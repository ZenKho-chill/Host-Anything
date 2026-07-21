// Copyright 2026 Host Anything Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/host-anything/hostanything/internal/config"
	"github.com/host-anything/hostanything/internal/logging"
	"github.com/host-anything/hostanything/internal/server"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/internal/crypto"
	"github.com/host-anything/hostanything/internal/runtime"
	"github.com/host-anything/hostanything/pkg/types"
)

// app holds all fully wired application components.
// It is created by [buildApp] and driven by [app.run].
type app struct {
	cfg    *types.SystemConfig
	logger *slog.Logger
	srv    *http.Server
}

// buildApp loads configuration, initializes all components in dependency order,
// and returns a ready-to-run application. It is the only place that wires
// dependencies together — no construction logic belongs in main.go.
func buildApp(configPath, version string) (*app, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("buildApp: load config: %w", err)
	}

	logger, err := logging.NewLogger(cfg.Server.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("buildApp: init logger: %w", err)
	}

	enabledRuntimes := collectEnabledRuntimes(cfg)

	reg, err := template.NewRegistry(cfg.Paths.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("buildApp: init template registry: %w", err)
	}

	masterKey, err := crypto.LoadOrCreateKey(cfg.Paths.DataDir)
	if err != nil {
		return nil, fmt.Errorf("buildApp: load master key: %w", err)
	}

	mgr := runtime.NewServiceManager(logger)
	// (Adapters would be registered here, e.g. docker.NewAdapter())

	srv, err := server.NewServer(server.Options{
		Config:          cfg,
		Logger:          logger,
		Version:         version,
		EnabledRuntimes: enabledRuntimes,
		Registry:        reg,
		Manager:         mgr,
		MasterKey:       masterKey,
	})
	if err != nil {
		return nil, fmt.Errorf("buildApp: create server: %w", err)
	}

	return &app{
		cfg:    cfg,
		logger: logger,
		srv:    srv,
	}, nil
}

// collectEnabledRuntimes returns the ordered list of runtime adapter names
// that are enabled in the system configuration.
func collectEnabledRuntimes(cfg *types.SystemConfig) []string {
	var runtimes []string
	if cfg.Runtimes.DockerEnabled {
		runtimes = append(runtimes, "docker")
	}
	if cfg.Runtimes.PodmanEnabled {
		runtimes = append(runtimes, "podman")
	}
	if cfg.Runtimes.K8sEnabled {
		runtimes = append(runtimes, "k8s")
	}
	if cfg.Runtimes.HostEnabled {
		runtimes = append(runtimes, "host")
	}
	return runtimes
}

// run starts the HTTP server and blocks until ctx is cancelled.
// On cancellation, it performs a graceful shutdown with a 15-second timeout.
func (a *app) run(ctx context.Context) error {
	const shutdownTimeout = 15 * time.Second

	listenErr := make(chan error, 1)

	go func() {
		a.logger.Info("hostanything started",
			"address", a.srv.Addr,
			"runtimes", collectEnabledRuntimes(a.cfg),
		)
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErr <- fmt.Errorf("app.run: listen: %w", err)
			return
		}
		close(listenErr)
	}()

	select {
	case err := <-listenErr:
		// Server exited before context was cancelled — propagate error.
		return err

	case <-ctx.Done():
		a.logger.Info("shutdown signal received, draining connections")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := a.srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("app.run: graceful shutdown: %w", err)
		}

		a.logger.Info("shutdown complete")
		return nil
	}
}

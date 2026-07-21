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
	"os"
	"time"

	"github.com/host-anything/hostanything/internal/config"
	"github.com/host-anything/hostanything/internal/crypto"
	"github.com/host-anything/hostanything/internal/logging"
	"github.com/host-anything/hostanything/internal/runtime"
	"github.com/host-anything/hostanything/internal/runtime/docker"
	"github.com/host-anything/hostanything/internal/runtime/host"
	"github.com/host-anything/hostanything/internal/runtime/kubernetes"
	"github.com/host-anything/hostanything/internal/runtime/podman"
	"github.com/host-anything/hostanything/internal/server"
	"github.com/host-anything/hostanything/internal/store"
	"github.com/host-anything/hostanything/internal/template"
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

	reg, err := template.NewRegistry(cfg.Paths.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("buildApp: init template registry: %w", err)
	}

	masterKey, err := crypto.LoadOrCreateKey(cfg.Paths.DataDir)
	if err != nil {
		return nil, fmt.Errorf("buildApp: load master key: %w", err)
	}

	db, err := store.Open(cfg.Paths.DataDir, logger)
	if err != nil {
		return nil, fmt.Errorf("buildApp: open db: %w", err)
	}

	// Seed admin if environment variables are provided during installation
	adminUser := os.Getenv("HA_ADMIN_USERNAME")
	adminPass := os.Getenv("HA_ADMIN_PASSWORD")
	if adminUser != "" && adminPass != "" {
		hash, err := crypto.HashPassword(adminPass)
		if err == nil {
			if err := db.SeedAdmin(context.Background(), adminUser, hash); err != nil {
				logger.Error("failed to seed admin user", "error", err)
			} else {
				logger.Info("admin user seeded successfully", "username", adminUser)
			}
		}
	}

	mgr := runtime.NewServiceManager(logger)

	var enabledRuntimes []string
	if cfg.Runtimes.DockerEnabled {
		if dAdapter, err := docker.NewAdapter(); err == nil {
			mgr.RegisterAdapter("docker", dAdapter)
			enabledRuntimes = append(enabledRuntimes, "docker")
		} else {
			logger.Warn("docker runtime disabled (auto-detect failed)", "error", err)
		}
	}
	if cfg.Runtimes.PodmanEnabled {
		if pAdapter, err := podman.NewAdapter(); err == nil {
			mgr.RegisterAdapter("podman", pAdapter)
			enabledRuntimes = append(enabledRuntimes, "podman")
		} else {
			logger.Warn("podman runtime disabled (auto-detect failed)", "error", err)
		}
	}
	if cfg.Runtimes.K8sEnabled {
		if kAdapter, err := kubernetes.NewAdapter(); err == nil {
			mgr.RegisterAdapter("k8s", kAdapter)
			enabledRuntimes = append(enabledRuntimes, "k8s")
		} else {
			logger.Warn("kubernetes runtime disabled (auto-detect failed)", "error", err)
		}
	}
	if cfg.Runtimes.HostEnabled {
		if hAdapter, err := host.NewAdapter(); err == nil {
			mgr.RegisterAdapter("host", hAdapter)
			enabledRuntimes = append(enabledRuntimes, "host")
		} else {
			logger.Warn("host runtime disabled (auto-detect failed)", "error", err)
		}
	}

	srv, err := server.NewServer(server.Options{
		Config:          cfg,
		Logger:          logger,
		Version:         version,
		EnabledRuntimes: enabledRuntimes,
		Registry:        reg,
		Manager:         mgr,
		MasterKey:       masterKey,
		DB:              db,
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

// collectEnabledRuntimes is no longer used, as we auto-detect during initialization.

// run starts the HTTP server and blocks until ctx is cancelled.
// On cancellation, it performs a graceful shutdown with a 15-second timeout.
func (a *app) run(ctx context.Context) error {
	const shutdownTimeout = 15 * time.Second

	listenErr := make(chan error, 1)

	go func() {
		a.logger.Info("hostanything started",
			"address", a.srv.Addr,
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

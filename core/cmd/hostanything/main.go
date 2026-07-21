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

/*
Package main provides the entry point for the hostanything binary.

It parses command-line flags, wires all dependencies (see wire.go),
and starts the API server with graceful shutdown on SIGINT/SIGTERM.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/host-anything/hostanything/internal/config"
)

// version is the binary version string. It is overridden at build time:
//
//	go build -ldflags "-X main.version=1.0.0"
var version = "dev"

func main() {
	fs := flag.NewFlagSet("hostanything", flag.ExitOnError)
	configPath := fs.String("config", config.DefaultConfigPath, "path to configuration file")
	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *showVersion {
		fmt.Printf("hostanything %s\n", version)
		os.Exit(0)
	}

	// Cancel context on SIGINT or SIGTERM, triggering graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := buildApp(*configPath, version)
	if err != nil {
		log.Fatalf("fatal: failed to initialize hostanything: %v", err)
	}

	if err := application.run(ctx); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

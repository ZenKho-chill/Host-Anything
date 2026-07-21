// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

/*
Package runtime defines the core interface for interacting with various
container and process runtimes. It uses an adapter pattern to support
different backends without tying the core logic to a specific runtime.

Sub-packages include:
  - docker: Adapter for Docker Engine
  - podman: Adapter for Podman
  - k8s: Adapter for Kubernetes
  - host: Adapter for raw host processes (systemd)
*/
package runtime

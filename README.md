# Host Anything

[![Build Status](https://img.shields.io/github/actions/workflow/status/hostanything/hostanything/build.yml?branch=main)](https://github.com/hostanything/hostanything/actions)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hostanything/hostanything)](https://golang.org/doc/devel/release.html)

**Deploy anything. Configure everything. Touch nothing else.**

## Overview

Host Anything is a unified hosting platform that allows you to manage the configuration of your services across various containerization tools and runtimes without interfering with the services themselves. 

We provide the control plane; you bring the runtimes.

## Key Features

- 📝 **Template-driven**: Manage configurations dynamically through flexible TOML templates.
- 🐳 **Multiple Runtimes**: Out-of-the-box support for Docker, Podman, Kubernetes, and bare-metal (via host mode).
- 🔐 **Web UI with Auth**: Built-in, fully authenticated web interface to manage your deployments visually.
- 🌍 **GitHub Marketplace Integration**: Explore and pull deployment templates right from a robust open-source ecosystem.
- 🛡️ **Fail2Ban Protection**: Deep integration with fail2ban to secure your exposed endpoints from malicious actors.
- 🐧 **Debian Native**: First-class support for Debian systems, distributed easily as `.deb` packages.

## Architecture

Host Anything employs an adapter pattern to communicate with different runtimes (Docker, Podman, K8s). It focuses exclusively on managing configurations (environment variables, ports, volumes, and resource limits), treating the underlying runtime logic as a black box.

## Quick Start

> **🚧 Coming Soon**
> 
> Detailed quick start instructions will be provided as we get closer to our v0.1.0 release.

## Roadmap

| Milestone | Focus | Expected |
| --------- | ----- | -------- |
| **M1** | Core architecture & Template engine | Q3 2026 |
| **M2** | Docker/Podman runtime adapters | Q4 2026 |
| **M3** | Basic Web UI & Auth | Q1 2027 |
| **M4** | Fail2Ban integration & Security | Q2 2027 |
| **M5** | K8s runtime adapter | Q3 2027 |
| **M6** | GitHub Marketplace Integration | Q4 2027 |

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started, set up your development environment, and submit pull requests.

## License

Host Anything is licensed under the [Apache License 2.0](LICENSE).

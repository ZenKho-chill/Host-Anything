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

// Package logging provides structured JSON logging for the hostanything daemon.
// It wraps the standard library's [log/slog] package with opinionated defaults
// suitable for a long-running system service on Debian.
//
// Logs are written in JSON format to stderr for machine-readability and
// fail2ban compatibility. The auth log format is defined in SPEC-030.
//
// Usage:
//
//	logger, err := logging.NewLogger("info")
//	if err != nil {
//	    log.Fatalf("failed to init logger: %v", err)
//	}
//	logger.Info("server started", "address", ":8080")
package logging

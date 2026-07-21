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

// Package types defines all shared interfaces and data types used across
// the hostanything codebase. It contains no business logic and must not
// import any internal/ packages.
//
// Key types:
//   - [RuntimeAdapter]: The central interface every runtime must implement.
//   - [ServiceSpec]: The normalized service description passed to adapters.
//   - [SystemConfig]: The parsed system-level configuration.
//   - [ServiceState]: Lifecycle states for deployed services (see SPEC-002).
package types

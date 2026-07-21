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

package types

// ServiceState represents the lifecycle state of a deployed service.
// Valid transitions are documented in SPEC-002.
type ServiceState string

// Service lifecycle states as defined in SPEC-002.
const (
	// ServiceStatePending means the service has been registered but not yet deployed.
	ServiceStatePending ServiceState = "PENDING"

	// ServiceStateDeploying means the service is currently being deployed.
	ServiceStateDeploying ServiceState = "DEPLOYING"

	// ServiceStateRunning means the service is active and healthy.
	ServiceStateRunning ServiceState = "RUNNING"

	// ServiceStateStopping means the service is gracefully shutting down.
	ServiceStateStopping ServiceState = "STOPPING"

	// ServiceStateStopped means the service has cleanly stopped.
	ServiceStateStopped ServiceState = "STOPPED"

	// ServiceStateError means the service encountered an unrecoverable error.
	ServiceStateError ServiceState = "ERROR"

	// ServiceStateUpdating means the service configuration is being applied.
	ServiceStateUpdating ServiceState = "UPDATING"
)

// ServiceStatus holds the current runtime status of a deployed service,
// as reported by a RuntimeAdapter.
type ServiceStatus struct {
	// State is the current lifecycle state.
	State ServiceState

	// UptimeSeconds is the number of seconds the service has been in RUNNING state.
	// Zero if the service is not currently running.
	UptimeSeconds int64

	// Message provides additional context, especially useful in ERROR state.
	Message string
}

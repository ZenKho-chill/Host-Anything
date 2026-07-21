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

package runtime

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/host-anything/hostanything/pkg/types"
)

const deployTimeout = 5 * time.Minute

// ServiceManager orchestrates the deployment lifecycle of services.
// It implements the finite state machine defined in SPEC-002 and acts
// as a safe proxy to the underlying RuntimeAdapters.
type ServiceManager struct {
	logger   *slog.Logger
	adapters map[string]types.RuntimeAdapter

	// In-memory state tracking for M3. In a real system (M4/M5),
	// this would be backed by the database.
	mu     sync.RWMutex
	states map[string]types.ServiceState
}

// NewServiceManager initializes the lifecycle manager.
func NewServiceManager(logger *slog.Logger) *ServiceManager {
	return &ServiceManager{
		logger:   logger,
		adapters: make(map[string]types.RuntimeAdapter),
		states:   make(map[string]types.ServiceState),
	}
}

// RegisterAdapter binds a runtime name (e.g., "docker") to an adapter implementation.
func (sm *ServiceManager) RegisterAdapter(name string, adapter types.RuntimeAdapter) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.adapters[name] = adapter
}

func (sm *ServiceManager) setState(serviceID string, state types.ServiceState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[serviceID] = state
	sm.logger.Info("service state changed", "service", serviceID, "state", state)
}

// GetState returns the current lifecycle state of a service.
func (sm *ServiceManager) GetState(serviceID string) types.ServiceState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if state, exists := sm.states[serviceID]; exists {
		return state
	}
	return "" // Unknown/Not Found
}

// DeployService executes the PENDING -> DEPLOYING -> RUNNING/ERROR transition.
func (sm *ServiceManager) DeployService(spec *types.ServiceSpec) error {
	sm.mu.RLock()
	adapter, exists := sm.adapters[spec.Template.Runtime.Preferred]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("manager.DeployService: adapter %q not registered", spec.Template.Runtime.Preferred)
	}

	sm.setState(spec.ServiceID, types.ServiceStateDeploying)

	ctx, cancel := context.WithTimeout(context.Background(), deployTimeout)
	defer cancel()

	err := adapter.Deploy(ctx, spec)
	if err != nil {
		sm.setState(spec.ServiceID, types.ServiceStateError)
		return fmt.Errorf("manager.DeployService: deploy failed: %w", err)
	}

	sm.setState(spec.ServiceID, types.ServiceStateRunning)
	return nil
}

// StopService executes the RUNNING -> STOPPING -> STOPPED/ERROR transition.
func (sm *ServiceManager) StopService(serviceID, runtimeName string) error {
	sm.mu.RLock()
	adapter, exists := sm.adapters[runtimeName]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("manager.StopService: adapter %q not registered", runtimeName)
	}

	// Simple transition guard (M3 simplified)
	if sm.GetState(serviceID) != types.ServiceStateRunning && sm.GetState(serviceID) != types.ServiceStateError {
		return fmt.Errorf("manager.StopService: service %q is not running or in error", serviceID)
	}

	sm.setState(serviceID, types.ServiceStateStopping)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := adapter.Stop(ctx, serviceID); err != nil {
		sm.setState(serviceID, types.ServiceStateError)
		return fmt.Errorf("manager.StopService: %w", err)
	}

	sm.setState(serviceID, types.ServiceStateStopped)
	return nil
}

// ListServices returns a list of all currently tracked services and their states.
func (sm *ServiceManager) ListServices() map[string]types.ServiceState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	copy := make(map[string]types.ServiceState, len(sm.states))
	for k, v := range sm.states {
		copy[k] = v
	}
	return copy
}

// LogsService retrieves the log stream for a service.
func (sm *ServiceManager) LogsService(ctx context.Context, serviceID, runtimeName string) (io.ReadCloser, error) {
	sm.mu.RLock()
	adapter, exists := sm.adapters[runtimeName]
	sm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("manager.LogsService: adapter %q not registered", runtimeName)
	}

	return adapter.Logs(ctx, serviceID)
}

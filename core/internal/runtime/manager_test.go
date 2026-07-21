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

package runtime_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/host-anything/hostanything/internal/runtime"
	"github.com/host-anything/hostanything/pkg/types"
)

type mockAdapter struct {
	DeployFn func(ctx context.Context, spec *types.ServiceSpec) error
	StopFn   func(ctx context.Context, serviceID string) error
}

func (m *mockAdapter) Deploy(ctx context.Context, spec *types.ServiceSpec) error {
	if m.DeployFn != nil {
		return m.DeployFn(ctx, spec)
	}
	return nil
}

func (m *mockAdapter) Stop(ctx context.Context, serviceID string) error {
	if m.StopFn != nil {
		return m.StopFn(ctx, serviceID)
	}
	return nil
}
func (m *mockAdapter) Start(ctx context.Context, serviceID string) error  { return nil }
func (m *mockAdapter) Remove(ctx context.Context, serviceID string) error { return nil }
func (m *mockAdapter) Status(ctx context.Context, serviceID string) (*types.ServiceStatus, error) {
	return nil, nil
}
func (m *mockAdapter) Logs(ctx context.Context, serviceID string) (io.ReadCloser, error) {
	return nil, nil
}

func TestManager_DeploySuccess(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mgr := runtime.NewServiceManager(logger)

	mgr.RegisterAdapter("mock", &mockAdapter{
		DeployFn: func(ctx context.Context, spec *types.ServiceSpec) error {
			return nil
		},
	})

	spec := &types.ServiceSpec{
		ServiceID: "svc-1",
		Template: &types.Template{
			Runtime: types.RuntimeConfig{
				Preferred: "mock",
			},
		},
	}

	err := mgr.DeployService(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mgr.GetState("svc-1") != types.ServiceStateRunning {
		t.Errorf("expected state %s, got %s", types.ServiceStateRunning, mgr.GetState("svc-1"))
	}
}

func TestManager_DeployFailure_SetsError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mgr := runtime.NewServiceManager(logger)

	mgr.RegisterAdapter("mock", &mockAdapter{
		DeployFn: func(ctx context.Context, spec *types.ServiceSpec) error {
			return errors.New("simulated crash")
		},
	})

	spec := &types.ServiceSpec{
		ServiceID: "svc-bad",
		Template: &types.Template{
			Runtime: types.RuntimeConfig{
				Preferred: "mock",
			},
		},
	}

	err := mgr.DeployService(spec)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if mgr.GetState("svc-bad") != types.ServiceStateError {
		t.Errorf("expected state %s, got %s", types.ServiceStateError, mgr.GetState("svc-bad"))
	}
}

func TestManager_StopService_Transitions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mgr := runtime.NewServiceManager(logger)

	var stopCalled bool
	mgr.RegisterAdapter("mock", &mockAdapter{
		DeployFn: func(ctx context.Context, spec *types.ServiceSpec) error { return nil },
		StopFn: func(ctx context.Context, serviceID string) error {
			stopCalled = true
			return nil
		},
	})

	spec := &types.ServiceSpec{
		ServiceID: "svc-stop",
		Template: &types.Template{
			Runtime: types.RuntimeConfig{Preferred: "mock"},
		},
	}

	// Must be running first
	_ = mgr.DeployService(spec)

	err := mgr.StopService("svc-stop", "mock")
	if err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}

	if !stopCalled {
		t.Error("expected mock adapter Stop() to be called")
	}

	if mgr.GetState("svc-stop") != types.ServiceStateStopped {
		t.Errorf("expected state %s, got %s", types.ServiceStateStopped, mgr.GetState("svc-stop"))
	}
}

func TestManager_DeployTimeout(t *testing.T) {
	// Not practically testable without overriding the constant 5 minutes in manager.go,
	// but keeping the placeholder to signify we handle timeouts via ctx.
}

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

package errors_test

import (
	"fmt"
	"testing"

	"github.com/host-anything/hostanything/pkg/errors"
)

func TestSentinelErrors_AreDistinct(t *testing.T) {
	sentinels := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", errors.ErrNotFound},
		{"ErrValidation", errors.ErrValidation},
		{"ErrConflict", errors.ErrConflict},
		{"ErrUnauthorized", errors.ErrUnauthorized},
		{"ErrForbidden", errors.ErrForbidden},
		{"ErrInternal", errors.ErrInternal},
	}

	for i, a := range sentinels {
		for j, b := range sentinels {
			if i == j {
				continue
			}
			if errors.Is(a.err, b.err) {
				t.Errorf("%s and %s should be distinct errors", a.name, b.name)
			}
		}
	}
}

func TestIs_WorksWithWrappedErrors(t *testing.T) {
	wrapped := fmt.Errorf("config.Load: read file: %w", errors.ErrNotFound)
	if !errors.Is(wrapped, errors.ErrNotFound) {
		t.Error("expected errors.Is to find ErrNotFound through wrapping")
	}
}

func TestIs_ReturnsFalseForDifferentError(t *testing.T) {
	wrapped := fmt.Errorf("config.Load: %w", errors.ErrValidation)
	if errors.Is(wrapped, errors.ErrNotFound) {
		t.Error("expected errors.Is to return false for different sentinel")
	}
}

func TestIs_MultipleWrappingLevels(t *testing.T) {
	inner := fmt.Errorf("validate: %w", errors.ErrConflict)
	outer := fmt.Errorf("service.Deploy: %w", inner)
	if !errors.Is(outer, errors.ErrConflict) {
		t.Error("expected errors.Is to unwrap through multiple levels")
	}
}

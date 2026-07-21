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

package logging_test

import (
	"testing"

	"github.com/host-anything/hostanything/internal/logging"
	"github.com/host-anything/hostanything/pkg/errors"
)

func TestNewLogger_AllValidLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			logger, err := logging.NewLogger(level)
			if err != nil {
				t.Errorf("expected no error for level %q, got: %v", level, err)
			}
			if logger == nil {
				t.Error("expected non-nil logger")
			}
		})
	}
}

func TestNewLogger_InvalidLevel_ReturnsValidationError(t *testing.T) {
	_, err := logging.NewLogger("verbose")
	if !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for invalid level, got: %v", err)
	}
}

func TestNewLogger_EmptyLevel_ReturnsValidationError(t *testing.T) {
	_, err := logging.NewLogger("")
	if !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for empty level, got: %v", err)
	}
}

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

package logging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/host-anything/hostanything/pkg/errors"
)

// levelMap maps string log level names to their [slog.Level] equivalent.
var levelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// NewLogger constructs a structured JSON logger writing to stderr at the given level.
//
// The level parameter must be one of: "debug", "info", "warn", "error".
// JSON format is used to ensure compatibility with log aggregators and fail2ban
// (see SPEC-030 for the required auth-event log schema).
//
// Returns [errors.ErrValidation] if the level string is not recognized.
func NewLogger(level string) (*slog.Logger, error) {
	slogLevel, ok := levelMap[level]
	if !ok {
		return nil, fmt.Errorf(
			"logging.NewLogger: unknown log level %q, must be one of: debug, info, warn, error: %w",
			level, errors.ErrValidation,
		)
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: slogLevel == slog.LevelDebug, // include source location in debug mode
	})

	return slog.New(handler), nil
}

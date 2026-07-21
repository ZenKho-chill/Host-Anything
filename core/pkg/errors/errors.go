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

// Package errors defines sentinel error types used throughout the
// hostanything codebase. All errors returned from internal packages
// must be wrapped with context using fmt.Errorf("package.Func: %w", err).
//
// Callers should use errors.Is() to test for specific error conditions
// rather than comparing error values directly.
package errors

import "errors"

// Sentinel errors for the hostanything application.
// These are intended to be tested with errors.Is() on wrapped errors.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("not found")

	// ErrValidation is returned when input fails validation checks.
	ErrValidation = errors.New("validation failed")

	// ErrConflict is returned when an operation conflicts with current state.
	ErrConflict = errors.New("conflict")

	// ErrUnauthorized is returned when a request lacks valid authentication.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a request lacks sufficient permissions.
	ErrForbidden = errors.New("forbidden")

	// ErrInternal is returned for unexpected internal failures.
	ErrInternal = errors.New("internal error")
)

// Is reports whether any error in err's chain matches target.
// Re-exported from the standard library for convenience.
var Is = errors.Is

// As finds the first error in err's chain that matches target.
// Re-exported from the standard library for convenience.
var As = errors.As

// Unwrap returns the wrapped error from err, if any.
// Re-exported from the standard library for convenience.
var Unwrap = errors.Unwrap

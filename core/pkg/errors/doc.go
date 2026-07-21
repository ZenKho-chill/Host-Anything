// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

/*
Package errors provides standard sentinel errors and error wrapping
conventions for the hostanything project.

Errors should be wrapped with context using fmt.Errorf with the %w verb
to preserve the original error trace.
*/
package errors

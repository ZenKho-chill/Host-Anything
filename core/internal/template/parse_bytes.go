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

package template

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/host-anything/hostanything/pkg/types"
)

// ParseBytes reads template TOML from a byte slice, decodes it into a
// types.Template struct, and performs schema validation according to SPEC-001.
// Unknown fields are strictly rejected.
func ParseBytes(data []byte) (*types.Template, error) {
	var t types.Template

	// Enforce strict unmarshalling — unknown keys in the TOML return an error.
	meta, err := toml.DecodeReader(bytes.NewReader(data), &t)
	if err != nil {
		return nil, fmt.Errorf("template.ParseBytes: decode TOML: %w", err)
	}
	if len(meta.Undecoded()) > 0 {
		return nil, fmt.Errorf("template.ParseBytes: unknown fields in TOML schema: %v", meta.Undecoded())
	}

	// Validate the decoded struct
	if err := Validate(&t); err != nil {
		return nil, fmt.Errorf("template.ParseBytes: validation failed: %w", err)
	}

	return &t, nil
}



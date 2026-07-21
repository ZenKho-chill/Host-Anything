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

package template_test

import (
	"strings"
	"testing"

	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

func TestValidateUserInput(t *testing.T) {
	tmpl := &types.Template{
		Config: []types.ConfigVar{
			{Name: "REQ_STR", Type: template.ConfigTypeString, Required: true},
			{Name: "DEF_STR", Type: template.ConfigTypeString, Required: true, Default: "default"},
			{Name: "INT_VAL", Type: template.ConfigTypeInt, Required: false},
			{Name: "BOOL_VAL", Type: template.ConfigTypeBoolean, Required: false},
			{Name: "ENUM_VAL", Type: template.ConfigTypeEnum, Required: false, Options: []string{"A", "B"}},
			{Name: "REGEX_VAL", Type: template.ConfigTypeString, Required: false, ValidationRegex: "^[A-Z]+$"},
		},
	}

	tests := []struct {
		name    string
		input   map[string]string
		wantErr string
	}{
		{
			name: "valid input all types",
			input: map[string]string{
				"REQ_STR":   "val",
				"INT_VAL":   "42",
				"BOOL_VAL":  "True",
				"ENUM_VAL":  "B",
				"REGEX_VAL": "VALID",
			},
			wantErr: "",
		},
		{
			name:  "missing required without default",
			input: map[string]string{
				// REQ_STR missing
			},
			wantErr: "missing required variable",
		},
		{
			name: "missing required with default is ok",
			input: map[string]string{
				"REQ_STR": "val", // satisfied
				// DEF_STR missing, but has default
			},
			wantErr: "",
		},
		{
			name: "invalid int",
			input: map[string]string{
				"REQ_STR": "val",
				"INT_VAL": "abc",
			},
			wantErr: "must be an integer",
		},
		{
			name: "invalid bool",
			input: map[string]string{
				"REQ_STR":  "val",
				"BOOL_VAL": "yes",
			},
			wantErr: "must be a boolean",
		},
		{
			name: "invalid enum",
			input: map[string]string{
				"REQ_STR":  "val",
				"ENUM_VAL": "C",
			},
			wantErr: "must be one of",
		},
		{
			name: "regex mismatch",
			input: map[string]string{
				"REQ_STR":   "val",
				"REGEX_VAL": "invalid123",
			},
			wantErr: "does not match regex",
		},
		{
			name: "unknown variable",
			input: map[string]string{
				"REQ_STR": "val",
				"HACKER":  "true",
			},
			wantErr: "unknown variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := template.ValidateUserInput(tmpl, tt.input)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
				}
			}
		})
	}
}

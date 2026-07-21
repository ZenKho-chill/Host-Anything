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
)

func TestParseBytes_Valid(t *testing.T) {
	tomlData := `
[meta]
name = "redis"
version = "1.0.0"
description = "In-memory cache"
author = "Host Anything"
license = "MIT"
tags = ["cache"]

[requirements]
min_memory = "256MB"
min_cpu = 0.5
disk_space = "1GB"

[[config]]
name = "REDIS_PASSWORD"
type = "secret"
required = true
validation_regex = "^.{8,}$"

[runtime]
preferred = "docker"
supported = ["docker"]
image = "redis:7"
command = ["redis-server", "--requirepass", "${REDIS_PASSWORD}"]

[[volumes]]
name = "data"
mount_path = "/data"

[[network]]
internal_port = 6379

[healthcheck]
command = "redis-cli ping"

[update]
strategy = "recreate"
`

	tmpl, err := template.ParseBytes([]byte(tomlData))
	if err != nil {
		t.Fatalf("unexpected error parsing valid template: %v", err)
	}

	if tmpl.Meta.Name != "redis" {
		t.Errorf("expected name 'redis', got %q", tmpl.Meta.Name)
	}
	if len(tmpl.Config) != 1 {
		t.Fatalf("expected 1 config var, got %d", len(tmpl.Config))
	}
}

func TestParseBytes_UnknownField(t *testing.T) {
	tomlData := `
[meta]
name = "redis"
version = "1.0.0"
description = "desc"
author = "author"
license = "MIT"
unknown_field = "should-fail"

[runtime]
supported = ["docker"]
image = "redis:7"
`
	_, err := template.ParseBytes([]byte(tomlData))
	if err == nil {
		t.Fatal("expected error when unknown field is present")
	}
	if !strings.Contains(err.Error(), "unknown fields") {
		t.Errorf("expected error to mention unknown fields, got: %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	tomlData := `
[meta]
# missing name
version = "1.0.0"
description = "desc"
author = "author"
license = "MIT"

[runtime]
supported = ["docker"]
image = "redis:7"
`
	_, err := template.ParseBytes([]byte(tomlData))
	if err == nil {
		t.Fatal("expected error due to missing required field")
	}
	if !strings.Contains(err.Error(), "missing required field 'name'") {
		t.Errorf("expected missing name error, got: %v", err)
	}
}

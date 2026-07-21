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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/host-anything/hostanything/pkg/types"
)

// Registry manages the local template store.
// Templates are expected to be in: {baseDir}/{name}/{version}/template.toml
type Registry struct {
	baseDir string
}

// NewRegistry creates a new registry targeting the given directory.
// It creates the directory if it does not exist.
func NewRegistry(baseDir string) (*Registry, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("template.NewRegistry: %w", err)
	}
	return &Registry{baseDir: baseDir}, nil
}

// List scans the registry and returns summaries of all installed templates.
func (r *Registry) List() ([]types.TemplateSummary, error) {
	var summaries []types.TemplateSummary

	// Walk the base directory to find template.toml files
	err := filepath.WalkDir(r.baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip unreadable paths
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "template.toml" {
			t, err := ParseTOML(path)
			if err != nil {
				// We log or ignore invalid templates in the registry during List
				return nil
			}
			summaries = append(summaries, types.TemplateSummary{
				Name:        t.Meta.Name,
				Version:     t.Meta.Version,
				Description: t.Meta.Description,
				Author:      t.Meta.Author,
				Tags:        t.Meta.Tags,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("template.Registry.List: %w", err)
	}

	return summaries, nil
}

// Get retrieves a specific template. If version is "latest" or empty,
// it returns the highest semver available (currently simplified to any available).
func (r *Registry) Get(name, version string) (*types.Template, error) {
	nameDir := filepath.Join(r.baseDir, name)

	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(nameDir), filepath.Clean(r.baseDir)) {
		return nil, fmt.Errorf("template.Registry.Get: invalid name %q", name)
	}

	// For M2, we simplify "latest" by just reading the first version dir we find.
	// A robust SemVer sorter is needed for a full implementation.
	targetVersion := version
	if targetVersion == "" || targetVersion == "latest" {
		entries, err := os.ReadDir(nameDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("template %q not found", name)
			}
			return nil, fmt.Errorf("template.Registry.Get: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				targetVersion = e.Name()
				break
			}
		}
	}

	if targetVersion == "" {
		return nil, fmt.Errorf("template %q not found", name)
	}

	targetPath := filepath.Join(nameDir, targetVersion, "template.toml")
	t, err := ParseTOML(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template %q version %q not found", name, targetVersion)
		}
		return nil, fmt.Errorf("template.Registry.Get: %w", err)
	}

	return t, nil
}

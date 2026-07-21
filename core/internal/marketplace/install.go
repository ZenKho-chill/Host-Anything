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

package marketplace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/host-anything/hostanything/pkg/types"
)

// Installer handles downloading and registering a remote template into the
// local template registry directory.
type Installer struct {
	client      *Client
	templateDir string
}

// NewInstaller creates a new Installer.
// templateDir is the local directory where templates are stored.
func NewInstaller(templateDir string) *Installer {
	return &Installer{
		client:      NewClient(),
		templateDir: templateDir,
	}
}

// Install fetches a template from the given GitHub repository, validates it,
// and writes it to the local template registry.
// The template is saved to {templateDir}/{name}/{version}/template.toml.
func (i *Installer) Install(ctx context.Context, owner, repo string) (*types.Template, error) {
	tmpl, err := i.client.FetchTemplate(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("marketplace.Installer.Install: fetch: %w", err)
	}

	if err := validateTemplate(tmpl); err != nil {
		return nil, fmt.Errorf("marketplace.Installer.Install: validation: %w", err)
	}

	destDir := filepath.Join(i.templateDir, tmpl.Meta.Name, tmpl.Meta.Version)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("marketplace.Installer.Install: create dir %q: %w", destDir, err)
	}

	destPath := filepath.Join(destDir, TemplateFileName)
	f, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("marketplace.Installer.Install: create file: %w", err)
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	if err := enc.Encode(tmpl); err != nil {
		return nil, fmt.Errorf("marketplace.Installer.Install: encode toml: %w", err)
	}

	return tmpl, nil
}

// validateTemplate performs basic structural validation on a remotely-fetched
// template before it is written to the local registry.
func validateTemplate(tmpl *types.Template) error {
	if tmpl.Meta.Name == "" {
		return fmt.Errorf("template.meta.name must not be empty")
	}
	if tmpl.Meta.Version == "" {
		return fmt.Errorf("template.meta.version must not be empty")
	}
	if tmpl.Runtime.Image == "" {
		return fmt.Errorf("template.runtime.image must not be empty")
	}
	return nil
}

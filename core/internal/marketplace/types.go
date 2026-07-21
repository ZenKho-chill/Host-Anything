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

// OfficialOrg is the GitHub organisation that owns all official templates.
// Repositories owned by this org are marked as IsOfficial = true.
const OfficialOrg = "host-anything"

// TemplateTopic is the GitHub topic used to discover Host Anything templates.
const TemplateTopic = "hostanything-template"

// DefaultBranch is the first branch tried when fetching template.toml.
const DefaultBranch = "main"

// FallbackBranch is tried if template.toml is not found on DefaultBranch.
const FallbackBranch = "master"

// TemplateFileName is the expected filename within the repository root.
const TemplateFileName = "template.toml"

// MarketplaceResult represents a single search result from the GitHub marketplace.
type MarketplaceResult struct {
	// Name is the repository name (e.g. "hostanything-template-redis").
	Name string `json:"name"`
	// Owner is the GitHub username or organisation (e.g. "host-anything").
	Owner string `json:"owner"`
	// Description is the repository's short description.
	Description string `json:"description"`
	// Stars is the number of GitHub stars.
	Stars int `json:"stars"`
	// RepoURL is the full HTML URL to the GitHub repository.
	RepoURL string `json:"repo_url"`
	// IsOfficial is true when Owner == OfficialOrg.
	IsOfficial bool `json:"is_official"`
}

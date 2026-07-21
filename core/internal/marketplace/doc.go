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

// Package marketplace provides integration with the GitHub-based Host Anything
// template marketplace. It allows users to search for community and official
// service templates, preview their metadata, and install them into the local
// template registry.
//
// Templates are discovered via the GitHub Search API by the topic
// "hostanything-template". Official templates are published under the
// "host-anything" GitHub organisation and receive a trust badge in the UI.
//
// The optional GITHUB_TOKEN environment variable can be set to authenticate
// requests and raise the API rate limit from 60 to 5,000 requests per hour.
package marketplace

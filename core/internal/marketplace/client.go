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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/host-anything/hostanything/pkg/types"
)

// Client handles all communication with GitHub APIs for the marketplace.
type Client struct {
	httpClient *http.Client
	token      string
}

// githubSearchResponse mirrors the relevant fields from GitHub's repository
// search API response.
type githubSearchResponse struct {
	Items []githubRepo `json:"items"`
}

// githubRepo mirrors the relevant fields of a repository object in GitHub API responses.
type githubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StarCount   int    `json:"stargazers_count"`
	HTMLURL     string `json:"html_url"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// NewClient creates a new marketplace Client.
// If the GITHUB_TOKEN environment variable is set, it is used to authenticate
// requests and raise the rate limit from 60 to 5,000 requests per hour.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		token:      os.Getenv("GITHUB_TOKEN"),
	}
}

// SearchTemplates queries the GitHub Search API for repositories matching the
// given query and the "hostanything-template" topic.
// Returns a slice of MarketplaceResult ordered by GitHub's relevance ranking.
func (c *Client) SearchTemplates(ctx context.Context, query string) ([]MarketplaceResult, error) {
	q := fmt.Sprintf("topic:%s %s", TemplateTopic, query)
	apiURL := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=30",
		url.QueryEscape(q),
	)

	resp, err := c.doRequest(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("marketplace.Client.SearchTemplates: %w", err)
	}
	defer resp.Body.Close()

	var searchResp githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("marketplace.Client.SearchTemplates: decode response: %w", err)
	}

	results := make([]MarketplaceResult, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		results = append(results, MarketplaceResult{
			Name:        item.Name,
			Owner:       item.Owner.Login,
			Description: item.Description,
			Stars:       item.StarCount,
			RepoURL:     item.HTMLURL,
			IsOfficial:  item.Owner.Login == OfficialOrg,
		})
	}

	return results, nil
}

// FetchTemplate fetches and parses the template.toml file from the given
// GitHub repository. It tries the main branch first, then master.
func (c *Client) FetchTemplate(ctx context.Context, owner, repo string) (*types.Template, error) {
	var (
		rawURL string
		body   []byte
		err    error
	)

	for _, branch := range []string{DefaultBranch, FallbackBranch} {
		rawURL = fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/%s/%s",
			owner, repo, branch, TemplateFileName,
		)
		body, err = c.fetchRaw(ctx, rawURL)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("marketplace.Client.FetchTemplate: fetch %s/%s: %w", owner, repo, err)
	}

	var tmpl types.Template
	if _, err := toml.Decode(string(body), &tmpl); err != nil {
		return nil, fmt.Errorf("marketplace.Client.FetchTemplate: parse toml: %w", err)
	}

	return &tmpl, nil
}

// fetchRaw downloads the raw content at the given URL and returns it as bytes.
// Returns an error if the status code is not 200.
func (c *Client) fetchRaw(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return data, nil
}

// doRequest performs an authenticated GET to the GitHub API.
// It returns a non-nil response only when the status code is 200.
func (c *Client) doRequest(ctx context.Context, apiURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GitHub API returned status %d for %s", resp.StatusCode, apiURL)
	}

	return resp, nil
}

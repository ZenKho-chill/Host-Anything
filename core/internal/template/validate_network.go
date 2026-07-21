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

	"github.com/host-anything/hostanything/pkg/types"
)

func validateNetwork(nets []types.NetworkConfig) error {
	for i, n := range nets {
		if n.InternalPort <= 0 || n.InternalPort > 65535 {
			return fmt.Errorf("network at index %d has invalid internal_port %d (must be 1-65535)", i, n.InternalPort)
		}
		if n.ExternalPort < 0 || n.ExternalPort > 65535 {
			return fmt.Errorf("network at index %d has invalid external_port %d (must be 0-65535)", i, n.ExternalPort)
		}
		if n.Protocol != "" {
			switch n.Protocol {
			case ProtocolTCP, ProtocolUDP, ProtocolHTTP:
				// valid
			default:
				return fmt.Errorf("network at index %d has invalid protocol %q", i, n.Protocol)
			}
		}
	}
	return nil
}

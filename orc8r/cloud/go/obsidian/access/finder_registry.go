/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package access

import (
	"strings"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/lib/go/protos"
)

const (
	MAGMA_ROOT_PART     = obsidian.RestRoot + obsidian.UrlSep
	MAGMA_ROOT_PART_LEN = len(MAGMA_ROOT_PART)
)

// RequestIdentityFinder Identity finder type
type RequestIdentityFinder func(c echo.Context) []*protos.Identity

// GetIdentityFinder returns an Identity finder for the request
func GetIdentityFinder(c echo.Context) RequestIdentityFinder {
	if c != nil {
		path := c.Path()
		if strings.HasPrefix(path, MAGMA_ROOT_PART) {
			parts := strings.Split(path[MAGMA_ROOT_PART_LEN:], obsidian.UrlSep)
			if len(parts) > 0 {
				p := parts[0]
				registry, ok := finderRegistries[p]
				if ok && len(parts) > 1 {
					p = parts[1]
				} else {
					// fall back to "versionless" V0
					registry, ok = finderRegistries[obsidian.V0]
				}
				if ok {
					fr, ok := registry.finderMap[p]
					if ok {
						return fr
					}
					return registry.defaultFinder
				}
			}
		}
	}
	return nil
}

// finderRegistries declares Versioned API Identity finders,
// add an entry for every new API Version
var finderRegistries = map[string]finderRegistryType{
	obsidian.V0: makeFinderRegistry(obsidian.V0),
	obsidian.V1: makeFinderRegistry(obsidian.V1),
}

type finderMap map[string]RequestIdentityFinder
type finderRegistryType struct {
	finderMap
	defaultFinder RequestIdentityFinder
}

func makeFinderRegistry(version string) finderRegistryType {
	magmaRoot := makeVersionedRoot(version, "")
	networkRoot := makeVersionedRoot(version, obsidian.MagmaNetworksUrlPart)
	operatorRoot := makeVersionedRoot(version, obsidian.MagmaOperatorsUrlPart)
	return finderRegistryType{
		finderMap: finderMap{
			obsidian.MagmaNetworksUrlPart:  func(c echo.Context) []*protos.Identity { return getNetworkIdentity(c, networkRoot) },
			obsidian.MagmaOperatorsUrlPart: func(c echo.Context) []*protos.Identity { return getOperatorIdentity(c, operatorRoot) },
		},
		defaultFinder: func(c echo.Context) []*protos.Identity { return getDefaultNetworkIdentity(c, magmaRoot) },
	}
}

func makeVersionedRoot(version, part string) string {
	if len(version) > 0 {
		return obsidian.RestRoot + obsidian.UrlSep + version + obsidian.UrlSep + part
	} else {
		return obsidian.RestRoot + obsidian.UrlSep + part
	}
}

// Network Identity Finder
func getNetworkIdentity(c echo.Context, networkRoot string) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), networkRoot) {
		nid, err := obsidian.GetNetworkId(c)
		if err == nil && len(nid) > 0 {
			// All checks pass - return a Network Identity
			return []*protos.Identity{identity.NewNetwork(nid)}
		}
		// No network ID -> requires wildcard access
		return []*protos.Identity{identity.NewNetworkWildcard()}
	}
	// We don't really know what resource is being requested - request all wildcards
	return SupervisorWildcards()
}

// Default Network Identity Finder, similar to getNetworkIdentity(), but returns SupervisorWildcards if :network_id
// is not found. To be used for default finders where we cannot be sure if the request is actually network scoped
func getDefaultNetworkIdentity(c echo.Context, versionRoot string) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), versionRoot) {
		if nid, err := obsidian.GetNetworkId(c); err == nil && len(nid) > 0 {
			return []*protos.Identity{identity.NewNetwork(nid)}
		}
	}
	return SupervisorWildcards()
}

// Operator Identity Finder
func getOperatorIdentity(c echo.Context, identityRoot string) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), identityRoot) {
		oid, err := obsidian.GetOperatorId(c)
		if err == nil && len(oid) > 0 {
			// All checks pass - return a Network Identity
			return []*protos.Identity{identity.NewOperator(oid)}
		}
		// No network ID -> requires wildcard access
		return []*protos.Identity{identity.NewOperatorWildcard()}
	}
	// We don't really know what resource is being requested - request all wildcards
	return SupervisorWildcards()
}

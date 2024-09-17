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

// Package registry defines API to register and fing EAP providers
package registry

import (
	"sync"

	"github.com/golang/glog"

	"magma/feg/gateway/services/eap/providers"
)

// registry is a "read mostly" map of eap providers
// to register provider, add it into init()

var (
	eapProviderRegistry               = map[uint8]providers.Method{}
	supportedTypes                    = []uint8{}
	registryMu          *sync.RWMutex = new(sync.RWMutex)
)

// Register adds (registers) the provider to the internal registry, if a provider for the same type is already
// registered it'll be overwritten.
// Register returns the previously registered provider for the type or nil if none was registered for the type before
func Register(p providers.Method) (oldProvider providers.Method) {
	typ := p.EAPType()
	registryMu.Lock()
	defer registryMu.Unlock()
	oldProvider, previousExists := eapProviderRegistry[typ]
	eapProviderRegistry[typ] = p
	if previousExists {
		glog.Errorf(
			"EAP Provider is already registered for type %d: %s. Will overwrite with: %s",
			typ, oldProvider, p)
	} else {
		supportedTypes = append(supportedTypes, typ)
	}
	return
}

// GetProvider returns registered Method provider for EAP type
func GetProvider(typ uint8) providers.Method {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, found := eapProviderRegistry[typ]
	if found {
		return p
	}
	return nil
}

// SupportedTypes returns sorted list (ascending, by type) of registered EAP Providers
// SupportedTypes makes copy of an internally maintained supported types list, so callers
// are advised to save the result locally and re-use it if needed
func SupportedTypes() []uint8 {
	registryMu.RLock()
	defer registryMu.RUnlock()
	res := make([]uint8, len(supportedTypes))
	copy(res, supportedTypes)
	return res
}

// Sort interface
type typesSlice []uint8

func (p typesSlice) Len() int           { return len(p) }
func (p typesSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p typesSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

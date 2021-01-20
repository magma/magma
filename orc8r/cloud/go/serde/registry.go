/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package serde

import (
	"fmt"

	"github.com/pkg/errors"
)

// Registry provides a serde registry.
type Registry interface {
	// GetSerde returns the serde registered for the passed type.
	// Returns an err iff not found.
	GetSerde(typ string) (Serde, error)

	// GetMap returns all registered serdes, keyed by their type.
	GetMap() map[string]Serde

	// MustMerge returns the union of two serde registries.
	// Panics if a serde type is found in both registries.
	MustMerge(rr Registry) Registry
}

func HasSerde(r Registry, typ string) bool {
	_, err := r.GetSerde(typ)
	return err == nil
}

type registry map[string]Serde

func NewRegistry(serdes ...Serde) Registry {
	r := registry{}
	for _, s := range serdes {
		r[s.GetType()] = s
	}
	return r
}

func (r registry) MustMerge(rr Registry) Registry {
	ret := registry{}
	for _, s := range r {
		ret[s.GetType()] = s
	}
	for _, s := range rr.GetMap() {
		if _, ok := ret[s.GetType()]; ok {
			panic(fmt.Sprintf("cannot merge serde registries when both contain serde for type %s", s.GetType()))
		}
		ret[s.GetType()] = s
	}
	return ret
}

func (r registry) GetSerde(typ string) (Serde, error) {
	serde, ok := r[typ]
	if !ok {
		return nil, errors.Errorf("no serde in registry for type %s", typ)
	}
	return serde, nil
}

func (r registry) GetMap() map[string]Serde {
	// Return copy
	ret := registry{}
	for k, v := range r {
		ret[k] = v
	}
	return ret
}

// Serialize an object by delegating to the passed serde registry.
func Serialize(data interface{}, typ string, registry Registry) ([]byte, error) {
	serde, err := registry.GetSerde(typ)
	if err != nil {
		return nil, err
	}
	return serde.Serialize(data)
}

// Deserialize a byte array by delegating to the passed serde registry.
func Deserialize(data []byte, typ string, registry Registry) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}
	serde, err := registry.GetSerde(typ)
	if err != nil {
		return nil, err
	}
	return serde.Deserialize(data)
}

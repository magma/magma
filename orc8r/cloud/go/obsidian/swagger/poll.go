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

package swagger

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sync"

	"magma/orc8r/cloud/go/obsidian/swagger/spec"
)

// GetCombinedSpec polls every servicer registered with
// a Swagger spec and merges them together to return a combined spec.
func GetCombinedSpec(yamlCommon string) (string, error) {
	servicers, err := GetSpecServicers()
	if err != nil {
		return "", err
	}

	yamlSpecMap := sync.Map{}
	var wg sync.WaitGroup
	for _, s := range servicers {
		wg.Add(1)
		s := s
		go func() {
			defer wg.Done()
			yamlSpec, err := s.GetPartialSpec()
			yamlSpecMap.Store(s.GetService(), yamlSpec)
			if err != nil {
				// Swallow error because the polling should continue
				// even if it fails to receive from a single servicer
				err = errors.Wrapf(err, "get Swagger spec from %s service", s.GetService())
				glog.Error(err)
			}
		}()
	}
	wg.Wait()

	var yamlSpecs []string
	for _, s := range servicers {
		yaml, ok := yamlSpecMap.Load(s.GetService())
		if !ok {
			glog.Errorf("failed to fetch Swagger spec from map for service %s", s.GetService())
		}
		yamlStr, ok := yaml.(string)
		if !ok {
			glog.Errorf("failed to cast Swagger spec to string for service %s", s.GetService())
		}
		yamlSpecs = append(yamlSpecs, yamlStr)
	}

	combined, warnings, err := spec.Combine(yamlCommon, yamlSpecs)
	if err != nil {
		return "", err
	}
	if warnings != nil {
		glog.Infof("Some Swagger spec traits were overwritten or unable to be read: %+v", warnings)
	}

	return combined, nil
}

// GetServiceSpec returns a service's standalone spec.
func GetServiceSpec(service string) (string, error) {
	remoteSpec := NewRemoteSpec(service)
	yamlSpec, err := remoteSpec.GetStandaloneSpec()
	if err != nil {
		return "", err
	}

	return yamlSpec, nil
}

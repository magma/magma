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

// File registry.go provides a stream provider registry by forwarding calls to
// the service registry.

package providers

import (
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// GetStreamProvider gets the stream provider for a stream name.
// Returns an error if no provider has been registered for the stream.
func GetStreamProvider(streamName string) (StreamProvider, error) {
	if len(streamName) == 0 {
		return nil, errors.New("stream name cannot be empty string")
	}

	services, err := getServicesForStream(streamName)
	if err != nil {
		return nil, err
	}

	n := len(services)
	if n == 0 {
		return nil, fmt.Errorf("no stream providers found for stream name %s", streamName)
	}
	if n != 1 {
		glog.Warningf("Found %d stream providers for stream name %s", n, streamName)
	}

	return NewRemoteProvider(services[0], streamName), nil
}

func getServicesForStream(streamName string) ([]string, error) {
	services, err := registry.FindServices(orc8r.StreamProviderLabel)
	if err != nil {
		return []string{}, err
	}
	var ret []string
	for _, s := range services {
		streams, err := registry.GetAnnotationList(s, orc8r.StreamProviderStreamsAnnotation)
		// Ignore annotation errors, since they indicate either
		//	- service registry contents were recently updated
		//	- this service has incorrect annotations given its label
		if err != nil {
			glog.Warningf("Received error getting annotation %s for service %s: %v", orc8r.StreamProviderStreamsAnnotation, s, err)
			continue
		}
		if funk.Contains(streams, streamName) {
			ret = append(ret, s)
		}
	}

	return ret, nil
}

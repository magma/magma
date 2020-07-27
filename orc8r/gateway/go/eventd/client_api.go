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

package eventd

import (
	"context"
	"strings"

	"magma/gateway/mconfig"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	mcfgprotos "magma/orc8r/lib/go/protos/mconfig"
	platformregistry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

const (
	ServiceName      = "EVENTD"
	DefaultVerbosity = 0
)

type Verbosity bool

// V generates a Verbosity boolean to determine whether an event should be
// logged based off of the provided verbosity and the configured verbosity.
func V(eventVerbosity int32) Verbosity {
	configuredVerbosity, err := getMconfigLogVerbosity()
	if err != nil {
		configuredVerbosity = DefaultVerbosity - 1
		glog.V(1).Infof("Could not load mconfig event verbosity: %s; Using verbosity: %d", err, configuredVerbosity)
	}
	return eventVerbosity <= configuredVerbosity
}

// Log sends an event to eventd if the event's Verbosity is true.
func (v Verbosity) Log(request *protos.Event) error {
	if !v {
		return nil
	}
	client, err := getEventdClient()
	if err != nil {
		return err
	}
	_, err = client.LogEvent(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return err
}

func getEventdClient() (protos.EventServiceClient, error) {
	conn, err := platformregistry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewEventServiceClient(conn), nil
}

func getMconfigLogVerbosity() (int32, error) {
	eventdCfg := &mcfgprotos.EventD{}
	err := mconfig.GetServiceConfigs(strings.ToLower(ServiceName), eventdCfg)
	if err != nil {
		return 0, err
	}
	return eventdCfg.GetEventVerbosity(), nil
}

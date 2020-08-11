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

package servicers

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"

	"magma/devmand/cloud/go/devmand"
	"magma/devmand/cloud/go/protos/mconfig"
	"magma/devmand/cloud/go/services/devmand/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
)

type builderServicer struct{}

func NewBuilderServicer() builder_protos.MconfigBuilderServer {
	return &builderServicer{}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (ret *builder_protos.BuildResponse, err error) {
	ret = &builder_protos.BuildResponse{ConfigsByKey: map[string]*any.Any{}, JsonConfigsByKey: map[string][]byte{}}

	// TODO(T71525030): revert defer, above changes, and fn signature changes
	defer func() { err = ret.FillJSONConfigs(err) }()

	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.Graph)
	if err != nil {
		return nil, err
	}
	devmandAgent, err := graph.GetEntity(devmand.SymphonyAgentType, request.GatewayId)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}

	devices, err := graph.GetAllChildrenOfType(devmandAgent, devmand.SymphonyDeviceType)
	if err != nil {
		return nil, err
	}

	managedDevices := map[string]*mconfig.ManagedDevice{}
	for _, device := range devices {
		deviceConfig := device.Config.(*models.SymphonyDeviceConfig)
		var channels *mconfig.Channels

		if deviceConfig.Channels != nil {
			var snmpConfig *mconfig.SNMPChannel
			if deviceConfig.Channels.SnmpChannel != nil {
				snmpModel := deviceConfig.Channels.SnmpChannel
				snmpConfig = &mconfig.SNMPChannel{
					Community: snmpModel.Community,
					Version:   snmpModel.Version,
				}
			}

			var frinxConfig *mconfig.FrinxChannel
			if deviceConfig.Channels.FrinxChannel != nil {
				frinxModel := deviceConfig.Channels.FrinxChannel
				frinxConfig = &mconfig.FrinxChannel{
					Authorization: frinxModel.Authorization,
					DeviceType:    frinxModel.DeviceType,
					DeviceVersion: frinxModel.DeviceVersion,
					FrinxPort:     frinxModel.FrinxPort,
					Host:          frinxModel.Host,
					Password:      frinxModel.Password,
					Port:          frinxModel.Port,
					TransportType: frinxModel.TransportType,
					Username:      frinxModel.Username,
				}
			}

			var cambiumConfig *mconfig.CambiumChannel
			if deviceConfig.Channels.CambiumChannel != nil {
				cambiumModel := deviceConfig.Channels.CambiumChannel
				cambiumConfig = &mconfig.CambiumChannel{
					ClientId:     cambiumModel.ClientID,
					ClientIp:     cambiumModel.ClientIP,
					ClientMac:    cambiumModel.ClientMac,
					ClientSecret: cambiumModel.ClientSecret,
				}
			}

			var otherConfig *mconfig.OtherChannel
			if deviceConfig.Channels.OtherChannel != nil {
				otherConfig = &mconfig.OtherChannel{
					ChannelProps: deviceConfig.Channels.OtherChannel.ChannelProps,
				}
			}

			channels = &mconfig.Channels{
				SnmpChannel:    snmpConfig,
				FrinxChannel:   frinxConfig,
				CambiumChannel: cambiumConfig,
				OtherChannel:   otherConfig,
			}
		}

		deviceMconfig := &mconfig.ManagedDevice{
			DeviceConfig: deviceConfig.DeviceConfig,
			Host:         deviceConfig.Host,
			DeviceType:   deviceConfig.DeviceType,
			Platform:     deviceConfig.Platform,
			Channels:     channels,
		}
		managedDevices[device.Key] = deviceMconfig
	}

	devmandMconfig := &mconfig.DevmandGatewayConfig{ManagedDevices: managedDevices}
	ret.ConfigsByKey["devmand"], err = ptypes.MarshalAny(devmandMconfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ret, nil
}

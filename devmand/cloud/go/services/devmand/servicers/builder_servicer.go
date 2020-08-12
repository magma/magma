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

	"github.com/golang/protobuf/proto"

	"magma/devmand/cloud/go/devmand"
	devmand_mconfig "magma/devmand/cloud/go/protos/mconfig"
	"magma/devmand/cloud/go/services/devmand/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
)

type builderServicer struct{}

func NewBuilderServicer() builder_protos.MconfigBuilderServer {
	return &builderServicer{}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (*builder_protos.BuildResponse, error) {
	ret := &builder_protos.BuildResponse{ConfigsByKey: map[string][]byte{}}

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

	managedDevices := map[string]*devmand_mconfig.ManagedDevice{}
	for _, device := range devices {
		deviceConfig := device.Config.(*models.SymphonyDeviceConfig)
		var channels *devmand_mconfig.Channels

		if deviceConfig.Channels != nil {
			var snmpConfig *devmand_mconfig.SNMPChannel
			if deviceConfig.Channels.SnmpChannel != nil {
				snmpModel := deviceConfig.Channels.SnmpChannel
				snmpConfig = &devmand_mconfig.SNMPChannel{
					Community: snmpModel.Community,
					Version:   snmpModel.Version,
				}
			}

			var frinxConfig *devmand_mconfig.FrinxChannel
			if deviceConfig.Channels.FrinxChannel != nil {
				frinxModel := deviceConfig.Channels.FrinxChannel
				frinxConfig = &devmand_mconfig.FrinxChannel{
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

			var cambiumConfig *devmand_mconfig.CambiumChannel
			if deviceConfig.Channels.CambiumChannel != nil {
				cambiumModel := deviceConfig.Channels.CambiumChannel
				cambiumConfig = &devmand_mconfig.CambiumChannel{
					ClientId:     cambiumModel.ClientID,
					ClientIp:     cambiumModel.ClientIP,
					ClientMac:    cambiumModel.ClientMac,
					ClientSecret: cambiumModel.ClientSecret,
				}
			}

			var otherConfig *devmand_mconfig.OtherChannel
			if deviceConfig.Channels.OtherChannel != nil {
				otherConfig = &devmand_mconfig.OtherChannel{
					ChannelProps: deviceConfig.Channels.OtherChannel.ChannelProps,
				}
			}

			channels = &devmand_mconfig.Channels{
				SnmpChannel:    snmpConfig,
				FrinxChannel:   frinxConfig,
				CambiumChannel: cambiumConfig,
				OtherChannel:   otherConfig,
			}
		}

		deviceMconfig := &devmand_mconfig.ManagedDevice{
			DeviceConfig: deviceConfig.DeviceConfig,
			Host:         deviceConfig.Host,
			DeviceType:   deviceConfig.DeviceType,
			Platform:     deviceConfig.Platform,
			Channels:     channels,
		}
		managedDevices[device.Key] = deviceMconfig
	}

	vals := map[string]proto.Message{
		"devmand": &devmand_mconfig.DevmandGatewayConfig{ManagedDevices: managedDevices},
	}
	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

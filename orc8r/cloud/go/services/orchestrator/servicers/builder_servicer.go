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
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	mconfig_protos "magma/orc8r/lib/go/protos/mconfig"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var localBuilders = []mconfig.Builder{
	&baseOrchestratorBuilder{},
}

type builderServicer struct{}

func NewBuilderServicer() builder_protos.MconfigBuilderServer {
	return &builderServicer{}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (*builder_protos.BuildResponse, error) {
	ret := &builder_protos.BuildResponse{ConfigsByKey: map[string][]byte{}}

	for _, b := range localBuilders {
		partialConfig, err := b.Build(request.Network, request.Graph, request.GatewayId)
		if err != nil {
			return nil, errors.Wrapf(err, "sub-builder %+v error", b)
		}
		for key, config := range partialConfig {
			_, ok := ret.ConfigsByKey[key]
			if ok {
				return nil, fmt.Errorf("builder received partial config for key %v from multiple sub-builders", key)
			}
			ret.ConfigsByKey[key] = config
		}
	}

	return ret, nil
}

type baseOrchestratorBuilder struct{}

func (b *baseOrchestratorBuilder) Build(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (mconfig.ConfigsByKey, error) {
	networkID := network.ID
	nativeGraph, err := (configurator.EntityGraph{}).FromProto(graph, serdes.Entity)
	if err != nil {
		return nil, err
	}

	// Gateway must be present in the graph
	gateway, err := nativeGraph.GetEntity(orc8r.MagmadGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return nil, errors.Errorf("could not find magmad gateway %s in graph", gatewayID)
	}
	if err != nil {
		return nil, err
	}

	vals := map[string]proto.Message{}
	if gateway.Config != nil {
		gatewayConfig := gateway.Config.(*models.MagmadGatewayConfigs)
		vals["magmad"], err = getMagmadMconfig(&gateway, &nativeGraph, gatewayConfig)
		if err != nil {
			return nil, err
		}
		vals["td-agent-bit"] = getFluentBitMconfig(networkID, gatewayID, gatewayConfig)
		vals["eventd"] = getEventdMconfig(gatewayConfig)
		vals["ovpn"] = getVpnMconfig(gatewayConfig)
	}
	vals["control_proxy"] = &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO}
	vals["metricsd"] = &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO}

	configs, err := mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func getMagmadMconfig(
	gateway *configurator.NetworkEntity, graph *configurator.EntityGraph, gatewayConfig *models.MagmadGatewayConfigs,
) (*mconfig_protos.MagmaD, error) {
	version, images, err := getPackageVersionAndImages(gateway, graph)
	if err != nil {
		return nil, err
	}

	ret := &mconfig_protos.MagmaD{
		LogLevel:                protos.LogLevel_INFO,
		CheckinInterval:         int32(gatewayConfig.CheckinInterval),
		CheckinTimeout:          int32(gatewayConfig.CheckinTimeout),
		AutoupgradeEnabled:      swag.BoolValue(gatewayConfig.AutoupgradeEnabled),
		AutoupgradePollInterval: gatewayConfig.AutoupgradePollInterval,
		PackageVersion:          version,
		Images:                  images,
		DynamicServices:         gatewayConfig.DynamicServices,
		FeatureFlags:            gatewayConfig.FeatureFlags,
	}

	return ret, nil
}

func getPackageVersionAndImages(magmadGateway *configurator.NetworkEntity, graph *configurator.EntityGraph) (string, []*mconfig_protos.ImageSpec, error) {
	tier, err := graph.GetFirstAncestorOfType(*magmadGateway, orc8r.UpgradeTierEntityType)
	if err == merrors.ErrNotFound {
		return "0.0.0-0", []*mconfig_protos.ImageSpec{}, nil
	}
	if err != nil {
		return "0.0.0-0", []*mconfig_protos.ImageSpec{}, errors.Wrap(err, "failed to load upgrade tier")
	}

	tierConfig := tier.Config.(*models.Tier)
	retImages := make([]*mconfig_protos.ImageSpec, 0, len(tierConfig.Images))
	for _, image := range tierConfig.Images {
		retImages = append(retImages, &mconfig_protos.ImageSpec{Name: swag.StringValue(image.Name), Order: swag.Int64Value(image.Order)})
	}
	return tierConfig.Version.ToString(), retImages, nil
}

func getFluentBitMconfig(networkID string, gatewayID string, mdGw *models.MagmadGatewayConfigs) *mconfig_protos.FluentBit {
	ret := &mconfig_protos.FluentBit{
		ExtraTags: map[string]string{
			"network_id": networkID,
			"gateway_id": gatewayID,
		},
		ThrottleRate:     1000,
		ThrottleWindow:   5,
		ThrottleInterval: "1m",
	}

	if mdGw.Logging != nil && mdGw.Logging.Aggregation != nil {
		ret.FilesByTag = mdGw.Logging.Aggregation.TargetFilesByTag
		if mdGw.Logging.Aggregation.ThrottleRate != nil {
			ret.ThrottleRate = *mdGw.Logging.Aggregation.ThrottleRate
		}
		if mdGw.Logging.Aggregation.ThrottleWindow != nil {
			ret.ThrottleWindow = *mdGw.Logging.Aggregation.ThrottleWindow
		}
		if mdGw.Logging.Aggregation.ThrottleInterval != nil {
			ret.ThrottleInterval = *mdGw.Logging.Aggregation.ThrottleInterval
		}
	}

	return ret
}

func getEventdMconfig(gatewayConfig *models.MagmadGatewayConfigs) *mconfig_protos.EventD {
	ret := &mconfig_protos.EventD{
		LogLevel:       protos.LogLevel_INFO,
		EventVerbosity: -1,
	}
	if gatewayConfig.Logging != nil && gatewayConfig.Logging.EventVerbosity != nil {
		ret.EventVerbosity = *gatewayConfig.Logging.EventVerbosity
	}
	return ret
}

func getVpnMconfig(gatewayConfig *models.MagmadGatewayConfigs) *mconfig_protos.OpenVPN {
	ret := &mconfig_protos.OpenVPN{
		EnableShellAccess: false,
	}
	if gatewayConfig.Vpn != nil {
		ret.EnableShellAccess = *gatewayConfig.Vpn.EnableShell
	}

	return ret
}

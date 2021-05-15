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

package servicers_test

import (
	"context"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/envoy_controller/control_plane"
	"magma/feg/gateway/services/envoy_controller/control_plane/mocks"
	"magma/feg/gateway/services/envoy_controller/servicers"
	lte_proto "magma/lte/cloud/go/protos"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	orig_src "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/original_src/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	IMSI1 = "IMSI00101"
	IMSI2 = "IMSI00102"
	UEIP1 = "3.3.33.3"
	UEIP2 = "2.2.2.2"
)

var (
	imsis   = []string{IMSI1, IMSI2}
	header1 = &protos.Header{
		Name:  "IMSI",
		Value: "024212312312",
	}
	header1_envoy = &core.HeaderValueOption{
		Header: &core.HeaderValue{
			Key:   "IMSI",
			Value: "024212312312",
		}}
	addUe1Rule1 = &protos.AddUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
		Websites: []string{"neverssl.com", "google.com"},
		Headers:  []*protos.Header{header1},
		RuleId:   "rule1",
	}
	addUe1Rule2 = &protos.AddUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
		Websites: []string{"newwebsite.com", "no_conlict.com"},
		Headers: []*protos.Header{{
			Name:  "TEST",
			Value: "12345",
		}},
		RuleId: "rule2",
	}
	addUe1Rule3Conflict = &protos.AddUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
		Websites: []string{"this_will_error.com", "neverssl.com"},
		Headers:  []*protos.Header{header1},
		RuleId:   "rule3",
	}
	addUe2Rule2 = &protos.AddUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP2),
		},
		Websites: []string{"magma.com", "qqq.com"},
		Headers: []*protos.Header{{
			Name:  "IMSI",
			Value: "111111",
		},
			{
				Name:  "MSISDN",
				Value: "THIS_IS_MSISDN",
			}},
		RuleId: "rule2",
	}
	ue1Entry = map[string]*control_plane.UEInfo{
		"rule1": {
			Websites: []string{"neverssl.com", "google.com"},
			Headers:  []*protos.Header{header1},
		}}
	ue1EntryRule12 = map[string]*control_plane.UEInfo{
		"rule1": {
			Websites: []string{"neverssl.com", "google.com"},
			Headers:  []*protos.Header{header1},
		},
		"rule2": {
			Websites: []string{"newwebsite.com", "no_conlict.com"},
			Headers: []*protos.Header{{
				Name:  "TEST",
				Value: "12345",
			}},
		}}
	ue2Entry = map[string]*control_plane.UEInfo{
		"rule2": {
			Websites: []string{"magma.com", "qqq.com"},
			Headers: []*protos.Header{{
				Name:  "IMSI",
				Value: "111111",
			},
				{
					Name:  "MSISDN",
					Value: "THIS_IS_MSISDN",
				}},
		}}
	emptyDict = control_plane.UEInfoMap{}
	ue1Dict   = control_plane.UEInfoMap{
		UEIP1: ue1Entry,
	}
	ue1Rule12Dict = control_plane.UEInfoMap{
		UEIP1: ue1EntryRule12,
	}
	ue_2_dict = control_plane.UEInfoMap{
		UEIP2: ue2Entry,
	}
	ue12Dict = control_plane.UEInfoMap{
		UEIP1: ue1Entry,
		UEIP2: ue2Entry,
	}
	deactivateUe1 = &protos.DeactivateUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
	}
	deactivateUe1Rule1 = &protos.DeactivateUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
		RuleId: "rule1",
	}
	deactivateUe1Rule2 = &protos.DeactivateUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte(UEIP1),
		},
		RuleId: "rule2",
	}
	addSuccess        = &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_SUCCESS}
	deactivateSuccess = &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_SUCCESS}

	routes = []*v2route.Route{{
		Match: &v2route.RouteMatch{
			PathSpecifier: &v2route.RouteMatch_Prefix{
				Prefix: "/",
			},
		},
		Action: &v2route.Route_Route{
			Route: &v2route.RouteAction{
				ClusterSpecifier: &v2route.RouteAction_Cluster{
					Cluster: "cluster1",
				},
			},
		},
	}}
	virtualHosts = []*v2route.VirtualHost{
		{
			Name:                "local_service",
			Domains:             []string{"*"},
			RequestHeadersToAdd: []*core.HeaderValueOption{},
			Routes:              routes,
		},
		{
			Name:                "local_service",
			Domains:             []string{"neverssl.com", "google.com"},
			RequestHeadersToAdd: []*core.HeaderValueOption{header1_envoy},
			Routes:              routes,
		},
	}
	pbst, _ = ptypes.MarshalAny(
		&hcm.HttpConnectionManager{
			CodecType:        hcm.HttpConnectionManager_AUTO,
			StatPrefix:       "ingress_http",
			UseRemoteAddress: &wrappers.BoolValue{Value: true},
			CommonHttpProtocolOptions: &core.HttpProtocolOptions{
				IdleTimeout: ptypes.DurationProto(3600 * time.Second),
				// TODO figure out why this doesn't work properly
				//HeadersWithUnderscoresAction: core.HttpProtocolOptions_REJECT_REQUEST,
			},
			Http2ProtocolOptions: &core.Http2ProtocolOptions{
				MaxConcurrentStreams:        &wrappers.UInt32Value{Value: 16},
				InitialStreamWindowSize:     &wrappers.UInt32Value{Value: 65536},
				InitialConnectionWindowSize: &wrappers.UInt32Value{Value: 1048576},
			},
			RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
				RouteConfig: &v2.RouteConfiguration{
					Name:         "matched_website_route",
					VirtualHosts: virtualHosts,
				},
			},
			HttpFilters: []*hcm.HttpFilter{{
				Name: wellknown.Router,
			}},
		},
	)

	mo_src, _    = ptypes.MarshalAny(&orig_src.OriginalSrc{})
	retListener1 = &v2.Listener{
		Name: "default_http",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address: "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 80,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				FilterChainMatch: &listener.FilterChainMatch{
					SourcePrefixRanges: []*core.CidrRange{{
						AddressPrefix: UEIP1,
						PrefixLen:     &wrappers.UInt32Value{Value: 32},
					}}},
				Filters: []*listener.Filter{{
					Name: wellknown.HTTPConnectionManager,
					ConfigType: &listener.Filter_TypedConfig{
						TypedConfig: pbst,
					},
				}},
			}},
		ListenerFilters: []*listener.ListenerFilter{

			{
				Name: "envoy.filters.listener.original_dst",
			},
			{
				Name: "envoy.filters.listener.original_src",
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: mo_src,
				},
			},
		},
		Transparent: &wrappers.BoolValue{Value: true},
	}
)

// ---- TESTS ----
func TestNormalCallFlow(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	listener, err := control_plane.GetListener(ue1Dict)
	assert.NoError(t, err)
	assert.Equal(t, listener, retListener1)

	assert.NoError(t, err)
	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

func TestAddRemoveFlow(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	cli.On("UpdateSnapshot", ue12Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")
	ret, err = srv.AddUEHeaderEnrichment(ctx, addUe2Rule2)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", ue_2_dict).Return()
	ret2, err := srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1)
	assert.NoError(t, err)
	assert.Equal(t, ret2, deactivateSuccess, "Rule should be added successfully")

	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

func TestDoubleActivation(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", ue1Rule12Dict).Return()
	ret, err = srv.AddUEHeaderEnrichment(ctx, addUe1Rule2)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

func TestMultiRemoval(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", ue1Rule12Dict).Return()
	ret, err = srv.AddUEHeaderEnrichment(ctx, addUe1Rule2)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret2, err := srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1Rule2)
	assert.NoError(t, err)
	assert.Equal(t, ret2, deactivateSuccess, "Rule should be deactivated successfully")

	cli.AssertExpectations(t)
}

func TestCompleteRemoval(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", emptyDict).Return()
	ret2, err := srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1Rule1)
	assert.Equal(t, ret2, deactivateSuccess, "Rule should be deactivated successfully")

	assert.NoError(t, err)
	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

func TestUERemoval(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	assert.Equal(t, ret, addSuccess, "Rule should be added successfully")

	cli.On("UpdateSnapshot", emptyDict).Return()
	ret2, err := srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1)
	assert.NoError(t, err)
	assert.Equal(t, ret2, deactivateSuccess, "Rule should be deactivated successfully")

	assert.NoError(t, err)
	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

func TestInvalidDeactivate(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	ret, err := srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1)
	assert.NoError(t, err)
	ue_not_found := &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_UE_NOT_FOUND}

	assert.Equal(t, ret, ue_not_found, "UE can't be deleted if it doesn't exist")

	cli.On("UpdateSnapshot", ue1Dict).Return()
	_, err = srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)

	ret, err = srv.DeactivateUEHeaderEnrichment(ctx, deactivateUe1Rule2)
	assert.NoError(t, err)
	rule_not_found := &protos.DeactivateUEHeaderEnrichmentResult{Result: protos.DeactivateUEHeaderEnrichmentResult_RULE_NOT_FOUND}

	assert.Equal(t, ret, rule_not_found, "Rule can't be deleted if it doesn't exist")

	cli.AssertExpectations(t)
}

func TestInvalidAdd(t *testing.T) {
	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", ue1Dict).Return()
	_, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)

	ret, err := srv.AddUEHeaderEnrichment(ctx, addUe1Rule1)
	assert.NoError(t, err)
	rule_id_conflict := &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_RULE_ID_CONFLICT}
	assert.Equal(t, ret, rule_id_conflict, "Can't insert duplicate rule")

	ret, err = srv.AddUEHeaderEnrichment(ctx, addUe1Rule3Conflict)
	assert.NoError(t, err)
	ip_conflict := &protos.AddUEHeaderEnrichmentResult{Result: protos.AddUEHeaderEnrichmentResult_WEBSITE_CONFLICT}
	assert.Equal(t, ret, ip_conflict, "Can't insert rule that will cause website collison")

	assert.NoError(t, err)
	assert.NoError(t, err)
	cli.AssertExpectations(t)
}

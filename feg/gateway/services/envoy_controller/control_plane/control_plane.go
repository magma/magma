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

/* This file contains modified code from github.com/envoyproxy/go-control-plane
module which is published under license Apache 2.0. See envoy_controller/README
for more details.
*/

/*
  Copyright 2018 Envoyproxy Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package control_plane

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	orig_src "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/original_src/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
)

const (
	XdsCluster                  = "xds_cluster"
	Ads                         = "ads"
	Xds                         = "xds"
	Rest                        = "rest"
	debug                       = true
	port                        = 18000
	gatewayPort                 = 18001
	mode                        = Ads
	any_addr                    = "0.0.0.0"
	maxConcurrentStreams        = 16
	initialStreamWindowSize     = 65536   // 64Kib
	initialConnectionWindowSize = 1048576 // 1 MiB
	connectTimeout              = 6 * time.Second
	idleTimeout                 = 3600 * time.Second
	httpPort                    = 80
	clusterName                 = "cluster1"
	targetPrefix                = "/"
	virtualHostName             = "local_service"
	routeConfigName             = "matched_website_route"
	grpcMaxConcurrentStreams    = 1000000
	listenerName                = "default_http"
	defaultRouteName            = "default_route"
)

type EnvoyController interface {
	UpdateSnapshot(UEInfoMap)
}

type UEInfo struct {
	Websites []string
	Headers  []*protos.Header
}

type UEInfoMap map[string]map[string]*UEInfo

type ControllerClient struct {
	version int32
	config  cache.SnapshotCache
}

type callbacks struct {
	signal         chan struct{}
	fetches        int
	requests       int
	deltaRequests  int
	deltaResponses int
	mu             sync.Mutex
}

// Hasher returns node ID as an ID
type Hasher struct{}

func (cb *callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	glog.V(2).Infof("cb.Report() fetches %d,  callbacks %d", cb.fetches, cb.requests)
}

// methods below (starting with On...) are included to satisfy the xds.Callbacks interface

func (cb *callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	glog.V(2).Infof("OnStreamOpen %d open for %s", id, typ)
	return nil
}

func (cb *callbacks) OnStreamClosed(id int64) {
	glog.V(2).Infof("OnStreamClosed %d closed", id)
}

func (cb *callbacks) OnStreamRequest(int64, *discovery.DiscoveryRequest) error {
	glog.V(2).Infof("OnStreamRequest")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnStreamResponse(context.Context, int64, *discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {
	glog.V(2).Infof("OnStreamResponse...")
	cb.Report()
}
func (cb *callbacks) OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error {
	glog.V(2).Infof("OnFetchRequest...")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}

func (cb *callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {
	glog.V(2).Infof("OnFetchResponse...")
}

func (cb *callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	glog.V(2).Infof("delta stream %d open for %s\n", id, typ)
	return nil
}

func (cb *callbacks) OnDeltaStreamClosed(id int64) {
	glog.V(2).Infof("delta stream %d closed\n", id)
}

func (cb *callbacks) OnStreamDeltaResponse(int64, *discovery.DeltaDiscoveryRequest, *discovery.DeltaDiscoveryResponse) {
	glog.V(2).Infof("OnStreamDeltaResponse...")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaResponses++
}

func (cb *callbacks) OnStreamDeltaRequest(int64, *discovery.DeltaDiscoveryRequest) error {
	glog.V(2).Infof("OnStreamDeltaRequest...")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaRequests++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}

// ID function
func (h Hasher) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

type HTTPGateway struct {
	// Server is the underlying gRPC server
	Server xds.Server
}

// RunManagementServer starts an xDS server at the given port.
func RunManagementServer(ctx context.Context, server xds.Server, port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.Fatalf("failed to listen %s", err)
	}

	// register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	glog.Infof("Management server listening on port %d", port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			glog.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

// RunManagementGateway starts an HTTP gateway to an xDS server.
func RunManagementGateway(ctx context.Context, srv xds.Server, port uint) {
	glog.Infof("Gateway listening HTTP/1.1 on port %d", port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &HTTPGateway{Server: srv}}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			glog.Error(err)
		}
	}()
}

func (h *HTTPGateway) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	gtw := xds.HTTPGateway{Server: h.Server}
	bytes, code, err := gtw.ServeHTTP(req)

	if err != nil {
		http.Error(resp, err.Error(), code)
		return
	}

	if bytes == nil {
		resp.WriteHeader(http.StatusNotModified)
		return
	}

	if _, err = resp.Write(bytes); err != nil {
		glog.Errorf("gateway error: %v", err)
	}
}

func newCallbacks(signal chan struct{}, fetches int, requests int, deltaRequests int, deltaResponses int) *callbacks {
	return &callbacks{
		signal:         signal,
		fetches:        fetches,
		requests:       requests,
		deltaRequests:  deltaRequests,
		deltaResponses: deltaResponses,
	}
}

func GetControllerClient() *ControllerClient {
	cli := ControllerClient{}
	ctx := context.Background()

	glog.Infof("Starting Envoy control plane")

	signal := make(chan struct{})
	cb := newCallbacks(signal, 0, 0, 0, 0)
	cli.config = cache.NewSnapshotCache(mode == Ads, Hasher{}, nil)

	srv := xds.NewServer(ctx, cli.config, cb)

	// start the xDS server
	go RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)

	cb.Report()

	return &cli
}

func getHttpConnectionManager(routeConfigName string, virtualHosts []*route.VirtualHost) *hcm.HttpConnectionManager {
	useRemoteAddress := &wrappers.BoolValue{Value: true}
	commonHttpProtocolOptions := &core.HttpProtocolOptions{
		IdleTimeout: ptypes.DurationProto(idleTimeout),
		// TODO figure out why this doesn't work properly
		//HeadersWithUnderscoresAction: core.HttpProtocolOptions_REJECT_REQUEST,
	}
	http2ProtocolOptions := &core.Http2ProtocolOptions{
		MaxConcurrentStreams:        &wrappers.UInt32Value{Value: maxConcurrentStreams},
		InitialStreamWindowSize:     &wrappers.UInt32Value{Value: initialStreamWindowSize},
		InitialConnectionWindowSize: &wrappers.UInt32Value{Value: initialConnectionWindowSize},
	}
	routeSpecifier := &hcm.HttpConnectionManager_RouteConfig{
		RouteConfig: &route.RouteConfiguration{
			Name:         routeConfigName,
			VirtualHosts: virtualHosts,
		},
	}
	httpFilters := []*hcm.HttpFilter{{
		Name: wellknown.Router,
	}}

	return &hcm.HttpConnectionManager{
		CodecType:                 hcm.HttpConnectionManager_AUTO,
		StatPrefix:                "ingress_http",
		UseRemoteAddress:          useRemoteAddress,
		CommonHttpProtocolOptions: commonHttpProtocolOptions,
		Http2ProtocolOptions:      http2ProtocolOptions,
		RouteSpecifier:            routeSpecifier,
		HttpFilters:               httpFilters,
	}
}

func getVirtualHost(virtualHostName string, domains []string, requestHeadersToAdd []*core.HeaderValueOption) *route.VirtualHost {
	routes := []*route.Route{{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: targetPrefix,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterName,
				},
			},
		},
	}}
	return &route.VirtualHost{
		Name:                virtualHostName,
		Domains:             domains,
		RequestHeadersToAdd: requestHeadersToAdd,
		Routes:              routes,
	}
}

func getHeadersToAdd(ueInfo *UEInfo) []*core.HeaderValueOption {
	requestHeadersToAdd := []*core.HeaderValueOption{}
	for _, header := range ueInfo.Headers {
		headerValueOption := &core.HeaderValueOption{
			Header: &core.HeaderValue{
				Key:   header.Name,
				Value: header.Value,
			},
		}
		requestHeadersToAdd = append(requestHeadersToAdd, headerValueOption)
	}
	return requestHeadersToAdd
}

func getUEFilterChains(ues UEInfoMap) ([]*listener.FilterChain, error) {
	filterChains := []*listener.FilterChain{}
	for ue_ip_addr, rule_map := range ues {
		glog.V(2).Infof("Adding UE - " + ue_ip_addr)
		virtualHosts := []*route.VirtualHost{getVirtualHost(virtualHostName, []string{"*"}, []*core.HeaderValueOption{})}

		for _, ueInfo := range rule_map {
			requestHeadersToAdd := getHeadersToAdd(ueInfo)
			virtualHosts = append(virtualHosts, getVirtualHost(virtualHostName, ueInfo.Websites, requestHeadersToAdd))
		}

		pbst, err := ptypes.MarshalAny(getHttpConnectionManager(routeConfigName, virtualHosts))
		if err != nil {
			glog.Errorf("Couldn't marshal UE HTTP connection manager")
			continue
		}
		filterChainMatch := &listener.FilterChainMatch{
			SourcePrefixRanges: []*core.CidrRange{{
				AddressPrefix: ue_ip_addr,
				PrefixLen:     &wrappers.UInt32Value{Value: 32},
			}}}
		filters := []*listener.Filter{{
			Name: wellknown.HTTPConnectionManager,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: pbst,
			},
		}}
		filterChains = append(filterChains, &listener.FilterChain{
			FilterChainMatch: filterChainMatch,
			Filters:          filters,
		})
		glog.V(2).Infof("Returning virtual hosts %s", virtualHosts)

	}

	return filterChains, nil
}

func GetListener(ues UEInfoMap) (*listener.Listener, error) {
	glog.V(2).Infof("Creating listener " + listenerName)
	filterChains, err := getUEFilterChains(ues)
	if err != nil {
		return nil, err
	}

	o_src := &orig_src.OriginalSrc{}
	mo_src, err := ptypes.MarshalAny(o_src)
	if err != nil {
		return nil, errors.New("Couldn't marshal OriginalSrc")
	}
	listenerFilters := []*listener.ListenerFilter{
		{
			Name: "envoy.filters.listener.original_dst",
		},
		{
			Name: "envoy.filters.listener.original_src",
			ConfigType: &listener.ListenerFilter_TypedConfig{
				TypedConfig: mo_src,
			},
		},
	}

	address := &core.Address{
		Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address: any_addr,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: httpPort,
				},
			},
		},
	}

	var listener = &listener.Listener{
		Name:            listenerName,
		Transparent:     &wrappers.BoolValue{Value: true},
		Address:         address,
		FilterChains:    filterChains,
		ListenerFilters: listenerFilters,
	}

	glog.V(2).Infof("Returning listener %s", listener)
	return listener, nil
}

func getDefaultReq() UEInfoMap {
	ret := UEInfoMap{}
	ret["0.0.0.0"] = map[string]*UEInfo{}
	ret["0.0.0.0"]["default"] = &UEInfo{
		Websites: []string{"0.0.0.0"},
	}
	return ret
}

func (cli *ControllerClient) UpdateSnapshot(ues UEInfoMap) {
	cluster := []types.Resource{
		&cluster.Cluster{
			Name:                 clusterName,
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
			ConnectTimeout:       ptypes.DurationProto(connectTimeout),
			LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		},
	}

	if len(ues) == 0 {
		ues = getDefaultReq()
	}
	listener, err := GetListener(ues)
	listener_resource := []types.Resource{listener}

	if err != nil {
		glog.Errorf("Get Listener error %s", err)
		return
	}
	nodeId := cli.config.GetStatusKeys()[0]

	atomic.AddInt32(&cli.version, 1)
	glog.Infof("Saved snapshot version " + fmt.Sprint(cli.version))
	snap, _ := cache.NewSnapshot(
		fmt.Sprint(cli.version),
		map[resource.Type][]types.Resource{
			resource.ClusterType:  cluster,
			resource.ListenerType: listener_resource,
		},
	)
	cli.config.SetSnapshot(context.Background(), nodeId, snap)
}

package control_plane

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"sync"
	"sync/atomic"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"

	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	orig_src "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/original_src/v3"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"magma/feg/cloud/go/protos"
)

var (
	any_addr = "0.0.0.0"

	version int32

	config cache.SnapshotCache
)

const (
	XdsCluster  = "xds_cluster"
	Ads         = "ads"
	Xds         = "xds"
	Rest        = "rest"
	debug       = true
	port        = 18000
	gatewayPort = 18001
	mode        = Ads
)

type logger struct{}

func (logger logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (logger logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
func (cb *callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.WithFields(log.Fields{"fetches": cb.fetches, "requests": cb.requests}).Info("cb.Report()  callbacks")
}
func (cb *callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	log.Infof("OnStreamOpen %d open for %s", id, typ)
	return nil
}
func (cb *callbacks) OnStreamClosed(id int64) {
	log.Infof("OnStreamClosed %d closed", id)
}
func (cb *callbacks) OnStreamRequest(int64, *v2.DiscoveryRequest) error {
	log.Infof("OnStreamRequest")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	log.Infof("OnStreamResponse...")
	cb.Report()
}
func (cb *callbacks) OnFetchRequest(ctx context.Context, req *v2.DiscoveryRequest) error {
	log.Infof("OnFetchRequest...")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {}

type callbacks struct {
	signal   chan struct{}
	fetches  int
	requests int
	mu       sync.Mutex
}

// Hasher returns node ID as an ID
type Hasher struct {
}

// ID function
func (h Hasher) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

const grpcMaxConcurrentStreams = 1000000

// RunManagementServer starts an xDS server at the given port.
func RunManagementServer(ctx context.Context, server xds.Server, port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	// register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	log.WithFields(log.Fields{"port": port}).Info("management server listening")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

// RunManagementGateway starts an HTTP gateway to an xDS server.
func RunManagementGateway(ctx context.Context, srv xds.Server, port uint) {
	log.WithFields(log.Fields{"port": port}).Info("gateway listening HTTP/1.1")
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &xds.HTTPGateway{Server: srv}}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
}

func Setup() {
	ctx := context.Background()

	log.Printf("Starting control plane")

	signal := make(chan struct{})
	cb := &callbacks{
		signal:   signal,
		fetches:  0,
		requests: 0,
	}
	config = cache.NewSnapshotCache(mode == Ads, Hasher{}, nil)

	srv := xds.NewServer(ctx, config, cb)

	// start the xDS server
	go RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)

	// what in the golang is this?
	<-signal

	cb.Report()

}

func UpdateSnapshot(ues []*protos.AddUEHeaderEnrichmentRequest) {

	nodeId := config.GetStatusKeys()[0]

	var clusterName = "cluster1"
	cluster := []cache.Resource{
		&v2.Cluster{
			Name:                 clusterName,
			ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_ORIGINAL_DST},
			ConnectTimeout:       ptypes.DurationProto(6 * time.Second),
			LbPolicy:             v2.Cluster_CLUSTER_PROVIDED,
		},
	}

	filterChains := []*listener.FilterChain{}
	var listenerName = "default_http"
	var targetPrefix = "/"
	var virtualHostName = "local_service"
	var routeConfigName = "local_route"

	for _, req := range ues {
		var ue_ip_addr = string(req.UeIp.Address)
		requestHeadersToAdd := []*core.HeaderValueOption{}

		for _, header := range req.Headers {
			requestHeadersToAdd = append(requestHeadersToAdd, &core.HeaderValueOption{
				Header: &core.HeaderValue{
					Key:   header.Name,
					Value: header.Value,
				},
			})
		}

		virtualHosts := []*v2route.VirtualHost{}
		virtualHosts = append(virtualHosts, &v2route.VirtualHost{
			Name:                virtualHostName,
			Domains:             req.Websites,
			RequestHeadersToAdd: requestHeadersToAdd,
			Routes: []*v2route.Route{{
				Match: &v2route.RouteMatch{
					PathSpecifier: &v2route.RouteMatch_Prefix{
						Prefix: targetPrefix,
					},
				},
				Action: &v2route.Route_Route{
					Route: &v2route.RouteAction{
						ClusterSpecifier: &v2route.RouteAction_Cluster{
							Cluster: clusterName,
						},
					},
				},
			}}})

		httpManager := &hcm.HttpConnectionManager{
			CodecType:        hcm.HttpConnectionManager_AUTO,
			StatPrefix:       "ingress_http",
			UseRemoteAddress: &wrappers.BoolValue{Value: true},
			CommonHttpProtocolOptions: &core.HttpProtocolOptions{
				IdleTimeout: ptypes.DurationProto(3600 * time.Second),
				//HeadersWithUnderscoresAction: core.HttpProtocolOptions_REJECT_REQUEST,
			},
			Http2ProtocolOptions: &core.Http2ProtocolOptions{
				MaxConcurrentStreams:        &wrappers.UInt32Value{Value: 100},
				InitialStreamWindowSize:     &wrappers.UInt32Value{Value: 65536},   // 64Kib
				InitialConnectionWindowSize: &wrappers.UInt32Value{Value: 1048576}, // 1 MiB
			},
			StreamIdleTimeout: ptypes.DurationProto(300 * time.Second), // 5 mins, must be disabled for long-lived and streaming requests
			RequestTimeout:    ptypes.DurationProto(300 * time.Second), // 5 mins, must be disabled for long-lived and streaming requests
			RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
				RouteConfig: &v2.RouteConfiguration{
					Name:         routeConfigName,
					VirtualHosts: virtualHosts,
				},
			},
			HttpFilters: []*hcm.HttpFilter{{
				Name: wellknown.Router,
			}},
		}

		pbst, err := ptypes.MarshalAny(httpManager)
		if err != nil {
			panic(err)
		}

		log.Infof(">>>>>>>>>>>>>>>>>>> adding UE " + ue_ip_addr)

		filterChains = append(filterChains, &listener.FilterChain{
			FilterChainMatch: &listener.FilterChainMatch{
				SourcePrefixRanges: []*core.CidrRange{{
					AddressPrefix: ue_ip_addr,
					PrefixLen:     &wrappers.UInt32Value{Value: 32},
				}},
			},
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		})

	}

	o_src := &orig_src.OriginalSrc{}
	mo_src, err := ptypes.MarshalAny(o_src)
	if err != nil {
		panic(err)
	}

	log.Infof(">>>>>>>>>>>>>>>>>>> creating listener " + listenerName)
	var listener = []cache.Resource{
		&v2.Listener{
			Name:        listenerName,
			Transparent: &wrappers.BoolValue{Value: true},
			Address: &core.Address{
				Address: &core.Address_SocketAddress{
					SocketAddress: &core.SocketAddress{
						Address: any_addr,
						PortSpecifier: &core.SocketAddress_PortValue{
							PortValue: 80,
						},
					},
				},
			},
			FilterChains: filterChains,
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
		}}

	// Save snapshot
	atomic.AddInt32(&version, 1)
	log.Infof(">>>>>>>>>>>>>>>>>>> creating snapshot Version " + fmt.Sprint(version))
	snap := cache.NewSnapshot(fmt.Sprint(version), nil, cluster, nil, listener, nil)
	config.SetSnapshot(nodeId, snap)
}

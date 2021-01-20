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
	"fmt"
	"io"
	"sync"
	"time"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// heartBeatInterval is the heart beat interval from cloud to gateway
const heartBeatInterval = time.Minute

type SyncRPCService struct {
	// hostName is the host at which this service instance is running on
	hostName string
	broker   broker.GatewayRPCBroker
}

func NewSyncRPCService(hostName string, broker broker.GatewayRPCBroker) (protos.SyncRPCServiceServer, error) {
	return &SyncRPCService{hostName: hostName, broker: broker}, nil
}

// SyncRPC exists for backwards compatibility.
//
// Deprecated: Use EstablishSyncRPCStream instead.
func (srv *SyncRPCService) SyncRPC(stream protos.SyncRPCService_SyncRPCServer) error {
	return srv.EstablishSyncRPCStream(stream)
}

// EstablishSyncRPCStream is the RPC call that will be called by gateways.
// It establishes a bidirectional stream between the gateway and the cloud,
// and the streams will close if it returns.
//
// Every active connection will run this function in its own goroutine.
func (srv *SyncRPCService) EstablishSyncRPCStream(stream protos.SyncRPCService_EstablishSyncRPCStreamServer) error {
	// Check if we can get a valid Gateway identity.
	gw, err := identity.GetStreamGatewayId(stream)
	if err != nil {
		return err
	}
	if gw == nil || len(gw.HardwareId) == 0 {
		return status.Errorf(codes.PermissionDenied, "Gateway hardware ID is nil")
	}
	return srv.serveGwId(stream, gw.HardwareId)
}

// streamCoordinator manages a SyncRPC bidirectional stream.
type streamCoordinator struct {
	GwID    string
	ErrChan chan error
	Wg      *sync.WaitGroup
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func newStreamCoordinator(gwId string, streamCtx context.Context) *streamCoordinator {
	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(streamCtx)
	return &streamCoordinator{gwId, errChan, wg, ctx, cancel}
}

// serveGwId handles the SyncRPC bidirectional stream for a particular gateway.
// It starts goroutines to manage SyncRPCRequest and SyncRPCResponse message streams.
//
// It is called directly by the test service.
func (srv *SyncRPCService) serveGwId(stream protos.SyncRPCService_EstablishSyncRPCStreamServer, gwId string) error {
	start := time.Now()

	coordinator := newStreamCoordinator(gwId, stream.Context())
	queue := srv.broker.InitializeGateway(gwId)
	glog.Infof("HWID %v: initialized gateway connection", gwId)

	coordinator.Wg.Add(2)
	go srv.receiveFromStream(stream, coordinator)
	go srv.sendToStream(stream, queue, coordinator)

	err := <-coordinator.ErrChan
	if err == nil {
		glog.Infof("HWID %v: client closed stream by sending EOF, cleaning up", gwId)
	} else {
		glog.Errorf("HWID %v: %v", gwId, err)
	}
	coordinator.Cancel()
	coordinator.Wg.Wait()
	srv.broker.CleanupGateway(gwId)
	glog.Infof("HWID %v: cleanup successful, stream was alive for %v seconds", gwId, time.Since(start).Seconds())

	return err
}

// sendToStream manages the SyncRPCRequest message stream.
// If messages are available in the SyncRPCRequest queue, it will send them
// to the gateway. Otherwise, it will send a heartbeat to the gateway after
// heartBeatInterval.
func (srv *SyncRPCService) sendToStream(
	stream protos.SyncRPCService_EstablishSyncRPCStreamServer,
	queue chan *protos.SyncRPCRequest,
	coordinator *streamCoordinator,
) {
	defer coordinator.Wg.Done()
	for {
		select {
		case <-coordinator.Ctx.Done():
			coordinator.sendErr(fmt.Errorf("HWID %v: sendToStream: context canceled: %v", coordinator.GwID, coordinator.Ctx.Err()))
			return
		case <-time.After(heartBeatInterval):
			glog.V(2).Infof("HWID %v: sending heartbeat", coordinator.GwID)
			err := stream.Send(&protos.SyncRPCRequest{HeartBeat: true})
			if err != nil {
				coordinator.sendErr(fmt.Errorf("HWID %v: sendToStream: error sending to stream (heartbeat): %v", coordinator.GwID, err))
				return
			}
		case reqToSend, ok := <-queue:
			if !ok {
				coordinator.sendErr(fmt.Errorf("HWID %v: sendToStream: requests queue is closed", coordinator.GwID))
				return
			}
			if reqToSend != nil {
				glog.V(2).Infof("HWID %v: sending req to stream", coordinator.GwID)
				err := stream.Send(reqToSend)
				if err != nil {
					coordinator.sendErr(fmt.Errorf("HWID %v: sendToStream: error sending to stream: %v", coordinator.GwID, err))
					return
				}
			}
		}
	}
}

// receiveFromStream manages the SyncRPCResponse stream and processes responses
// that it receives.
func (srv *SyncRPCService) receiveFromStream(
	stream protos.SyncRPCService_EstablishSyncRPCStreamServer,
	coordinator *streamCoordinator,
) {
	defer coordinator.Wg.Done()
	for {
		// recv() can be canceled via ctx
		syncRPCResp, err := RecvWithContext(coordinator.Ctx, func() (*protos.SyncRPCResponse, error) { return stream.Recv() })
		if err == io.EOF {
			coordinator.sendErr(nil)
			return
		} else if err != nil {
			coordinator.sendErr(fmt.Errorf("HWID %v: receiveFromStream: error receiving from stream: %v", coordinator.GwID, err))
			return
		} else {
			glog.V(2).Infof("HWID %v: processing response", coordinator.GwID)
			err := srv.processSyncRPCResp(syncRPCResp, coordinator.GwID)
			if err != nil {
				coordinator.sendErr(fmt.Errorf("HWID %v: receiveFromStream: error processing stream response: %v", coordinator.GwID, err))
				return
			}
		}
	}
}

// processSyncRPCResp processes a SyncRPC response. It will either handle a
// heartbeat or call upon the broker to send the response to the HTTP server.
//
// Returning err indicates to end the bidirectional stream.
func (srv *SyncRPCService) processSyncRPCResp(resp *protos.SyncRPCResponse, hwId string) error {
	if resp.HeartBeat {
		err := directoryd.MapHWIDToHostname(hwId, srv.hostName)
		if err != nil {
			// Cannot persist <gwId, hostName> so nobody can send things to this
			// gateway use the stream, therefore return err to end the stream.
			return err
		}
	} else if resp.ReqId > 0 {
		err := srv.broker.ProcessGatewayResponse(resp)
		if err != nil {
			// No need to end the stream, just log the error.
			glog.Errorf("HWID %v: error processing gateway response: %v", hwId, err)
		}
	} else {
		glog.Errorf("HWID %v: cannot send a non-heartbeat with invalid ReqId", hwId)
	}
	return nil
}

// sendErr tries to send the err to stream coordinator's error channel.
// After timing out trying to send, instead log the error and return.
func (streamCoordinator *streamCoordinator) sendErr(err error) {
	select {
	case streamCoordinator.ErrChan <- err:
		return
	case <-time.After(time.Second):
		if err == nil {
			glog.Infof("HWID %v: client closed stream by sending EOF, cleaning up", streamCoordinator.GwID)
		} else {
			glog.Errorf(err.Error())
		}
		return
	}
}

type WrappedSyncResponse struct {
	Resp *protos.SyncRPCResponse
	Err  error
}

// RecvWithContext runs f and returns its error. If the context is canceled or
// times out first, it returns the context's error instead.
// See https://github.com/grpc/grpc-go/issues/1229#issuecomment-300938770
func RecvWithContext(ctx context.Context, f func() (*protos.SyncRPCResponse, error)) (*protos.SyncRPCResponse, error) {
	wrappedRespChan := make(chan WrappedSyncResponse, 1)
	go func() {
		resp, err := f()
		wrappedRespChan <- WrappedSyncResponse{resp, err}
		close(wrappedRespChan)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case wrappedResp := <-wrappedRespChan:
		return wrappedResp.Resp, wrappedResp.Err
	}
}

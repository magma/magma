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

package servicers

import (
	"context"
	"fmt"
	"sync"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/multiplex"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/errors"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

// How CentralSessionControllers works
//
// CentralSessionControllers holds an slice of N Controllers. Each controller uses an particular diameter
// client to a specific PCRF and OCS. Each diameter client is configured with its own src and dst port
// Subscribers are forwarded to a different controller using algorithm defined in multiplex object
//
// CreateSession and Terminate session gRPC message come per subscriber so they are forwarded to
// the controller directly
//
// UpdateSession gRPC message comes in groups, so before sending the request to each controller we look
// into each Credit and Policy request and forward it to the right controller depending on multiplex
//
// Health, Enable and Disable returns error if any of the controllers return errors. No partial results are give.

// CentralSessionControllerServerWithHealth is an interface just to group CentralSessionControllerServer
// and ServiceHealthServer. This is used by NewCentralSessionControllerWithHealth to be able to return
// either CentralSessionControllers or CentralSessionController (without S)
type CentralSessionControllerServerWithHealth interface {
	protos.CentralSessionControllerServer
	fegprotos.ServiceHealthServer
}

type CentralSessionControllers struct {
	centralControllers []*CentralSessionController
	multiplexor        multiplex.Multiplexor
}

type ControllerParam struct {
	CreditClient gy.CreditClient
	PolicyClient gx.PolicyClient
	Config       *SessionControllerConfig
}

// NewCentralSessionControllers creates centralControllers which is a slice of centralController.
// This should be used only if more than one server is configred
func NewCentralSessionControllers(
	controlParam []*ControllerParam,
	dbClient policydb.PolicyDBClient,
	mux multiplex.Multiplexor,
) *CentralSessionControllers {
	totalLen := len(controlParam)
	controllers := make([]*CentralSessionController, 0, totalLen)
	for _, cp := range controlParam {
		singleController := NewCentralSessionController(cp.CreditClient, cp.PolicyClient, dbClient, cp.Config)
		controllers = append(controllers, singleController)
	}
	return &CentralSessionControllers{
		centralControllers: controllers,
		multiplexor:        mux,
	}
}

// NewCentralSessionControllerDefaultMultiplesWithHealth returns a different type of controller depending on the amount
// of servers configured. In case only one server is configured, there is no need to calculate where this
// subscriber should be sent, so in that case we return CentralSessionController (without S). In case of multiple servers
// configured, it creates a CentralSessionControllers and uses a **StaticMultiplexByIMSI** as a multiplexor
func NewCentralSessionControllerDefaultMultiplexWithHealth(
	controlParam []*ControllerParam,
	dbClient policydb.PolicyDBClient,
) (CentralSessionControllerServerWithHealth, error) {
	if len(controlParam) == 1 {
		cp := controlParam[0]
		return NewCentralSessionController(cp.CreditClient, cp.PolicyClient, dbClient, cp.Config), nil
	}
	mux, err := multiplex.NewStaticMultiplexByIMSI(len(controlParam))
	if err != nil {
		return nil, err
	}
	return NewCentralSessionControllers(controlParam, dbClient, mux), nil
}

// CreateSession begins a UE session by requesting rules from PCEF
// and credit from OCS (if RatingGroup is present) and returning them.
func (srv *CentralSessionControllers) CreateSession(
	ctx context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	subs := request.GetCommonContext().GetSid()
	if subs == nil || len(subs.GetId()) == 0 {
		return nil, fmt.Errorf("Create Session Request Request malformed. Missing Subscriber.id")
	}
	controller, err := getControllerPerKey(
		srv.centralControllers,
		srv.multiplexor,
		multiplex.NewContext().WithIMSI(subs.GetId()),
	)
	if err != nil {
		return nil, err
	}
	return controller.CreateSession(ctx, request)
}

// UpdateSession handles periodic updates from gateways that include quota
// exhaustion and terminations
func (srv *CentralSessionControllers) UpdateSession(
	ctx context.Context,
	request *protos.UpdateSessionRequest,
) (*protos.UpdateSessionResponse, error) {
	requestsByController, err := getUpdateSessionRequestPerController(request, srv.centralControllers, srv.multiplexor)
	if err != nil {
		return nil, err
	}
	jobs := make(chan *protos.UpdateSessionResponse)
	wg := sync.WaitGroup{}
	// Create and run N producers (N controllers)
	for controller, request := range requestsByController {
		wg.Add(1)
		controllerShadow, requestShadow := controller, request
		go func() {
			defer wg.Done()
			singleUpdateSessionResponse, err := controllerShadow.UpdateSession(ctx, requestShadow)
			if err != nil {
				glog.Errorf("UpdateSession returned and error: %s", err)
				return
			}
			jobs <- singleUpdateSessionResponse
		}()
	}

	// Create One consumer to collect the responses from the producers
	done := make(chan *protos.UpdateSessionResponse)
	go func() {
		mergedResponse := &protos.UpdateSessionResponse{
			Responses:             make([]*protos.CreditUpdateResponse, 0),
			UsageMonitorResponses: make([]*protos.UsageMonitoringUpdateResponse, 0),
		}
		for singleResponse := range jobs {
			mergedResponse.Responses =
				append(mergedResponse.Responses, singleResponse.Responses...)
			mergedResponse.UsageMonitorResponses =
				append(mergedResponse.UsageMonitorResponses, singleResponse.UsageMonitorResponses...)
		}
		done <- mergedResponse
	}()

	wg.Wait()
	close(jobs)
	responseUpdatetSession, ok := <-done
	close(done)
	if !ok {
		return nil, fmt.Errorf("Couldnt read from channel")
	}
	return responseUpdatetSession, nil
}

// TerminateSession handles a session termination by sending single ccr-t on gx sending ccr-t per
// rating group on gy.
func (srv *CentralSessionControllers) TerminateSession(
	ctx context.Context,
	request *protos.SessionTerminateRequest,
) (*protos.SessionTerminateResponse, error) {
	if request == nil || len(request.GetSessionId()) == 0 {
		return nil, fmt.Errorf("Could not terminate session")
	}
	// be aware that this is sessionID format, not IMSI!!
	controller, err := getControllerPerKey(
		srv.centralControllers,
		srv.multiplexor,
		multiplex.NewContext().WithSessionId(request.GetSessionId()),
	)
	if err != nil {
		return nil, err
	}
	return controller.TerminateSession(ctx, request)
}

// Disable closes all existing diameter connections and disables
// connection creation for the time specified in the request
func (srv *CentralSessionControllers) Disable(
	ctx context.Context,
	req *fegprotos.DisableMessage,
) (*orcprotos.Void, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil Disable Request")
	}
	for _, controller := range srv.centralControllers {
		// this will never error. Error was check on req == nil
		controller.Disable(ctx, req)
	}
	return &orcprotos.Void{}, nil
}

// Enable enables diameter connection creation and gets a connection to the
// diameter server(s). if creation is already enabled and a connection already
// exists, enable has no effect
func (srv *CentralSessionControllers) Enable(
	ctx context.Context,
	void *orcprotos.Void,
) (*orcprotos.Void, error) {
	multiError := errors.NewMulti()
	for i, controller := range srv.centralControllers {
		_, err := controller.Enable(ctx, void)
		multiError = multiError.AddFmt(err, "error(%d):", i+1)
	}
	return &orcprotos.Void{}, multiError.AsError()
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
func (srv *CentralSessionControllers) GetHealthStatus(
	ctx context.Context,
	void *orcprotos.Void,
) (*fegprotos.HealthStatus, error) {
	for _, controller := range srv.centralControllers {
		healthMessage, err := controller.GetHealthStatus(ctx, void)
		if err != nil || healthMessage.Health == fegprotos.HealthStatus_UNHEALTHY {
			return healthMessage, err
		}
	}
	return &fegprotos.HealthStatus{
		Health:        fegprotos.HealthStatus_HEALTHY,
		HealthMessage: "All metrics appear healthy",
	}, nil
}

// getUpdateSessionRequestPerController creates a new UpdateSessionRequest per
// controller depending on the IMSIS of each request
func getUpdateSessionRequestPerController(
	request *protos.UpdateSessionRequest,
	controllers []*CentralSessionController,
	mux multiplex.Multiplexor,
) (map[*CentralSessionController]*protos.UpdateSessionRequest, error) {
	controllersToRequest := make(map[*CentralSessionController]*protos.UpdateSessionRequest)

	//Gy - Credit
	for _, creditUpdate := range request.GetUpdates() {
		controller, err := getControllerPerKey(
			controllers, mux,
			multiplex.NewContext().WithSessionId(creditUpdate.GetSessionId()),
		)
		if err != nil {
			return nil, err
		}
		fillMapWithUpdateSessionRequestIfEmpty(controllersToRequest, controller)
		controllersToRequest[controller].Updates = append(controllersToRequest[controller].Updates, creditUpdate)
	}

	//Gx - Policy
	for _, usageM := range request.GetUsageMonitors() {
		controller, err := getControllerPerKey(
			controllers, mux,
			multiplex.NewContext().WithSessionId(usageM.GetSessionId()),
		)
		if err != nil {
			return nil, err
		}
		fillMapWithUpdateSessionRequestIfEmpty(controllersToRequest, controller)
		controllersToRequest[controller].UsageMonitors = append(controllersToRequest[controller].UsageMonitors, usageM)
	}
	return controllersToRequest, nil
}

func fillMapWithUpdateSessionRequestIfEmpty(
	controllersToRequest map[*CentralSessionController]*protos.UpdateSessionRequest,
	controller *CentralSessionController) {
	_, found := controllersToRequest[controller]
	if !found {
		controllersToRequest[controller] = &protos.UpdateSessionRequest{}
	}
}

// getControllerPerKey provides the controllerId on a given selector (selector may include IMSI numeric, IMSIstr and Session ID)
func getControllerPerKey(
	controllers []*CentralSessionController,
	mux multiplex.Multiplexor,
	muxCtx *multiplex.Context,
) (*CentralSessionController, error) {
	index, err := mux.GetIndex(muxCtx)
	if err != nil {
		return nil, err
	}
	if index >= len(controllers) {
		return nil, fmt.Errorf("Index %d is bigger than the amount of controllers %d", index, len(controllers))
	}
	return controllers[index], nil
}

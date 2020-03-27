/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

/*
* ##########################################
* How CentralSessionControllers works
*
* CentralSessionControllers holds an slice of N Controllers. Each controller uses an particular diameter
* client to an specific PCRF and OCS. Each diameter client is configured with its own src and dst port
*
* Subscribers are forwarded to a different controller based on their IMSI that comes either in IMSI form on
* CreateSession or as SessionId (IMSI######-1234) for UpdateSession and Terminate session
*
* CreateSession and Terminate session gRPC message come per IMSI so they are forwarded to the controller directly
*
* UpdateSession gRPC message comes in groups, so before sending to each controller we look into each Credit and Policy
* request and forward it to the right controller depending on IMSI
*
* Health, Enable and Disable returns error if any of the controllers return errors. No partial results are give.
* ##########################################
 */

type CentralSessionControllers struct {
	centralControllers []*CentralSessionController
}

func NewCentralSessionControllers(
	creditClient []gy.CreditClient,
	policyClient []gx.PolicyClient,
	dbClient policydb.PolicyDBClient,
	cfg []*SessionControllerConfig,
) *CentralSessionControllers {
	totalLen := len(creditClient)
	if totalLen != len(policyClient) || totalLen != len(cfg) {
		panic("Same size required for CreditClient (Gy), PolicyClient (Gx) and SessionControllersConfig")
	}

	controllers := make([]*CentralSessionController, 0)
	centralControllers := CentralSessionControllers{}
	for n := 0; n < len(creditClient); n++ {
		singleController := NewCentralSessionController(creditClient[n], policyClient[n], dbClient, cfg[n])
		controllers = append(controllers, singleController)
	}
	centralControllers.centralControllers = controllers
	return &centralControllers
}

func NewCentralSessionControllers_SingleServer(
	creditClient gy.CreditClient,
	policyClient gx.PolicyClient,
	dbClient policydb.PolicyDBClient,
	cfg *SessionControllerConfig,
) *CentralSessionControllers {
	return NewCentralSessionControllers(
		[]gy.CreditClient{creditClient}, []gx.PolicyClient{policyClient},
		dbClient, []*SessionControllerConfig{cfg})
}

// CreateSession begins a UE session by requesting rules from PCEF
// and credit from OCS (if RatingGroup is present) and returning them.
func (srv *CentralSessionControllers) CreateSession(
	ctx context.Context,
	request *protos.CreateSessionRequest,
) (*protos.CreateSessionResponse, error) {
	subs := request.GetSubscriber()
	if subs == nil || len(subs.GetId()) == 0 {
		return nil, fmt.Errorf("Create Session Request Request malformed. Missing Subscriber.id")
	}
	imsi := subs.GetId()
	controller, err := getControllerFromImsi(imsi, srv.centralControllers)
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
	requestsByController, err := getUpdateSessionRequestPerController(request, srv.centralControllers)
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
	imsi, err := parseImsiFromSessionId(request.GetSessionId())
	if err != nil {
		return nil, err
	}
	controller, err := getControllerFromImsi(imsi, srv.centralControllers)
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
	var hasErrors []string
	for _, controller := range srv.centralControllers {
		_, err := controller.Disable(ctx, req)
		if err != nil {
			hasErrors = append(hasErrors, err.Error())
		}
	}
	if hasErrors != nil {
		return nil, fmt.Errorf("Errors found while disabling SessionProxy: %s",
			fmt.Errorf(strings.Join(hasErrors, "\n")))
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
	var hasErrors []string
	for _, controller := range srv.centralControllers {
		_, err := controller.Enable(ctx, void)
		if err != nil {
			hasErrors = append(hasErrors, err.Error())
		}
	}
	if hasErrors != nil {
		return nil, fmt.Errorf("Errors found while disabling SessionProxy: %s",
			fmt.Errorf(strings.Join(hasErrors, "\n")))
	}
	return &orcprotos.Void{}, nil
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the service
// TODO: report per each service
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

// splitUpdateSessionRequest creates a new UpdateSessionRequest per controller depending on the IMSIS of
// each request
func getUpdateSessionRequestPerController(
	request *protos.UpdateSessionRequest,
	controllers []*CentralSessionController,
) (map[*CentralSessionController]*protos.UpdateSessionRequest, error) {
	controllersToRequest := make(map[*CentralSessionController]*protos.UpdateSessionRequest)
	//Gy - Credit
	for _, creditUpdate := range request.GetUpdates() {
		controller, err := getControllerFromSessionId(controllers, creditUpdate.GetSessionId())
		if err != nil {
			return nil, err
		}
		_, found := controllersToRequest[controller]
		if !found {
			controllersToRequest[controller] = &protos.UpdateSessionRequest{}
		}
		controllersToRequest[controller].Updates = append(controllersToRequest[controller].Updates, creditUpdate)
	}
	//Gx - Policy
	for _, usageM := range request.GetUsageMonitors() {
		controller, err := getControllerFromSessionId(controllers, usageM.GetSessionId())
		if err != nil {
			return nil, err
		}
		_, found := controllersToRequest[controller]
		if found == false {
			controllersToRequest[controller] = &protos.UpdateSessionRequest{}
		}
		controllersToRequest[controller].UsageMonitors = append(controllersToRequest[controller].UsageMonitors, usageM)
	}
	return controllersToRequest, nil
}

// getControllerFromSessionId provides the controllerId on a given SessionID (note session ID contains the IMSI)
func getControllerFromSessionId(controllers []*CentralSessionController, sessionId string) (*CentralSessionController, error) {
	imsiStr, err := parseImsiFromSessionId(sessionId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	controller, err := getControllerFromImsi(imsiStr, controllers)
	if err != nil {
		return nil, err
	}
	return controller, nil
}

// getControllerFromImsi returns the controller that should serve that subscriber
func getControllerFromImsi(imsi string, controllers []*CentralSessionController) (*CentralSessionController, error) {
	if controllers == nil {
		return nil, fmt.Errorf("No controllers available.")
	}
	size := len(controllers)
	index, err := GetControllerIndexFromImsi(imsi, size)
	if err != nil {
		return nil, err
	}
	return controllers[index], nil
}

// parseImsiFromSessionId extracts IMSI from a sessionId. SessionId format is is considered
// to be IMMSIxxxxxx-1234, where xxxxx is the imsi to be extracted
func parseImsiFromSessionId(sessionId string) (string, error) {
	sessionId = strings.TrimPrefix(sessionId, "IMSI")
	data := strings.Split(sessionId, "-")
	if len(data) != 2 {
		return "", fmt.Errorf("Couldn't parse Subscrier ID from sessionID. Format should be IMISxxxxx-RandomNumber")
	}
	return data[0], nil
}

// GetControllerIndexFromImsi describes how we allocate the subscriber on the controllers
func GetControllerIndexFromImsi(imsi string, numberControlers int) (int, error) {
	//Remove the prefix if any
	imsi = strings.TrimPrefix(imsi, "IMSI")
	//IMSI must be parsed with 64 bitSize
	imsiUint, err := strconv.ParseUint(imsi, 10, 64)
	if err != nil {
		return -1, err
	}
	index := int(imsiUint % uint64(numberControlers))
	return index, nil
}

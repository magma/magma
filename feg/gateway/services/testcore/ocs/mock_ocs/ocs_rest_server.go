/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mock_ocs

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"
)

// OCSRestServer is an OCS queryable through a REST API to add accounts, set
// credits, and control settings. This is to be used for integration testing
type OCSRestServer struct {
	addr       string
	diamServer *OCSDiamServer
}

type AccountJSON struct {
	Imsi   string       `json:"imsi"`
	Credit []CreditJSON `json:"credit"`
}

type CreditJSON struct {
	ChargingKey uint32                     `json:"charging_key"`
	Unit        protos.CreditInfo_UnitType `json:"unit_type"`
	Volume      uint64                     `json:"volume"`
}

type OCSSettingsJSON struct {
	MaxUsageBytes *uint32 `json:"max_usage_bytes,omitempty"`
	MaxUsageTime  *uint32 `json:"max_usage_time,omitempty"`
	ValidityTime  *uint32 `json:"validity_time,omitempty"`
}

// NewOCSRestServer initializes a new REST server and diam server.
// Input: string address to start the REST server on
// 				 *sm.Settings to pass to the diameter server
//				 *DiameterServerConfig to send to the diameter server
// Output: *OCSRestServer
func NewOCSRestServer(
	addr string,
	diameterSettings *diameter.DiameterClientConfig,
	serverConfig *diameter.DiameterServerConfig,
) *OCSRestServer {
	return &OCSRestServer{
		addr: addr,
		diamServer: NewOCSDiamServer(
			diameterSettings,
			&OCSConfig{
				ServerConfig:  serverConfig,
				MaxUsageBytes: 2048, // default
				MaxUsageTime:  1000, // seconds default
				ValidityTime:  60,   // seconds default
				GyInitMethod:  gy.PerSessionInit,
			},
		),
	}
}

// Start initializes the REST server endpoints, begins the diameter server, and
// listens for HTTP requests, blocking
func (srv *OCSRestServer) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/account", postAccountHandler(srv)).Methods(http.MethodPost)
	router.HandleFunc("/account/{imsi}/credits", getAccountCreditsHandler(srv)).Methods(http.MethodGet)
	router.HandleFunc("/reset", postResetHandler(srv)).Methods(http.MethodPost)
	router.HandleFunc("/settings", postSettingsHandler(srv)).Methods(http.MethodPost)

	lis, err := srv.diamServer.StartListener()
	if err != nil {
		return err
	}

	go srv.diamServer.Start(lis)

	return http.ListenAndServe(srv.addr, router)
}

// postAccountHandler handles POST /account, adding a new account to the OCS
// Body:
// {
// 	"imsi": "1234",
// 	"credit": [
// 		{
// 			"charging_key": 1,
// 			"unit_type": 1, // 1 = bytes, 2 = time
// 			"volume": 1024, // bytes
// 		},
// 		...
// 	]
// }
// Return: 200 success, 404 error if error
func postAccountHandler(srv *OCSRestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account AccountJSON
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		srv.diamServer.CreateAccount(
			context.Background(),
			&lteprotos.SubscriberID{Id: account.Imsi},
		)
		for _, credit := range account.Credit {
			srv.diamServer.SetCredit(
				context.Background(),
				&protos.CreditInfo{
					Imsi:        account.Imsi,
					ChargingKey: credit.ChargingKey,
					Volume:      credit.Volume,
					UnitType:    credit.Unit,
				},
			)
		}
		glog.V(2).Infof("Added new account for subscriber %s", account.Imsi)
	}
}

// getAccountCreditsHandler handles GET /accounts/{imsi}/credits, returning the
// amount of credit for an account
// Body: None
// Returns:
// {
// 	"imsi": "1234",
// 	"credit": [
// 		{
// 			"charging_key": 1,
// 			"unit_type": 1, // 1 = bytes, 2 = time
// 			"volume": 1024, // bytes
// 		},
// 		...
// 	]
// }
func getAccountCreditsHandler(srv *OCSRestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		imsi := params["imsi"]
		credits, err := srv.diamServer.GetCredits(imsi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		creditJSONList := []CreditJSON{}
		for key, credit := range credits {
			creditJSONList = append(creditJSONList, CreditJSON{
				ChargingKey: key,
				Unit:        credit.Unit,
				Volume:      credit.Volume,
			})
		}
		account := AccountJSON{
			Imsi:   imsi,
			Credit: creditJSONList,
		}
		json.NewEncoder(w).Encode(account)
		glog.V(2).Infof("Updated credits for subscriber %s", imsi)
	}
}

// postResetHandler handles POST /reset, removing all the accounts tracked
// Body: None
// Returns: 200 success
func postResetHandler(srv *OCSRestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.diamServer.ClearSubscribers(context.Background(), &orcprotos.Void{})
		glog.V(2).Infof("Reset OCS")
	}
}

// postSettingsHandler handles POST /settings, changing the OCS return settings
// like max usage and validity time for credits
// Body:
// {
// 	"max_usage_bytes": 2048, // bytes, optional
// 	"max_usage_time": 1000, // seconds, optional
// 	"validity_time": 3600, // seconds, optional
// }
func postSettingsHandler(srv *OCSRestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings OCSSettingsJSON
		err := json.NewDecoder(r.Body).Decode(&settings)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		srv.diamServer.SetOCSSettings(
			context.Background(),
			&protos.OCSConfig{
				MaxUsageBytes: *settings.MaxUsageBytes,
				MaxUsageTime:  *settings.MaxUsageTime,
				ValidityTime:  *settings.ValidityTime,
			},
		)
		glog.V(2).Infof("Updated OCS settings")
	}
}

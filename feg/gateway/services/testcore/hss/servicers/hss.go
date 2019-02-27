/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"time"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/testcore/hss/crypto"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/protos"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/sm"
	"golang.org/x/net/context"
)

const hssProductName = "magma"

// HomeSubscriberServer tracks all the accounts needed for authenticating users.
type HomeSubscriberServer struct {
	store    storage.SubscriberStore
	Config   *mconfig.HSSConfig
	Milenage *crypto.MilenageCipher

	// authSqnInd is an index used in the array scheme described by 3GPP TS 33.102 Appendix C.1.2 and C.2.2.
	// SQN consists of two parts (SQN = SEQ||IND).
	AuthSqnInd uint64
}

// NewHomeSubscriberServer initializes a HomeSubscriberServer with an empty accounts map.
// Output: a new HomeSubscriberServer
func NewHomeSubscriberServer(store storage.SubscriberStore, config *mconfig.HSSConfig) (*HomeSubscriberServer, error) {
	milenage, err := crypto.NewMilenageCipher(config.LteAuthAmf)
	if err != nil {
		return nil, err
	}
	return &HomeSubscriberServer{
		store:    store,
		Config:   config,
		Milenage: milenage,
	}, nil
}

// AddSubscriber tries to add this subscriber to the server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func (srv *HomeSubscriberServer) AddSubscriber(ctx context.Context, req *lteprotos.SubscriberData) (*protos.Void, error) {
	err := srv.store.AddSubscriber(req)
	err = storage.ConvertStorageErrorToGrpcStatus(err)
	return &protos.Void{}, err
}

// GetSubscriberData looks up a subscriber by their Id.
// If the subscriber cannot be found, an error is returned instead.
// Input: The id of the subscriber to be looked up.
// Output: The data of the corresponding subscriber.
func (srv *HomeSubscriberServer) GetSubscriberData(ctx context.Context, req *lteprotos.SubscriberID) (*lteprotos.SubscriberData, error) {
	data, err := srv.store.GetSubscriberData(req.Id)
	err = storage.ConvertStorageErrorToGrpcStatus(err)
	return data, err
}

// UpdateSubscriber changes the data stored for an existing subscriber.
// If the subscriber cannot be found, an error is returned instead.
// Input: The new subscriber data to store
func (srv *HomeSubscriberServer) UpdateSubscriber(ctx context.Context, req *lteprotos.SubscriberData) (*protos.Void, error) {
	err := srv.store.UpdateSubscriber(req)
	err = storage.ConvertStorageErrorToGrpcStatus(err)
	return &protos.Void{}, err
}

// DeleteSubscriber deletes a subscriber by their Id.
// If the subscriber is not found, then this call is ignored.
// Input: The id of the subscriber to be deleted.
func (srv *HomeSubscriberServer) DeleteSubscriber(ctx context.Context, req *lteprotos.SubscriberID) (*protos.Void, error) {
	err := srv.store.DeleteSubscriber(req.Id)
	err = storage.ConvertStorageErrorToGrpcStatus(err)
	return &protos.Void{}, err
}

// Start begins the server and blocks, listening to the network
// Input: the address to start listening on
// Output: error if the server could not be started
func (srv *HomeSubscriberServer) Start() error {
	serverCfg := srv.Config.Server
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(serverCfg.DestHost),
		OriginRealm:      datatype.DiameterIdentity(serverCfg.DestRealm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(hssProductName),
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	}

	mux := sm.New(settings)
	mux.HandleFunc("ALL", handleUnknownMessage) // default handler
	mux.Handle(diam.AIR, srv.handleMessage(NewAIA))
	mux.Handle(diam.ULR, srv.handleMessage(NewULA))
	mux.Handle(diam.MAR, srv.handleMessage(NewMAA))
	mux.Handle(diam.SAR, srv.handleMessage(NewSAA))

	server := &diam.Server{
		Network: serverCfg.Protocol,
		Addr:    serverCfg.Address,
		Handler: mux,
	}
	listener, err := diam.MultistreamListen(serverCfg.Protocol, serverCfg.Address)
	if err != nil {
		return err
	}
	// If the port is 0 or not specified, overwriting the config allows the
	// chosen port to be known by the application.
	serverCfg.Address = listener.Addr().String()
	return server.Serve(listener)
}

/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"

	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
)

type directoryLookupServicer struct {
	store storage.DirectorydStorage
}

func NewDirectoryLookupServicer(store storage.DirectorydStorage) (protos.DirectoryLookupServer, error) {
	srv := &directoryLookupServicer{store: store}
	return srv, nil
}

func (d *directoryLookupServicer) GetHostnameForHWID(
	ctx context.Context, req *protos.GetHostnameForHWIDRequest,
) (*protos.GetHostnameForHWIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	hostname, err := d.store.GetHostnameForHWID(req.Hwid)
	res := &protos.GetHostnameForHWIDResponse{Hostname: hostname}

	return res, err
}

func (d *directoryLookupServicer) MapHWIDsToHostnames(ctx context.Context, req *protos.MapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	err = d.store.MapHWIDsToHostnames(req.HwidToHostname)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetIMSIForSessionID(
	ctx context.Context, req *protos.GetIMSIForSessionIDRequest,
) (*protos.GetIMSIForSessionIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	imsi, err := d.store.GetIMSIForSessionID(req.NetworkID, req.SessionID)
	res := &protos.GetIMSIForSessionIDResponse{Imsi: imsi}

	return res, err
}

func (d *directoryLookupServicer) MapSessionIDsToIMSIs(ctx context.Context, req *protos.MapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	err = d.store.MapSessionIDsToIMSIs(req.NetworkID, req.SessionIDToIMSI)

	return &protos.Void{}, err
}

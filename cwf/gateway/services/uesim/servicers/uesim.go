/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
	"os/exec"
	"strconv"

	"fbc/lib/go/radius"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	networkIDPlaceholder = "magma"
	blobTypePlaceholder  = "uesim"
	trafficMSS           = "1300"
	trafficSrvIP         = "192.168.129.42"
)

// UESimServer tracks all the UEs being simulated.
type UESimServer struct {
	store blobstore.BlobStorageFactory
	cfg   *UESimConfig
}

type UESimConfig struct {
	op                []byte
	amf               []byte
	radiusAuthAddress string
	radiusAcctAddress string
	radiusSecret      string
	brMac             string
}

// NewUESimServer initializes a UESimServer with an empty store map.
// Output: a new UESimServer
func NewUESimServer(factory blobstore.BlobStorageFactory) (*UESimServer, error) {
	config, err := getUESimConfig()
	if err != nil {
		return nil, err
	}
	return &UESimServer{
		store: factory,
		cfg:   config,
	}, nil
}

// AddUE tries to add this UE to the server.
// Input: The UE data which will be added.
func (srv *UESimServer) AddUE(ctx context.Context, ue *cwfprotos.UEConfig) (ret *protos.Void, err error) {
	ret = &protos.Void{}

	err = validateUEData(ue)
	if err != nil {
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	blob, err := ueToBlob(ue)
	store, err := srv.store.StartTransaction(nil)
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				err = errors.Wrap(err, "Error while rolling back transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		}
	}()

	err = store.CreateOrUpdate(networkIDPlaceholder, []blobstore.Blob{blob})
	return
}

// Authenticate triggers an authentication for the UE with the specified IMSI.
// Input: The IMSI of the UE to try to authenticate.
// Output: The resulting Radius packet returned by the Radius server.
func (srv *UESimServer) Authenticate(ctx context.Context, id *cwfprotos.AuthenticateRequest) (*cwfprotos.AuthenticateResponse, error) {
	eapIDResp, err := srv.CreateEAPIdentityRequest(id.GetImsi())
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaIDReq, err := radius.Exchange(context.Background(), eapIDResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaIDResp, err := srv.HandleRadius(id.GetImsi(), akaIDReq)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaChalReq, err := radius.Exchange(context.Background(), akaIDResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaChalResp, err := srv.HandleRadius(id.GetImsi(), akaChalReq)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	result, err := radius.Exchange(context.Background(), akaChalResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	resultBytes, err := result.Encode()
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, errors.Wrap(err, "Error encoding Radius packet")
	}
	radiusPacket := &cwfprotos.AuthenticateResponse{RadiusPacket: resultBytes}

	return radiusPacket, nil
}

func (srv *UESimServer) Disconnect(ctx context.Context, id *cwfprotos.DisconnectRequest) (*cwfprotos.DisconnectResponse, error) {
	radiusP, err := srv.MakeAccountingStopRequest()
	if err != nil {
		return nil, errors.Wrap(err, "Error making Accounting Stop Radius message")
	}
	response, err := radius.Exchange(context.Background(), radiusP, srv.cfg.radiusAcctAddress)
	if err != nil {
		return nil, errors.Wrap(err, "Error exchanging Radius message")
	}
	encoded, err := response.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "Error encoding Radius packet")
	}
	return &cwfprotos.DisconnectResponse{RadiusPacket: encoded}, nil
}

func (srv *UESimServer) GenTraffic(ctx context.Context, req *cwfprotos.GenTrafficRequest) (*cwfprotos.GenTrafficResponse, error) {
	if req == nil {
		return &cwfprotos.GenTrafficResponse{}, fmt.Errorf("Nil GenTrafficRequest provided")
	}
	var cmd *exec.Cmd

	argList := []string{"--json", "--get-server-output", "-c", trafficSrvIP, "-M", trafficMSS}
	if req.Volume != nil {
		argList = append(argList, []string{"-n", req.Volume.Value}...)
	}

	if req.ReverseMode {
		argList = append(argList, "-R")
	}

	if req.Bitrate != nil {
		argList = append(argList, []string{"-b", req.Bitrate.Value}...)
	}

	if req.TimeInSecs != 0 {
		argList = append(argList, []string{"-t", strconv.FormatUint(req.TimeInSecs, 10)}...)
	}

	if req.ReportingIntervalInSecs != 0 {
		argList = append(argList, []string{"-i", strconv.FormatUint(req.ReportingIntervalInSecs, 10)}...)
	}

	cmd = exec.Command("iperf3", argList...)
	cmd.Dir = "/usr/bin"
	output, err := cmd.Output()
	if err != nil {
		glog.Info("args = ", argList)
		glog.Info("error = ", err)
		err = errors.Wrap(err, fmt.Sprintf("argList %v\n output %v", argList, string(output)))
	}
	return &cwfprotos.GenTrafficResponse{Output: output}, err
}

// Converts UE data to a blob for storage.
func ueToBlob(ue *cwfprotos.UEConfig) (blobstore.Blob, error) {
	marshaledUE, err := protos.Marshal(ue)
	if err != nil {
		return blobstore.Blob{}, err
	}
	return blobstore.Blob{
		Type:  blobTypePlaceholder,
		Key:   ue.GetImsi(),
		Value: marshaledUE,
	}, nil
}

// Converts a blob back into a UE config
func blobToUE(blob blobstore.Blob) (*cwfprotos.UEConfig, error) {
	ue := &cwfprotos.UEConfig{}
	err := protos.Unmarshal(blob.Value, ue)
	if err != nil {
		return nil, err
	}
	return ue, nil
}

// getUE gets the UE with the specified IMSI from the blobstore.
func getUE(blobStoreFactory blobstore.BlobStorageFactory, imsi string) (ue *cwfprotos.UEConfig, err error) {
	store, err := blobStoreFactory.StartTransaction(nil)
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				glog.Errorf("Error while rolling back transaction: %s", err)
			}
		}
	}()

	blob, err := store.Get(networkIDPlaceholder, storage.TypeAndKey{Type: blobTypePlaceholder, Key: imsi})
	if err != nil {
		err = errors.Wrap(err, "Error getting UE with specified IMSI")
		return
	}
	ue, err = blobToUE(blob)
	return
}

// ConvertStorageErrorToGrpcStatus converts a UE error into a gRPC status error.
func ConvertStorageErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Unknown, err.Error())
}

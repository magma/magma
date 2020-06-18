/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// session_manager package defines local session manager client API
package session_manager

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"magma/feg/gateway/registry"
	"magma/lte/cloud/go/protos"
)

type sessionManagerClient struct {
	protos.LocalSessionManagerClient
}

// getSessionManagerClient is a utility function to get a RPC connection to the
// Local SessionManager service
func getSessionManagerClient() (*sessionManagerClient, error) {
	conn, err := registry.GetConnection(registry.SESSION_MANAGER)
	if err != nil {
		errMsg := fmt.Sprintf("Local SessionManager client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &sessionManagerClient{protos.NewLocalSessionManagerClient(conn)}, err
}

func ReportRuleStats(in *protos.RuleRecordTable) error {
	if in == nil {
		return errors.New("Nil RuleRecordTable Request")
	}
	cli, err := getSessionManagerClient()
	if err != nil {
		return err
	}
	_, err = cli.ReportRuleStats(context.Background(), in)
	return err
}

func CreateSession(in *protos.LocalCreateSessionRequest) (*protos.LocalCreateSessionResponse, error) {
	if in == nil {
		return nil, errors.New("Nil LocalCreateSessionRequest")
	}
	cli, err := getSessionManagerClient()
	if err != nil {
		return nil, err
	}
	return cli.CreateSession(context.Background(), in)
}

func EndSession(in *protos.LocalEndSessionRequest) (*protos.LocalEndSessionResponse, error) {
	if in == nil {
		return nil, errors.New("Nil LocalEndSessionRequest")
	}
	cli, err := getSessionManagerClient()
	if err != nil {
		return nil, err
	}
	return cli.EndSession(context.Background(), in)
}

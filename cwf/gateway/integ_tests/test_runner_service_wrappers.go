/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"context"
	"fmt"

	"magma/cwf/gateway/registry"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/hss"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/protos"
	registryTestUtils "magma/orc8r/cloud/go/test_utils"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Wrapper for GRPC Client functionality
type hssClient struct {
	fegprotos.HSSConfiguratorClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type pcrfClient struct {
	fegprotos.MockPCRFClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type ocsClient struct {
	fegprotos.MockOCSClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type pipelinedClient struct {
	lteprotos.PipelinedClient
	cc *grpc.ClientConn
}

/**  ========== HSS Helpers ========== **/
// getHSSClient is a utility function to getHSSClient a RPC connection to a
// remote HSS service.
func getHSSClient() (*hssClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockHSSRemote)
	if err != nil {
		errMsg := fmt.Sprintf("HSS client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &hssClient{
		fegprotos.NewHSSConfiguratorClient(conn),
		conn,
	}, err
}

// addSubscriber tries to add this subscriber to the HSS server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func addSubscriberToHSS(sub *lteprotos.SubscriberData) error {
	err := hss.VerifySubscriberData(sub)
	if err != nil {
		errMsg := fmt.Errorf("Invalid AddSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.AddSubscriber(context.Background(), sub)
	return err
}

func deleteSubscribersFromHSS(subscriberID string) error {
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.DeleteSubscriber(context.Background(), &lteprotos.SubscriberID{Id: subscriberID})
	return err
}

/**  ========== PCRF Helpers ========== **/
// getPCRFClient is a utility function to get an RPC connection to a
// remote PCRF service.
func getPCRFClient() (*pcrfClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockPCRFRemote)
	if err != nil {
		errMsg := fmt.Sprintf("PCRF client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &pcrfClient{
		MockPCRFClient: fegprotos.NewMockPCRFClient(conn),
		cc:             conn,
	}, err
}

// addSubscriber tries to add this subscriber to the PCRF server.
// Input: The subscriber data which will be added.
func addSubscriberToPCRF(sub *lteprotos.SubscriberID) error {
	cli, err := getPCRFClient()
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

func clearSubscribersFromPCRF() error {
	cli, err := getPCRFClient()
	if err != nil {
		return err
	}
	_, err = cli.ClearSubscribers(context.Background(), &protos.Void{})
	return err
}

func addPCRFRules(rules *fegprotos.AccountRules) error {
	cli, err := getPCRFClient()
	if err != nil {
		return err
	}
	_, err = cli.SetRules(context.Background(), rules)
	return err
}

func addPCRFUsageMonitors(monitorInfo *fegprotos.UsageMonitorInfo) error {
	cli, err := getPCRFClient()
	if err != nil {
		return err
	}
	_, err = cli.SetUsageMonitors(context.Background(), monitorInfo)
	return err
}

/**  ========== OCS Helpers ========== **/
// getOCSClient is a utility function to an RPC connection to a
// remote OCS service.
func getOCSClient() (*ocsClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockOCSRemote)
	if err != nil {
		errMsg := fmt.Sprintf("PCRF client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &ocsClient{
		fegprotos.NewMockOCSClient(conn),
		conn,
	}, err
}

// addSubscriber tries to add this subscriber to the OCS server.
// Input: The subscriber data which will be added.
func addSubscriberToOCS(sub *lteprotos.SubscriberID) error {
	cli, err := getOCSClient()
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

func clearSubscribersFromOCS() error {
	cli, err := getOCSClient()
	if err != nil {
		return err
	}
	_, err = cli.ClearSubscribers(context.Background(), &protos.Void{})
	return err
}

/**  ========== Pipelined Helpers ========== **/
// getPipelinedClient is a utility function to an RPC connection to a
// remote OCS service.
func getPipelinedClient() (*pipelinedClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registryTestUtils.GetConnectionWithAuthority(PipelinedRemote)
	if err != nil {
		errMsg := fmt.Sprintf("Pipelined client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &pipelinedClient{
		lteprotos.NewPipelinedClient(conn),
		conn,
	}, err
}

func deactivateSubscriberFlows(imsi string) error {
	cli, _ := getPipelinedClient()
	_, err := cli.DeactivateFlows(context.Background(), &lteprotos.DeactivateFlowsRequest{
		Sid:     &lteprotos.SubscriberID{Id: imsi},
		RuleIds: []string{makeRuleIDFromIMSI(imsi)},
	})
	return err
}

func getPolicyUsage() (*lteprotos.RuleRecordTable, error) {
	client, _ := getPipelinedClient()
	return client.GetPolicyUsage(context.Background(), &protos.Void{})
}

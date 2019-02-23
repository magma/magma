/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package client provides a thin client for contacting the subscriberdb service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package subscriberdb

import (
	"fmt"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const ServiceName = "SUBSCRIBERDB"

// Utility function to get a RPC connection to the subscriberdb service
func getSubscriberdbClient() (
	lteprotos.SubscriberDBControllerClient, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		glog.Errorf("Subscriberdb client initialization error: %s", err)
		return nil, nil, fmt.Errorf(
			"Subscriberdb client initialization error: %s", err)
	}
	return lteprotos.NewSubscriberDBControllerClient(conn), conn, err
}

// Add a new subscriber.
// The subscriber must not be existing already.
func AddSubscriber(networkId string, sd *lteprotos.SubscriberData) error {
	sd.NetworkId = &protos.NetworkID{Id: networkId}
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = client.AddSubscriber(context.Background(), sd); err != nil {
		glog.Errorf("[Network: %s] AddSubscriber error: %s", networkId, err)
		return err
	}
	return nil
}

// Get the subscriber data.
func GetSubscriber(networkId string, subscriberId string) (
	*lteprotos.SubscriberData, error) {
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lookup := &lteprotos.SubscriberLookup{
		NetworkId: &protos.NetworkID{Id: networkId},
		Sid:       lteprotos.SidFromString(subscriberId)}
	data, err := client.GetSubscriberData(context.Background(), lookup)
	if err != nil {
		glog.Errorf("[Network: %s, Sub: %s] GetSubscriberData error: %s",
			networkId, subscriberId, err)
		return nil, err
	}
	return data, nil
}

// Update the subscriber info.
func UpdateSubscriber(networkId string, sd *lteprotos.SubscriberData) error {
	sd.NetworkId = &protos.NetworkID{Id: networkId}
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = client.UpdateSubscriber(context.Background(), sd); err != nil {
		glog.Errorf("[Network: %s] UpdateSubscriber error: %s", networkId, err)
		return err
	}
	return nil
}

// Delete the subscriber.
func DeleteSubscriber(networkId string, subscriberId string) error {
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	lookup := &lteprotos.SubscriberLookup{
		NetworkId: &protos.NetworkID{Id: networkId},
		Sid:       lteprotos.SidFromString(subscriberId)}
	if _, err := client.DeleteSubscriber(context.Background(), lookup); err != nil {
		glog.Errorf("[Network: %s, Sub: %s] DeleteSubscribererror: %s",
			networkId, subscriberId, err)
		return err
	}
	return nil
}

// List all existing subscribers.
// Returns an array of subscriber ids.
func ListSubscribers(networkId string) ([]string, error) {
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	subs, err := client.ListSubscribers(
		context.Background(),
		&protos.NetworkID{Id: networkId})
	if err != nil {
		glog.Errorf("ListSubscribers error: %s", err)
		return nil, err
	}
	sids := subs.GetSids()
	ret := make([]string, len(sids))
	for i := range sids {
		ret[i] = lteprotos.SidString(sids[i])
	}
	return ret, nil
}

func GetAllSubscriberData(networkId string) ([]*lteprotos.SubscriberData, error) {
	client, conn, err := getSubscriberdbClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := client.GetAllSubscriberData(context.Background(), &protos.NetworkID{Id: networkId})
	if err != nil {
		return []*lteprotos.SubscriberData{}, err
	}
	return response.Subscribers, nil
}

/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"fmt"
	"sort"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SubscribersTable = "subscriberdb"

type SubscriberDBStorage struct {
	store datastore.Api
}

func getSubscriberDBTableName(networkID string) string {
	return datastore.GetTableName(networkID, SubscribersTable)
}

func NewSubscriberDBStorage(ds datastore.Api) (*SubscriberDBStorage, error) {
	if ds == nil {
		return nil, fmt.Errorf("Nil SubscriberDBStorage datastore")
	}
	s := &SubscriberDBStorage{ds}
	return s, nil
}

func (s *SubscriberDBStorage) AddSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error) {
	sid := lteprotos.SidString(subs.Sid)
	table := getSubscriberDBTableName(subs.NetworkId.Id)
	glog.V(2).Info("Adding subscriber ", sid, " to ", table)

	if _, _, err := s.store.Get(table, sid); err == nil {
		errMsg := fmt.Sprintf("Subscriber %s already exists ", sid)
		glog.Error(errMsg)
		return nil, status.Error(codes.AlreadyExists, errMsg)
	}

	// Marshal the protobuf and store the byte stream in the Datastore
	value, err := proto.Marshal(subs)
	if err != nil {
		errMsg := fmt.Sprintf("Marshalling error on sid: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}

	// Add the subscriber to the Datastore
	if err = s.store.Put(table, sid, value); err != nil {
		errMsg := fmt.Sprintf("Error adding subscriber: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}
	return &protos.Void{}, nil
}

func (s *SubscriberDBStorage) DeleteSubscriber(lookup *lteprotos.SubscriberLookup) (*protos.Void, error) {
	sid := lteprotos.SidString(lookup.Sid)
	table := getSubscriberDBTableName(lookup.NetworkId.Id)
	glog.V(2).Info("Deleting subscriber ", sid, " from ", table)

	if err := s.store.Delete(table, sid); err != nil {
		errMsg := fmt.Sprintf("Error deleting subscriber: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}
	return &protos.Void{}, nil
}

func (s *SubscriberDBStorage) UpdateSubscriber(subs *lteprotos.SubscriberData) (*protos.Void, error) {
	sid := lteprotos.SidString(subs.Sid)
	table := getSubscriberDBTableName(subs.NetworkId.Id)
	glog.V(2).Info("Updating subscriber ", sid, " in ", table)

	// Marshal the protobuf and store the byte stream in the Datastore
	value, err := proto.Marshal(subs)
	if err != nil {
		errMsg := fmt.Sprintf("Marshalling error on sid: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}

	// Update the subscriber in the Datastore
	if err = s.store.Put(table, sid, value); err != nil {
		errMsg := fmt.Sprintf("Error updating subscriber: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}
	return &protos.Void{}, nil
}

func (s *SubscriberDBStorage) GetSubscriberData(lookup *lteprotos.SubscriberLookup) (*lteprotos.SubscriberData, error) {
	sid := lteprotos.SidString(lookup.Sid)
	subs := lteprotos.SubscriberData{}
	table := getSubscriberDBTableName(lookup.NetworkId.Id)

	value, _, err := s.store.Get(table, sid)
	if err != nil {
		errMsg := fmt.Sprintf("Error fetching subscriber: %s, %s", sid, err)
		glog.Error(errMsg)

		if err == datastore.ErrNotFound {
			return nil, status.Error(codes.NotFound, errMsg)
		}
		return nil, status.Error(codes.Aborted, errMsg)
	}
	if err = proto.Unmarshal(value, &subs); err != nil {
		errMsg := fmt.Sprintf("Unmarshalling error on sid: %s, %s", sid, err)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}
	return &subs, nil
}

func (s *SubscriberDBStorage) ListSubscribers(networkID *protos.NetworkID) (*lteprotos.SubscriberIDSet, error) {
	table := getSubscriberDBTableName(networkID.Id)
	keys, err := s.store.ListKeys(table)
	if err != nil {
		errMsg := fmt.Sprintf("Error listing Subscribers %s, for network: %s", err, networkID.Id)
		glog.Error(errMsg)
		return nil, status.Error(codes.Aborted, errMsg)
	}

	sids := make([]*lteprotos.SubscriberID, 0, len(keys))
	for _, key := range keys {
		if sid, err := lteprotos.SidProto(key); err == nil {
			sids = append(sids, sid)
		} else {
			glog.Warningf("Unable to convert sid %s string to proto struct", key)
		}
	}
	return &lteprotos.SubscriberIDSet{Sids: sids}, nil
}

func (s *SubscriberDBStorage) GetAllSubscriberData(networkID *protos.NetworkID) (*lteprotos.GetAllSubscriberDataResponse, error) {
	table := getSubscriberDBTableName(networkID.Id)
	sids, err := s.store.ListKeys(table)
	if err != nil {
		msg := fmt.Sprintf("Error listing all subscriber IDs: %s", err)
		glog.Error(msg)
		return nil, status.Error(codes.Aborted, msg)
	}

	subscribersBySid, err := s.store.GetMany(table, sids)
	if err != nil {
		msg := fmt.Sprintf("Error getting all subscribers: %s", err)
		glog.Error(msg)
		return nil, status.Error(codes.Aborted, msg)
	}

	allSubProtos := make([]*lteprotos.SubscriberData, 0, len(subscribersBySid))
	for _, sid := range *getSortedKeys(&subscribersBySid) {
		marshaledDatum := subscribersBySid[sid].Value
		sub := &lteprotos.SubscriberData{}
		if err = proto.Unmarshal(marshaledDatum, sub); err != nil {
			msg := fmt.Sprintf("Could not unmarshal subscriber data: %s", err)
			glog.Error(msg)
			return nil, status.Error(codes.Aborted, msg)
		}
		allSubProtos = append(allSubProtos, sub)
	}
	return &lteprotos.GetAllSubscriberDataResponse{Subscribers: allSubProtos}, nil
}

func getSortedKeys(in *map[string]datastore.ValueWrapper) *[]string {
	ret := make([]string, 0, len(*in))
	for k := range *in {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return &ret
}

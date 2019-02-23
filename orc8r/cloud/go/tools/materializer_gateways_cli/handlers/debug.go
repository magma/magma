/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

var debugGatewayID string

func init() {
	cmd := Commands.Add(
		"debug",
		"Print all the events for a given gateway in order",
		GetRecordedEventsForGateway,
	)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: envdir /var/opt/magma/envdir %s debug [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\tRecommended to redirect output to a file and perform further analysis using jq\n")
		f.PrintDefaults()
	}
	f.StringVar(&debugGatewayID, "gateway", "", "Gateway ID to debug")
}

func GetRecordedEventsForGateway(cmd *commands.Command, args []string) int {
	consumer, err := newConsumer()
	if err != nil {
		glog.Errorf("Error initializing stream consumer: %s", err)
		return 1
	}
	defer consumer.Close()

	err = consumer.Subscribe("magma.public.gatewaystates", nil)
	if err != nil {
		glog.Errorf("Error subscribing to recorded state topic: %s", err)
		return 1
	}

	// All the updates for a given gateway will be on the same partition so we
	// will have ordering here
	decoder := streaming.NewDecoder()
	recordCount := 0
	matchedCount := 0
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			glog.Errorf("Error reading message from consumer: %s", err)
			return 1
		}

		recordCount++
		if recordCount%1000 == 0 {
			glog.Infof("Processed %d records", recordCount)
		}

		update, err := decoder.GetUpdateFromStateRecorderMessage(message)
		if err != nil {
			glog.Errorf("Error decoding update: %s", err)
			return 1
		}

		id := getAssociatedId(update.Payload)
		if id != debugGatewayID {
			continue
		}

		matchedCount++
		if matchedCount%100 == 0 {
			glog.Infof("Found %d matching records", matchedCount)
		}

		marshaled, err := json.MarshalIndent(update, "", "  ")
		if err != nil {
			glog.Errorf("Error marshaling update to print: %s", err)
			return 1
		}
		fmt.Println(string(marshaled))
		fmt.Println()
	}
}

func newConsumer() (streaming.Consumer, error) {
	clusterServers, ok := os.LookupEnv("CLUSTER_SERVERS")
	if !ok {
		return nil, fmt.Errorf("CLUSTER_SERVERS was not defined")
	}
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": clusterServers,
		"group.id":          "materializer-cli-debugger",
		// reset earliest and auto commit disabled without manual commit
		// means that we will always start from the beginning of the log
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
}

func getAssociatedId(payload streaming.UpdatePayload) string {
	switch payload.(type) {
	case *streaming.GatewayConfigUpdate:
		return payload.(*streaming.GatewayConfigUpdate).ConfigKey
	case *streaming.GatewayStatusUpdate:
		return payload.(*streaming.GatewayStatusUpdate).GatewayID
	case *streaming.GatewayRecordUpdate:
		return payload.(*streaming.GatewayRecordUpdate).GatewayID
	default:
		return ""
	}
}

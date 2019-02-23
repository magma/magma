/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/tools/commands"
)

const (
	defaultAuthKey = "\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"
	defaultAuthOpc = "\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"
)

var (
	cmdRegistry           = new(commands.Map) // manages the commands which this CLI supports
	subscriberID          string
	gsmSubscriptionActive bool
	lteSubscriptionActive bool
	networkID             string
	lteAuthNextSeq        uint64
	subProfile            string
	authKey               string
	authOpc               string
)

func main() {
	flag.Parse()
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s command [OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		cmdRegistry.Usage()
	}
	cmdName := flag.Arg(0)
	if len(flag.Args()) < 1 || cmdName == "" || cmdName == "help" {
		flag.Usage()
		os.Exit(1)
	}

	cmd := cmdRegistry.Get(cmdName)
	if cmd == nil {
		fmt.Println("\nInvalid Command: ", cmdName)
		flag.Usage()
		os.Exit(1)
	}
	args := os.Args[2:]
	cmd.Flags().Parse(args)
	if len(subscriberID) == 0 {
		println("Error: Subscriber ID missing")
		cmd.Usage()
		os.Exit(1)
	}
	os.Exit(cmd.Handle(args))
}

// getSubscriber handles the GET command (retrieves subscriber data from the hss)
func getSubscriber(_ *commands.Command, _ []string) int {
	client, err := connectToHss()
	if err != nil {
		fmt.Printf("Failed to connect to hss: %v\n", err)
		return 1
	}

	id := &lteprotos.SubscriberID{Id: subscriberID}
	data, err := client.GetSubscriberData(context.Background(), id)
	if err != nil {
		fmt.Printf("Failed to get subscriber data: %v\n", err)
		return 1
	}

	fmt.Printf("Retreived subscriber data: %v\n", data)
	return 0
}

// addSubscriber handles the ADD command (adds a new subscriber to the hss)
func addSubscriber(_ *commands.Command, _ []string) int {
	client, err := connectToHss()
	if err != nil {
		fmt.Printf("Failed to connect to hss: %v\n", err)
		return 1
	}

	_, err = client.AddSubscriber(context.Background(), getSubscriberData())
	if err != nil {
		fmt.Printf("Failed to add subscriber: %v\n", err)
		return 1
	}

	return 0
}

// updateSubscriber handles the UPDATE command (updates an existing subscriber in the hss)
func updateSubscriber(_ *commands.Command, _ []string) int {
	client, err := connectToHss()
	if err != nil {
		fmt.Printf("Failed to connect to hss: %v\n", err)
		return 1
	}

	_, err = client.UpdateSubscriber(context.Background(), getSubscriberData())
	if err != nil {
		fmt.Printf("Failed to update subscriber: %v\n", err)
		return 1
	}

	return 0
}

// deleteSubscriber handles the DEL command (deletes a subscriber from the hss)
func deleteSubscriber(_ *commands.Command, _ []string) int {
	client, err := connectToHss()
	if err != nil {
		fmt.Printf("Failed to connect to hss: %v\n", err)
		return 1
	}

	id := &lteprotos.SubscriberID{Id: subscriberID}
	_, err = client.DeleteSubscriber(context.Background(), id)
	if err != nil {
		fmt.Printf("Failed to delete subscriber: %v\n", err)
		return 1
	}

	return 0
}

func init() {
	getCmd := cmdRegistry.Add(
		"GET",
		"Retrieve subscriber data by id",
		getSubscriber)
	getFlags := getCmd.Flags()
	getFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], getCmd.Name(), getCmd.Name())
		getFlags.PrintDefaults()
	}
	getFlags.StringVar(&subscriberID, "subscriber_id", subscriberID, "IMSI of the subscriber to look up")

	addCmd := cmdRegistry.Add(
		"ADD",
		"Add a new subscriber",
		addSubscriber)
	addFlags := addCmd.Flags()
	addFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], addCmd.Name(), addCmd.Name())
		addFlags.PrintDefaults()
	}
	addSubscriberDataFlags(addFlags)

	updateCmd := cmdRegistry.Add(
		"UPDATE",
		"Update an existing subscriber",
		updateSubscriber)
	updateFlags := updateCmd.Flags()
	updateFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], updateCmd.Name(), updateCmd.Name())
		updateFlags.PrintDefaults()
	}
	addSubscriberDataFlags(updateFlags)

	delCmd := cmdRegistry.Add(
		"DEL",
		"Delete subscriber data",
		deleteSubscriber)
	delFlags := delCmd.Flags()
	delFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], delCmd.Name(), delCmd.Name())
		delFlags.PrintDefaults()
	}
	delFlags.StringVar(&subscriberID, "subscriber_id", subscriberID, "IMSI of the subscriber to delete")
}

// addSubscriberDataFlags adds all of the flags needed to fill a SubscriberData proto.
func addSubscriberDataFlags(flags *flag.FlagSet) {
	flags.StringVar(&subscriberID, "subscriber_id", subscriberID, "IMSI of the subscriber")
	flags.BoolVar(&gsmSubscriptionActive, "gsm_subscription_active", gsmSubscriptionActive, "Whether the gsm subscription is active")
	flags.BoolVar(&lteSubscriptionActive, "lte_subscription_active", lteSubscriptionActive, "Whether the lte subscription is active")
	flags.StringVar(&authKey, "auth_key", defaultAuthKey, "authentication key")
	flags.StringVar(&authOpc, "auth_opc", defaultAuthOpc, "Operator configuration field signed with authentication key")
	flags.StringVar(&networkID, "network_id", networkID, "Uniquely identifies the network")
	flags.Uint64Var(&lteAuthNextSeq, "lte_auth_next_seq", lteAuthNextSeq, "Next SEQ to be used for calculating the AUTN")
	flags.StringVar(&subProfile, "sub_profile", subProfile, "Subscription profile")
}

// connectToHss establishes a grpc connection to the hss configurator service.
func connectToHss() (protos.HSSConfiguratorClient, error) {
	conn, err := registry.GetConnection(testcore.MockHSSServiceName)
	if err != nil {
		return nil, err
	}
	client := protos.NewHSSConfiguratorClient(conn)
	return client, nil
}

// getSubscriberData uses the command line flag values to create a SubscriberData proto.
func getSubscriberData() *lteprotos.SubscriberData {
	return &lteprotos.SubscriberData{
		Sid: &lteprotos.SubscriberID{Id: subscriberID},
		Gsm: &lteprotos.GSMSubscription{State: getGSMSubscriptionState()},
		Lte: &lteprotos.LTESubscription{
			State:    getLTESubscriptionState(),
			AuthKey:  []byte(authKey),
			AuthOpc:  []byte(authOpc),
			AuthAlgo: lteprotos.LTESubscription_MILENAGE,
		},
		NetworkId:  &orcprotos.NetworkID{Id: networkID},
		State:      &lteprotos.SubscriberState{LteAuthNextSeq: lteAuthNextSeq},
		SubProfile: subProfile,
	}
}

// getGSMSubscriptionState uses the command line args to determine whether the gsm subscription is active
func getGSMSubscriptionState() lteprotos.GSMSubscription_GSMSubscriptionState {
	if gsmSubscriptionActive {
		return lteprotos.GSMSubscription_ACTIVE
	}
	return lteprotos.GSMSubscription_INACTIVE
}

// getLTESubscriptionState uses the command line args to determine whether the lte subscription is active
func getLTESubscriptionState() lteprotos.LTESubscription_LTESubscriptionState {
	if lteSubscriptionActive {
		return lteprotos.LTESubscription_ACTIVE
	}
	return lteprotos.LTESubscription_INACTIVE
}

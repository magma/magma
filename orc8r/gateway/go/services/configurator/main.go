/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package main - implementation of a stand alone configurator
package main

import (
	"flag"

	"github.com/golang/glog"

	"magma/gateway/services/configurator/service"
)

func main() {
	flag.Parse() // for glog
	updateNotifier := make(chan interface{})
	cfg := service.NewConfigurator(updateNotifier)
	go func() {
		for i := range updateNotifier {
			switch u := i.(type) {
			case service.UpdateCompletion:
				glog.Infof("mconfigs updated successfully for services: %v", u)
			default:
				glog.Errorf("unknown completion type: %T", u)
			}
		}
	}()

	if err := cfg.Start(); err != nil {
		glog.Fatalf("configurator start error: %v", err)
	}
}

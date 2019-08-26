/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package oc

import (
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

func ExampleConfig_Build() {
	var cfg struct {
		Census Config `group:"oc" namespace:"oc" env-namespace:"OC"`
	}

	err := os.Setenv("OC_VIEWS", "proc,http")
	if err != nil {
		log.Fatalf("settings census environ: %v", err)
	}

	if _, err := flags.Parse(&cfg); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("creating logger: %v", err)
	}

	census, err := cfg.Census.Build(WithLogger(logger))
	if err != nil {
		log.Fatalf("building census: %v", err)
	}
	defer census.Close()

	http.Handle("/metrics", census.StatsHandler)
	_ = http.ListenAndServe(":9100", nil)
	// - Output:
}

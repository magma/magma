// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	stdlog "log"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/store/sign/s3"

	"github.com/jessevdk/go-flags"
)

type cliFlags struct {
	Addr   string     `env:"ADDR" long:"addr" description:"the address to listen on"`
	S3     s3.Config  `group:"s3" namespace:"s3" env-namespace:"S3"`
	Log    log.Config `group:"log" namespace:"log" env-namespace:"LOG"`
	Census oc.Options `group:"oc" namespace:"oc" env-namespace:"OC"`
}

func main() {
	var cf cliFlags
	if _, err := flags.Parse(&cf); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	srv, _, err := NewServer(&cf)
	if err != nil {
		stdlog.Fatal(err)
	}
	go func() { _ = srv.ListenAndServe(cf.Addr) }()

	<-ctxutil.WithSignal(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	).Done()
	_ = srv.Shutdown(context.Background())
}

// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	stdlog "log"
	"net/url"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/jessevdk/go-flags"

	_ "github.com/facebookincubator/symphony/frontier/ent/runtime"
)

type (
	cliFlags struct {
		Addr            string     `env:"ADDR" long:"addr" description:"the address to listen on"`
		InventoryTarget target     `env:"INVENTORY_TARGET" long:"inventory-target" required:"true" description:"the url of the inventory (static files) to proxy to"`
		ProxyTarget     target     `env:"PROXY_TARGET" long:"proxy-target" required:"true" description:"the url to proxy to"`
		KeyPairs        []key      `env:"KEY_PAIRS" env-delim:"," long:"key-pairs" required:"true" description:"authentication / encryption key pairs"`
		Log             log.Config `group:"log" namespace:"log" env-namespace:"LOG"`
		Census          oc.Options `group:"oc" namespace:"oc" env-namespace:"OC"`
	}

	// target attaches flags methods to url.URL.
	target struct{ *url.URL }
	// key attaches flags methods to []byte.
	key []byte
)

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

func (t *target) UnmarshalFlag(value string) (err error) {
	if t.URL, err = url.Parse(value); err != nil {
		return &flags.Error{
			Type:    flags.ErrMarshal,
			Message: err.Error(),
		}
	}
	if t.Scheme == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "proxy target requires url schema",
		}
	}
	return nil
}

// nolint:unparam // implements flags.Unmarshaler
func (k *key) UnmarshalFlag(value string) error {
	*k = []byte(value)
	return nil
}

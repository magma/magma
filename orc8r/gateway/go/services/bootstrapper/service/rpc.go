/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of bootstrapper
package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"magma/orc8r/lib/go/registry"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

// GetBootstrapperCloudConnection initializes and returns Bootstrapper cloud grpc connection
func (b *Bootstrapper) GetBootstrapperCloudConnection() (*grpc.ClientConn, error) {
	// Do not use proxied connection for bootstrapper

	addrPieces := strings.Split(b.CpConfig.BootstrapAddr, ":")
	addr := fmt.Sprintf("%s:%d", addrPieces[0], b.CpConfig.BootstrapPort)

	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxLocalTimeoutSec*time.Second)
	defer cancel()
	opts := b.getGrpcOpts(b.CpConfig.ProxyCloudConnection)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	// if the proxy is not present, we should fail fast & retry with direct TLS connection
	if err != nil {
		firstErr := fmt.Errorf("Bootstrapper dial failure for address: %s; GRPC Dial error: %s", addr, err)
		if b.CpConfig.ProxyCloudConnection {
			log.Printf("%v; trying direct TLS cloud connection", firstErr)
			// Try to call cloud directly
			ctxTls, cancelTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
			defer cancelTls()
			addr = fmt.Sprintf("%s:%d", addrPieces[0], DefaultTLSBootstrapPort)
			opts = b.getGrpcOpts(false)
			conn, err = grpc.DialContext(ctxTls, addr, opts...)
			if err != nil {
				return conn, fmt.Errorf(
					"Bootstrapper TLS dial failure for address: %s; GRPC Dial error: %s", addr, err)
			}
		} else {
			return conn, firstErr // already direct cloud conn, fail
		}
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return conn, nil
}

func (b *Bootstrapper) getGrpcOpts(useProxy bool) []grpc.DialOption {
	var opts = []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 30 * time.Second,
		}),
		grpc.WithBlock(),
		grpc.WithAuthority(b.CpConfig.BootstrapAddr),
	}
	if useProxy {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// always try to add OS certs
		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("OS Cert Pool initialization error: %v", err)
			certPool = x509.NewCertPool()
		}
		// Add magma RootCA
		if rootCa, err := ioutil.ReadFile(b.CpConfig.RootCaFile); err == nil {
			if !certPool.AppendCertsFromPEM(rootCa) {
				log.Printf("Failed to append certificates from %s", b.CpConfig.RootCaFile)
			}
		} else {
			log.Printf("Cannot load Root CA from '%s': %v", b.CpConfig.RootCaFile, err)
		}
		var tlsCfg *tls.Config
		if len(certPool.Subjects()) > 0 {
			tlsCfg = &tls.Config{
				InsecureSkipVerify: false, // last resort - do not verify the server cert, but rely only on
				RootCAs:            certPool,
			}
		} else {
			tlsCfg = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	}
	return opts
}

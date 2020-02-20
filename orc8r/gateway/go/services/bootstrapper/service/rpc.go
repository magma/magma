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
	addrPieces := strings.Split(b.CpConfig.BootstrapAddr, ":")
	addr := fmt.Sprintf("%s:%d", addrPieces[0], b.CpConfig.BootstrapPort)

	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxLocalTimeoutSec*time.Second)
	defer cancel()
	proxied := b.CpConfig.ProxyCloudConnection
	opts := b.getGrpcOpts(proxied)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err == nil {
		return conn, ctx.Err()
	}
	// if the proxy is not present, we should with direct TLS connection to default and configured TLS ports
	firstErr := fmt.Errorf("Bootstrapper dial failure for address: %s; GRPC Dial error: %s", addr, err)
	if proxied {
		addr = fmt.Sprintf("%s:%d", addrPieces[0], DefaultTLSBootstrapPort)
		log.Printf("%v; trying secure connection to: %s", firstErr, addr)
		// Try to call cloud directly
		ctxTls, cancelTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
		defer cancelTls()
		opts = b.getGrpcOpts(false)
		conn, err = grpc.DialContext(ctxTls, addr, opts...)
		if err == nil {
			return conn, ctxTls.Err()
		}
		err = fmt.Errorf("Bootstrapper TLS dial failure for address: %s; GRPC Dial error: %s", addr, err)
		// final attempt, use direct cloud connection and configured bootstrapper port instead of default TLS port
		if b.CpConfig.BootstrapPort != DefaultTLSBootstrapPort {
			addr = fmt.Sprintf("%s:%d", addrPieces[0], b.CpConfig.BootstrapPort)
			log.Printf("%v; trying: %s", err, addr)
			ctx2Tls, cance2lTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
			defer cance2lTls()
			conn, err = grpc.DialContext(ctx2Tls, addr, opts...)
			if err == nil {
				return conn, ctx2Tls.Err()
			}
			err = fmt.Errorf("final Bootstrapper TLS dial failure for: %s; GRPC Dial error: %s", addr, err)
		}
		log.Print(err)
	}
	return conn, firstErr
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

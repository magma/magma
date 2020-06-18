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
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"

	"magma/gateway/config"
	"magma/orc8r/lib/go/registry"
)

// GetBootstrapperCloudConnection initializes and returns Bootstrapper cloud grpc connection
func (b *Bootstrapper) GetBootstrapperCloudConnection() (*grpc.ClientConn, error) {
	cfg := config.GetControlProxyConfigs()
	addrPieces := strings.Split(cfg.BootstrapAddr, ":")
	addr := fmt.Sprintf("%s:%d", addrPieces[0], cfg.BootstrapPort)

	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxLocalTimeoutSec*time.Second)
	defer cancel()
	proxied := cfg.ProxyCloudConnection
	opts := b.getGrpcOpts(proxied, cfg)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err == nil {
		return conn, ctx.Err()
	}
	// in case the proxy is not present, we should try again with direct TLS connection to the default as well as
	// the configured TLS ports
	firstErr := fmt.Errorf("Bootstrapper dial failure for address: %s; GRPC Dial error: %s", addr, err)
	if proxied {
		addr = fmt.Sprintf("%s:%d", addrPieces[0], DefaultTLSBootstrapPort)
		glog.Warningf("%v; trying secure connection to: %s", firstErr, addr)
		// Try to call cloud directly
		ctxTls, cancelTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
		defer cancelTls()
		opts = b.getGrpcOpts(false, cfg)
		conn, err = grpc.DialContext(ctxTls, addr, opts...)
		if err == nil {
			return conn, ctxTls.Err()
		}
		err = fmt.Errorf("Bootstrapper TLS dial failure for address: %s; GRPC Dial error: %s", addr, err)
		// final attempt, use direct cloud connection and configured bootstrapper port instead of default TLS port
		if cfg.BootstrapPort != DefaultTLSBootstrapPort {
			addr = fmt.Sprintf("%s:%d", addrPieces[0], cfg.BootstrapPort)
			glog.Warningf("%v; trying: %s", err, addr)
			ctx2Tls, cance2lTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
			defer cance2lTls()
			conn, err = grpc.DialContext(ctx2Tls, addr, opts...)
			if err == nil {
				return conn, ctx2Tls.Err()
			}
			err = fmt.Errorf("final Bootstrapper TLS dial failure for: %s; GRPC Dial error: %s", addr, err)
		}
		glog.Error(err)
	}
	return conn, firstErr
}

func (b *Bootstrapper) getGrpcOpts(useProxy bool, cfg *config.ControlProxyCfg) []grpc.DialOption {
	var opts = []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 30 * time.Second,
		}),
		grpc.WithBlock(),
		grpc.WithAuthority(cfg.BootstrapAddr),
	}
	if useProxy {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// always try to add OS certs
		certPool, err := x509.SystemCertPool()
		if err != nil {
			glog.Warningf("OS Cert Pool initialization error: %v", err)
			certPool = x509.NewCertPool()
		}
		// Add magma RootCA
		if rootCa, err := ioutil.ReadFile(cfg.RootCaFile); err == nil {
			if !certPool.AppendCertsFromPEM(rootCa) {
				glog.Warningf("Failed to append certificates from %s", cfg.RootCaFile)
			}
		} else {
			glog.Warningf("Cannot load Root CA from '%s': %v", cfg.RootCaFile, err)
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

/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
func (b *Bootstrapper) GetBootstrapperCloudConnection() (conn *grpc.ClientConn, err error) {
	cfg := config.GetControlProxyConfigs()
	addrPieces := strings.Split(cfg.BootstrapAddr, ":")
	configuredAddr := fmt.Sprintf("%s:%d", addrPieces[0], cfg.BootstrapPort)
	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxLocalTimeoutSec*time.Second)
	defer cancel()
	opts := b.getGrpcOpts(false, cfg)
	conn, err = grpc.DialContext(ctx, configuredAddr, opts...)
	if err == nil {
		err = ctx.Err()
	}
	if err == nil {
		return
	}
	glog.Errorf("Bootstrapper dial failure %v for address: %s", err, configuredAddr)
	// in case of an error, try again with direct TLS connection to the default TLS port
	if cfg.BootstrapPort != DefaultTLSBootstrapPort {
		addr := fmt.Sprintf("%s:%d", addrPieces[0], DefaultTLSBootstrapPort)
		glog.Infof("trying default bootstrapper TLS port, address: %s", addr)
		// Try to call cloud directly
		ctxTls, cancelTls := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
		defer cancelTls()
		opts = b.getGrpcOpts(false, cfg)
		conn, err = grpc.DialContext(ctxTls, addr, opts...)
		if err == nil {
			err = ctxTls.Err()
		}
		if err != nil {
			glog.Errorf("Bootstrapper TLS dial failure for address: %s; GRPC Dial error: %s", addr, err)
		}
	}
	return
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
				InsecureSkipVerify: false, // last resort - do not verify the server cert, but rely only on challenge
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

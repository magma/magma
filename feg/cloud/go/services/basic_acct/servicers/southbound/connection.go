/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/golang/glog"
	"github.com/magma/augmented-networks/accounting/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const (
	grpcMaxTimeoutSec = 60
	grpcRemoteTimeout = time.Second * 20
	grpcMaxDelaySec   = 20
)

var (
	keepaliveParams = keepalive.ClientParameters{
		Time:                59 * time.Second,
		Timeout:             20 * time.Second,
		PermitWithoutStream: true,
	}
)

// GetAcctClient returns a new AN accounting client & CTX with timeout
func (s *BaseAccService) GetAcctClient() (protos.AccountingClient, context.Context, context.CancelFunc, error) {
	conn, err := s.GetConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcRemoteTimeout)
	return protos.NewAccountingClient(conn), ctx, cancel, nil
}

// GetConnection returns existing or new connection to the remote client
func (s *BaseAccService) GetConnection() (*grpc.ClientConn, error) {
	if s == nil {
		return nil, status.Error(codes.InvalidArgument, "nil BaseAccService")
	}
	s.RLock()
	conn := s.remoteConn
	if conn != nil && validateConnState(conn.GetState()) {
		s.RUnlock()
		return conn, nil
	}
	s.RUnlock()

	newConn, err := s.CreateConnection()
	if err != nil {
		return nil, err
	}
	s.Lock()
	defer s.Unlock()

	if conn = s.remoteConn; conn != nil && validateConnState(conn.GetState()) {
		go newConn.Close()
	} else {
		conn = newConn
		s.remoteConn = conn
	}
	return conn, nil
}

// CreateConnection creates a new gRPC connection
func (s *BaseAccService) CreateConnection() (*grpc.ClientConn, error) {
	if s == nil {
		return nil, status.Error(codes.InvalidArgument, "nil BaseAccService")
	}
	s.RLock()
	cfg := s.cfg
	addr := cfg.RemoteAddr
	clientCrt, clientCrtKey, rootCa, notls, insecure :=
		cfg.ClientCrt, cfg.ClientCrtKey, cfg.RootCaCert, cfg.NoTls, cfg.Insecure
	s.RUnlock()

	if len(addr) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty target address")
	}

	opts := getDialOptions(clientCrt, clientCrtKey, rootCa, notls, insecure)

	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxTimeoutSec*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to dial %s: %v", addr, err)
	}
	return conn, err
}

func validateConnState(st connectivity.State) bool {
	if st >= connectivity.Idle && st <= connectivity.Ready {
		return true
	}
	return false
}

func getDialOptions(clientCrt, clientCrtKey, rootCa string, notls, insecure bool) []grpc.DialOption {
	bckoff := backoff.DefaultConfig
	bckoff.MaxDelay = grpcMaxDelaySec * time.Second
	var opts = []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bckoff,
			MinConnectTimeout: grpcMaxTimeoutSec * time.Second,
		}),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepaliveParams),
	}
	if notls {
		return append(opts, grpc.WithInsecure())
	}
	tlsCfg := &tls.Config{}
	if insecure {
		tlsCfg.InsecureSkipVerify = true
	} else {
		// always try to add OS certs
		certPool, err := x509.SystemCertPool()
		if err != nil {
			glog.Warningf("OS Cert Pool initialization error: %v", err)
			certPool = x509.NewCertPool()
		}
		if len(rootCa) > 0 {
			// Add server RootCA
			if rootCa, err := ioutil.ReadFile(rootCa); err == nil {
				if !certPool.AppendCertsFromPEM(rootCa) {
					glog.Errorf("Failed to append certificates from %s", rootCa)
				}
			} else {
				glog.Errorf("Cannot load Root CA from '%s': %v", rootCa, err)
			}
		}
		if len(certPool.Subjects()) > 0 {
			tlsCfg.RootCAs = certPool
		} else {
			glog.Warning("Empty server certificate pool, using TLS InsecureSkipVerify")
			tlsCfg.InsecureSkipVerify = true
		}
		if len(clientCrt) > 0 {
			if len(clientCrtKey) > 0 {
				clientCert, err := tls.LoadX509KeyPair(clientCrt, clientCrtKey)
				if err == nil {
					tlsCfg.Certificates = []tls.Certificate{clientCert}
				} else {
					glog.Errorf("failed to load Client Certificate & Key from '%s', '%s': %v",
						clientCrt, clientCrtKey, err)
				}
			} else {
				glog.Errorf("failed to get gateway certificate key location: %v", err)
			}
		} else {
			glog.Errorf("failed to get gateway certificate location: %v", err)
		}
	}
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	return opts
}

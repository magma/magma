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

package main

import (
	"context"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy"
)

// s8Cli is a helper interface so so_client can use either a its own s8_proxy or use the s8_proxy running on that vm
type s8Cli interface {
	CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error)
	DeleteSession(req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error)
	SendEcho(req *protos.EchoRequest) (*protos.EchoResponse, error)
}

//s8CliImpl will represent the client that is running on the vm (probably on a docker container. It will use
// s8_proxy API send the messages
type s8CliImpl struct{}

func (s8CliImpl) CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	return s8_proxy.CreateSession(req)
}

func (s8CliImpl) DeleteSession(req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	return s8_proxy.DeleteSession(req)
}

func (s8CliImpl) SendEcho(req *protos.EchoRequest) (*protos.EchoResponse, error) {
	return s8_proxy.SendEcho(req)
}

// s8BuiltIn is a actual s8 client, separated from the one running on the VM.
type s8BuiltIn struct {
	server protos.S8ProxyServer
}

func (s s8BuiltIn) CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	return s.server.CreateSession(context.Background(), req)
}

func (s s8BuiltIn) DeleteSession(req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	return s.server.DeleteSession(context.Background(), req)
}

func (s s8BuiltIn) SendEcho(req *protos.EchoRequest) (*protos.EchoResponse, error) {
	return s.server.SendEcho(context.Background(), req)
}

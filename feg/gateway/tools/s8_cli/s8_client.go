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

type s8Cli interface {
	CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error)
}

type s8CliImpl struct{}

func (s8CliImpl) CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	return s8_proxy.CreateSession(req)
}

type s8BuiltIn struct {
	server protos.S8ProxyServer
}

func (s s8BuiltIn) CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	return s.server.CreateSession(context.Background(), req)
}

// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capture

import (
	"context"
	"fmt"
	"strings"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/protos/magma/capture"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Middleware defines middleware's functions - used for generating mocks.
type Middleware interface {
	GetDialOptions() []grpc.DialOption
	GetServerOptions() []grpc.ServerOption
}

type middleware struct {
	config.Configer
	*Buffer
}

// NewMiddleware returns a configured middleware.
func NewMiddleware(cfgr config.Configer, b *Buffer) middleware {
	return middleware{
		Configer: cfgr,
		Buffer:   b,
	}
}

// GetDialOptions returns all dial options attached to middleware.
func (m *middleware) GetDialOptions() []grpc.DialOption {
	var options []grpc.DialOption
	return append(options, m.withCaptureUnaryClientInterceptor())
}

// GetServerOptions returns all server options attached to middleware.
func (m *middleware) GetServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	return append(options, m.withCaptureUnaryServerInterceptor())
}

// withCaptureUnaryServerInterceptor returns a UnaryServerInterceptor that captures targeted traffic.
func (m *middleware) withCaptureUnaryServerInterceptor() grpc.ServerOption {
	captureInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {
		resp, retErr := handler(ctx, req)
		if m.isTargeted(info.FullMethod) {
			rq, err := anypb.New(req.(proto.Message))
			if err != nil {
				panic(err)
			}
			rp, err := anypb.New(resp.(proto.Message))
			if err != nil {
				panic(err)
			}

			captured := &capture.UnaryCall{
				Method:   info.FullMethod,
				Request:  rq,
				Response: rp,
			}
			if retErr != nil {
				captured.Err = err.Error()
			}
			m.Buffer.Write(captured)
		}
		return resp, retErr
	}
	return grpc.UnaryInterceptor(captureInterceptor)
}

// withCaptureUnaryClientInterceptor returns an interceptor that captures and buffers all rpc calls that go through it.
func (m *middleware) withCaptureUnaryClientInterceptor() grpc.DialOption {
	captureInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		retErr := invoker(ctx, method, req, reply, cc, opts...)
		if m.isTargeted(method) {
			rq, err := anypb.New(req.(proto.Message))
			if err != nil {
				panic(err)
			}
			rp, err := anypb.New(reply.(proto.Message))
			if err != nil {
				panic(err)
			}

			captured := &capture.UnaryCall{
				Method:   method,
				Request:  rq,
				Response: rp,
			}
			if retErr != nil {
				captured.Err = err.Error()
			}
			m.Buffer.Write(captured)
		}
		return retErr
	}
	return grpc.WithUnaryInterceptor(captureInterceptor)
}

// isTargeted compares the intercepted method to the match specs in the config.
// Supports wild cards.
func (m *middleware) isTargeted(interceptedMethod string) bool {
	captureConfig := m.Config()
	for _, spec := range captureConfig.CaptureConfig.MatchSpecs {
		if spec.GetService() == "*" && spec.GetMethod() == "*" {
			return true
		}
		if spec.GetService() == "*" {
			targetSuffix := fmt.Sprintf("/%s", spec.GetMethod())
			if strings.HasSuffix(interceptedMethod, targetSuffix) {
				return true
			}
		}
		if spec.GetMethod() == "*" {
			targetPrefix := fmt.Sprintf("/%s", spec.GetService())
			if strings.HasPrefix(interceptedMethod, targetPrefix) {
				return true
			}
		}
		target := fmt.Sprintf("/%s/%s", spec.GetService(), spec.GetMethod())
		if interceptedMethod == target {
			return true
		}
	}
	return false
}

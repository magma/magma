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

package config

import (
	"context"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/log"
	pb "github.com/magma/magma/src/go/protos/magma/config"
)

// ConfigServer handles ConfigServer RPCs.
type ConfigServer struct {
	log.Logger
	pb.ConfigServer
	config.Configer
}

// NewConfigServer returns a ConfigServer injected with the provided logger and
// Configer.
func NewConfigServer(logger log.Logger, cm config.Configer) *ConfigServer {
	return &ConfigServer{Logger: logger, Configer: cm}
}

// GetConfig calls Config.GetConfig.
func (c *ConfigServer) GetConfig(ctx context.Context, req *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	c.Logger.Debug().Print("GetConfig")
	config := c.Configer.Config()
	return &pb.GetConfigResponse{Config: config}, nil
}

// UpdateConfig calls Config.UpdateConfig.
func (c *ConfigServer) UpdateConfig(ctx context.Context, req *pb.UpdateConfigRequest) (*pb.UpdateConfigResponse, error) {
	c.Logger.
		With("config", req.GetConfig()).
		Debug().Print("UpdateConfig")
	if err := c.Configer.UpdateConfig(req.GetConfig()); err != nil {
		return nil, err
	}
	return &pb.UpdateConfigResponse{Config: c.Configer.Config()}, nil
}

func (c *ConfigServer) ReplaceConfig(ctx context.Context, req *pb.ReplaceConfigRequest) (*pb.ReplaceConfigResponse, error) {
	c.Logger.
		With("config", req.GetConfig()).
		Debug().Print("ReplaceConfig")
	if err := c.Configer.ReplaceConfig(req.GetConfig()); err != nil {
		return nil, err
	}
	return &pb.ReplaceConfigResponse{Config: c.Configer.Config()}, nil
}

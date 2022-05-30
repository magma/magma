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
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/magma/magma/src/go/agwd/config/mock_config"
	"github.com/magma/magma/src/go/internal/testutil"
	pb "github.com/magma/magma/src/go/protos/magma/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCfgr := config.NewMockConfiger(ctrl)
	logger, _ := testutil.NewTestLogger()
	NewConfigServer(logger, mockCfgr)
}

func TestConfigServer_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedConfig := &pb.AgwD{
		LogLevel:                        pb.AgwD_DEBUG,
		SctpdDownstreamServiceTarget:    "sctpd_down",
		SctpdUpstreamServiceTarget:      "sctpd_up",
		MmeSctpdDownstreamServiceTarget: "mme_down",
		MmeSctpdUpstreamServiceTarget:   "mme_up",
		SentryDsn:                       "",
		ConfigServicePort:               "1234",
		VagrantPrivateNetworkIp:         "0.0.0.0",
		CaptureServicePort:              "12345",
		CaptureConfig: &pb.CaptureConfig{
			MatchSpecs: []*pb.CaptureConfig_MatchSpec{
				{
					Service: "service",
					Method:  "method",
				},
				{
					Service: "service",
					Method:  "method2",
				},
			}},
	}

	ctx := context.Background()
	req := &pb.GetConfigRequest{}
	res := &pb.GetConfigResponse{Config: expectedConfig}

	mockCfgr := config.NewMockConfiger(ctrl)
	mockCfgr.EXPECT().Config().Return(expectedConfig)

	logger, logBuffer := testutil.NewTestLogger()

	cs := NewConfigServer(logger, mockCfgr)

	got, err := cs.GetConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tGetConfig\n",
		logBuffer.String())

}

func TestConfigServer_UpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedConfig := &pb.AgwD{
		LogLevel:                        pb.AgwD_DEBUG,
		SctpdDownstreamServiceTarget:    "sctpd_down",
		SctpdUpstreamServiceTarget:      "sctpd_up",
		MmeSctpdDownstreamServiceTarget: "mme_down",
		MmeSctpdUpstreamServiceTarget:   "mme_up",
		SentryDsn:                       "",
		ConfigServicePort:               "1234",
		VagrantPrivateNetworkIp:         "0.0.0.0",
		CaptureServicePort:              "12345",
		CaptureConfig: &pb.CaptureConfig{
			MatchSpecs: []*pb.CaptureConfig_MatchSpec{
				{
					Service: "service",
					Method:  "method",
				},
				{
					Service: "service",
					Method:  "method2",
				},
			}},
	}

	ctx := context.Background()
	req := &pb.UpdateConfigRequest{Config: expectedConfig}
	res := &pb.UpdateConfigResponse{Config: expectedConfig}

	mockCfgr := config.NewMockConfiger(ctrl)
	mockCfgr.EXPECT().UpdateConfig(req.Config)
	mockCfgr.EXPECT().Config().Return(expectedConfig)

	logger, _ := testutil.NewTestLogger()

	cs := NewConfigServer(logger, mockCfgr)

	got, err := cs.UpdateConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
}

func TestConfigServer_ReplaceConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	replacementConfig := &pb.AgwD{
		LogLevel:                        pb.AgwD_DEBUG,
		SctpdDownstreamServiceTarget:    "sctpd_down",
		SctpdUpstreamServiceTarget:      "sctpd_up",
		MmeSctpdDownstreamServiceTarget: "mme_down",
		MmeSctpdUpstreamServiceTarget:   "mme_up",
		SentryDsn:                       "",
		ConfigServicePort:               "1234",
		VagrantPrivateNetworkIp:         "0.0.0.0",
		CaptureServicePort:              "12345",
		CaptureConfig: &pb.CaptureConfig{
			MatchSpecs: []*pb.CaptureConfig_MatchSpec{
				{
					Service: "service",
					Method:  "method",
				},
				{
					Service: "service",
					Method:  "method2",
				},
			}},
	}

	ctx := context.Background()
	req := &pb.ReplaceConfigRequest{Config: replacementConfig}
	res := &pb.ReplaceConfigResponse{Config: replacementConfig}

	mockCfgr := config.NewMockConfiger(ctrl)
	mockCfgr.EXPECT().ReplaceConfig(req.Config)
	mockCfgr.EXPECT().Config().Return(replacementConfig)

	logger, _ := testutil.NewTestLogger()

	cs := NewConfigServer(logger, mockCfgr)

	got, err := cs.ReplaceConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
}

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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/agwd/config"
	mock_cfg "github.com/magma/magma/src/go/agwd/config/mock_config"
	configpb "github.com/magma/magma/src/go/protos/magma/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware_Success(t *testing.T) {
	cfgr := config.NewConfigManager()
	buf := NewBuffer()
	mw := NewMiddleware(cfgr, buf)
	assert.Equal(t, mw.Config(), cfgr.Config())
	assert.Equal(t, mw.Buffer, buf)
}

func TestMiddleware_GetDialOptions(t *testing.T) {
	cfgr := config.NewConfigManager()
	buf := NewBuffer()
	mw := NewMiddleware(cfgr, buf)

	dialOptions := mw.GetDialOptions()
	assert.Equal(t, 1, len(dialOptions))
}

func TestMiddleware_GetServerOptions(t *testing.T) {
	cfgr := config.NewConfigManager()
	buf := NewBuffer()
	mw := NewMiddleware(cfgr, buf)

	serverOptions := mw.GetServerOptions()
	assert.Equal(t, 1, len(serverOptions))
}

func TestNewMiddleware_isTargeted(t *testing.T) {
	t.Parallel()

	spec1 := &configpb.CaptureConfig_MatchSpec{
		Service: "magma.sctpd.SctpdDownlink",
		Method:  "SendDl",
	}
	spec2 := &configpb.CaptureConfig_MatchSpec{
		Service: "magma.sctpd.service2",
		Method:  "SendDl2",
	}
	wildCardServiceSpec := &configpb.CaptureConfig_MatchSpec{
		Service: "*",
		Method:  "SendDl",
	}
	wildCardMethodSpec := &configpb.CaptureConfig_MatchSpec{
		Service: "magma.sctpd.SctpdDownlink",
		Method:  "*",
	}

	allWildCardSpec := &configpb.CaptureConfig_MatchSpec{
		Service: "*",
		Method:  "*",
	}

	tests := []struct {
		specs []*configpb.CaptureConfig_MatchSpec
		call  string
		want  bool
	}{
		{},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  false,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{spec1},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{spec1, spec2},
			call:  "/magma.sctpd.service2/SendDl2",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{wildCardServiceSpec},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{wildCardMethodSpec},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{allWildCardSpec},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{spec1},
			call:  "/magma.sctpd.SctpdDownlink/noMatch",
			want:  false,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{spec1, spec2},
			call:  "/magma.sctpd.SctpdDownlink/SendDl",
			want:  true,
		},
		{
			specs: []*configpb.CaptureConfig_MatchSpec{spec1, spec2},
			call:  "/magma.sctpd.SctpdDownlink/a",
			want:  false,
		},
	}

	for _, test := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cfg := &configpb.AgwD{
			CaptureConfig: &configpb.CaptureConfig{MatchSpecs: test.specs}}

		mockCfgr := mock_cfg.NewMockConfiger(ctrl)
		mockCfgr.EXPECT().Config().Return(cfg)

		mw := NewMiddleware(
			mockCfgr,
			NewBuffer())
		got := mw.isTargeted(test.call)
		assert.Equal(
			t,
			test.want,
			got,
			"isTargeted(%s) = %v, want %v with config: %v",
			test.call,
			got,
			test.want,
			test.specs)
	}
}

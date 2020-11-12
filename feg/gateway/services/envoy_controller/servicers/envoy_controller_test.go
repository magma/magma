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

package servicers_test

import (
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/envoy_controller/control_plane/mocks"
	"magma/feg/gateway/services/envoy_controller/servicers"
	lte_proto "magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	IMSI1 = "IMSI00101"
	IMSI2 = "IMSI00102"
)

var (
	imsis       = []string{IMSI1, IMSI2}
	add_ue_reqs = []*protos.AddUEHeaderEnrichmentRequest{
		{
			UeIp: &lte_proto.IPAddress{
				Version: lte_proto.IPAddress_IPV4,
				Address: []byte("3.3.33.3"),
			},
			Websites: []string{"neverssl.com", "google.com"},
			Headers: []*protos.Header{{
				Name:  "IMSI",
				Value: "024212312312",
			}},
		},
		{
			UeIp: &lte_proto.IPAddress{
				Version: lte_proto.IPAddress_IPV4,
				Address: []byte("2.2.2.2"),
			},
			Websites: []string{"magma.com", "qqq.com"},
			Headers: []*protos.Header{{
				Name:  "IMSI",
				Value: "111111",
			},
				{
					Name:  "MSISDN",
					Value: "THIS_IS_MSISDN",
				}},
		},
	}
	overwrtie_ue_req = []*protos.AddUEHeaderEnrichmentRequest{
		{
			UeIp: &lte_proto.IPAddress{
				Version: lte_proto.IPAddress_IPV4,
				Address: []byte("2.2.2.2"),
			},
			Websites: []string{"magma.com", "qqq.com"},
			Headers: []*protos.Header{{
				Name:  "NEW_IMSI",
				Value: "212412",
			}},
		},
	}
	deactivate_req = &protos.DeactivateUEHeaderEnrichmentRequest{
		UeIp: &lte_proto.IPAddress{
			Version: lte_proto.IPAddress_IPV4,
			Address: []byte("3.3.33.3"),
		},
	}
)

// ---- TESTS ----
func TestEnvoyControllerInit(t *testing.T) {

	cli := new(mocks.EnvoyController)
	srv := servicers.NewEnvoyControllerService(cli)

	ctx := context.Background()

	cli.On("UpdateSnapshot", add_ue_reqs[:1]).Return()
	cli.On("UpdateSnapshot", add_ue_reqs).Return()
	_, err := srv.AddUEHeaderEnrichment(ctx, add_ue_reqs[0])
	_, err = srv.AddUEHeaderEnrichment(ctx, add_ue_reqs[1])

	cli.On("UpdateSnapshot", add_ue_reqs[1:]).Return()
	_, err = srv.DeactivateUEHeaderEnrichment(ctx, deactivate_req)

    // Make sure duplicate doesn't get reinserted but gets overwritten
	cli.On("UpdateSnapshot", overwrtie_ue_req).Return()
	_, err = srv.AddUEHeaderEnrichment(ctx, overwrtie_ue_req[0])

	assert.NoError(t, err)

	cli.AssertExpectations(t)
}

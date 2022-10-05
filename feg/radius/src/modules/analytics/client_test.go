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

package analytics

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"testing"

	"fbc/cwf/radius/modules/analytics/graphql"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	accessToken, ok := os.LookupEnv("INTEGRATION_ACCESS_TOKEN")
	if !ok {
		t.SkipNow()
	}
	u, err := user.Current()
	require.NoError(t, err, "failed getting user")
	c := graphql.NewClient(graphql.ClientConfig{
		Token:    accessToken,
		Endpoint: fmt.Sprintf("https://graph.%s.sb.expresswifi.com/graphql", u.Username),
		HTTPClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
	})

	t.Log("creating session")
	cop := NewCreateSessionOp(&RadiusSession{
		NASIPAddress:         "10.0.0.1",
		NASIdentifier:        "10.0.0.2",
		AcctSessionID:        "10.0.0.3:10.0.0.4",
		CalledStationID:      "10.0.0.5",
		FramedIPAddress:      "10.0.0.6",
		CallingStationID:     "10.0.0.7",
		NormalizedMacAddress: "10.0.0.8",
		RADIUSServerID:       1,
	})
	require.NoError(t, c.Do(cop), "failed creating session")
	t.Logf("radius_session_id:%d\n", cop.Response().FBID)

	t.Logf("updating session with id %d\n", cop.Response().FBID)
	uop := NewUpdateSessionOp(&RadiusSession{
		FBID:                 cop.Response().FBID,
		NASIPAddress:         "10.0.1.1",
		NASIdentifier:        "10.0.1.2",
		AcctSessionID:        "10.0.1.3:10.0.1.4",
		CalledStationID:      "10.0.1.5",
		FramedIPAddress:      "10.0.1.6",
		CallingStationID:     "10.0.1.7",
		NormalizedMacAddress: "10.0.1.8",
		RADIUSServerID:       2,
		UploadBytes:          1234,
		Vendor:               int64(Ruckus),
	})
	require.NoError(t, c.Do(uop), "failed updating session")
}

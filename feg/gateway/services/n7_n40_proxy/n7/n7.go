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

package n7

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
)

// NewN7Client creates a N7 oapi client and sets the OAuth2 client credentiatls for authorizing requests
func NewN7Client(cfg *PCFConfig) (*n7_sbi.ClientWithResponses, error) {
	tokenConfig := clientcredentials.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenUrl,
	}
	tokenCtxt := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	// Create new N7 client object and assosiate it with oAuth2 HTTP client
	client, err := n7_sbi.NewClientWithResponses(
		fmt.Sprintf("%s://%s", cfg.ApiRoot.Scheme, cfg.ApiRoot.Host),
		n7_sbi.WithHTTPClient(tokenConfig.Client(tokenCtxt)),
	)
	return client, err
}

func removeIMSIPrefix(imsi string) string {
	return strings.TrimPrefix(imsi, "IMSI")
}

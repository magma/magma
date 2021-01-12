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

package ofpanalytics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type (
	Config struct {
		URI         string
		AccessToken string
		DryRun      bool
	}

	moduleContext struct {
		URI         string
		AccessToken string
		HTTPClient  *http.Client
	}

	EndpointResponse struct {
		Auth []string `json:"config:Auth-Type,omitempty"`
	}
)

var (
	defaultTimeout   = 5 * time.Second
	acceptCode       = "254"
	rejectCode       = "5"
	analyticsVersion = "v1.1"
)

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var mCtx moduleContext
	var cfg Config

	err := mapstructure.Decode(config, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.URI == "" {
		return nil, errors.New("OFPAnalytics module cannot be initialized with an empty URI value")
	}
	if cfg.AccessToken == "" {
		return nil, errors.New("OFPAnalytics module cannot be initialized with an empty access token value")
	}

	mCtx.URI = cfg.URI
	mCtx.AccessToken = cfg.AccessToken
	mCtx.HTTPClient = &http.Client{
		Timeout: defaultTimeout,
	}
	// For DryRun we're going to allow connection to an unauthenticated end
	if cfg.DryRun {
		mCtx.HTTPClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	logger.Info("OFPAnalytics module initialized successfully")
	return mCtx, nil
}

// Normalizing fields
func normalize(radiusAvp string) string {
	var replacer = strings.NewReplacer("-", ":")
	return strings.ToLower(replacer.Replace(radiusAvp))
}

// Handle module interface implementation
func Handle(m modules.Context, rc *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	ctx, ok := m.(moduleContext)
	if !ok {
		return nil, fmt.Errorf("unable to obtain context")
	}

	// Currently only handling authorization requests - we have roadmap tasks to support full v2 integration T64414814
	// When we'll have v2 support we can remove the hand crafted json packet
	jsonPacket := map[string]map[string]interface{}{
		"Called-Station-Id":  {"type": "string", "value": []string{normalize(rfc2865.CalledStationID_GetString(r.Packet))}},
		"Calling-Station-Id": {"type": "string", "value": []string{normalize(rfc2865.CallingStationID_GetString(r.Packet))}},
		"NAS-Identifier":     {"type": "string", "value": []string{rfc2865.NASIdentifier_GetString(r.Packet)}},
		"XWF-C-Version":      {"type": "string", "value": []string{analyticsVersion}},
	}
	// If no nas ip address is specified then no field will be sent
	if rfc2865.NASIPAddress_Get(r.Packet) != nil {
		jsonPacket["NAS-IP-Address"] =
			map[string]interface{}{"type": "string", "value": []string{rfc2865.NASIPAddress_Get(r.Packet).String()}}
	}
	encodedMsg, err := json.Marshal(jsonPacket)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal radius packet")
	}
	req, err := http.NewRequest(http.MethodPost, ctx.URI, bytes.NewReader(encodedMsg))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+ctx.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending response to endpoint")
	}
	defer resp.Body.Close()
	rc.Logger.Debug("got response", zap.String("status", resp.Status),
		zap.String("url", resp.Request.URL.String()),
		zap.Any("request", r.Packet.Attributes),
		zap.String("Called-Station-Id", rfc2865.CalledStationID_GetString(r.Packet)),
		zap.String("Calling-Station-Id", rfc2865.CallingStationID_GetString(r.Packet)),
		zap.String("NAS-Identifier", rfc2865.NASIdentifier_GetString(r.Packet)),
		zap.String("XWF-C-Version", analyticsVersion))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error processing message by endpoint. Response status %d", resp.StatusCode)
	}
	var endPointResponse EndpointResponse
	if err := json.NewDecoder(resp.Body).Decode(&endPointResponse); err != nil {
		return nil, errors.Wrap(err, "unable to decode endpoint response")
	}

	if len(endPointResponse.Auth) == 0 {
		return nil, fmt.Errorf("malformed auth response: no acceptance code")
	}
	var p *radius.Packet
	switch endPointResponse.Auth[0] {
	case acceptCode:
		p = r.Response(radius.CodeAccessAccept)
	case rejectCode:
		p = r.Response(radius.CodeAccessReject)
	}

	response := &modules.Response{
		Code:       p.Code,
		Attributes: p.Attributes,
	}
	return response, nil
}

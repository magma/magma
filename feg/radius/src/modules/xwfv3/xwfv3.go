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

package xwfv3

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fbc/cwf/radius/modules/xwfv3/xwfhttp2"
	"fmt"
	"net/http"
	"strings"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/machine"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// ExpressWiFiVendorSpecificServerRADIUSAttributeType ...
const ExpressWiFiVendorSpecificServerRADIUSAttributeType uint32 = 99999

// Config configuration structure for restproxy module
type Config struct {
	URI           string
	AccessToken   string
	Method        string
	SSEMacAddress string
}

// ModuleCtx ...
type ModuleCtx struct {
	uri         string
	http2client *xwfhttp2.Client
	method      string
	sseMACAddr  string
}

type wwwResp struct {
	Data string `json:"data,omitempty"`
}

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var mCtx ModuleCtx
	var cfg Config
	err := mapstructure.Decode(config, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.URI == "" {
		return nil, errors.New("XWFv3 module cannot be initialized with an empty URI value")
	}
	if cfg.AccessToken == "" {
		return nil, errors.New("XWFv3 module cannot be initialized with an empty access token value")
	}

	mCtx.uri = cfg.URI
	mCtx.method, err = normalizeHTTPMethod(cfg.Method)
	if err != nil {
		return nil, err
	}
	mCtx.sseMACAddr = cfg.SSEMacAddress
	if mCtx.sseMACAddr == "" {
		mCtx.sseMACAddr = machine.GetMachineMACAddressID()
	}
	logger.Info(
		"Generating sessions with SSE MAC Address",
		zap.String("sse_mac_address", mCtx.sseMACAddr),
	)

	mCtx.http2client = xwfhttp2.NewClient(cfg.AccessToken)
	logger.Info("XWFv3 module initialized successfully")
	return mCtx, nil
}

// Handle module interface implementation
func Handle(m modules.Context, rc *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)
	var res *radius.Packet
	if !strings.EqualFold(http.MethodPost, mCtx.method) {
		return nil, errors.New("XWFv3 only supports POST method at this point")
	}

	// Add XWFv3 version header
	xwfVersionAttr, err := radius.NewVendorSpecific(
		ExpressWiFiVendorSpecificServerRADIUSAttributeType,
		radius.Attribute([]byte{4, 4, 'v', '3', '.', '0'}),
	)
	if err != nil {
		return nil, errors.New("Failed encoding XWFv3 Version")
	}
	r.Packet.Add(rfc2865.VendorSpecific_Type, xwfVersionAttr)

	// Serialize the packet
	data, err := r.Packet.Encode()
	if err != nil {
		return nil, err
	}

	respBody, err := mCtx.http2client.PostJSON(mCtx.uri, map[string]string{
		// Transform the radius request to a json suitable body for www
		"data": base64.StdEncoding.EncodeToString(data),
	}, map[string]string{
		"sse-client-mac-address": mCtx.sseMACAddr,
	})

	if err != nil {
		return nil, err
	}

	// Parsing the json response
	decoder := json.NewDecoder(bytes.NewReader(respBody))
	encodedRadius := &wwwResp{}
	if decoder.Decode(encodedRadius) != nil {
		return nil, err
	}

	// Decoding the base64 string to binary form
	radiusResponse, err := base64.StdEncoding.DecodeString(encodedRadius.Data)
	if err != nil {
		return nil, err
	}

	res, err = radius.Parse(radiusResponse, r.Secret)
	if err != nil {
		return nil, err
	}

	response := &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}
	return response, nil
}

func normalizeHTTPMethod(method string) (string, error) {
	switch strings.ToUpper(method) {
	case http.MethodPost:
		return http.MethodPost, nil
	case http.MethodGet:
		return http.MethodGet, nil
	default:
		return "", fmt.Errorf("unsupported http method %s", method)
	}
}

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
	"fbc/lib/go/radius"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

func normalizeMethod(method string) (string, error) {
	switch strings.ToUpper(method) {
	case http.MethodPost:
		return http.MethodPost, nil
	case http.MethodGet:
		return http.MethodGet, nil
	default:
		return "", fmt.Errorf("unsupported http method %s", method)
	}
}

// Config configuration structure for restproxy module
type Config struct {
	URI         string
	AccessToken string
	Method      string
}

var uri string
var http2client *xwfhttp2.Client
var method string

type wwwResp struct {
	Data string `json:"data,omitempty"`
}

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) error {
	var postProxyConfig Config
	err := mapstructure.Decode(config, &postProxyConfig)
	if err != nil {
		return err
	}

	if postProxyConfig.URI == "" {
		return errors.New("rest proxy module cannot be initialized with an empty URI value")
	}
	if postProxyConfig.AccessToken == "" {
		return errors.New("rest proxy module cannot be initialized with an empty access token value")
	}

	uri = postProxyConfig.URI
	method, err = normalizeMethod(postProxyConfig.Method)
	if err != nil {
		return err
	}

	http2client = xwfhttp2.NewClient(postProxyConfig.AccessToken)
	logger.Info("rest proxy module initialized successfully")
	return nil
}

// Handle module interface implementation
func Handle(_ *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {

	var res *radius.Packet
	if strings.EqualFold(http.MethodPost, method) {

		data, err := r.Packet.Encode()
		if err != nil {
			return nil, err
		}

		respBody, err := http2client.PostJSON(uri, map[string]string{
			// Transform the radius request to a json suitable body for www
			"data": base64.StdEncoding.EncodeToString(data),
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

	} else if strings.EqualFold(http.MethodGet, method) {
		// Task: T45993664
		return nil, fmt.Errorf("unimplemented method: %s", http.MethodGet)
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}, nil
}

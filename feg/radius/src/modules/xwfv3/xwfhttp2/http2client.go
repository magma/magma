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

package xwfhttp2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// DefaultTimeout the default timeout of a request
const DefaultTimeout = 5 * time.Second

//Client struct encapsulates an http2 connection to www
type Client struct {
	http2client http.Client
	accessToken string
	Timeout     time.Duration
}

//isSuccessful
func isSuccessful(resp *http.Response) bool {
	statusCode := resp.StatusCode

	// A list of successful status codes can be viewed at: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
	successCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
	}

	for _, code := range successCodes {
		if statusCode == code {
			return true
		}
	}

	return false
}

// NewClient returns an initialized Client struct.
func NewClient(accessToken string) *Client {
	return &Client{
		http2client: http.Client{
			Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			Timeout:   DefaultTimeout,
		},
		accessToken: accessToken,
	}
}

func (c *Client) addAccessToken(req *http.Request) {
	q := req.URL.Query()
	q.Add("access_token", c.accessToken)
	req.URL.RawQuery = q.Encode()
}

//Get method will perform an http get request for the specified url and will return the response bytes
func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.addAccessToken(req)

	resp, err := c.http2client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccessful(resp) {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err

}

//PostJSON method will perform an http post to the url specified with a json body
func (c *Client) PostJSON(url string, bodyToSend map[string]string, headers map[string]string) ([]byte, error) {

	// Turning the map into a json
	body, err := json.Marshal(bodyToSend)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("radius-packet-encoding", "base64/binary")
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	c.addAccessToken(req)

	resp, err := c.http2client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccessful(resp) {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, err
}

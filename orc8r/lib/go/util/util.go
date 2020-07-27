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

package util

import (
	"crypto/x509/pkix"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func GetFreeTcpPort(preferredPort int) int {
	l, err := net.Listen("tcp", ":"+string(preferredPort))
	if err != nil {
		l, _ = net.Listen("tcp", "")
	}
	addr, _ := net.ResolveTCPAddr("tcp", l.Addr().String())
	l.Close()
	return addr.Port
}

func SendHttpRequest(
	method, url,
	payload string,
	// optional headers, each header is a slice of strings in the form:
	//   [<header name>, value1, value2...]
	headers ...[]string,
) (int, string, error) {

	var body io.Reader = nil
	if len(payload) > 0 {
		body = strings.NewReader(payload)
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, "", err
	}
	for _, h := range headers {
		if l := len(h); l > 0 {
			k := h[0]
			if l > 1 {
				request.Header.Set(k, h[1])
				if l > 2 {
					for _, v := range h[2:] {
						request.Header.Add(k, v)
					}
				}
			} else {
				request.Header.Set(k, "")
			}
		}
		request.Header.Set(h[0], h[1])
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	var client = &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, string(contents), nil
}

func CompareJSON(j1, j2 string) bool {
	var struct1, struct2 interface{}

	if len(j1) < 8 || len(j2) < 8 {
		return j1 == j2
	}
	err := json.Unmarshal([]byte(j1), &struct1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(j2), &struct2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(struct1, struct2)
}

// FormatPkixSubject returns a string representing x509 cert subject/issuer
func FormatPkixSubject(s *pkix.Name) string {
	if s == nil {
		return ""
	}
	res, _ := json.Marshal(struct {
		Country            []string `json:"C,omitempty"`
		Organization       []string `json:"O,omitempty"`
		OrganizationalUnit []string `json:"OU,omitempty"`
		Locality           []string `json:"L,omitempty"`
		Province           []string `json:"S,omitempty"`
		StreetAddress      []string `json:"STREET,omitempty"`
		PostalCode         []string `json:",omitempty"`
		SerialNumber       string   `json:",omitempty"`
		CommonName         string   `json:"CN,omitempty"`
	}{
		s.Country,
		s.Organization,
		s.OrganizationalUnit,
		s.Locality,
		s.Province,
		s.StreetAddress,
		s.PostalCode,
		s.SerialNumber,
		s.CommonName,
	})
	return string(res)
}

// IsTruthyEnv returns true for any value not "false", "0", "no..."
func IsTruthyEnv(envName string) bool {
	value := os.Getenv(envName)
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	if value == "0" || strings.HasPrefix(value, "false") || strings.HasPrefix(value, "no") {
		return false
	}
	return true
}

// GetEnvBool returns value of the environment
// variable if it exists, or defaultValue if not
func GetEnvBool(envVariable string, defaultValue ...bool) bool {
	if len(envVariable) > 0 {
		if envValue := os.Getenv(envVariable); len(envValue) > 0 {
			envValueBool, _ := strconv.ParseBool(envValue)
			return envValueBool
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

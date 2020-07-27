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

package libgraphql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	// ClientMutationID field used as a unique identifier for the outgoing request.
	ClientMutationID = "client_mutation_id"
	// defaultClientTimeout holds the default timeout for graphql http client.
	defaultClientTimeout = 10 * time.Second
)

// Op is the interface that wraps the GraphQL HTTP operation.
type Op interface {
	json.Unmarshaler
	Doc() string
	Vars() (string, error)
}

// Client is the libgraphql client.
type Client struct {
	ClientConfig
	// client state.
	auth string
}

// ClientConfig holds the configuration for the GraphQL client.
type ClientConfig struct {
	// Token is the bearer token for the authorization header.
	Token string
	// Endpoint is the GraphQL endpoint. defaults to "graph.expresswifi.com/graphql"
	Endpoint string
	// HTTPClient is an optional HTTP client. defaults to http.DefaultClient.
	HTTPClient *http.Client
}

// NewClient creates a new libgraphql.Client by the given config.
func NewClient(c ClientConfig) *Client {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: defaultClientTimeout}
	}
	return &Client{ClientConfig: c, auth: "Bearer " + c.Token}
}

// Do executes the given libgraphql.Op interface.
func (c *Client) Do(op Op) error {
	v := url.Values{}
	v.Add("doc", op.Doc())
	vars, err := op.Vars()
	if err != nil {
		return err
	}
	v.Add("variables", vars)
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("libgraphql: invalid status code: %s", res.Status)
	}
	if err := json.NewDecoder(res.Body).Decode(op); err != nil {
		return err
	}
	return nil
}

// Vars is the variables container for GraphQL mutations.
type Vars map[string]interface{}

// String returns the string representation of GraphQL variables.
func (v Vars) String() (string, error) {
	v[ClientMutationID] = uuid.New().String()
	b := new(strings.Builder)
	// variables wrapped in a "data" object.
	if err := json.NewEncoder(b).Encode(map[string]Vars{"data": v}); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Errors wraps a list of errors.
type Errors []Error

func (e Errors) Error() string {
	s := new(strings.Builder)
	for i := range e {
		s.WriteString(e[i].Error())
	}
	return s.String()
}

// Error is a GraphQL error from WWW.
type Error struct {
	Code     uint   `json:"code,omitempty"`
	Desc     string `json:"description,omitempty"`
	Message  string `json:"message,omitempty"`
	Summary  string `json:"summary,omitempty"`
	Severity string `json:"severity,omitempty"`
	TraceID  string `json:"fbtrace_id,omitempty"`
}

func (e *Error) Error() string { return fmt.Sprintf("libgraphql: %s", e.Message) }

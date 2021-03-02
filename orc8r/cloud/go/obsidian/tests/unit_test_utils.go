/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Method           string
	URL              string
	Payload          encoding.BinaryMarshaler
	MalformedPayload bool
	Handler          echo.HandlerFunc

	ParamNames  []string
	ParamValues []string

	ExpectedStatus int
	ExpectedResult encoding.BinaryMarshaler

	ExpectedError          string
	ExpectedErrorSubstring string
}

// RunUnitTest runs a test case using the given Echo instance.
// Does not start an obsidian server.
func RunUnitTest(t *testing.T, e *echo.Echo, test Test) {
	var req *http.Request
	if test.Payload != nil {
		payloadBytes, err := test.Payload.MarshalBinary()
		if !assert.NoError(t, err) {
			return
		}
		if test.MalformedPayload {
			payloadBytes = append([]byte{'x'}, payloadBytes...)
		}
		req = httptest.NewRequest(test.Method, test.URL, bytes.NewReader(payloadBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(test.Method, test.URL, bytes.NewReader([]byte{}))
	}

	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)
	c.SetParamNames(test.ParamNames...)
	c.SetParamValues(test.ParamValues...)
	handlerErr := test.Handler(c)
	if handlerErr != nil {
		// In prod, this is handled by CollectStats middleware
		c.Error(handlerErr)
	}
	assert.Equal(t, test.ExpectedStatus, recorder.Code)

	if test.ExpectedError != "" {
		if httpErr, ok := handlerErr.(*echo.HTTPError); ok {
			assert.Equal(t, test.ExpectedError, httpErr.Message)
		} else {
			assert.EqualError(t, handlerErr, test.ExpectedError)
		}
	} else if test.ExpectedErrorSubstring != "" {
		if handlerErr == nil {
			assert.Fail(t, "unexpected nil error", "error was nil but was expecting %s", test.ExpectedErrorSubstring)
		} else {
			assert.Contains(t, handlerErr.Error(), test.ExpectedErrorSubstring)
		}
	} else {
		if assert.NoError(t, handlerErr) && test.ExpectedResult != nil {
			expectedBytes, err := test.ExpectedResult.MarshalBinary()
			if assert.NoError(t, err) {
				// Convert to string for more readable assert failure messages.
				//
				// json.Marshal returns the serialized value as-is, without
				// appending a newline.
				//
				// The echo.Context's JSON method uses a json.Encoder to encode
				// its object. The json.Encoder object always appends a newline
				// to the end of the serialized value.
				//
				// To handle this mismatch, trim a newline from both values.
				expected := strings.TrimSuffix(string(expectedBytes), "\n")
				actual := strings.TrimSuffix(recorder.Body.String(), "\n")
				assert.Equal(t, expected, actual)
			}
		}
	}
}

// GetHandlerByPathAndMethod fetches the first obsidian.Handler that matches the
// given path and method from a list of handlers. If no such handler exists, it
// will fail.
func GetHandlerByPathAndMethod(t *testing.T, handlers []obsidian.Handler, path string, method obsidian.HttpMethod) obsidian.Handler {
	for _, handler := range handlers {
		if handler.Path == path && handler.Methods == method {
			return handler
		}
	}
	assert.Fail(t, fmt.Sprintf("no handler registered for path %s", path))
	return obsidian.Handler{}
}

func JSONMarshaler(v interface{}) encoding.BinaryMarshaler {
	return &jsonMarshaler{v: v}
}

type jsonMarshaler struct {
	v interface{}
}

func (j *jsonMarshaler) MarshalBinary() (data []byte, err error) {
	return json.Marshal(j.v)
}

func StringMarshaler(s string) encoding.BinaryMarshaler {
	return &stringMarshaler{s}
}

type stringMarshaler struct {
	s string
}

func (m *stringMarshaler) MarshalBinary() (data []byte, err error) {
	return []byte(m.s), nil
}

func ByteIdentityMarshaler(v []byte) encoding.BinaryMarshaler {
	return &byteIdentityMarshaler{v: v}
}

type byteIdentityMarshaler struct {
	v []byte
}

func (j *byteIdentityMarshaler) MarshalBinary() (data []byte, err error) {
	return j.v, nil
}

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

package obsidian

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectHandlerFuncCalled bool

	handlerFuncError         string
	migratedHandlerFuncError string

	expectedError string
}

func TestRegister(t *testing.T) {
	runCase(t, testCase{
		expectHandlerFuncCalled: true,
	})
	runCase(t, testCase{expectHandlerFuncCalled: true})
	runCase(t, testCase{
		expectHandlerFuncCalled: true,
		handlerFuncError:        "foo",
		expectedError:           "foo",
	})
}

func runCase(t *testing.T, tc testCase) {
	registries = map[HttpMethod]handlerRegistry{
		GET:    {},
		POST:   {},
		PUT:    {},
		DELETE: {},
	}
	handlerFuncCalled := false
	mockHandler := Handler{
		Methods: GET,
		Path:    "/foo/",
		HandlerFunc: func(c echo.Context) error {
			handlerFuncCalled = true
			if tc.handlerFuncError != "" {
				return errors.New(tc.handlerFuncError)
			}
			return nil
		},
	}
	err := Register(mockHandler)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/foo/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = registries[GET]["/foo/"](c)
	assert.Equal(t, 200, rec.Code)
	if tc.expectedError != "" {
		assert.EqualError(t, err, tc.expectedError)
	} else {
		assert.NoError(t, err)
	}
	assert.Equal(t, tc.expectHandlerFuncCalled, handlerFuncCalled)
}

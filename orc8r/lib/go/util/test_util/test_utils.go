/**
 * Copyright 2022 The Magma Authors.
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
package test_util

import (
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// AssertErrorsEqual checks if actualError contains relevant parts of expectedError.
// This method of comparison is necessary since string representations of proto messages contain
// non-deterministic whitespaces (https://github.com/golang/protobuf/issues/1269)
func AssertErrorsEqual(t *testing.T, expectedError, actualError, separator string) {
	origExp, origActual := expectedError, actualError
	if !strings.Contains(separator, " ") {
		// remove all spaces
		expectedError = strings.ReplaceAll(expectedError, " ", "")
		actualError = strings.ReplaceAll(actualError, " ", "")
	}
	// split expected error string at `separator`
	expectedErrorParts := strings.Split(expectedError, separator)
	// check if actual message contains relevant parts of expected error
	for _, messagePart := range expectedErrorParts {
		assert.Contains(t, actualError, messagePart, "\n\texpected: %s\n\tactual:   %s", origExp, origActual)
	}
}

// AssertMapsEqual compares two maps of with proto messages.
func AssertMapsEqual(t *testing.T, expected, actual map[string]proto.Message) {
	assert.Equal(t, len(expected), len(actual))
	for key, actualVal := range actual {
		if proto.Equal(expected[key], actualVal) {
			continue
		}
		assert.Equal(t, expected[key], actualVal)
	}
}

// AssertListsEqual compares two lists of proto messages.
func AssertListsEqual(t *testing.T, expected, actual interface{}) {
	e, a := reflect.ValueOf(expected), reflect.ValueOf(actual)
	if ekind, akind := e.Kind(), a.Kind(); ekind != akind || ekind != reflect.Slice {
		assert.Equal(t, expected, actual)
	}
	assert.Equal(t, e.Len(), a.Len())
	for i := 0; i < e.Len(); i++ {
		expVal, actVal := e.Index(i).Interface(), a.Index(i).Interface()
		expMsg, msgOk := expVal.(proto.Message)
		actualMsg, actOk := actVal.(proto.Message)
		if msgOk && actOk && proto.Equal(expMsg, actualMsg) {
			continue
		}
		assert.Equal(t, expVal, actVal)
	}
}

// AssertEqual is a wrapper for assert.Equal which properly handles proto.Message comparisons
func AssertEqual(t *testing.T, expected, actual interface{}) {
	em, emOk := expected.(proto.Message)
	am, amOk := actual.(proto.Message)
	if emOk && amOk && proto.Equal(em, am) {
		return
	}
	assert.Equal(t, expected, actual)
}

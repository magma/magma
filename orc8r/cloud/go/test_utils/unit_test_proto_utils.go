/*
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

package test_utils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

const (
	Separator = "$$"
)

// AssertErrorsEqual checks if actualError contains relevant parts of expectedError.
// This method of comparison is necessary since string representations of proto messages contain
// non-deterministic whitespaces (https://github.com/golang/protobuf/issues/1269)
func AssertErrorsEqual(t *testing.T, expectedError, actualError string) {
	// split expected error string at `separator`
	expectedErrorParts := strings.Split(expectedError, Separator)
	// check if actual message contains relevant parts of expected error
	for _, messagePart := range expectedErrorParts {
		assert.Contains(t, actualError, messagePart, "actual error does not contain all parts of "+
			"expected error\nexpected: %v\nactual:   %v", strings.Join(expectedErrorParts, " "), actualError)
	}
}

// AssertMapsEqual compares two maps of proto messages.
func AssertMapsEqual(t *testing.T, expected, actual interface{}) {
	expValue, actValue := reflect.ValueOf(expected), reflect.ValueOf(actual)
	ekind, akind := expValue.Kind(), actValue.Kind()
	if ekind != reflect.Map && akind != reflect.Map {
		assert.Fail(t, fmt.Sprintf("Type error: Both arguments are not Maps \nexpected: %+v\nactual : %+v", ekind, akind))
	} else if ekind != reflect.Map {
		assert.Fail(t, fmt.Sprintf("Type error: expected argument is not a a Map \nexpected:  %+v", ekind))
	} else if akind != reflect.Map {
		assert.Fail(t, fmt.Sprintf("Type error: second argument is not a Map \nactual  : %+v", akind))
	} else if assert.Equal(t, expValue.Len(), actValue.Len()) {
		for _, key := range expValue.MapKeys() {
			if assert.True(t, actValue.MapIndex(key).IsValid(), "Key %v is missing in actual:%v", key, actValue.MapKeys()) {
				AssertMessagesEqual(t, expValue.MapIndex(key).Interface(), actValue.MapIndex(key).Interface())
			}
		}
	}
}

// AssertListsEqual compares two lists of proto messages.
func AssertListsEqual(t *testing.T, expected, actual interface{}) {
	expValue, actValue := reflect.ValueOf(expected), reflect.ValueOf(actual)
	ekind, akind := expValue.Kind(), actValue.Kind()
	if ekind != reflect.Slice && akind != reflect.Slice {
		assert.Fail(t, fmt.Sprintf("Type error: Both arguments are not Slices \nexpected: %+v\nactual : %+v", ekind, akind))
	} else if ekind != reflect.Slice {
		assert.Fail(t, fmt.Sprintf("Type error: expected argument is not a a Slice \nexpected:  %+v", ekind))
	} else if akind != reflect.Slice {
		assert.Fail(t, fmt.Sprintf("Type error: second argument is not a Slice \nactual  : %+v", akind))
	} else {
		assert.Equal(t, expValue.Len(), actValue.Len())
		for i := 0; i < expValue.Len(); i++ {
			AssertMessagesEqual(t, expValue.Index(i).Interface(), actValue.Index(i).Interface())
		}
	}
}

// AssertMessagesEqual compares two proto.Messages with proto.Equal.
// Prints string representations of messages upon inequality.
// String representations provided by protobuf are not stable.
// Please only pay attention to the field values themselves, not to whitespaces
// in between.
func AssertMessagesEqual(t *testing.T, expected, actual interface{}) {
	expMsg, expOk := expected.(proto.Message)
	actMsg, actOk := actual.(proto.Message)
	if !expOk || !actOk {
		assert.Equal(t, expected, actual)
	} else if !proto.Equal(expMsg, actMsg) {
		assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %+v\n"+
			"actual  : %+v\n"+
			"Note that string representations provided by protobuf are not guaranteed to be stable and the comparison "+
			"between messages does not rely on them. Please pay attention to differences in the field values only.",
			expMsg, actMsg))
	}
}

// AssertMessagesNotEqual compares two proto.Messages with proto.Equal.
// Fails if messages are equal and prints string representations of messages.
// String representations provided by protobuf are not stable.
// Please only pay attention to the field values themselves, not to whitespaces
// in between.
func AssertMessagesNotEqual(t *testing.T, expected, actual interface{}) {
	expMsg, expOk := expected.(proto.Message)
	actMsg, actOk := actual.(proto.Message)
	if !expOk || !actOk {
		assert.NotEqual(t, expected, actual)
	} else if proto.Equal(expMsg, actMsg) {
		assert.Fail(t, "Messages are equal but should be different.")
	}
}

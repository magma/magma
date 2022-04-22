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
	"strings"
	"testing"

	protoV1 "github.com/golang/protobuf/proto"
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
		assert.Contains(t, actualError, messagePart, "actual error does not contain all parts of expected error\nexpected: %v\nactual:   %v", strings.Join(expectedErrorParts, " "), actualError)
	}
}

// AssertMapsEqual compares two maps of with proto messages.
func AssertMapsEqual(t *testing.T, expected, actual map[string]protoV1.Message) {
	assert.Equal(t, len(expected), len(actual))
	for key := range actual {
		AssertMessagesEqual(t, protoV1.MessageV2(expected[key]), protoV1.MessageV2(actual[key]))
	}
}

// AssertMessagesEqual compares two proto.Message's with proto.Equal.
// Prints string representations of messages upon inequality.
// Note, that string representations provided by protobuf are not stable.
// Please only pay attention to the field values themselves, not to whitespaces
// in between.
func AssertMessagesEqual(t *testing.T, expected, actual proto.Message) bool {
	equal := proto.Equal(expected, actual)
	if !equal {
		assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %+v\n"+
			"actual: %+v\n"+
			"Note, that string representations provided by protobuf are not guaranteed to be stable and the comparison between messages does not rely on them. Please pay attention to differences in the field values only.",
			expected, actual))
	}
	return true
}

// AssertMessagesNotEqual compares two proto.Message's with proto.Equal.
// Fails if messages are equal and prints string representations of messages.
// Note, that string representations provided by protobuf are not stable.
// Please only pay attention to the field values themselves, not to whitespaces
// in between.
func AssertMessagesNotEqual(t *testing.T, expected, actual proto.Message) bool {
	equal := proto.Equal(expected, actual)
	if equal {
		assert.Fail(t, "Messages are equal but should be different.")
	}
	return true
}

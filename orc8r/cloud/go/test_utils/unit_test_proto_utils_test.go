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

package test_utils

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"magma/orc8r/lib/go/protos"
	mconfig_protos "magma/orc8r/lib/go/protos/mconfig"
)

type testCase struct {
	expected    interface{}
	actual      interface{}
	test_result bool
}

func TestAssertMapsEqual(t *testing.T) {
	testCases := []testCase{
		{
			actual:      1,
			expected:    1,
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13},
			expected:    1,
			test_result: false,
		}, {
			actual:      1,
			expected:    map[string]int{"k1": 7, "k2": 13},
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13},
			expected:    1,
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13},
			expected:    map[string]int{},
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13},
			expected:    map[string]int{"k1": 7, "k3": 13},
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13, "k3": 13},
			expected:    map[string]int{"k1": 8, "k2": 13, "k3": 13},
			test_result: false,
		}, {
			actual:      map[string]int{"k1": 7, "k2": 13, "k3": 13},
			expected:    map[string]int{"k1": 7, "k2": 13, "k3": 13},
			test_result: true,
		}, {
			actual: map[string]int{"k1": 7, "k2": 13, "k3": 13},
			expected: map[string]proto.Message{
				"eventd":         &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
				"state":          &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				"shared_mconfig": &mconfig_protos.SharedMconfig{SentryConfig: nil},
			},
			test_result: false,
		}, {
			actual: map[string]proto.Message{
				"eventd":         &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
				"state":          &mconfig_protos.State{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				"shared_mconfig": &mconfig_protos.SharedMconfig{SentryConfig: nil},
			},
			expected: map[string]proto.Message{
				"eventd":         &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
				"state":          &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				"shared_mconfig": &mconfig_protos.SharedMconfig{SentryConfig: nil},
			},
			test_result: false,
		}, {
			actual: map[string]proto.Message{
				"eventd":         &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
				"state":          &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				"shared_mconfig": &mconfig_protos.SharedMconfig{SentryConfig: nil},
			},
			expected: map[string]proto.Message{
				"eventd":         &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
				"state":          &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				"shared_mconfig": &mconfig_protos.SharedMconfig{SentryConfig: nil},
			},
			test_result: true,
		},
	}
	runTestCases(t, testCases, AssertMapsEqual)
}

func TestAssertListsEqual(t *testing.T) {
	testCases := []testCase{
		{
			actual:      1,
			expected:    1,
			test_result: false,
		}, {
			actual: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			expected:    1,
			test_result: false,
		}, {
			actual: 1,
			expected: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			test_result: false,
		}, {
			actual: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 5, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			expected: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			test_result: false,
		}, {
			actual: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			expected: []*mconfig_protos.State{
				{SyncInterval: 2, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
				{SyncInterval: 4, LogLevel: protos.LogLevel_INFO},
			},
			test_result: true,
		},
	}
	runTestCases(t, testCases, AssertListsEqual)
}

func TestAssertMessageEqual(t *testing.T) {
	testCases := []testCase{
		{
			actual:      1,
			expected:    1,
			test_result: true,
		}, {
			actual:      1,
			expected:    &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			test_result: false,
		}, {
			actual:      &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			expected:    1,
			test_result: false,
		}, {
			actual:      &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			expected:    &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			test_result: false,
		}, {
			actual:      &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			expected:    &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			test_result: true,
		},
	}
	runTestCases(t, testCases, AssertMessagesEqual)

}

func TestAssertMessageNotEqual(t *testing.T) {
	testCases := []testCase{
		{
			actual:      1,
			expected:    1,
			test_result: false,
		}, {
			actual:      1,
			expected:    &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			test_result: true,
		}, {
			actual:      &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			expected:    1,
			test_result: true,
		}, {
			actual:      &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			expected:    &mconfig_protos.State{SyncInterval: 3, LogLevel: protos.LogLevel_INFO},
			test_result: true,
		}, {
			actual:      &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			expected:    &mconfig_protos.EventD{LogLevel: protos.LogLevel_INFO, EventVerbosity: 0},
			test_result: false,
		},
	}
	runTestCases(t, testCases, AssertMessagesNotEqual)

}

func runTestCases(t *testing.T, testCases []testCase, assertFunction func(t *testing.T, expected, actual interface{})) {
	_t := new(testing.T)
	for _, testCase := range testCases {
		if testCase.test_result {
			assertFunction(t, testCase.expected, testCase.actual)
		} else {
			assertFunction(_t, testCase.expected, testCase.actual)
			if !_t.Failed() {
				t.Error("Error: Arguments are asserted to be equal but are not equal")
			}
		}
	}
}

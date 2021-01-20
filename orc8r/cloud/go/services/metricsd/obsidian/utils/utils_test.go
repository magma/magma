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

package utils_test

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/obsidian/utils"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	exampleUnixTimeString := "1547751342"
	exampleUnixFloatString := "1547751342.23"
	exampleUnixTime := time.Unix(1547751342, 0)

	exampleRFCTimeString := "2018-07-01T20:10:30.781Z"
	exampleBadRFCTimeString := "2018-07-01T20:10.781Z"
	exampleRFCTime, err := time.Parse(time.RFC3339, "2018-07-01T20:10:30.781Z")
	assert.NoError(t, err)

	defaultTime := time.Time{}

	time, err := utils.ParseTime(exampleUnixTimeString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleUnixTime, time)

	time, err = utils.ParseTime(exampleUnixFloatString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleUnixTime, time)

	time, err = utils.ParseTime(exampleRFCTimeString, &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, exampleRFCTime, time)

	_, err = utils.ParseTime(exampleBadRFCTimeString, &defaultTime)
	assert.Error(t, err)

	time, err = utils.ParseTime("", &defaultTime)
	assert.NoError(t, err)
	assert.Equal(t, time, defaultTime)

	_, err = utils.ParseTime("", nil)
	assert.Error(t, err)
}

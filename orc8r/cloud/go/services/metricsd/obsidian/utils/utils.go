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

package utils

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

const (
	ParamQuery      = "query"
	ParamRangeStart = "start"
	ParamRangeEnd   = "end"
	ParamStepWidth  = "step"
	ParamTime       = "time"

	ParamMetric      = "metric"
	ParamMatchTarget = "match_target"
	ParamLimit       = "limit"

	StatusSuccess = "success"
)

func ParseTime(timeString string, defaultTime *time.Time) (time.Time, error) {
	if timeString == "" {
		if defaultTime != nil {
			return *defaultTime, nil
		}
		return time.Time{}, fmt.Errorf("time parameter not provided")
	}
	time, err := parseUnixTime(timeString)
	if err == nil {
		return time, nil
	}
	return parseRFCTime(timeString)
}

func parseUnixTime(timeString string) (time.Time, error) {
	timeNum, err := strconv.ParseFloat(timeString, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(timeNum), 0), nil
}

func parseRFCTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}

func ParseDuration(durationString, defaultDuration string) (time.Duration, error) {
	if durationString == "" {
		durationString = defaultDuration
	}
	// If last char is a digit, append 's' since number of seconds is assumed
	if unicode.IsDigit(rune(durationString[len(durationString)-1])) {
		durationString += "s"
	}
	return time.ParseDuration(durationString)
}

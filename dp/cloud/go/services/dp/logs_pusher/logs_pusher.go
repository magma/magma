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

package logs_pusher

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type DPLog struct {
	EventTimestamp   int64  `json:"event_timestamp"`
	LogFrom          string `json:"log_from"`
	LogTo            string `json:"log_to"`
	LogName          string `json:"log_name"`
	LogMessage       string `json:"log_message"`
	CbsdSerialNumber string `json:"cbsd_serial_number"`
	NetworkId        string `json:"network_id"`
	FccId            string `json:"fcc_id"`
}

type LogPusher func(ctx context.Context, log *DPLog, consumerUrl string) error

func PushDPLog(ctx context.Context, log *DPLog, consumerUrl string) error {
	client := http.Client{}
	body, _ := json.Marshal(log)
	req, _ := http.NewRequest(http.MethodPost, consumerUrl, strings.NewReader(string(body)))
	req.Header.Set("contentType", "application/json")
	_, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	return nil
}

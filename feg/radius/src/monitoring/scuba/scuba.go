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

package scuba

import (
	"encoding/json"
	"errors"
	"fbc/cwf/radius/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ANY_SCUBA_CATEGORY json_to_any_scuba category
// See https://our.internmc.facebook.com/intern/wiki/Internet.org/Express_Wi-Fi_(XWF)/Scuba_Tailer/
// for more details
const ANY_SCUBA_CATEGORY string = "xwf_json_to_any_scuba"

type scubaWriteSyncer struct {
	disabled bool
	config   *config.Scuba
	url      url.URL
	table    string
	msgQ     chan string
}

func (s *scubaWriteSyncer) Write(p []byte) (int, error) {
	if s.disabled {
		return 0, errors.New("Logger is already closed, cannot write")
	}
	s.msgQ <- string(p)
	return len(p), nil
}

func (s *scubaWriteSyncer) Sync() error {
	return nil
}

func (s *scubaWriteSyncer) Close() error {
	s.disabled = true
	return nil
}

// ScribeEntry ...
type ScribeEntry struct {
	Category string `json:"category"`
	Message  string `json:"message"`
}

func (s *scubaWriteSyncer) makeScribeEntry(msg string) ScribeEntry {
	return ScribeEntry{
		Category: ANY_SCUBA_CATEGORY,
		Message:  fmt.Sprintf("perfpipe_%s %s", s.table, msg),
	}
}

func (s *scubaWriteSyncer) serve() {
	for {
		// Grab messages from the queue
		var messages []ScribeEntry
		firstMsg := <-s.msgQ
		messages = append(messages, s.makeScribeEntry(firstMsg))
		flush := time.NewTimer(time.Second * time.Duration(s.config.FlushIntervalSec))

	Remaining:
		for i := 0; i < s.config.BatchSize-1; i++ {
			select {
			case msg := <-s.msgQ:
				messages = append(
					messages,
					s.makeScribeEntry(msg),
				)
			case <-flush.C:
				break Remaining
			default:
				break Remaining
			}
		}

		flush.Stop()

		// Break or go back to wait
		if len(messages) == 0 {
			if s.disabled {
				break // exit the go routine
			}
			continue
		}

		// Build the message
		msgs, err := json.Marshal(messages)
		if err != nil {
			fmt.Printf("ERROR serializing %d log(s): %s\n", len(messages), err.Error())
			continue
		}

		form := url.Values{
			"access_token": []string{s.config.AccessToken},
			"logs":         []string{string(msgs)},
		}

		// Do Post
		res, err := http.Post(
			s.config.GraphURL,
			"application/x-www-form-urlencoded",
			strings.NewReader(form.Encode()),
		)
		if err != nil {
			fmt.Printf("ERROR sending %d log(s) to Scuba: %s (target: %s)\n", len(messages), err.Error(), s.config.GraphURL)
			continue
		}

		if res.StatusCode != 200 {
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf(
					"ERROR sending %d log(s) to Scuba: Got status code %d (%s)\n",
					len(messages),
					res.StatusCode,
					res.Status,
				)
			} else {
				fmt.Printf(
					"ERROR sending %d logs to Scuba: Got status code %d (%s): %s\n",
					len(messages),
					res.StatusCode,
					res.Status,
					string(bodyBytes),
				)
			}
		}
	}
}

// Initialize ...
func Initialize(config *config.Scuba, logger *zap.Logger) {
	zap.RegisterSink(
		"scuba",
		func(url *url.URL) (zap.Sink, error) {
			result := &scubaWriteSyncer{
				disabled: false,
				config:   config,
				url:      *url,
				table:    url.Hostname(),
				msgQ:     make(chan string, config.MessageQueueSize),
			}
			go result.serve()
			return result, nil
		},
	)
}

// NewLogger creates a new Scuba logger
func NewLogger(table string, options ...zap.Option) (*zap.Logger, error) {
	// Create configuration
	c := zap.NewProductionConfig()
	c.Level.SetLevel(zap.DebugLevel)
	if c.OutputPaths == nil {
		c.OutputPaths = []string{}
	}
	c.OutputPaths = append(c.OutputPaths, fmt.Sprintf("scuba://%s", table))
	return c.Build(options...)
}

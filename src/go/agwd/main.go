// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/agwd/server"
	"github.com/magma/magma/src/go/log"
	"github.com/magma/magma/src/go/log/zap"
)

func main() {
	configFlag := flag.String(
		"c", "/etc/magma/agwd.json", "Path to config file")
	flag.Parse()

	cfgr := config.NewConfigManager()
	cfgr_err := config.LoadConfigFile(cfgr, *configFlag)

	lm := log.NewManager(zap.NewLogger())
	lm.
		LoggerFor("").
		SetLevel(config.LogLevel(cfgr.Config().GetLogLevel()))

	if cfgr_err != nil {
		lm.LoggerFor("").Warning().Printf("using default configuration as LoadConfigFile failed with %q", cfgr_err)
	}

	server.Start(cfgr, lm.LoggerFor("server"))

	stopper := make(chan struct{})
	<-stopper
}

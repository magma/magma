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
	uber_zap "go.uber.org/zap"

	"github.com/magma/magma/accessd/server"
	"github.com/magma/magma/log"
	"github.com/magma/magma/log/zap"
)

func main() {
	lm := log.NewManager(zap.NewLogger(uber_zap.NewDevelopmentConfig()))

	if err := server.Start(lm.LoggerFor("server")); err != nil {
		panic(err)
	}
}

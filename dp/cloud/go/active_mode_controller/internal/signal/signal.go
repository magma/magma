/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type app interface {
	Run(ctx context.Context) error
}

func Run(ctx context.Context, app app) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	appCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()
	return app.Run(appCtx)
}

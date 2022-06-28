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

package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"magma/dp/cloud/go/active_mode_controller/config"
	"magma/dp/cloud/go/active_mode_controller/internal/app"
	"magma/dp/cloud/go/active_mode_controller/internal/signal"
	"magma/dp/cloud/go/active_mode_controller/internal/time"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Printf("failed to read config: %s", err)
		os.Exit(1)
	}
	clock := &time.Clock{}
	seed := rand.NewSource(clock.Now().Unix())
	a := app.NewApp(
		app.WithConfig(cfg),
		app.WithClock(clock),
		app.WithRNG(rand.New(seed)),
	)
	ctx := context.Background()
	if err := signal.Run(ctx, a); err != nil && err != context.Canceled {
		log.Printf("failed to stop app: %s", err)
		os.Exit(1)
	}
}

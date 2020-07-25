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

package lbcanary

import (
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

const maxCanariesWeight = 50

var r *rand.Rand

var totalCanariesWeight int
var lbConfig *config.LoadBalanceConfig

// Init filter interface implementation
func Init(c *config.ServerConfig) error {
	lbConfig = &c.LoadBalance
	r = rand.New(rand.NewSource(time.Now().Unix()))
	totalCanariesWeight = 0
	for _, canary := range lbConfig.Canaries {
		if canary.Name == config.LiveTier {
			return errors.New("reserved canary name 'live' specified")
		}
		totalCanariesWeight += canary.TrafficSlicePercent
	}
	if totalCanariesWeight > maxCanariesWeight {
		return errors.New("canaries are over allocated")
	}
	return nil
}

// Process filter interface implementation
func Process(c *modules.RequestContext, _ string, _ *radius.Request) error {
	// Register Upstream host
	if err := allocateTier(c); err != nil {
		return err
	}

	return nil
}

func allocateTier(c *modules.RequestContext) error {
	// Load session state
	state, err := c.SessionStorage.Get()
	if err != nil {
		c.Logger.Debug(
			"New session detected",
			zap.Error(err),
		)
		state = &session.State{}
	}

	if state.UpstreamHost != "" {
		return nil
	}

	counter := monitoring.NewOperation("pick_canary_tier").Start()
	tier := pickRandomTier()

	state.Tier = tier
	err = c.SessionStorage.Set(*state)
	if err != nil {
		c.Logger.Error(
			"Error persisting allocated upstream host",
			zap.Error(err),
		)
		counter.Failure("storage_error")
		return err
	}
	counter.Success()
	return nil
}

func pickRandomTier() string {
	selection := 1 + r.Intn(100)
	tier := config.LiveTier
	for _, canary := range lbConfig.Canaries {
		selection -= canary.TrafficSlicePercent
		if selection <= 0 {
			tier = canary.Name
			break
		}
	}
	return tier
}

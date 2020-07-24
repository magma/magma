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

package lballocate

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

const errRequiredTierNotSpecifiedText = "required Tier or DefaultTier not specified"
const errCanaryNotFoundText = "canary not found"
const errServiceTierNotFoundText = "service tier definition not found"
const errListenerRoutingNotFoundText = "listener routing definition not found"
const errNoUpstreamHostsText = "upstream hosts not specified"

var errRequiredTierNotSpecified = errors.New(errRequiredTierNotSpecifiedText)
var errCanaryNotFound = errors.New(errCanaryNotFoundText)
var errServiceTierNotFound = errors.New(errServiceTierNotFoundText)
var errListenerRoutingNotFound = errors.New(errListenerRoutingNotFoundText)
var errNoUpstreamHosts = errors.New(errNoUpstreamHostsText)

var r *rand.Rand

var lbConfig *config.LoadBalanceConfig
var tierRoutings map[string]map[string]string
var serviceTiers map[string][]string

// Init module interface implementation
//nolint:unparam
func Init(c *config.ServerConfig) error {
	lbConfig = &c.LoadBalance

	// prepare map: tier -> (map: listener -> service tier)
	tierRoutings = make(map[string]map[string]string)
	tierRoutings[config.LiveTier] = getListenerServiceTierLookup(&lbConfig.LiveTier)
	for _, canary := range lbConfig.Canaries {
		tierRoutings[canary.Name] = getListenerServiceTierLookup(&canary.Routing)
	}

	// prepare ServiceTier.Name->UpstreamServers lookup
	serviceTiers = make(map[string][]string)
	for _, serviceTier := range lbConfig.ServiceTiers {
		serviceTiers[serviceTier.Name] = serviceTier.UpstreamHosts
	}

	r = rand.New(rand.NewSource(time.Now().Unix()))
	return nil
}

func getListenerServiceTierLookup(tierRouting *config.TierRouting) map[string]string {
	listenerServiceTierMap := make(map[string]string)
	for _, route := range tierRouting.Routes {
		listenerServiceTierMap[route.Listener] = route.ServiceTier
	}
	return listenerServiceTierMap
}

// Process module interface implementation
func Process(c *modules.RequestContext, listenerName string, _ *radius.Request) error {
	// Register Upstream host
	if err := allocateUpstreamHost(c, listenerName); err != nil {
		return err
	}

	return nil
}

func allocateUpstreamHost(c *modules.RequestContext, listenerName string) error {
	// Load session state
	state, err := c.SessionStorage.Get()
	if err != nil {
		c.Logger.Error(
			"error loading session state, unable to register upstream host(s)",
			zap.Error(err),
		)
		return err
	}

	if state.UpstreamHost != "" {
		return nil
	}

	counter := monitoring.NewOperation("pick_upstream_host").Start()

	upstreamHost, err := pickRandomUpstreamHost(c, state, listenerName)
	if err != nil {
		counter.Failure("allocation_error")
		return err
	}

	state.UpstreamHost = upstreamHost
	err = c.SessionStorage.Set(*state)
	if err != nil {
		c.Logger.Error(
			"error persisting allocated upstream host",
			zap.Error(err),
		)
		counter.Failure("storage_error")
		return err
	}
	counter.Success()
	return nil
}

func pickRandomUpstreamHost(c *modules.RequestContext, state *session.State, listenerName string) (string, error) {

	tier, err := getTier(state)
	if err != nil {
		return "", err
	}

	var upstreamHost string
	upstreamHosts, hostCount, err := getUpstreamHosts(c, tier, listenerName)
	if err != nil {
		return "", err
	}

	selection := r.Intn(hostCount)
	upstreamHost = upstreamHosts[selection]
	return upstreamHost, nil
}

func getTier(state *session.State) (string, error) {
	if state.Tier != "" {
		return state.Tier, nil
	}
	if lbConfig.DefaultTier != "" {
		return lbConfig.DefaultTier, nil
	}
	return "", errRequiredTierNotSpecified
}

func getUpstreamHosts(c *modules.RequestContext, tier string, listenerName string) ([]string, int, error) {
	listenerServiceTierLookup, found := tierRoutings[tier]
	if !found {
		c.Logger.Error(errCanaryNotFoundText,
			zap.String("canary", tier))
		return nil, 0, errCanaryNotFound
	}

	serviceTier, found := listenerServiceTierLookup[listenerName]
	if !found {
		c.Logger.Error(errListenerRoutingNotFoundText,
			zap.String("canary", tier),
			zap.String("listener", listenerName))
		return nil, 0, errListenerRoutingNotFound
	}

	upstreamHosts, found := serviceTiers[serviceTier]
	if !found {
		c.Logger.Error(errServiceTierNotFoundText,
			zap.String("service_tier", serviceTier))
		return nil, 0, errServiceTierNotFound
	}

	hostCount := len(upstreamHosts)
	if hostCount == 0 {
		c.Logger.Error(errNoUpstreamHostsText,
			zap.String("service_tier", serviceTier))
		return nil, 0, errNoUpstreamHosts
	}

	return upstreamHosts, hostCount, nil
}

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

package active_mode_controller

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/glog"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator"
	"magma/dp/cloud/go/services/dp/storage"
)

type App struct {
	db                    *sql.DB
	clock                 Clock
	rng                   action_generator.RNG
	heartbeatSendTimeout  time.Duration
	pollingInterval       time.Duration
	cbsdInactivityTimeout time.Duration
	amcManager            storage.AmcManager
}

func NewApp(options ...Option) *App {
	a := &App{}
	for _, o := range options {
		o(a)
	}
	return a
}

type Clock interface {
	Now() time.Time
	Tick(duration time.Duration) *time.Ticker
}

type Option func(*App)

func WithDb(db *sql.DB) Option {
	return func(a *App) { a.db = db }
}

func WithAmcManager(manager storage.AmcManager) Option {
	return func(a *App) { a.amcManager = manager }
}

func WithRNG(rng action_generator.RNG) Option {
	return func(a *App) { a.rng = rng }
}

func WithClock(clock Clock) Option {
	return func(a *App) { a.clock = clock }
}

func WithHeartbeatSendTimeout(sendTimeout time.Duration, sendInterval time.Duration) Option {
	return func(a *App) { a.heartbeatSendTimeout = sendTimeout + sendInterval }
}

func WithPollingInterval(interval time.Duration) Option {
	return func(a *App) { a.pollingInterval = interval }
}

func WithCbsdInactivityTimeout(timeout time.Duration) Option {
	return func(a *App) { a.cbsdInactivityTimeout = timeout }
}

func (a *App) Run(ctx context.Context) error {
	ticker := a.clock.Tick(a.pollingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, err := storage.WithinTx(a.db, a.getStateAndProcessData)
			if err != nil {
				glog.Errorf("failed to process data: %s", err)
			}
		}
	}
}

// TODO add context
func (a *App) getStateAndProcessData(tx *sql.Tx) (any, error) {
	state, err := a.amcManager.GetState(tx)
	if err != nil {
		return nil, err
	}
	generator := &action_generator.ActionGenerator{
		HeartbeatTimeout:  a.heartbeatSendTimeout + a.pollingInterval,
		InactivityTimeout: a.cbsdInactivityTimeout,
		Rng:               a.rng,
	}
	now := a.clock.Now()
	actions := generator.GenerateActions(state, now)
	for _, act := range actions {
		if err := act.Do(tx, a.amcManager); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

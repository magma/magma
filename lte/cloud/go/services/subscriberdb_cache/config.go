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

package subscriberdb_cache

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	// SleepIntervalSecs is the time interval between each service worker loop.
	SleepIntervalSecs int `yaml:"sleepIntervalSecs"`
	// UpdateIntervalSecs is the target time interval to update each digest.
	UpdateIntervalSecs int `yaml:"updateIntervalSecs"`
}

func (config Config) Validate() error {
	errs := &multierror.Error{}
	if config.SleepIntervalSecs <= 0 {
		errs = multierror.Append(errs, fmt.Errorf("invalid worker sleep interval"))
	}
	if config.UpdateIntervalSecs < 60 {
		errs = multierror.Append(errs, fmt.Errorf("worker update interval smaller than 1 minute"))
	}
	return errs.ErrorOrNil()
}

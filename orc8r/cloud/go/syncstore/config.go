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

package syncstore

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	// CacheWriterValidIntervalSecs specifies the time duration (in secs) after
	// which a cacheWriter object is subject to garbage collection.
	//
	// NOTE: the caller should enforce that this value is smaller than the service
	// worker loop interval, to prevent workers with "older" cache writers overwriting
	// concurrent workers with "newer" ones.
	CacheWriterValidIntervalSecs int64
	// TableNamePrefix is used to namespace all syncstore tables to prevent
	// naming collisions among different services using syncstore.
	TableNamePrefix string
}

func (config *Config) Validate(writer bool) error {
	errs := &multierror.Error{}
	if config.TableNamePrefix == "" {
		errs = multierror.Append(errs, errors.New("empty table name prefix for syncstore"))
	}
	if writer && config.CacheWriterValidIntervalSecs <= 0 {
		errs = multierror.Append(errs, fmt.Errorf("invalid cache writer valid interval: %+v", config.CacheWriterValidIntervalSecs))
	}
	return errs.ErrorOrNil()
}

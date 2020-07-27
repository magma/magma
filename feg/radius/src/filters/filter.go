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

package filters

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
)

type (
	// Filter represents a request filter action
	Filter interface {
		Init(c *config.ServerConfig) error
		Process(c *modules.RequestContext, l string, r *radius.Request) error
	}

	// FilterInitFunc type for filter's Init function
	FilterInitFunc func(c *config.ServerConfig) error

	// FilterProcessFunc type for filter's Process function
	FilterProcessFunc func(c *modules.RequestContext, l string, r *radius.Request) error
)

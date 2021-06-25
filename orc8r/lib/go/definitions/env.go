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

package definitions

import (
	"os"

	"github.com/pkg/errors"
)

// GetEnvWithDefault returns the string value of the environment variable,
// defaulting to a specified value if it doesn't exist.
func GetEnvWithDefault(variable string, defaultValue string) string {
	value := os.Getenv(variable)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

// MustGetEnv returns the string value of the environment variable,
// panics it doesn't exist
func MustGetEnv(variable string) string {
	value := os.Getenv(variable)
	if len(value) == 0 {
		panic(errors.Errorf("%s env not found", variable))
	}
	return value
}

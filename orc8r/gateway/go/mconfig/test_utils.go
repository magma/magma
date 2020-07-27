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

package mconfig

import (
	"io/ioutil"
	"os"
	"time"
)

func CreateLoadTempConfig(configJSON string) error {
	StopRefreshTicker()

	tmpfile, err := ioutil.TempFile("", "mconfig_test")
	if err != nil {
		return err
	}
	// Write marshaled configs
	if _, err = tmpfile.Write([]byte(configJSON)); err != nil {
		return err
	}
	mcpath := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(mcpath)

	time.Sleep(time.Second) // give extra time for in-flight RefreshTicker to complete

	return RefreshConfigsFrom(mcpath)
}

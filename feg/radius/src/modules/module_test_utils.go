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

package modules

import (
	"context"
	"fmt"
	"os"
	"time"

	"layeh.com/radius"
)

func WaitForRadiusServerToBeReady(secret []byte, addr string) (err error) {
	// stop printing custom server messages until the server is up
	temp := os.Stdout
	os.Stdout = nil
	defer func() {
		os.Stdout = temp
	}()
	MaxRetries := 20
	for r := 0; r < MaxRetries; r++ {
		_, err = radius.Exchange(
			context.Background(),
			radius.New(radius.CodeStatusServer, secret),
			addr,
		)
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fmt.Errorf(
		"radius server failed to be ready after %d retries: %v",
		MaxRetries, err,
	)
}

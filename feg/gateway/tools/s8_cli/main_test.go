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

//go:build cli_test
// +build cli_test

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestS8Cli(t *testing.T) {
	// replace the arguments
	os.Args = []string{"s8_cli", "cs", "-test", "-delete", "0"}
	fmt.Printf("\ncommand to run: %+v\n", os.Args)
	main()
}

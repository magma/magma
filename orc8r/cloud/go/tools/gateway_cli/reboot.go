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

package main

import (
	"os"

	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdReboot := &cobra.Command{
		Use:   "reboot --network=<network-id> --gateway=<gateway-id>",
		Short: "reboot gateway device",
		Run:   rebootCmd,
	}

	rootCmd.AddCommand(cmdReboot)
}

func rebootCmd(cmd *cobra.Command, args []string) {
	err := magmad.GatewayReboot(networkId, gatewayId)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}

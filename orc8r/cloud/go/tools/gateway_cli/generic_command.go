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
	"fmt"
	"os"

	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/spf13/cobra"
)

func init() {
	cmdGenericCommand := &cobra.Command{
		Use:   "generic_command <command> <params>",
		Short: "Execute generic command on gateway",
		Args:  cobra.ExactArgs(2),
		Run:   genericCommandCmd,
	}

	rootCmd.AddCommand(cmdGenericCommand)
}

func genericCommandCmd(cmd *cobra.Command, args []string) {
	paramsStruct := structpb.Struct{}
	err := jsonpb.UnmarshalString(args[1], &paramsStruct)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	genericCommandParams := protos.GenericCommandParams{
		Command: args[0],
		Params:  &paramsStruct,
	}

	response, err := magmad.GatewayGenericCommand(networkId, gatewayId, &genericCommandParams)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", response)
}

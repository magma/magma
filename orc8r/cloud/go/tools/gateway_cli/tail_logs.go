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
	"io"
	"os"
	"os/signal"
	"syscall"

	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdTailLogs := &cobra.Command{
		Use:   "tail_logs [service]",
		Short: "tail gateway logs",
		Args:  cobra.MaximumNArgs(1),
		Run:   tailLogsCmd,
	}

	rootCmd.AddCommand(cmdTailLogs)
}

func tailLogsCmd(cmd *cobra.Command, args []string) {
	var service string
	if len(args) == 1 {
		service = args[0]
	}
	stream, err := magmad.TailGatewayLogs(networkId, gatewayId, service)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	// https://stackoverflow.com/q/11268943
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-term
	}()
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			break
		}
		fmt.Print(line.Line)
	}
}

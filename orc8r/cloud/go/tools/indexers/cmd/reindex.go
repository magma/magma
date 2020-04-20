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

package cmd

import (
	"context"
	"io"
	"log"
	"time"

	"magma/orc8r/cloud/go/services/state/protos"

	"github.com/spf13/cobra"
)

var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Kick off required reindex jobs, blocking till complete",
	Run:   runReindex,
}

func init() {
	rootCmd.AddCommand(reindexCmd)
	reindexCmd.Flags().StringVarP(&reindexID, "id", "i", "", "reindex specific indexer")
	reindexCmd.Flags().BoolVarP(&reindexForce, "force", "f", false, "force reindex even if automatic reindexing is enabled")
}

func runReindex(cmd *cobra.Command, args []string) {
	printAlive()
	stream, err := getClient().StartReindex(context.Background(), &protos.StartReindexRequest{IndexerId: reindexID, Force: reindexForce})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res.Update)
	}
}

// printAlive to convey that the CLI isn't hanging during long-running reindex operations.
func printAlive() {
	time.AfterFunc(10*time.Second, func() {
		stderrln("This may take a while")
	})
}

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
	"fmt"
	"log"

	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/protos"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all indexer IDs and versions",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listID, "id", "i", "", "restrict to specific indexer ID")
	listCmd.Flags().BoolVarP(&listShort, "name", "n", false, "only print indexer IDs")
}

func runList(cmd *cobra.Command, args []string) {
	res, err := getClient().GetIndexers(context.Background(), &protos.GetIndexersRequest{})
	if err != nil {
		log.Fatal(err)
	}

	if listID != "" {
		v, ok := res.IndexersById[listID]
		if !ok {
			log.Fatalf("No indexer found for ID %s", listID)
		}
		printVersions(indexer.MakeVersion(v))
		return
	}

	printVersions(indexer.MakeVersions(res.IndexersById)...)
}

func printVersions(versions ...*indexer.Versions) {
	for _, v := range versions {
		if listShort {
			fmt.Println(v.IndexerID)
			continue
		}
		fmt.Println(v)
	}
}

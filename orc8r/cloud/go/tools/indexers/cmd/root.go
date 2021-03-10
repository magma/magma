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
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"magma/orc8r/cloud/go/services/state"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/spf13/cobra"
)

var (
	// Global flag vars
	rootSilent   bool
	listID       string
	listShort    bool
	reindexID    string
	reindexForce bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootSilent, "silent", "s", false, "silence log output from loading Magma plugins")
}

var rootCmd = &cobra.Command{
	Use:              "indexers",
	Short:            "indexers CLI provides methods for viewing and managing state indexers",
	PersistentPreRun: globalPre,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func globalPre(cmd *cobra.Command, args []string) {
	if rootSilent {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	registry.MustPopulateServices()
}

func getClient() indexer_protos.IndexerManagerClient {
	conn, err := registry.GetConnection(state.ServiceName)
	if err != nil {
		log.Fatal(err)
	}
	return indexer_protos.NewIndexerManagerClient(conn)
}

func stderrln(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

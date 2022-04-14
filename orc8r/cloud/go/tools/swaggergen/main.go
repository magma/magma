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

/*
	Package main holds the swaggergen executable combined with combine-swagger
	The combine_swagger and swaggergen tools are strongly related and share a lot of functionality.
	Merging these tools as separate sub-commands under a single executable would reduce final container image size,
	as well as providing improved conceptual cohesion.

	combine_swagger command is now accessible by argument:
	swaggergen --combine
*/

package main

import (
	"flag"
	"os"

	combine_swagger "magma/orc8r/cloud/go/tools/swaggergen/combine_swagger"
	swaggergen "magma/orc8r/cloud/go/tools/swaggergen/swaggergen"
)

func main() {
	cmdCombine := flag.Bool("combine", false, "calls combine_swagger command")

	//swaggergen commands
	targetFilepath := flag.String("target", "", "Target swagger spec to generate code from")
	configFilepath := flag.String("config", "", "Config file for go-swagger command")
	rootDir := flag.String("root", os.Getenv("MAGMA_ROOT"), "Root path to resolve dependency and output directories based on")

	//combine_swagger commands
	inDir := flag.String("in", "", "Input directory")
	commonFilepath := flag.String("common", "", "Common definitions filepath")
	outFilepath := flag.String("out", "", "Output directory")
	generateStandAloneSpec := flag.Bool("standalone", true, "Generate standalone specs")

	flag.Parse()

	if *cmdCombine {
		combine_swagger.Run(inDir, commonFilepath, outFilepath, generateStandAloneSpec)
	} else {
		swaggergen.Run(targetFilepath, configFilepath, rootDir)
	}
}

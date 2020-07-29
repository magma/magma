// +build all qos
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

package integration

import "fmt"

func dumpPipelinedState(tr *TestRunner) {
	fmt.Println("******************* Dumping Pipelined State *******************")
	cmdList := [][]string{
		{"pipelined_cli.py", "debug", "qos"},
		{"pipelined_cli.py", "debug", "display_flows"},
	}
	cmdOutputList, err := tr.RunCommandInContainer("pipelined", cmdList)
	if err != nil {
		fmt.Printf("error dumping pipelined state %v", err)
		return
	}
	for _, cmdOutput := range cmdOutputList {
		fmt.Printf("command : \n%v\n", cmdOutput.cmd)
		fmt.Printf("output : \n%v\n", cmdOutput.output)
		fmt.Printf("error : \n%v\n", cmdOutput.err)
		fmt.Printf("\n\n")
	}
}

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
// package service_manager defines and implements API for service management
package service_manager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"

	"github.com/golang/glog"

	"magma/gateway/config"
)

var (
	registry = map[string]ServiceController{
		strings.ToLower(DockerController{}.Name()):  DockerController{},
		strings.ToLower(SystemdController{}.Name()): SystemdController{},
		strings.ToLower(RunitController{}.Name()):   RunitController{},
		"initd": InitdController{TailLogsCmd: DefaultInitdLogTailCmd},
		"procd": InitdController{TailLogsCmd: "logread -f"},
	}
	defaultController = DockerController{}
)

// Get returns Service Controller for configured init system or the default controller if match cannot be found
func Get() ServiceController {
	initSystem := strings.ToLower(config.GetMagmadConfigs().InitSystem)
	if contr, ok := registry[initSystem]; ok {
		return contr
	}
	glog.Warningf("process controller for '%s' cannot be found, using '%s' controller",
		initSystem, defaultController.Name())
	return defaultController
}

// StartCmdWithStderrStdoutTailer starts given command and waits for its completion,
// StartCmdWithStderrStdoutTailer creates the chan where all stderr & stdout strings are sent
// The chan will be closed when both stdout & stderr streams are closed or have errors
// To terminate the running command - use Process.Kill()
func StartCmdWithStderrStdoutTailer(cmd *exec.Cmd) (chan string, *os.Process, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to cerate stderr pipe for command '%s': %v", cmd.String(), err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to cerate stdout pipe for command '%s': %v", cmd.String(), err)
	}
	if err = cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed start command '%s': %v", cmd.String(), err)
	}
	c := make(chan string)
	var streamsDone = new(int32)
	go outputReader(stderr, c, streamsDone)
	go outputReader(stdout, c, streamsDone)
	go cmd.Wait()
	return c, cmd.Process, nil
}

func outputReader(rdr io.ReadCloser, outChan chan string, done *int32) {
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		outChan <- scanner.Text()
	}
	rdr.Close()
	if atomic.AddInt32(done, 1) == 2 {
		close(outChan)
	}
}

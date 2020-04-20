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
	"os"
	"os/exec"
)

// InitdController - initd based controller implementation
type InitdController struct {
	TailLogsCmd string
}

// DefaultInitdLogTailCmd logs tail command if none is specified
const DefaultInitdLogTailCmd = "tail -f /var/log/syslog"

// Name returns runit controller type name
func (InitdController) Name() string {
	return "initd"
}

// CmdName returns the base name of the start/stop/status command
func (InitdController) CmdName(service string) string {
	return "/etc/init.d/" + service
}

// Start starts service and returns error if unsuccessful
func (c InitdController) Start(service string) error {
	return exec.Command(c.CmdName(service), "start").Run()
}

// Stop stops service and returns error if unsuccessful
func (c InitdController) Stop(service string) error {
	return exec.Command(c.CmdName(service), "stop").Run()
}

// Restart restarts service and returns error if unsuccessful
func (c InitdController) Restart(service string) error {
	return exec.Command(c.CmdName(service), "restart").Run()
}

// GetState returns the given service state or error if unsuccessful
func (c InitdController) GetState(service string) (ServiceState, error) {
	err := exec.Command(c.CmdName(service), "status").Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr != nil && exitErr.ExitCode() != 0 {
			return Inactive, nil
		}
		return Error, err
	}
	return Active, nil
}

// TailLogs executes command to start tailing service logs and returns string chan to receive log strings
// closing the chan will terminate tailing
func (c InitdController) TailLogs(service string) (chan string, *os.Process, error) {
	var cmd *exec.Cmd
	tailCmd := c.TailLogsCmd
	if len(tailCmd) == 0 {
		tailCmd = DefaultInitdLogTailCmd
	}
	if len(service) == 0 {
		cmd = exec.Command("sh", "-c", tailCmd)
	} else {
		cmd = exec.Command("sh", "-c", tailCmd+" | grep \" "+service+"\"")
	}
	return StartCmdWithStderrStdoutTailer(cmd)
}

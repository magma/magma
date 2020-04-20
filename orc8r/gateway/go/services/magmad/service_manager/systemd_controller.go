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
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SystemdController - systemd based controller implementation
type SystemdController struct{}

// Name returns systemd controller type name
func (SystemdController) Name() string {
	return "systemd"
}

func (SystemdController) CmdName() string {
	return "systemctl"
}

func (SystemdController) ServiceName(service string) string {
	return fmt.Sprintf("magma@%s", service)
}

// Start starts service and returns error if unsuccessful
func (c SystemdController) Start(service string) error {
	return exec.Command(c.CmdName(), "start", c.ServiceName(service)).Run()
}

// Stop stops service and returns error if unsuccessful
func (c SystemdController) Stop(service string) error {
	return exec.Command(c.CmdName(), "stop", c.ServiceName(service)).Run()
}

// Restart restarts service and returns error if unsuccessful
func (c SystemdController) Restart(service string) error {
	return exec.Command(c.CmdName(), "reload-or-restart", c.ServiceName(service)).Run()
}

// GetState returns the given service state or error if unsuccessful
func (c SystemdController) GetState(service string) (ServiceState, error) {
	out, err := exec.Command(c.CmdName(), "is-active", c.ServiceName(service)).Output()
	if err != nil {
		return Error, err
	}
	state, err := parseSystemdInspectResult(out)
	if err != nil {
		err = fmt.Errorf("%v for service '%s', raw output: %s", err, service, string(out))
	}
	return state, err
}

// TailLogs executes command to start tailing service logs and returns string chan to receive log strings
// closing the chan will terminate tailing
func (c SystemdController) TailLogs(service string) (chan string, *os.Process, error) {
	var cmd *exec.Cmd
	if len(service) == 0 {
		cmd = exec.Command("sudo", "tail", "-f", "/var/log/syslog")
	} else {

		cmd = exec.Command("sudo", "journalctl", "-fu", c.ServiceName(service))
	}
	return StartCmdWithStderrStdoutTailer(cmd)
}

func parseSystemdInspectResult(out []byte) (ServiceState, error) {
	res := Error
	if len(out) == 0 {
		return res, fmt.Errorf("Empty returned status")
	}
	stateName := strings.ToLower(strings.TrimSpace(string(out)))
	if returnedState, ok := systemdStates[stateName]; ok {
		res = returnedState
	}
	return res, nil
}

var systemdStates = map[string]ServiceState{
	"inactive":     Inactive,
	"activating":   Activating,
	"active":       Active,
	"paused":       Inactive,
	"deactivating": Deactivating,
	"exited":       Inactive,
	"failed":       Failed,
	"unknown":      Unknown,
}

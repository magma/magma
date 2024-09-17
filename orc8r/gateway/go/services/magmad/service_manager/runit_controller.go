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

// RunitController - runit based controller implementation
type RunitController struct{}

// Name returns runit controller type name
func (RunitController) Name() string {
	return "runit"
}

func (RunitController) CmdName() string {
	return "sv"
}

func (RunitController) ServiceName(service string) string {
	return service
}

// Start starts service and returns error if unsuccessful
func (c RunitController) Start(service string) error {
	return exec.Command(c.CmdName(), "start", service).Run()
}

// Stop stops service and returns error if unsuccessful
func (c RunitController) Stop(service string) error {
	return exec.Command(c.CmdName(), "stop", c.ServiceName(service)).Run()
}

// Restart restarts service and returns error if unsuccessful
func (c RunitController) Restart(service string) error {
	return exec.Command(c.CmdName(), "reload-or-restart", c.ServiceName(service)).Run()
}

// GetState returns the given service state or error if unsuccessful
func (c RunitController) GetState(service string) (ServiceState, error) {
	out, err := exec.Command(c.CmdName(), "status", c.ServiceName(service)).Output()
	if err != nil {
		return Error, err
	}
	state, err := parseRunitStatusResult(out)
	if err != nil {
		err = fmt.Errorf("%v for service '%s', raw output: %s", err, service, string(out))
	}
	return state, err
}

// TailLogs executes command to start tailing service logs and returns string chan to receive log strings
// closing the chan will terminate tailing
func (c RunitController) TailLogs(service string) (chan string, *os.Process, error) {
	var cmd *exec.Cmd
	if len(service) == 0 {
		cmd = exec.Command("logread", "-f")
	} else {
		cmd = exec.Command("sh", "-c", "logread -f | grep "+service)
	}
	return StartCmdWithStderrStdoutTailer(cmd)
}

func parseRunitStatusResult(out []byte) (ServiceState, error) {
	res := Error
	if len(out) == 0 {
		return res, fmt.Errorf("Empty returned status")
	}
	statuses := strings.ToLower(strings.TrimSpace(string(out)))
	status := strings.TrimSpace(strings.Split(statuses, ":")[0])
	if returnedState, ok := runitStates[status]; ok {
		res = returnedState
	}
	return res, nil
}

var runitStates = map[string]ServiceState{
	"run":  Active,
	"down": Inactive,
	"fail": Failed,
}

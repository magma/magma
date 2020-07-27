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

// package generic_command implements magmad shell command execution functionality
package generic_command

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/shlex"

	"magma/gateway/config"
	"magma/orc8r/lib/go/protos"
)

const (
	Timeout = time.Minute
)

// Execute runs the command specified by req and returns its results
func Execute(ctx context.Context, req *protos.GenericCommandParams) (*protos.GenericCommandResponse, error) {
	var err error
	commands := config.GetMagmadConfigs().GenericCommandConfig.CommandsMap
	cmd, ok := commands[strings.ToLower(req.GetCommand())]
	res := &protos.GenericCommandResponse{}
	if !ok {
		return res, fmt.Errorf("unregistered command: %s", req.GetCommand())
	}
	cmdParams := []interface{}{}
	if cmd.AllowParams {
		if fields := req.GetParams().GetFields(); fields != nil {
			if param, ok := fields["shell_params"]; ok && param != nil {
				cmdParams = addValue(cmdParams, param)
			}
		}
	}
	formatters := strings.Count(cmd.CommandFmt, "%") - strings.Count(cmd.CommandFmt, "%%")*2
	cmdStr := cmd.CommandFmt
	if formatters > 0 {
		if len(cmdParams) < formatters {
			fillIn := make([]interface{}, formatters-len(cmdParams))
			for i, _ := range fillIn {
				fillIn[i] = ""
			}
			cmdParams = append(cmdParams, fillIn...)
		}
		cmdStr = fmt.Sprintf(cmd.CommandFmt, cmdParams[:formatters]...)
	} else {
		for _, p := range cmdParams {
			cmdStr += fmt.Sprintf(" %v", p)
		}
	}
	cmdList, splitErr := shlex.Split(cmdStr)
	if splitErr != nil || len(cmdList) == 0 {
		glog.Errorf("invalid command format: %v", splitErr)
		cmdList = []string{"sh", "-c", cmdStr}
	}
	execCtx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	exeCmd := exec.CommandContext(execCtx, cmdList[0], cmdList[1:]...)
	glog.Infof("executing command '%s'", exeCmd.String())
	var errBuff, outBuff bytes.Buffer
	exeCmd.Stderr = &errBuff
	exeCmd.Stdout = &outBuff
	err = exeCmd.Start()
	if err == nil {
		exeCmd.Wait()
		res = &protos.GenericCommandResponse{
			Response: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"result": {Kind: &structpb.Value_NumberValue{NumberValue: float64(exeCmd.ProcessState.ExitCode())}},
					"stdout": {Kind: &structpb.Value_StringValue{StringValue: outBuff.String()}},
					"stderr": {Kind: &structpb.Value_StringValue{StringValue: errBuff.String()}},
				},
			},
		}
	}
	return res, err
}

func addValue(params []interface{}, val *structpb.Value) []interface{} {
	if val != nil {
		switch val.GetKind().(type) {
		case *structpb.Value_StringValue:
			params = append(params, val.GetStringValue())
		case *structpb.Value_BoolValue:
			params = append(params, val.GetBoolValue())
		case *structpb.Value_NumberValue:
			params = append(params, val.GetNumberValue())
		case *structpb.Value_ListValue:
			for _, v := range val.GetListValue().GetValues() {
				params = addValue(params, v)
			}
		case *structpb.Value_StructValue:
			for name, v := range val.GetStructValue().GetFields() {
				params = addValue(append(params, name), v)
			}
		}
	}
	return params
}

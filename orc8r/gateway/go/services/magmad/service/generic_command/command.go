/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package generic_command implements magmad shell command execution functionality
package generic_command

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

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
	if ok && cmd != nil {
		cmdParams := []interface{}{}
		if cmd.AllowParams {
			if fields := req.GetParams().GetFields(); fields != nil {
				if param, ok := fields["shell_params"]; ok && param != nil {
					cmdParams = addValue(cmdParams, param)
				}
			}
		}
		cmdStr := fmt.Sprintf(cmd.CommandFmt, cmdParams)
		log.Printf("executing command '%s'", cmdStr)
		cmdList, splitErr := shlex.Split(cmdStr)
		if splitErr != nil || len(cmdList) == 0 {
			log.Printf("invalid command format: %v", splitErr)
			cmdList = []string{"sh", "-c", cmdStr}
		}
		execCtx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()
		exeCmd := exec.CommandContext(execCtx, cmdList[0], cmdList[1:]...)
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

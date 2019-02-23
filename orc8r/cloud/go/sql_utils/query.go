/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_utils

import (
	"fmt"
	"strings"
)

// GetPlaceholderArgList returns a string "($1, $2, ..., $numArgs)" for use
// in a SELECT ... IN query or INSERT query.
func GetPlaceholderArgList(startIdx int, numArgs int) string {
	retBuilder := strings.Builder{}
	retBuilder.WriteString("(")

	endIdx := startIdx + numArgs
	for i := startIdx; i < endIdx; i++ {
		retBuilder.WriteString(fmt.Sprintf("$%d", i))
		if i < endIdx-1 {
			retBuilder.WriteString(", ")
		}
	}
	retBuilder.WriteString(")")
	return retBuilder.String()
}

// GetUpdateClauseString returns a string "args[0] = $1, args[1] = $2, ..."
// for use in a UPDATE ... SET query.
func GetUpdateClauseString(startIdx int, argNames ...string) string {
	argsToJoin := make([]string, 0, len(argNames))
	for i, arg := range argNames {
		argsToJoin = append(argsToJoin, fmt.Sprintf("%s = $%d", arg, startIdx+i))
	}
	return strings.Join(argsToJoin, ", ")
}

// GetInsertArgListString returns a string "(args[0], args[1], ...)" for use
// in an INSERT query.
func GetInsertArgListString(args ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(args, ", "))
}

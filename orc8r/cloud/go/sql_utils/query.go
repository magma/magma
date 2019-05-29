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

const (
	postgresDialect = "psql"
	mysqlDialect    = "mysql"
)

// GetPlaceholderArgList returns a string
// "(${startIdx}, ${startIdx+1}, ..., ${startIdx+numArgs-1})"
func GetPlaceholderArgList(startIdx int, numArgs int) string {
	return GetPlaceholderArgListWithSuffix(startIdx, numArgs, "")
}

// GetPlaceholderArgListWithSuffix returns a string
// "(${startIdx}, ${startIdx+1}, ..., ${startIdx+numArgs-1}, {suffix})"
//
// The suffix argument is typically used for a field that's being updated
// in-place in an UPDATE query.
func GetPlaceholderArgListWithSuffix(startIdx int, numArgs int, suffix string) string {
	if numArgs == 0 {
		return fmt.Sprintf("(%s)", suffix)
	}

	retBuilder := strings.Builder{}
	retBuilder.WriteString("(")

	endIdx := startIdx + numArgs
	for i := startIdx; i < endIdx; i++ {
		retBuilder.WriteString(fmt.Sprintf("$%d", i))
		if i < endIdx-1 {
			retBuilder.WriteString(", ")
		}
	}

	if suffix != "" {
		retBuilder.WriteString(", ")
		retBuilder.WriteString(suffix)
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

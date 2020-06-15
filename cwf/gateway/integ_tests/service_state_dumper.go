// +build all qos
/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import "fmt"

func dumpPipelinedState(tr *TestRunner) {
	fmt.Println("******************* Dumping Pipelined State *******************")
	cmdList := [][]string{
		{"ovs-vsctl", "show"},
		{"ovs-ofctl", "-O", "OpenFlow15", "meter-stats", "cwag_br0"},
		{"ovs-ofctl", "-O", "OpenFlow15", "dump-meters", "cwag_br0"},
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

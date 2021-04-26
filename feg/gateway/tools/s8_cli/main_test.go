// +build cli_test

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestS8Cli(t *testing.T) {
	// replace the arguments
	os.Args = []string{"s8_cli", "cs", "-test", "-delete", "0"}
	fmt.Printf("\ncommand to run: %+v\n", os.Args)
	main()
}

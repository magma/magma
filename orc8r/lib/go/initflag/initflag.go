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
// package initflag initializes (parses) Go flag if needed, it allows the noise free use of golog & other packages
// relying on flags being parsed
package initflag

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	flag.Var(
		&syslogDest,
		syslogFlag,
		"Redirect stderr to syslog, optional syslog destination in network::address format (system default otherwise)",
	)

	if shouldParse() {
		// Save original settings
		orgUsage := flag.CommandLine.Usage
		origOut := flag.CommandLine.Output()
		origErrorHandling := flag.CommandLine.ErrorHandling()

		// Set to 'silent'
		flag.CommandLine.Init(flag.CommandLine.Name(), flag.ContinueOnError)
		flag.CommandLine.Usage = func() {}
		flag.CommandLine.SetOutput(ioutil.Discard)
		flag.Parse()

		// Restore original settings
		flag.CommandLine.Init(flag.CommandLine.Name(), origErrorHandling)
		flag.CommandLine.Usage = orgUsage
		flag.CommandLine.SetOutput(origOut)
	}
	// Check if the process needs to redirect stderr to syslog
	if f := flag.Lookup(syslogFlag); f != nil {
		if fTarget, ok := f.Value.(*syslogTarget); ok && fTarget.IsSet() {
			// Cannot use glog here, it should not be initialized yet
			fmt.Fprintf(os.Stderr, "INFO redirecting stderr to syslog\n")
			if err := redirectToSyslog(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR redirecting to syslog: %v\n", err)
			}
		}
	}
	// Check if the process needs to redirect stdout to stderr
	if *stdoutToStderr {
		stdout, os.Stdout = os.Stdout, os.Stderr
	}
}

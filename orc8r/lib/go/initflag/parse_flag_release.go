// +build !debug_build

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

package initflag

import (
	"flag"
	"os"
	"strings"
)

// shouldParse returns true if initflag should parse flags.
// this hack works around the fact that initflags breaks test tool outputs
// to completely disable this logic & preliminary flag parsing, use 'debug_build' Go build tag
func shouldParse() bool {
	isTest := strings.HasSuffix(os.Args[0], ".test") ||
		strings.HasSuffix(os.Args[0], "_test.go") ||
		strings.HasSuffix(os.Args[0], "_test_go") ||
		isInArgs("-test.v")
	return !flag.Parsed() && !isTest
}

// isInArgs returns true if any of the argunents passed to the command maches
// with the passed match stirng
func isInArgs(match string) bool {
	for _, arg := range os.Args {
		if arg == match {
			return true
		}
	}
	return false
}

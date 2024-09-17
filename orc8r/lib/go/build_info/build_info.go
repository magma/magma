/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// package build_info provides API to get information on the binary build/commit version
package build_info

import (
	"fmt"
)

const unknown = "UNDEFINED"

// the following variables (buildBranch, buildTag, buildCommitHash, buildCommitDate & buildDate)
// should be set by an external Makefile/build script/etc. via Go ldflags. For example:
//
//	export MAGMA_BUILD_COMMIT_HASH=$(git rev-parse HEAD)
//	export LD_FLAGS="-ldflags=-X 'magma/orc8r/lib/go/build_info.buildCommitHash=$MAGMA_BUILD_COMMIT_HASH'"
//	go build -o ./magmad "$LD_FLAGS" magma/gateway/services/magmad
var (
	buildBranch     = unknown
	buildTag        = unknown
	buildCommitHash = unknown
	buildCommitDate = unknown
	buildDate       = unknown
)

// Branch returns git branch of the current build
func Branch() string {
	return buildBranch
}

// Tag returns git tag of the current build
func Tag() string {
	return buildTag
}

// Commit returns git commit hash of the current build
func Commit() string {
	return buildCommitHash
}

// CommitDate returns git commit date of the current build
func CommitDate() string {
	return buildCommitDate
}

// BuildDate returns date & time of the current build
func BuildDate() string {
	return buildDate
}

// String returns formatted string of the current build's build & git commit information
func String() string {
	return fmt.Sprintf(
		"\nBuild Info:\n-----------\n"+
			"\tCommit Branch: %s\n"+
			"\tCommit Tag:    %s\n"+
			"\tCommit Hash:   %s\n"+
			"\tCommit Date:   %s\n"+
			"\tBuild  Date:   %s\n",
		buildBranch, buildTag, buildCommitHash, buildCommitDate, buildDate)
}

#!/usr/bin/env bash

################################################################################
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

set -euo pipefail

###############################################################################
# SCRIPT SECTION
###############################################################################

# This is a wrapper script for the bazel buildifier. Further information can be
# found at https://github.com/bazelbuild/buildtools/tree/master/buildifier#readme

BUILDIFIER_URL="https://github.com/bazelbuild/buildtools/releases/download/5.1.0/buildifier-linux-amd64"
BUILDIFIER_PATH="/tmp"
BUILDIFIER_EXECUTABLE="$BUILDIFIER_PATH/buildifier-linux-amd64"

# Command line argument, can be either 'check' or 'format'.
BUILDIFIER_MODE="${1:-}"

case "$BUILDIFIER_MODE" in
    "check")
        BUILDIFIER_OPTIONS=("-mode=diff" "--lint=warn" "--warnings=-unused-variable" "-v=true")
        ;;
    "format")
        BUILDIFIER_OPTIONS=("-mode=fix" "--lint=fix" "--warnings=all")
        ;;
    *)
        echo "Invalid argument '$BUILDIFIER_MODE'."
        echo "Valid arguments are 'check' or 'format'."
        exit 1
        ;;
esac

# Check if the Buildifier executable is already present on the system.
if [[ -s "$BUILDIFIER_EXECUTABLE" ]];
then
    echo "Buildifier executable is already present, skipping download."
else
    echo "Downloading pre-built buildifier ..."
    wget --quiet --directory-prefix "$BUILDIFIER_PATH" "$BUILDIFIER_URL"
    chmod +x "$BUILDIFIER_EXECUTABLE"
    echo "Download successful."
fi

WORKING_DIR="./"
echo "Info: Only the subfolders of the current directory are checked or formatted."

echo "Running bazel buildifier with the following command:"
set -x
# The '-r' option is used to find starlark files recursively in the WORKING_DIR.
"$BUILDIFIER_EXECUTABLE" "${BUILDIFIER_OPTIONS[@]}" -r "$WORKING_DIR"

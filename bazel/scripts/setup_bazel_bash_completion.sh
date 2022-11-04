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

########################################################
# Configures Bazel Bash-completion for the invoking user
########################################################

set -euo pipefail

function log() {
  echo "$@" >&2
}

function generate_completion() {
  # outputs the completion script to stdout
  #
  # inspired by https://github.com/bazelbuild/bazel/blob/5f3f59ba367158b6c2e811f68cc946537b9b74d6/scripts/BUILD#L6-L26
  # which is referenced in https://bazel.build/install/completion
  curl https://raw.githubusercontent.com/bazelbuild/bazel/release-5.3.0/scripts/bazel-complete-header.bash
  curl https://raw.githubusercontent.com/bazelbuild/bazel/release-5.3.0/scripts/bazel-complete-template.bash
  bazel help completion
}

if ! [ -f ~/.bash_completion.d/bazel-complete.bash ]
then
  log "Creating completion script"
  mkdir -p ~/.bash_completion.d
  generate_completion > ~/.bash_completion.d/bazel-complete.bash
else
  log "Completion script already exists"
fi

if ! grep --quiet bazel-complete.bash ~/.bashrc
then
  log "Adapting ~/.bashrc to source completion script"
  echo "source ~/.bash_completion.d/bazel-complete.bash" > ~/.bashrc
else
  log ".bashrc already sources the completion script"
fi

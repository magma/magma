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
# FUNCTION DECLARATIONS
###############################################################################

get_python_scripts() {
    echo "Collecting script targets..."
    mapfile -t PYTHON_SCRIPTS < <(bazel query "attr(tags, 'util_script', kind(.*_binary, \
        //orc8r/gateway/python/scripts/... union \
        //lte/gateway/python/scripts/... ))")
}

format_targets_to_paths() {
    for INDEX in "${!PYTHON_SCRIPTS[@]}"
    do
        # Strip leading '//'
        PYTHON_SCRIPTS[INDEX]="${PYTHON_SCRIPTS[INDEX]/\/\//}"
        # Replace ':' with '/'
        PYTHON_SCRIPTS[INDEX]="${PYTHON_SCRIPTS[INDEX]/://}"
    done
}

create_links() {
    echo "Linking bazel-built script executables to '/usr/local/bin/'..."
    for PYTHON_SCRIPT in "${PYTHON_SCRIPTS[@]}"
    do
        sudo ln -sf "/home/vagrant/magma/bazel-bin/${PYTHON_SCRIPT}" "/usr/local/bin/$(basename "${PYTHON_SCRIPT}").py"
    done
    echo "Linking finished."
}

mock_virtualenv() {
    # The virtualenv is not needed with bazel. Until the switchover
    # to bazel is complete, this creates an empty file that can
    # be sourced, without failure, in the LTE integration tests.
    # See https://github.com/magma/magma/issues/13807
    mkdir -p /home/vagrant/build/python/bin/
    touch /home/vagrant/build/python/bin/activate
}

###############################################################################
# SCRIPT SECTION
###############################################################################

PYTHON_SCRIPTS=()

get_python_scripts
format_targets_to_paths
create_links
mock_virtualenv

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

# Bazel commands have to be executed inside the repository.
cd "${MAGMA_ROOT}"
printf "\n###############################################################################\n"
printf "Linking the Python scripts needed for the integration tests."
printf "\n###############################################################################\n"
"${MAGMA_ROOT}"/bazel/scripts/link_scripts_for_bazel_integ_tests.sh

printf "\n###############################################################################\n"
printf "Copying Bazel specific systemd files to '/etc/systemd/system/'."
printf "\n###############################################################################\n"
sudo cp "${MAGMA_ROOT}"/lte/gateway/deploy/roles/magma/files/systemd_bazel/* /etc/systemd/system/
sudo systemctl daemon-reload

printf "\n###############################################################################\n"
printf "Building bazelified Magma AGW services and util scripts."
printf "\n###############################################################################\n"
# shellcheck disable=SC2046
bazel build $(bazel query "attr(tags, 'service|util_script', kind(.*_binary, //orc8r/... union //lte/... union //feg/... except //lte/gateway/c/core:mme_oai))")

printf "\n###############################################################################\n"
printf "Restarting the Magma AGW services." 
printf "\n###############################################################################\n"
sudo service "magma@*" stop && sudo service sctpd stop && sudo service magma_dp@envoy stop
sudo service magma@magmad start

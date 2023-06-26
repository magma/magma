# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os

import pytest


@pytest.hookimpl()
def pytest_sessionstart(session):
    if os.geteuid() != 0:
        raise Exception(
            "\n\n" + \
            "################################################################################\n" * 3 + \
            "To execute tests tagged as 'manual' you need to use the relevant shell script in\n" + \
            "'$MAGMA_ROOT/bazel/scripts/'!\n" + \
            "For Python sudo tests: '$MAGMA_ROOT/bazel/scripts/run_sudo_tests.sh'!\n" + \
            "For LTE integration tests: '$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh'!\n" + \
            "The scripts can be used to execute individual tests. To have the full overview\n" + \
            "of the available options add '--help' to them.\n" + \
            "Note: You got this error, because the test is tagged as manual and the user id\n" + \
            "is not zero, which indicates that you are not running the test as root.\n" + \
            "Disclaimer: Do not execute bazel commands as root, because it can destroy your\n" + \
            "local bazel setup, you have to use the above mentioned scripts." + \
            "\n################################################################################" * 3 + \
            "\n\n",
        )

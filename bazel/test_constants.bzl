# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
This file holds constants that are used for bazel test definitions. Especially
for test-tagging, this is meant to ensure that tags are used consistently.

Constants can be imported in a build file via
load("//bazel:test_constants.bzl", "CONSTANT1", "CONSTANT2")

Tags are defined as lists, i.e., you can add them in your build file via
list concatenation:
    load("//bazel:test_constants.bzl", "TAG_FOO", "TAG_BAR")
    ...
    tags = TAG_FOO + TAG_BAR
"""

# Used for a test that should be executed manually (read explicitly). This
# is a bazel keyword and causes "*" or "all" to not expand the
# respective target.
# See https://docs.bazel.build/versions/main/be/common-definitions.html
TAG_MANUAL = ["manual"]

# Used for a "sudo test" that needs to be executed by a user with sudo
# privileges. This is a restriction mainly known for some Python tests.
# To run sudo tests execute: $MAGMA_ROOT/bazel/scripts/run_sudo_tests.sh
# Note: for now a sudo test is also tagged as "manual".
TAG_SUDO_TEST = ["sudo_test"] + TAG_MANUAL

# Used for integration tests. These tests need to be executed manually
# by a user with sudo privileges. These tags represent test categories,
# which are used to determine the appropriate environment for them.
# To run the LTE integration tests execute: $MAGMA_ROOT/bazel/scripts/run_integ_tests.sh
TAG_INTEGRATION_TEST = ["integration_test"]
TAG_PRECOMMIT_TEST = ["precommit_test"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_EXTENDED_TEST = ["extended_test"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_EXTENDED_TEST_SETUP = ["extended_setup"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_EXTENDED_TEST_TEARDOWN = ["extended_teardown"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_NON_SANITY_TEST = ["nonsanity_test"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_NON_SANITY_TEST_SETUP = ["nonsanity_setup"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_NON_SANITY_TEST_TEARDOWN = ["nonsanity_teardown"] + TAG_MANUAL + TAG_INTEGRATION_TEST
TAG_TRAFFIC_SERVER_TEST = ["traffic_server_test"]

# Used for load tests. These tests need to be executed by the
# '$MAGMA_ROOT/bazel/scripts/run_load_tests.sh' script to preserve
# the ordering of the test cases.
TAG_LOAD_TEST = ["load_test"] + TAG_MANUAL

# Tag for utility scripts that are used in the Magma VM.
TAG_UTIL_SCRIPT = ["util_script"]

# Tag for Magma services.
TAG_SERVICE = ["service"]

# Tag to easily exclude OAI MME from builds
TAG_MME_OAI = ["mme_oai"]

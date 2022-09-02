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
# Note: for now a sudo test is also tagged as "manual".
TAG_SUDO_TEST = ["sudo_test"] + TAG_MANUAL

TAG_PRECOMMIT_TEST = ["precommit_test"] + TAG_MANUAL
TAG_EXTENDED_TEST = ["extended_test"] + TAG_MANUAL
TAG_EXTENDED_TEST_SETUP = ["extended_setup"] + TAG_MANUAL
TAG_EXTENDED_TEST_TEARDOWN = ["extended_teardown"] + TAG_MANUAL
TAG_NON_SANITY_TEST = ["nonsanity_test"] + TAG_MANUAL
TAG_NON_SANITY_TEST_SETUP = ["nonsanity_setup"] + TAG_MANUAL
TAG_NON_SANITY_TEST_TEARDOWN = ["nonsanity_teardown"] + TAG_MANUAL
TAG_TRAFFIC_SERVER_TEST = ["traffic_server_test"]

# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Tests for runfiles.bzl.
See https://bazel.build/rules/testing for general bazel rule testing documentation.
"""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load("@rules_pkg//:providers.bzl", "PackageFilesInfo")
load("//bazel:runfiles.bzl", "expand_runfiles")

# test suite

def runfiles_test(name):
    _setup_empty_targets_returns_empty_providers_test()
    _setup_targets_are_correctly_expanded_test()

    native.test_suite(
        name = name,
        tests = [
            ":empty_targets_returns_empty_providers_test",
            ":targets_are_correctly_expanded_test",
        ],
    )

# setup for rule to be tested

def _setup_empty_targets_returns_empty_providers_test():
    expand_runfiles(
        name = "expand_empty_targets",
        tags = ["manual"],  # should only be build here
    )

    rule_empty_targets_returns_empty_providers_test(
        name = "empty_targets_returns_empty_providers_test",
        target_under_test = ":expand_empty_targets",
    )

def _setup_targets_are_correctly_expanded_test():
    # testing an actually magma target instead of an artificial one
    # mconfigs proto should be sufficiently stable
    expand_runfiles(
        name = "expand_targets",
        tags = ["manual"],  # should only be build here
        targets = ["//lte/protos:mconfigs_python_proto"],
    )

    rule_targets_are_correctly_expanded_test(
        name = "targets_are_correctly_expanded_test",
        target_under_test = ":expand_targets",
    )

# asserts

def _empty_targets_returns_empty_providers_test_impl(ctx):
    env = analysistest.begin(ctx)

    target_under_test = analysistest.target_under_test(env)

    asserts.equals(
        env,
        expected = {"mode": "0755"},
        actual = target_under_test[PackageFilesInfo].attributes,
    )
    asserts.equals(
        env,
        expected = {},
        actual = target_under_test[PackageFilesInfo].dest_src_map,
    )
    asserts.equals(
        env,
        expected = depset([]),
        actual = target_under_test[DefaultInfo].files,
    )

    return analysistest.end(env)

expected_mapping = (
    "{" +
    '"orc8r/protos/common_pb2.py": <generated file orc8r/protos/common_python_proto_pb/orc8r/protos/common_pb2.py>, ' +
    '"lte/protos/mconfig/mconfigs_pb2.py": <generated file lte/protos/mconfigs_python_proto_pb/lte/protos/mconfig/mconfigs_pb2.py>' +
    "}"
)

expected_depset = (
    "depset([" +
    "<generated file orc8r/protos/common_python_proto_pb/orc8r/protos/common_pb2.py>, " +
    "<generated file lte/protos/mconfigs_python_proto_pb/lte/protos/mconfig/mconfigs_pb2.py>" +
    "])"
)

def _targets_are_correctly_expanded_test_impl(ctx):
    env = analysistest.begin(ctx)

    target_under_test = analysistest.target_under_test(env)

    asserts.equals(
        env,
        expected = expected_mapping,
        actual = str(target_under_test[PackageFilesInfo].dest_src_map),
    )
    asserts.equals(
        env,
        expected = expected_depset,
        actual = str(target_under_test[DefaultInfo].files),
    )

    return analysistest.end(env)

# creating rules for asserts

rule_empty_targets_returns_empty_providers_test = analysistest.make(_empty_targets_returns_empty_providers_test_impl)

rule_targets_are_correctly_expanded_test = analysistest.make(_targets_are_correctly_expanded_test_impl)

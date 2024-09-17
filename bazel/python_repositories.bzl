# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" Load external python repositories (that cannot be imported via pip). """

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

def python_repositories(name = ""):
    new_git_repository(
        name = "aioh2_repo",
        build_file = "//bazel/external:aioh2.BUILD",
        commit = "8c1b5ab2399443087795fe52b71e43b652b1031f",
        shallow_since = "1548652954 +0800",
        remote = "https://github.com/URenko/aioh2.git",
    )
    new_git_repository(
        name = "ryu_repo",
        build_file = "//bazel/external:ryu.BUILD",
        commit = "c776e4cb68600b2ee0a4f38364f4a355502777f1",  # Corresponds to: tag = "v4.34"
        shallow_since = "1569926530 +0900",
        remote = "https://github.com/faucetsdn/ryu.git",
        patches = [
            "//lte/gateway/deploy/roles/magma/files/patches:ryu_ipfix_args.patch",
            "//lte/gateway/deploy/roles/magma/files/patches:0001-Set-unknown-dpid-ofctl-log-to-debug.patch",
            "//lte/gateway/deploy/roles/magma/files/patches:0002-QFI-value-set-in-Openflow-controller-using-RYU.patch",
            "//lte/gateway/deploy/roles/magma/files/patches:0003-QFI-value-set-in-Openflow-controller-using-RYU.patch",
        ],
        patch_args = ["-p1"],
    )
    new_git_repository(
        name = "aioeventlet_repo",
        remote = "https://github.com/magma/deb-python-aioeventlet.git",
        build_file = "//bazel/external:aioeventlet.BUILD",
        commit = "86130360db113430370ed6c64d42aee3b47cd619",
        shallow_since = "1656345625 +0200",
    )

    # TODO: This is not a nice solution, because it is not really hermetic.
    # bcc is installd via apt bcc-tools from a magma repository
    native.new_local_repository(
        name = "bcc_repo",
        build_file = "//bazel/external:bcc.BUILD",
        path = "/usr/lib/python3/dist-packages/",
    )

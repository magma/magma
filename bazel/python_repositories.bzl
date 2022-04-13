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

def python_repositories():
    new_git_repository(
        name = "aioh2_repo",
        build_file = "//bazel/external:aioh2.BUILD",
        commit = "8c1b5ab2399443087795fe52b71e43b652b1031f",
        shallow_since = "1548652954 +0800",
        remote = "https://github.com/URenko/aioh2.git",
    )

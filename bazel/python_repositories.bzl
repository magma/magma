# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""All external repositories used for Python dependencies"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//bazel:python_packages.bzl", "PYTHON_PACKAGE")

def python_repositories():

    http_archive(
        name = PYTHON_PACKAGE.name,
        urls = PYTHON_PACKAGE.urls,
        sha256 = PYTHON_PACKAGE.sha256,
        strip_prefix = PYTHON_PACKAGE.strip_prefix,
        patch_cmds = PYTHON_PACKAGE.patch_cmds,
        build_file = PYTHON_PACKAGE.build_file,
    )


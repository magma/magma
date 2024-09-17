# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Python Toolchain and PIP Dependencies"""

load("@rules_python//python:pip.bzl", "pip_parse")

def configure_python_dependencies(name = None):
    native.register_toolchains("//bazel:py_toolchain")

    pip_parse(
        name = "python_deps",
        extra_pip_args = ["--require-hashes"],
        python_interpreter = "python3",
        requirements_lock = "//bazel/external:requirements.txt",
        visibility = ["//visibility:public"],
    )

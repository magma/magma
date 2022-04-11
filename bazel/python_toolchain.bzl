# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Python toolchain configuration"""

load("@rules_python//python:defs.bzl", "py_runtime", "py_runtime_pair")

def configure_python_toolchain(name = None):
    py_runtime(
        name = "python3",
        interpreter_path = "/usr/bin/python3.8",
        python_version = "PY3",
        visibility = ["//visibility:public"],
    )

    py_runtime_pair(
        name = "py_runtime_pair",
        py2_runtime = None,
        py3_runtime = ":python3",
        visibility = ["//visibility:public"],
    )

    native.toolchain(
        name = "py_toolchain",
        toolchain = ":py_runtime_pair",
        toolchain_type = "@bazel_tools//tools/python:toolchain_type",
        visibility = ["//visibility:public"],
    )

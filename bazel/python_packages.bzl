# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Hermetic Python interpreter configuration"""

load("@rules_python//python:pip.bzl", "pip_parse")

PY_VERSION = "3.8.5"
PY_HASH = "e3003ed57db17e617acb382b0cade29a248c6026b1bd8aad1f976e9af66a83b0"

BUILD_DIR = "/tmp/bazel/external/python_{0}".format(PY_VERSION)

def _patch_cmds():
    return [
        "mkdir -p {0}".format(BUILD_DIR),
        "cp -r * {0}".format(BUILD_DIR),
        "cd {0} && ./configure --prefix={0}/bazel_install".format(BUILD_DIR),
        "cd {0} && make install".format(BUILD_DIR),
        "rm -rf * && mv {0}/* .".format(BUILD_DIR),
        "ln -s bazel_install/bin/python3 python_bin",
    ]

PYTHON_PACKAGE = struct(
    name = "python_interpreter",
    sha256 = PY_HASH,
    strip_prefix = "Python-{0}".format(PY_VERSION),
    urls = ["https://www.python.org/ftp/python/{0}/Python-{0}.tar.xz".format(PY_VERSION)],
    build_file =  "//bazel/external:python_interpreter.BUILD",
    patch_cmds = _patch_cmds(),
    )

def configure_python_interpreter(name=None):

    native.register_toolchains("//bazel:py_toolchain")

    pip_parse(
        name = "python_deps",
        extra_pip_args = ["--require-hashes"],
        python_interpreter_target = "@python_interpreter//:python_bin",
        requirements_lock = "//bazel/external:requirements.txt",
        visibility = ["//visibility:public"],
    )
    

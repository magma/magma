# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@rules_cc//cc:defs.bzl", "cc_library")

# Manually generate config.h from config.h.in
genrule(
    name = "generate_config_h",
    srcs = ["include/tins/config.h.in"],
    outs = ["include/tins/config.h"],
    # Substitute #cmakedefine with #define
    cmd = 'sed "s/\\#cmakedefine/\\#define/" "$<" > "$@"',
)

cc_library(
    name = "libtins",
    srcs = glob(
        [
            "src/*.cpp",
            "src/detail/*.cpp",
            "src/tcp_ip/*.cpp",
            "src/utils/*.cpp",
            "src/dot11/*.cpp",
        ],
    ) + [":generate_config_h"],
    hdrs = glob(
        ["include/tins/**"],
    ),
    includes = ["include"],
    linkopts = ["-lpcap"],
    visibility = ["//visibility:public"],
)

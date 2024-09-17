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
    outs = ["config.h"],
    cmd = "\n".join([
        "touch $@",
        "echo '#define HAVE_STRINGS_H 1' >> $@",
        "echo '#define HAVE_STRING_H 1' >> $@",
        "echo '#define HAVE_SYS_STAT_H 1' >> $@",
        "echo '#define HAVE_SYS_TYPES_H 1' >> $@",
        "echo '#define HAVE_TLS 1' >> $@",
        "echo '#define HAVE_UNISTD_H 1' >> $@",
    ]),
)

cc_library(
    name = "fluid_base",
    srcs = glob(
        [
            "fluid/*.cc",
            "fluid/base/*.cc",
        ],
    ),
    hdrs = glob(
        [
            "fluid/*.hh",
            "fluid/base/*.hh",
        ],
    ) + [
        ":generate_config_h",
    ],
    includes = [""],
    linkopts = [
        "-lssl",
        "-lcrypto",
        "-levent",
        "-levent_pthreads",
        "-levent_openssl",
    ],
    visibility = ["//visibility:public"],
)

# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_skylib//rules:native_binary.bzl", "native_binary")
load("@rules_cc//cc:defs.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "folly",
    srcs = ["usr/local/lib/libfolly.so"],
    linkopts = [
        "-ldl",
        "-levent",
        "-ldouble-conversion",
        "-lgflags",
    ],
)

cc_library(
    name = "libmnl",
    linkopts = ["-lmnl"],
)

cc_library(
    name = "libpcap",
    linkopts = ["-lpcap"],
)

cc_library(
    name = "libuuid",
    linkopts = ["-luuid"],
)

cc_library(
    name = "sctp",
    linkopts = ["-lsctp"],
)

cc_library(
    name = "czmq",
    linkopts = ["-lczmq"],
)

cc_library(
    name = "libconfig",
    linkopts = ["-lconfig"],
)

cc_library(
    name = "libfd",
    srcs = [
        "usr/local/lib/libfdcore.so",
        "usr/local/lib/libfdproto.so",
    ],
    linkopts = [
        "-lfdcore",
        "-lfdproto",
    ],
)

cc_library(
    name = "libnettle",
    srcs = ["usr/lib/libnettle.so"],
    linkopts = ["-lnettle"],
)

cc_library(
    name = "libglog",
    srcs = glob(
        ["usr/lib/*-linux-gnu/libglog.so.0"],
        allow_empty = False,
    ),
)

cc_library(
    name = "libgnutls",
    srcs = ["usr/lib/libgnutls.so"],
    linkopts = ["-lgnutls"],
)

# TODO(GH9710): Generate asn1c with Bazel
native_binary(
    name = "asn1c",
    src = "usr/local/bin/asn1c",
    out = "asn1c",
)

cc_library(
    name = "libsqlite3-dev",
    linkopts = ["-lsqlite3"],
)

cc_library(
    name = "libsystemd",
    linkopts = ["-lsystemd"],
)

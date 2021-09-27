# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cc_library(
    name = "sentry_vendor",
    srcs = [
        "vendor/acutest.h",
        "vendor/jsmn.h",
        "vendor/mpack.c",
        "vendor/mpack.h",
        "vendor/stb_sprintf.c",
        "vendor/stb_sprintf.h",
    ],
)

cc_library(
    name = "sentry_core",
    srcs = [
        "src/sentry_alloc.c",
        "src/sentry_backend.c",
        "src/sentry_core.c",
        "src/sentry_database.c",
        "src/sentry_envelope.c",
        "src/sentry_json.c",
        "src/sentry_logger.c",
        "src/sentry_options.c",
        "src/sentry_os.c",
        "src/sentry_random.c",
        "src/sentry_ratelimiter.c",
        "src/sentry_scope.c",
        "src/sentry_session.c",
        "src/sentry_slice.c",
        "src/sentry_string.c",
        "src/sentry_sync.c",
        "src/sentry_transport.c",
        "src/sentry_utils.c",
        "src/sentry_uuid.c",
        "src/sentry_value.c",
    ],
    hdrs = [
        "src/sentry_alloc.h",
        "src/sentry_backend.h",
        "src/sentry_boot.h",
        "src/sentry_core.h",
        "src/sentry_database.h",
        "src/sentry_envelope.h",
        "src/sentry_json.h",
        "src/sentry_logger.h",
        "src/sentry_options.h",
        "src/sentry_os.h",
        "src/sentry_path.h",
        "src/sentry_random.h",
        "src/sentry_ratelimiter.h",
        "src/sentry_scope.h",
        "src/sentry_session.h",
        "src/sentry_slice.h",
        "src/sentry_string.h",
        "src/sentry_symbolizer.h",
        "src/sentry_sync.h",
        "src/sentry_transport.h",
        "src/sentry_unix_pageallocator.h",
        "src/sentry_unix_spinlock.h",
        "src/sentry_utils.h",
        "src/sentry_uuid.h",
        "src/sentry_value.h",
    ],
    strip_include_prefix = "src",
    deps = [":sentry_vendor"],
)

cc_library(
    name = "sentry_unwind",
    srcs = [
        "src/unwinder/sentry_unwinder.c",
        "src/unwinder/sentry_unwinder_libbacktrace.c",
    ],
    deps = [":sentry_core"],
)

cc_library(
    name = "sentry_transport",
    srcs = [
        "src/transports/sentry_disk_transport.c",
        "src/transports/sentry_function_transport.c",
        # for using curl
        "src/transports/sentry_transport_curl.c",
    ],
    hdrs = [
        "src/transports/sentry_disk_transport.h",
    ],
    strip_include_prefix = "src",
    deps = [":sentry_core"],
)

cc_library(
    name = "sentry_modulefinder",
    srcs = ["src/modulefinder/sentry_modulefinder_linux.c"],
    hdrs = ["src/modulefinder/sentry_modulefinder_linux.h"],
    deps = [":sentry_core"],
)

# WIP
cc_library(
    name = "sentry_breakpad",
    srcs = ["src/backends/sentry_backend_breakpad.cpp"],
    deps = [
        ":sentry_core",
        ":sentry_transport",
        # we need a dependency on breakpad here :'(
        # maybe useful: https://github.com/google/gapid/blob/master/tools/build/third_party/breakpad/breakpad.BUILD
    ],
)

# WIP
cc_library(
    name = "sentry",
    hdrs = ["include/sentry.h"],
    copts = [
        "-ldl",
        "-lrt",
        "-lpthread",
    ],
    includes = ["include"],
    visibility = ["//visibility:public"],
    deps = [":sentry_breakpad"],
)

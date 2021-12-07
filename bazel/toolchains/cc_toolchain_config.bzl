# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_tools//tools/build_defs/cc:action_names.bzl", "ACTION_NAMES")
load(
    "@bazel_tools//tools/cpp:cc_toolchain_config_lib.bzl",
    "action_config",
    "feature",
    "flag_group",
    "flag_set",
    "tool",
    "tool_path",
)

all_link_actions = [
    ACTION_NAMES.cpp_link_executable,
    ACTION_NAMES.cpp_link_dynamic_library,
    ACTION_NAMES.cpp_link_nodeps_dynamic_library,
]

all_compiler_actions = [
    ACTION_NAMES.cc_flags_make_variable,
    ACTION_NAMES.c_compile,
    ACTION_NAMES.cpp_compile,
    ACTION_NAMES.cpp_header_parsing,
]

gcc_path = "/usr/bin/gcc"

action_configs = [
    action_config(
        action_name = ACTION_NAMES.cpp_compile,
        tools = [
            tool(
                path = gcc_path,
            ),
        ],
    ),
    action_config(
        action_name = ACTION_NAMES.c_compile,
        tools = [
            tool(
                path = gcc_path,
            ),
        ],
    ),
]

def _impl(ctx):
    tool_paths = [
        tool_path(
            name = "gcc",
            path = "/usr/bin/gcc-9",
        ),
        tool_path(
            name = "ld",
            path = "/usr/bin/ld",
        ),
        tool_path(
            name = "ar",
            path = "/usr/bin/ar",
        ),
        tool_path(
            name = "cpp",
            path = "/usr/bin/cpp",
        ),
        tool_path(
            name = "gcov",
            path = "/usr/bin/gcov",
        ),
        tool_path(
            name = "nm",
            path = "/usr/bin/nm",
        ),
        tool_path(
            name = "objdump",
            path = "/usr/bin/objdump",
        ),
        tool_path(
            name = "strip",
            path = "/usr/bin/strip",
        ),
        tool_path(
            name = "cmake",
            path = "/usr/bin/cmake",
        ),
    ]
    features = [
        feature(
            name = "default_linker_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = all_link_actions,
                    flag_groups = ([
                        flag_group(
                            flags = [
                                "-lstdc++",
                                "-lm",
                            ],
                        ),
                    ]),
                ),
            ],
        ),
        feature(
            name = "default_compiler_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = all_link_actions,
                    flag_groups = ([
                        flag_group(
                            flags = [
                                "-std=c++14",
                            ],
                        ),
                    ]),
                ),
            ],
        ),
    ]
    return cc_common.create_cc_toolchain_config_info(
        ctx = ctx,
        features = features,
        cxx_builtin_include_directories = [
            "/usr/lib/gcc/x86_64-linux-gnu/9/include",
            "/usr/include",
            "/usr/local/include",
        ],
        toolchain_identifier = "local",
        host_system_name = "local",
        target_system_name = "local",
        target_cpu = "linux-x86_64",
        target_libc = "unknown",
        compiler = "gcc",
        abi_version = "unknown",
        abi_libc_version = "unknown",
        tool_paths = tool_paths,
    )

cc_toolchain_config = rule(
    attrs = {},
    provides = [CcToolchainConfigInfo],
    implementation = _impl,
)

# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Inspired by https://github.com/aspect-build/rules_container/blob/main/language/runfiles.bzl by @thesayyn
# Contributed to https://github.com/aspect-build/rules_container under Apache-2.0

"""
This file represents a workaround for the current state of https://github.com/bazelbuild/rules_pkg.
Dependencies are not packaged (even internal dependencies). The rule here expands all
internal and external dependencies into a PackageFilesInfo so that the files can be used in
package rules.

The dependencies are not put into service specific paths. This means, that, e.g., pip dependencies and
proto files that are used by multiple services are only added once in the packaging process.

This rule is currently only applied to python dependencies (go and c/c++ dependencies are linked statically
into the binaries).

Additionally this rule
* renames the relative path of files so that they can be found correctly in the target system (usually
  packaged into the "dist-packages" folder of the used python interpreter)
* excludes files that are not needed during runtime
"""

load("@rules_pkg//:providers.bzl", "PackageFilesInfo")

STRIP_PATHS = [
    "lte/gateway/python/",
    "orc8r/gateway/python/",
    "cwf/swagger/specs_root/",
    "feg/swagger/specs_root/",
    "lte/swagger/specs_root/",
    "orc8r/swagger/specs_root/",
]

# beware: order matters here, e.g., "lte/protos/oai/" needs to be before "lte/protos/"
STRIP_PATHS_PROTOS = [
    "cwf/protos/",
    "dp/protos/",
    "feg/protos/",
    "lte/protos/oai/",
    "lte/protos/",
    "orc8r/protos/",
    "orc8r/swagger/magmad_events_v1",
]

EXCLUDES = [
    # external protobuf is only needed during compile time
    "../com_google_protobuf",
    # external grpc is only needed during compile time
    "../com_github_grpc_grpc",
    # bazel compiled grpc library
    "_solib_k8/libexternal_Scom_Ugithub_Ugrpc_Ugrpc_Slibgrpc.so",
]

def _is_excluded(file):
    for exclude in EXCLUDES:
        if file.short_path.startswith(exclude):
            return True
    return False

def _runfile_path(file):
    path = file.short_path
    if path.startswith("../"):
        return _strip_external(path)
    return _strip_internal(path, file)

def _strip_external(path):
    path_clean = path.replace("../", "")

    # removes the first folder
    path_wo_first_folder = path_clean.partition("/")[2]

    # special case: grpc is packaged in subfolders (stripped here)
    if path_wo_first_folder.startswith("src/python/grpcio/"):
        return path_wo_first_folder.replace("src/python/grpcio/", "")

    return path_wo_first_folder

def _strip_internal(path, file):
    for prefix in STRIP_PATHS:
        if path.startswith(prefix):
            # lte/gateway/python/magma/foo/bar.py -> magma/foo/bar.py
            return path.replace(prefix, "", 1)

    for prefix in STRIP_PATHS_PROTOS:
        if path.startswith(prefix):
            # lte/protos/target_name/lte/protos/foo_pb2.py -> lte/protos/foo_pb2.py
            return path.replace(prefix, "", 1).replace(file.owner.name + "/", "", 1)

    print("Unhandled path: " + path)  # buildifier: disable=print

    return "FIXME"  # needs to be handled

def _runfiles_impl(ctx):
    py_infos = [target[PyInfo] for target in ctx.attr.targets]
    def_infos = [target[DefaultInfo] for target in ctx.attr.targets]

    files = depset(transitive = [py_info.transitive_sources for py_info in py_infos] + [def_info.default_runfiles.files for def_info in def_infos])
    file_map = {}
    mapped_files = []

    for file in files.to_list():
        if not _is_excluded(file):
            file_map[_runfile_path(file)] = file
            mapped_files = mapped_files + [file]

    files = depset(transitive = [files])

    return [
        PackageFilesInfo(
            dest_src_map = file_map,
            attributes = {"mode": "0755"},
        ),
        DefaultInfo(files = depset(mapped_files)),
    ]

expand_runfiles = rule(
    implementation = _runfiles_impl,
    attrs = {
        "targets": attr.label_list(providers = [PyInfo]),
    },
)

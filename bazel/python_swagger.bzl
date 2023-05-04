# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This file contains parts derived from [1] and [2].
# [1] Licensed under Apache License 2.0
#     Copyright 2019 OpenAPI-Generator-Bazel Contributors
#     https://github.com/OpenAPITools/openapi-generator-bazel/blob/fb7e302de4597277bea12757836f2ce988c805ee/internal/openapi_generator.bzl
# [2] Licensed under MIT License
#     Copyright (c) 2021 Meetup, Inc.
#     https://github.com/meetup/rules_openapi/blob/86fa11d0a8a8188ceecb74b5674af3f7363701e8/openapi/openapi.bzl
# We thank the authors of [1] and [2].

""" Provides support for building python code from swagger specifications. """

load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

SWAGGER_CLI_VERSION = "2.4.16"

# How to update SWAGGER_CLI_SHA256 in case of version upgrade:
#  - Run bazel build without the artifact_sha256 argument.
#  - Find line in output
#    "INFO: SHA256 (https://repo.maven.apache.org/maven2/io/swagger/swagger-codegen-cli/2.4.16/swagger-codegen-cli-2.4.16.jar) = 154b5a37254a3021a8cb669a1c57af78b45bb97e89e0425e3f055b1c79f74a93".
#  - Update SWAGGER_CLI_SHA256 accordingly.
#  - Alternatively download the .jar file and calculate the SHA256 by running sha256sum swagger-codegen-cli-2.4.16.jar.
SWAGGER_CLI_SHA256 = "154b5a37254a3021a8cb669a1c57af78b45bb97e89e0425e3f055b1c79f74a93"

# Loads the swagger codegen_cli dependency - to be called in WORKSPACE.bazel.
def load_swagger_repositories(name = "load_swagger_repositories"):
    jvm_maven_import_external(
        name = "maven_swagger_codegen_cli",
        artifact = "io.swagger:swagger-codegen-cli:" + SWAGGER_CLI_VERSION,
        artifact_sha256 = SWAGGER_CLI_SHA256,
        server_urls = ["https://repo.maven.apache.org/maven2"],
        licenses = ["notice"],  # Apache 2.0 License
    )

def _generator_command(ctx, gen_dir):
    java_path = ctx.attr._jdk[java_common.JavaRuntimeInfo].java_executable_exec_path

    gen_cmd = str(java_path)
    gen_cmd += " -cp {cli_jar}".format(cli_jar = ctx.file._swagger_cli.path)
    gen_cmd += " io.swagger.codegen.SwaggerCodegen generate -i {spec} -l python -o {output}".format(
        spec = ctx.file.spec.path,
        output = gen_dir,
    )
    gen_cmd += " -Dmodels"

    return gen_cmd

def _py_swagger_impl(ctx):
    declared_dir = ctx.actions.declare_directory(ctx.attr.name)

    # The prefix of the .yml path determines the component, e.g. orc8r, lte, feg, ...
    component = ctx.file.spec.path.split("/")[0]
    gen_dir = "{}/{}/swagger".format(declared_dir.path, component)

    commands = [
        "mkdir -p {}".format(gen_dir),
        _generator_command(ctx, gen_dir),
        # Build required folder structure and remove unused build outputs.
        "rm -rf {}/test".format(gen_dir),
        "mv {}/swagger_client/models {}".format(gen_dir, gen_dir),
        "rm -rf {}/swagger_client".format(gen_dir),
    ]

    inputs = ctx.files._jdk + [
        ctx.file._swagger_cli,
        ctx.file.spec,
    ]

    ctx.actions.run_shell(
        inputs = inputs,
        outputs = [declared_dir],
        command = " && ".join(commands),
    )

    return [
        DefaultInfo(
            files = depset([declared_dir]),
            runfiles = ctx.runfiles([declared_dir]),
        ),
        # Target can be consumed as python dependency.
        PyInfo(
            transitive_sources = depset([declared_dir]),
            # __main__ is the relative root of the runfile folder from the calling target
            imports = depset(["__main__/{}/swagger/{}".format(component, ctx.attr.name)]),
        ),
    ]

py_swagger = rule(
    attrs = {
        "spec": attr.label(
            mandatory = True,
            allow_single_file = [".yml"],
        ),
        "_jdk": attr.label(
            default = Label("@bazel_tools//tools/jdk:current_java_runtime"),
            providers = [java_common.JavaRuntimeInfo],
        ),
        "_swagger_cli": attr.label(
            cfg = "exec",
            default = Label("@maven_swagger_codegen_cli//jar"),
            allow_single_file = True,
        ),
    },
    implementation = _py_swagger_impl,
)

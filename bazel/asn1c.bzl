# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Custom library and macros to generate files with asn1c"""

load("@rules_cc//cc:defs.bzl", "cc_library")

def _get_dir_name(name):
    # we need to postfix the directory name with .c to trick Bazel into thinking this is a valid input
    # Related GH issue: https://github.com/bazelbuild/bazel/issues/10552
    return name + ".c"

def _contruct_substitution_commands(ctx, dir_path):
    """Returns a list of commands that replaces certain strings with another. 

    The replacement configuration is taken from ctx.attr.substitutions.
    """
    substitution_commands = []
    substitution_template = "egrep -lRZ \'{before}\' {dir} | xargs --no-run-if-empty -0 -l sed -i -e \'s/{before}/{after}/g\'"
    for before in ctx.attr.substitutions:
        after = ctx.attr.substitutions[before]
        substitution_commands.append(
            substitution_template.format(
                dir = dir_path,
                before = before,
                after = after,
            ),
        )
    return substitution_commands

def _construct_file_filter_commands(ctx, dir_path):
    """Returns a list of commands that removes any files that fit the condition described below.

    1. any file that does not end with .c or .h
    2. any file that is specified in ctx.attr.files_to_remove
    """
    filter_commands = []

    # cc_library will complain if there are files that are not one of (.c, .cc, .cpp, .h, ...)
    filter_files_template = "ls -d {dir}/* | grep --invert-match --extended-regexp \'{choose}\' | xargs  --no-run-if-empty rm "
    choose = "**.[.][ch]$"
    filter_commands.append(filter_files_template.format(dir = dir_path, choose = choose))

    # remove any files specified in ctx.attr
    for file_name in ctx.attr.files_to_remove:
        filter_commands.append("rm {dir}/{file} || true".format(dir = dir_path, file = file_name))
    return filter_commands

def _construct_asn1c_commands(ctx, dir_path):
    """Returns a string of command that runs asn1c"""
    prefix = ctx.attr.prefix
    flags = ctx.attr.flags
    asn1_file = ctx.attr.asn1_file.files.to_list()[0].path
    asn1c_command_template = "{asn1c} {flags} -D {dir} {input}"
    return [
        asn1c_command_template.format(
            asn1c = ctx.executable._asn1c.path,
            prefix = prefix,
            flags = flags,
            dir = dir_path,
            input = asn1_file,
        ),
    ]

def _gen_with_asn1c_impl(ctx):
    """generate files by running asn1c

    Args:
        ctx: passed through Bazel
    """
    name = ctx.attr.name
    gen_tree = ctx.actions.declare_directory(_get_dir_name(name))
    gen_path = gen_tree.path

    # Run asn1c to generate files and apply some modifications
    commands = _construct_asn1c_commands(ctx, gen_path) + _construct_file_filter_commands(ctx, gen_path) + _contruct_substitution_commands(ctx, gen_path)
    ctx.actions.run_shell(
        inputs = ctx.attr.asn1_file.files.to_list() + [ctx.executable._asn1c],
        outputs = [gen_tree],
        command = " && ".join(commands),
        progress_message = "Generating ASN1 files into: {dir}".format(dir = gen_path),
        env = {
            "ASN1C_PREFIX": ctx.attr.prefix,
        },
    )

    return [
        DefaultInfo(
            files = depset([gen_tree]),
        ),
        # Include additional information so the cc_library can understand where to look for headers/includes
        CcInfo(
            compilation_context = cc_common.create_compilation_context(
                headers = depset([gen_tree]),
                system_includes = depset([gen_path]),
            ),
        ),
    ]

def _get_attrs():
    return {
        "asn1_file": attr.label(
            mandatory = True,
            allow_single_file = [".asn1"],
            doc = """The asn file asn1c should use to generate files""",
        ),
        "files_to_remove": attr.string_list(
            doc = """A list of file names to remove after code generation""",
        ),
        "flags": attr.string(
            doc = """Command line flags passed to asn1c""",
        ),
        "prefix": attr.string(
            mandatory = True,
            doc = """Value that is set to ASN1C_PREFIX""",
        ),
        "substitutions": attr.string_dict(
            doc = """A string map of any substitions to be made""",
        ),
        # This makes it so that the executable asn1c is treated as an input to this rule
        "_asn1c": attr.label(
            executable = True,
            default = Label("@system_libraries//:asn1c"),
            cfg = "host",
        ),
    }

gen_with_asn1c = rule(
    implementation = _gen_with_asn1c_impl,
    attrs = _get_attrs(),
    output_to_genfiles = True,
)

def cc_asn1_library(
        name,
        asn1_file,
        prefix):
    """Create a CC library of generated asn1 files.

    This library wraps up generated files from these 4 actions:
      1. generate .c/.h files by running asn1c with the given input: asn1_file
      2. remove all non .c/.h files
      3. remove converter-example.c
      4. apply string substitutions

    Args:
        name: the name of rule
        asn1_file: relative path to the .asn1 file that will be passed to asn1c
        prefix: value that is set to ASN1C_PREFIX
    """
    gen_name = name + "_genrule"

    flags = "-pdu=all -fcompound-names -fno-include-deps -gen-PER"

    # Taken from https://github.com/magma/magma/blob/14c1cf643a61d576b3d24642e17ed3911d19210d/lte/gateway/c/core/oai/tasks/s1ap/CMakeLists.txt#L35
    # The original PR (PR2707) doesn't give an explanation on why this is necessary.
    # I'm guessing it is to avoid the following GCC error: To avoid the following GCC warning: integer constant is so large that it is unsigned
    substitutions = {
        "18446744073709551615": "18446744073709551615u",
    }

    gen_with_asn1c(
        name = gen_name,
        asn1_file = asn1_file,
        flags = flags,
        prefix = prefix,
        substitutions = substitutions,
        # We want to remove converter-example.c as it does not compile
        files_to_remove = ["converter-example.c"],
    )

    cc_library(
        name = name,
        srcs = [gen_name],
        # This is needed so that the CCInfo (header/include info) can be used
        deps = [gen_name],
    )

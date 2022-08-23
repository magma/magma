#!/usr/bin/env python3

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script was modified from github.com/envoyproxy/envoy
# See https://github.com/envoyproxy/envoy/blob/dca672af515029d36ddfc378480e3041f4d3a9f3/tools/gen_compilation_database.py

# The original script does not work for our repository, so I've made some changes to the compile_commands combination logic.
# Other than that, I also removed logic that was specific to envoyproxy.

import argparse
import json
import os
import subprocess
from pathlib import Path


# This method is equivalent to https://github.com/grailbio/bazel-compilation-database/blob/master/generate.py
def generate_compilation_database(args):
    subprocess.check_call(
        ["bazel", "build"] + [
            "--aspects=@bazel_compdb//:aspects.bzl%compilation_database_aspect",
            "--output_groups=compdb_files,header_files",
        ] + args.bazel_targets,
    )

    execroot = subprocess.check_output(
        ["bazel", "info", "execution_root"],
    ).decode().strip()

    db_entries = []
    for db in Path(execroot).glob('**/*.compile_commands.json'):
        raw_commands = db.read_text().splitlines()
        for raw_command in raw_commands:
            db_entries.append(json.loads(raw_command.rstrip(',')))

    def replace_execroot_marker(db_entry):
        if 'directory' in db_entry and db_entry['directory'] == '__EXEC_ROOT__':
            db_entry['directory'] = execroot
        if 'command' in db_entry:
            db_entry['command'] = (
                db_entry['command'].replace('-isysroot __BAZEL_XCODE_SDKROOT__', '')
            )
        return db_entry

    return list(map(replace_execroot_marker, db_entries))


def is_header(filename):
    for ext in (".h", ".hh", ".hpp", ".hxx"):
        if filename.endswith(ext):
            return True
    return False


def modify_compile_command(target, args):
    cc, options = target["command"].split(" ", 1)

    # # Workaround for bazel added C++11 options, those doesn't affect build itself but
    # # clang-tidy will misinterpret them.
    options = options.replace("-std=c++0x ", "")
    options = options.replace("-std=c++14 ", "")
    # clang-tidy does not recognize this flag
    options = options.replace("-fno-canonical-system-headers ", "")

    # Visual Studio Code doesn't seem to like "-iquote". Replace it with
    # old-style "-I".
    options = options.replace("-iquote ", "-I ")

    if is_header(target["file"]):
        options += " -Wno-pragma-once-outside-header -Wno-unused-const-variable"
        options += " -Wno-unused-function"

    target["command"] = " ".join([cc, options])
    return target


def fix_compilation_database(args, db):
    db = [modify_compile_command(target, args) for target in db]
    out = args.output_dir + "/compile_commands.json"
    print("Generated compile db: " + out)
    with open(out, "w") as db_file:
        json.dump(db, db_file, indent=2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Generate JSON compilation database')
    parser.add_argument(
        'bazel_targets',
        nargs='*',
        default=[
            "//lte/gateway/c/...",
            "//orc8r/gateway/c/...",
        ],
    )
    parser.add_argument(
        '--output_dir',
        default=os.getenv('MAGMA_ROOT'),
        help="Directory where compile_commands.json should be generated. Default value is $MAGMA_ROOT",
        action='store',
    )
    args = parser.parse_args()
    fix_compilation_database(args, generate_compilation_database(args))

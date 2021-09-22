#!/usr/bin/env python3

import argparse
import json
import os
import shlex
import subprocess
from pathlib import Path

# This file was taken and modified from envoyproxy/envoy
# https://github.com/envoyproxy/envoy/blob/e385e016e941b6c2d91096b0b792aa6940d4b38c/tools/gen_compilation_database.py


# This method is equivalent to https://github.com/grailbio/bazel-compilation-database/blob/master/generate.sh
def generate_compilation_database(args):
    # We need to download all remote outputs for generated source code. This option lives here to override those
    # specified in bazelrc.
    bazel_options = shlex.split(os.environ.get("BAZEL_BUILD_OPTIONS", "")) + [
        "--config=vm",
        "--remote_download_outputs=all",
    ]

    subprocess.check_call(
        ["bazel", "build"] + bazel_options + [
            "--aspects=@bazel_compdb//:aspects.bzl%compilation_database_aspect",
            "--output_groups=compdb_files,header_files",
        ] + args.bazel_targets,
    )

    execroot = subprocess.check_output(
        ["bazel", "info", "execution_root"]
        + bazel_options,
    ).decode().strip()

    compdb = []
    for compdb_file in Path(execroot).glob("**/*.compile_commands.json"):
        compdb.extend(
            json.loads("[" + compdb_file.read_text().replace("__EXEC_ROOT__", execroot) + "]"),
        )
    return compdb


def is_header(filename):
    for ext in (".h", ".hh", ".hpp", ".hxx"):
        if filename.endswith(ext):
            return True
    return False


def is_compile_target(target, args):
    filename = target["file"]
    if not args.include_headers and is_header(filename):
        return False

    if not args.include_genfiles:
        if filename.startswith("bazel-out/"):
            return False

    if not args.include_external:
        if filename.startswith("external/"):
            return False

    return True


def modify_compile_command(target, args):
    cc, options = target["command"].split(" ", 1)

    # Workaround for bazel added C++11 options, those doesn't affect build itself but
    # clang-tidy will misinterpret them.
    options = options.replace("-std=c++0x ", "")
    options = options.replace("-std=c++11 ", "")

    if args.vscode:
        # Visual Studio Code doesn't seem to like "-iquote". Replace it with
        # old-style "-I".
        options = options.replace("-iquote ", "-I ")

    if is_header(target["file"]):
        options += " -Wno-pragma-once-outside-header -Wno-unused-const-variable"
        options += " -Wno-unused-function"
        if not target["file"].startswith("external/"):
            options = "-x c++ -std=c++14 -fexceptions " + options

    target["command"] = " ".join([cc, options])
    return target


def fix_compilation_database(args, db):
    db = [modify_compile_command(target, args) for target in db if is_compile_target(target, args)]

    with open("compile_commands.json", "w") as db_file:
        json.dump(db, db_file, indent=2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Generate JSON compilation database')
    parser.add_argument('--include_external', action='store_true')
    parser.add_argument('--include_genfiles', action='store_true')
    parser.add_argument('--include_headers', action='store_true')
    parser.add_argument('--vscode', action='store_true')
    parser.add_argument(
        'bazel_targets',
        nargs='*',
        default=[
            "//orc8r/gateway/c/...",
            "//orc8r/protos/...",
            "//lte/gateway/c/...",
            "//lte/protos/...",
        ],
    )
    args = parser.parse_args()
    fix_compilation_database(args, generate_compilation_database(args))

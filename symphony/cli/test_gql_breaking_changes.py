#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
import argparse
import os
import sys
from distutils.version import LooseVersion
from zipfile import ZipFile

from graphql import build_ast_schema
from graphql.language.parser import parse
from graphql.utilities import find_breaking_changes
from graphql_compiler.gql.utils_schema import compile_schema_library


ERR_MSG_FORMAT = "Diff breaks schema {} version with changes: {}"


def extract_zip(input_zip_filepath):
    with ZipFile(input_zip_filepath) as input_zip:
        return {name: input_zip.read(name) for name in input_zip.namelist()}


def get_minimal_supported_version(schema_versions_library):
    minimal_version_filepath = os.path.join(
        schema_versions_library, "minimal_supported_version"
    )
    with open(minimal_version_filepath) as minimal_version_file:
        return minimal_version_file.read().strip()


def test_graphql_breaking_changes(schema_library: str) -> None:
    graphql_schema_versions = "graphql_schema_versions"
    zip_filepath = os.path.join(graphql_schema_versions, "old_schemas.zip")
    schema = compile_schema_library(schema_library)
    schemas = extract_zip(zip_filepath)
    minimal_version = get_minimal_supported_version(graphql_schema_versions)
    for ver, schema_str in schemas.items():
        if LooseVersion(ver) < LooseVersion(minimal_version):
            continue
        schema_str = schema_str.decode("utf-8")
        old_schema = build_ast_schema(parse(schema_str))
        breaking_changes = find_breaking_changes(old_schema, schema)
        assert len(breaking_changes) == 0, ERR_MSG_FORMAT.format(ver, breaking_changes)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_library_path", help="email to connect to inventory with", type=str
    )
    args: argparse.Namespace = parser.parse_args()

    schema_library_path: str = args.schema_library_path

    test_graphql_breaking_changes(schema_library_path)

    sys.exit(0)

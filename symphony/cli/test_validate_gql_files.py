#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
import argparse
import glob
import os
import sys

from graphql_compiler.gql.cli import process_file
from graphql_compiler.gql.query_parser import QueryParser
from graphql_compiler.gql.renderer_dataclasses import DataclassesRenderer
from graphql_compiler.gql.utils_schema import (
    compile_schema_library,
    read_fragment_queries,
)


def test_graphql_compilation(schema_library: str, graphql_library: str) -> None:
    schema = compile_schema_library(schema_library)
    fragment_library = read_fragment_queries(graphql_library)
    file_names = glob.glob(
        os.path.join(graphql_library, "**/*.graphql"), recursive=True
    )

    query_parser = QueryParser(schema)
    query_renderer = DataclassesRenderer(schema)

    for file_name in file_names:
        process_file(
            file_name, schema, query_parser, query_renderer, fragment_library, True
        )


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_library_path", help="email to connect to inventory with", type=str
    )
    parser.add_argument(
        "graphql_library_path", help="inventory connection password", type=str
    )
    args: argparse.Namespace = parser.parse_args()

    schema_library_path: str = args.schema_library_path
    graphql_library_path: str = args.graphql_library_path

    test_graphql_compilation(schema_library_path, graphql_library_path)

    sys.exit(0)

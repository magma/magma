#!/usr/bin/env python3
import glob
import os
from unittest.mock import MagicMock

import mock
import pkg_resources
from graphql_compiler.gql.cli import process_file
from graphql_compiler.gql.query_parser import QueryParser
from graphql_compiler.gql.renderer_dataclasses import DataclassesRenderer
from graphql_compiler.gql.utils_schema import (
    compile_schema_library,
    read_fragment_queries,
)
from libfb.py import testutil


class PyinventoryGraphqlCompilationTest(testutil.BaseFacebookTestCase):
    @mock.patch("os.path.exists", return_value=True)
    @mock.patch("subprocess.Popen")
    def test_graphql_compilation(self, mock_popen, mock_path_exists: MagicMock) -> None:
        schema_library = pkg_resources.resource_filename(__name__, "schema")
        schema = compile_schema_library(schema_library)
        graphql_library = pkg_resources.resource_filename(__name__, "graphql")
        fragment_library = read_fragment_queries(graphql_library)

        query_parser = QueryParser(schema)
        query_renderer = DataclassesRenderer(schema)

        filenames = glob.glob(
            os.path.join(graphql_library, "**/*.graphql"), recursive=True
        )

        for filename in filenames:
            process_file(
                filename, schema, query_parser, query_renderer, fragment_library, True
            )

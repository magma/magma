#!/usr/bin/env python3
import glob
import os
from graphql import build_ast_schema
from graphql.language.parser import parse
from fbc.symphony.cli.graphql_compiler.gql.query_parser import QueryParser
from fbc.symphony.cli.graphql_compiler.gql.renderer_dataclasses \
    import DataclassesRenderer
from libfb.py import testutil
import mock
import pkg_resources


class PyinventoryGraphqlCompilationTest(testutil.BaseFacebookTestCase):

    def verify_file(
            self,
            filename: str,
            parser: QueryParser,
            renderer: DataclassesRenderer):
        root, _s = os.path.splitext(filename)
        target_filename = root + '.py'

        with open(filename, 'r') as fin:
            query = fin.read()
            parsed = parser.parse(query)
            rendered = renderer.render(parsed)
            assert rendered == open(target_filename, 'r').read(), \
                f"Generated file name {target_filename} does " \
                "not match compilation result"

    @mock.patch("os.path.exists", return_value=True)
    @mock.patch("subprocess.Popen")
    def test_graphql_compilation(self, mock_popen, mock_path_exists):
        graphql_library = pkg_resources.resource_filename(
            __name__, "graphql")
        schema_filepath = pkg_resources.resource_filename(
            __name__, "schema/symphony.graphql")

        schema = build_ast_schema(parse((open(schema_filepath).read())))
        filenames = glob.glob(
            os.path.join(graphql_library, "*.graphql"), recursive=True)
        query_parser = QueryParser(schema)
        query_renderer = DataclassesRenderer(schema)

        for filename in filenames:
            self.verify_file(filename, query_parser, query_renderer)

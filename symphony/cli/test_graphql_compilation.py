#!/usr/bin/env python3
import glob
import os

import mock
import pkg_resources
from fbc.symphony.cli.graphql_compiler.gql.query_parser import QueryParser
from fbc.symphony.cli.graphql_compiler.gql.renderer_dataclasses import (
    DataclassesRenderer,
)
from fbc.symphony.cli.graphql_compiler.gql.utils_codegen import (
    get_enum_filename,
    get_input_filename,
)
from graphql import build_ast_schema
from graphql.language.parser import parse
from libfb.py import testutil


class PyinventoryGraphqlCompilationTest(testutil.BaseFacebookTestCase):
    def assert_rendered_file(self, file_name, file_content, rendered):
        assert (
            rendered == file_content
        ), f"""Generated file name {file_name} does
            not match compilation result:
            exising file:
            {file_content}
            compilation result:
            {rendered}"""

    def verify_file(
        self, filename: str, parser: QueryParser, renderer: DataclassesRenderer
    ) -> None:
        root, _s = os.path.splitext(filename)
        target_filename = "".join([root, ".py"])
        dir_name = os.path.dirname(target_filename)

        with open(filename, "r") as fin:
            query = fin.read()
            parsed = parser.parse(query)
            rendered = renderer.render(parsed)
            with open(target_filename, "r") as f:
                file_content = f.read()
                self.assert_rendered_file(target_filename, file_content, rendered)
            enums = renderer.render_enums(parsed)
            for enum_name, code in enums.items():
                target_enum_filename = os.path.join(
                    dir_name, "".join([get_enum_filename(enum_name), ".py"])
                )
                with open(target_enum_filename, "r") as f:
                    file_content = f.read()
                    self.assert_rendered_file(target_enum_filename, file_content, code)
            input_objects = renderer.render_input_objects(parsed)
            for input_object_name, code in input_objects.items():
                target_input_object_filename = os.path.join(
                    dir_name, "".join([get_input_filename(input_object_name), ".py"])
                )
                with open(target_input_object_filename, "r") as f:
                    file_content = f.read()
                    self.assert_rendered_file(
                        target_input_object_filename, file_content, code
                    )

    @mock.patch("os.path.exists", return_value=True)
    @mock.patch("subprocess.Popen")
    def test_graphql_compilation(self, mock_popen, mock_path_exists):
        graphql_library = pkg_resources.resource_filename(__name__, "graphql")
        schema_filepath = pkg_resources.resource_filename(
            __name__, "schema/symphony.graphql"
        )

        schema = build_ast_schema(parse((open(schema_filepath).read())))
        filenames = glob.glob(
            os.path.join(graphql_library, "**/*.graphql"), recursive=True
        )
        query_parser = QueryParser(schema)
        query_renderer = DataclassesRenderer(schema)

        for filename in filenames:
            self.verify_file(filename, query_parser, query_renderer)

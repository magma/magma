#!/usr/bin/env python3
import argparse
import glob
import os
from typing import Dict

from graphql import GraphQLSchema
from graphql.language.parser import parse
from graphql.utilities.find_deprecated_usages import find_deprecated_usages

from .query_parser import AnonymousQueryError, InvalidQueryError, QueryParser
from .renderer_dataclasses import DataclassesRenderer
from .utils_codegen import get_enum_filename, get_input_filename
from .utils_schema import compile_schema_library, read_fragment_queries


DEFAULT_CONFIG_FNAME = ".gql.json"


def assert_rendered_file(file_name: str, file_content: str, rendered: str) -> None:
    assert (
        rendered == file_content
    ), f"""Generated file name {file_name} does
            not match compilation result:
            exising file:
            {file_content}
            compilation result:
            {rendered}"""


def safe_remove(fname: str) -> None:
    try:
        os.remove(fname)
    except BaseException:
        pass


def verify_or_write_rendered(filename: str, rendered: str, verify: bool) -> None:
    if verify:
        with open(filename, "r") as f:
            file_content = f.read()
            assert_rendered_file(filename, file_content, rendered)
    else:
        with open(filename, "w") as outfile:
            outfile.write(rendered)


def process_file(
    filename: str,
    schema: GraphQLSchema,
    parser: QueryParser,
    renderer: DataclassesRenderer,
    fragment_library: Dict[str, str],
    verify: bool = False,
) -> None:
    full_fragments = "".join(
        [
            fragment_code
            for fragment_filename, fragment_code in fragment_library.items()
            if fragment_filename != filename
        ]
    )
    root, _s = os.path.splitext(filename)
    target_filename = "".join([root, ".py"])
    dir_name = os.path.dirname(target_filename)

    try:
        with open(filename, "r") as fin:
            query = fin.read()
            parsed_query = parse(query)
            usages = find_deprecated_usages(schema, parsed_query)
            assert (
                len(usages) == 0
            ), f"Graphql file name {filename} uses deprecated fields {usages}"
            is_fragment = filename.endswith("_fragment.graphql")
            parsed = parser.parse(query, full_fragments, is_fragment=is_fragment)
            rendered = renderer.render(parsed)
            verify_or_write_rendered(target_filename, rendered, verify)

            enums = renderer.render_enums(parsed)
            for enum_name, code in enums.items():
                target_enum_filename = os.path.join(
                    dir_name, "".join([get_enum_filename(enum_name), ".py"])
                )
                verify_or_write_rendered(target_enum_filename, code, verify)
            input_objects = renderer.render_input_objects(parsed)
            for input_object_name, code in input_objects.items():
                target_input_object_filename = os.path.join(
                    dir_name, "".join([get_input_filename(input_object_name), ".py"])
                )
                verify_or_write_rendered(target_input_object_filename, code, verify)
    except (AnonymousQueryError, InvalidQueryError, AssertionError):
        if verify:
            print(f"Failed to verify graphql file {filename}")
        else:
            print(f"Failed to process graphql file {filename}")
        safe_remove(target_filename)
        raise


def run(schema_library: str, graphql_library: str) -> None:
    schema = compile_schema_library(schema_library)
    fragment_library = read_fragment_queries(graphql_library)
    filenames = glob.glob(os.path.join(graphql_library, "**/*.graphql"), recursive=True)

    query_parser = QueryParser(schema)
    query_renderer = DataclassesRenderer(schema)

    py_filenames = glob.glob(os.path.join(graphql_library, "**/*.py"), recursive=True)
    for py_filename in py_filenames:
        if os.path.basename(py_filename) != "__init__.py":
            os.unlink(py_filename)

    for filename in filenames:
        process_file(filename, schema, query_parser, query_renderer, fragment_library)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_library", help="the graphql schemas storage path", type=str
    )
    parser.add_argument(
        "graphql_library", help="path where all queries files are stored", type=str
    )
    args = parser.parse_args()
    run(args.schema_library, args.graphql_library)

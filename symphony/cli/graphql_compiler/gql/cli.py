#!/usr/bin/env python3
import argparse
import glob
import os

from graphql.language.parser import parse
from graphql.utilities.find_deprecated_usages import find_deprecated_usages

from .query_parser import AnonymousQueryError, InvalidQueryError, QueryParser
from .renderer_dataclasses import DataclassesRenderer
from .utils_codegen import get_enum_filename, get_input_filename
from .utils_schema import compile_schema_library


DEFAULT_CONFIG_FNAME = ".gql.json"


def safe_remove(fname):
    try:
        os.remove(fname)
    except BaseException:
        pass


def process_file(
    filename: str, parser: QueryParser, renderer: DataclassesRenderer
) -> None:
    root, _s = os.path.splitext(filename)
    target_filename = "".join([root, ".py"])
    dir_name = os.path.dirname(target_filename)

    with open(filename, "r") as fin:
        query = fin.read()
        try:
            parsed = parser.parse(query)
            rendered = renderer.render(parsed)
            with open(target_filename, "w") as outfile:
                outfile.write(rendered)

            enums = renderer.render_enums(parsed)
            for enum_name, code in enums.items():
                target_enum_filename = os.path.join(
                    dir_name, "".join([get_enum_filename(enum_name), ".py"])
                )
                with open(target_enum_filename, "w") as outfile:
                    outfile.write(code)
            input_objects = renderer.render_input_objects(parsed)
            for input_object_name, code in input_objects.items():
                target_input_object_filename = os.path.join(
                    dir_name, "".join([get_input_filename(input_object_name), ".py"])
                )
                with open(target_input_object_filename, "w") as outfile:
                    outfile.write(code)

        except (AnonymousQueryError, InvalidQueryError):
            safe_remove(target_filename)
            raise


def run(schema_library: str, graphql_library: str) -> None:
    schema = compile_schema_library(schema_library)
    filenames = glob.glob(os.path.join(graphql_library, "**/*.graphql"), recursive=True)

    query_parser = QueryParser(schema)
    query_renderer = DataclassesRenderer(schema)

    py_filenames = glob.glob(os.path.join(graphql_library, "**/*.py"), recursive=True)
    for py_filename in py_filenames:
        if os.path.basename(py_filename) != "__init__.py":
            os.unlink(py_filename)

    for filename in filenames:
        with open(filename) as f:
            query = parse(f.read())
            usages = find_deprecated_usages(schema, query)
            assert (
                len(usages) == 0
            ), f"Graphql file name {filename} uses deprecated fields {usages}"

        process_file(filename, query_parser, query_renderer)


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

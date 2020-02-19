#!/usr/bin/env python3
import argparse
import glob
import os
from graphql import build_ast_schema
from graphql.language.parser import parse

from .query_parser import QueryParser, AnonymousQueryError, InvalidQueryError
from .renderer_dataclasses import DataclassesRenderer
from .utils_codegen import get_enum_filename

DEFAULT_CONFIG_FNAME = '.gql.json'


def safe_remove(fname):
    try:
        os.remove(fname)
    except BaseException:
        pass


def process_file(
        filename: str,
        parser: QueryParser,
        renderer: DataclassesRenderer,
        verify: bool):
    root, _s = os.path.splitext(filename)
    target_filename = root + '.py'

    with open(filename, 'r') as fin:
        query = fin.read()
        try:
            parsed = parser.parse(query)
            rendered = renderer.render(parsed)
            if verify:
                assert rendered == open(target_filename, 'r').read(), \
                    f"Generated file name {target_filename} does " \
                    "not match compilation result"
            else:
                with open(target_filename, 'w') as outfile:
                    outfile.write(rendered)

        except (AnonymousQueryError, InvalidQueryError):
            safe_remove(target_filename)
            raise
    enums = renderer.render_enums(parsed)
    for enum, code in enums.items():
        target_enum_filename = os.path.join(
            os.path.dirname(target_filename), "".join([get_enum_filename(enum), '.py']))
        if verify:
            assert code == open(target_enum_filename, 'r').read(), \
                f"Generated file name {target_enum_filename} does " \
                "not match compilation result"
        else:
            with open(target_enum_filename, 'w') as outfile:
                outfile.write(code)


def run(schema_filepath: str, graphql_library: str, verify: bool):
    schema = build_ast_schema(parse((open(schema_filepath).read())))

    filenames = glob.glob(os.path.join(graphql_library, "*.graphql"), recursive=True)

    query_parser = QueryParser(schema)
    query_renderer = DataclassesRenderer(schema)

    for filename in filenames:
        process_file(filename, query_parser, query_renderer, verify)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_filepath", help="the path of grahql schema file", type=str)
    parser.add_argument(
        "graphql_library", help="path where all queries files are stored", type=str
    )
    parser.add_argument('--verify', action='store_true')
    args = parser.parse_args()
    run(args.schema_filepath, args.graphql_library, args.verify)

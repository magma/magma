#!/usr/bin/env python3

import argparse

from graphql import (
    GraphQLEnumType,
    GraphQLInputObjectType,
    GraphQLScalarType,
    build_ast_schema,
)
from graphql.language.parser import parse


DEPRECATED_DOC = """---
id: graphql-breaking-changes
title: Graphql API Breaking Changes
---

[//]: <> ({})

## Deprecated Queries
{}

## Deprecated Mutations
{}

## Deprecated Fields
{}

"""


def document_all_deprecated_functions(schema_filepath: str, doc_filepath: str) -> None:
    deprecated_queries = []
    deprecated_mutations = []
    deprecated_fields = []
    with open(schema_filepath) as f:
        schema = build_ast_schema(parse(f.read()))
        for query_name, query_field in schema.query_type.fields.items():
            if query_field.is_deprecated:
                deprecated_queries.append(
                    " - ".join([f"`{query_name}`", query_field.deprecation_reason])
                )
        for mutation_name, mutation_field in schema.mutation_type.fields.items():
            if query_field.is_deprecated:
                deprecated_mutations.append(
                    " - ".join(
                        [f"`{mutation_name}`", mutation_field.deprecation_reason]
                    )
                )
        for type_name, type_object in schema.type_map.items():
            if isinstance(
                type_object,
                (GraphQLScalarType, GraphQLEnumType, GraphQLInputObjectType),
            ) or type_name in ("Query", "Mutation"):
                continue
            for field_name, field_object in type_object.fields.items():
                if field_object.is_deprecated:
                    field_long_name = ".".join([type_name, field_name])
                    deprecated_fields.append(
                        " - ".join(
                            [f"`{field_long_name}`", field_object.deprecation_reason]
                        )
                    )
    with open(doc_filepath, "w") as f:
        f.write(
            DEPRECATED_DOC.format(
                "@"
                + "generated This file was created by cli/extract_graphql_deprecations.py"
                + "do not change it manually",
                "\n".join(map(lambda s: "* " + s, deprecated_queries)),
                "\n".join(map(lambda s: "* " + s, deprecated_mutations)),
                "\n".join(map(lambda s: "* " + s, deprecated_fields)),
            )
        )


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_filepath", help="the path of grahql schema file", type=str
    )
    parser.add_argument(
        "doc_filepath", help="the path of grahql documentation file", type=str
    )
    args = parser.parse_args()
    document_all_deprecated_functions(args.schema_filepath, args.doc_filepath)

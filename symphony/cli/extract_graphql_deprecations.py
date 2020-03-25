#!/usr/bin/env python3

import argparse

from graphql import (
    GraphQLEnumType,
    GraphQLInputObjectType,
    GraphQLScalarType,
    build_ast_schema,
)
from graphql.language.parser import parse
from graphql_compiler.gql.utils_schema import compile_schema_library


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

## Deprecated Input Fields
{}

"""


def document_all_deprecated_functions(schema_library: str, doc_filepath: str) -> None:
    deprecated_queries = []
    deprecated_mutations = []
    deprecated_inputs = []
    deprecated_fields = []
    schema = compile_schema_library(schema_library)
    for query_name, query_field in schema.query_type.fields.items():
        if query_field.is_deprecated:
            deprecated_queries.append(
                " - ".join([f"`{query_name}`", query_field.deprecation_reason])
            )
    for mutation_name, mutation_field in schema.mutation_type.fields.items():
        if query_field.is_deprecated:
            deprecated_mutations.append(
                " - ".join([f"`{mutation_name}`", mutation_field.deprecation_reason])
            )
    for type_name, type_object in schema.type_map.items():
        if isinstance(
            type_object, (GraphQLScalarType, GraphQLEnumType)
        ) or type_name in ("Query", "Mutation"):
            continue
        if isinstance(type_object, GraphQLInputObjectType):
            for field_name, field_object in type_object.fields.items():
                directives = field_object.ast_node.directives
                for directive in directives:
                    if directive.name.value == "deprecatedInput":
                        field_long_name = ".".join([type_name, field_name])
                        for argument in directive.arguments:
                            if argument.name.value == "reason":
                                deprecation_reason = (
                                    argument.value.value
                                    if directive.arguments
                                    else None
                                )
                                deprecated_inputs.append(
                                    " - ".join(
                                        [f"`{field_long_name}`", deprecation_reason]
                                    )
                                )
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
                "\n".join(map(lambda s: "* " + s, deprecated_inputs)),
            )
        )


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "schema_library", help="the graphql schemas storage path", type=str
    )
    parser.add_argument(
        "doc_filepath", help="the path of grahql documentation file", type=str
    )
    args = parser.parse_args()
    document_all_deprecated_functions(args.schema_library, args.doc_filepath)

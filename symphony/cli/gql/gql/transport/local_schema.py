#!/usr/bin/env python3

from typing import Any, Dict, cast

from graphql.execution import ExecutionResult, execute
from graphql.language.ast import DocumentNode
from graphql.type.schema import GraphQLSchema

from .transport import ExtendedExecutionResult, Transport


class LocalSchemaTransport(Transport):
    def __init__(self, schema: GraphQLSchema) -> None:
        self.schema = schema

    def execute(
        self, document: DocumentNode, variable_values: Dict[str, Any] = {}  # noqa: B006
    ) -> ExtendedExecutionResult:
        result = cast(ExecutionResult, execute(self.schema, document))
        return ExtendedExecutionResult(
            response="", data=result.data, errors={}, extensions={}
        )

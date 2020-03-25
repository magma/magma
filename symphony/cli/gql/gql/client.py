#!/usr/bin/env python3

import warnings
from logging import Logger, getLogger
from typing import Any, Dict, Optional, cast

from graphql import (
    build_ast_schema,
    build_client_schema,
    get_introspection_query,
    parse,
)
from graphql.language.ast import DocumentNode
from graphql.type.schema import GraphQLSchema
from graphql.utilities.find_deprecated_usages import find_deprecated_usages
from graphql.validation import validate

from .transport.local_schema import LocalSchemaTransport
from .transport.transport import ExtendedExecutionResult, Transport


log: Logger = getLogger(__name__)


class OperationException(Exception):
    def __init__(self, err_msg: str, err_id: str) -> None:
        message = "Operation failed: %s (id:%s)" % (err_msg, err_id)
        super(OperationException, self).__init__(message)
        self.err_msg = err_msg
        self.err_id = err_id


class RetryError(Exception):
    """Custom exception thrown when retry logic fails"""

    def __init__(self, retries_count: int, last_exception: Optional[Exception]) -> None:
        message = "Failed %s retries: %s" % (retries_count, last_exception)
        super(RetryError, self).__init__(message)
        self.last_exception = last_exception


class GraphqlDeprecationWarning(DeprecationWarning):
    pass


class Client(object):
    schema: Optional[GraphQLSchema]
    introspection: Optional[Dict[str, Any]]
    transport: Transport
    retries: int

    def __init__(
        self,
        schema: Optional[GraphQLSchema] = None,
        introspection: Optional[Dict[str, Any]] = None,
        type_def: Optional[str] = None,
        transport: Optional[Transport] = None,
        fetch_schema_from_transport: bool = False,
        retries: int = 0,
    ) -> None:
        assert not (
            type_def and introspection
        ), "Cant provide introspection type definition at the same time"
        if transport and fetch_schema_from_transport:
            assert (
                not schema
            ), "Cant fetch the schema from transport if is already provided"
            introspection = transport.execute(
                parse(get_introspection_query(descriptions=True))
            ).data
        if introspection:
            assert not schema, "Cant provide introspection and schema at the same time"
            schema = build_client_schema(introspection)
        elif type_def:
            assert (
                not schema
            ), "Cant provide Type definition and schema at the same time"
            type_def_ast = parse(type_def)
            schema = build_ast_schema(type_def_ast)
        elif schema and not transport:
            transport = LocalSchemaTransport(schema)

        self.schema = schema
        self.introspection = introspection
        self.transport = cast(Transport, transport)
        self.retries = retries

    def validate(self, document: DocumentNode) -> None:
        schema = self.schema
        if not schema:
            raise Exception(
                "Cannot validate locally the document, you need to pass a schema."
            )
        validation_errors = validate(schema, document)
        if validation_errors:
            raise validation_errors[0]
        usages = find_deprecated_usages(schema, document)
        for usage in usages:
            message = (
                f"Query of deprecated grapqhl field in {usage}"
                "Consider upgrading to newer API version."
            )
            warnings.warn(message, GraphqlDeprecationWarning)

    def execute(self, document: DocumentNode, variable_values: Dict[str, Any]) -> str:
        if self.schema:
            self.validate(document)

        result = self._get_result(document, variable_values)
        if result.errors:
            raise OperationException(
                str(cast(Dict[int, str], result.errors)[0]),
                result.extensions.get("trace_id", ""),
            )

        return result.response

    def _get_result(
        self, document: DocumentNode, variable_values: Dict[str, Any]
    ) -> ExtendedExecutionResult:
        if not self.retries:
            return self.transport.execute(document, variable_values)

        last_exception = None
        retries_count = 0
        while retries_count < self.retries:
            try:
                result = self.transport.execute(document, variable_values)
                return result
            except Exception as e:
                last_exception = e
                log.warn(
                    "Request failed with exception %s. Retrying for the %s time...",
                    e,
                    retries_count + 1,
                    exc_info=True,
                )
            finally:
                retries_count += 1

        raise RetryError(retries_count, last_exception)

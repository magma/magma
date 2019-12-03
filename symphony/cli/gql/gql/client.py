#!/usr/bin/env python3

from logging import getLogger, Logger
from graphql import build_ast_schema, build_client_schema, \
    get_introspection_query, parse
from graphql.validation import validate
from .transport.local_schema import LocalSchemaTransport

log: Logger = getLogger(__name__)


class OperationException(Exception):
    def __init__(self, err_msg: str, err_id: str) -> None:
        message = "Operation failed: %s (id:%s)" % (err_msg, err_id)
        super(OperationException, self).__init__(message)
        self.err_msg = err_msg
        self.err_id = err_id


class RetryError(Exception):
    """Custom exception thrown when retry logic fails"""

    def __init__(self, retries_count: int, last_exception: Exception) -> None:
        message = "Failed %s retries: %s" % (retries_count, last_exception)
        super(RetryError, self).__init__(message)
        self.last_exception = last_exception


class Client(object):
    def __init__(
        self,
        schema=None,
        introspection=None,
        type_def=None,
        transport=None,
        fetch_schema_from_transport=False,
        retries=0,
    ):
        assert not (
            type_def and introspection
        ), "Cant provide introspection type definition at the same time"
        if transport and fetch_schema_from_transport:
            assert (
                not schema
            ), "Cant fetch the schema from transport if is already provided"
            introspection = transport.execute(parse(
                get_introspection_query(descriptions=True))).data
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
        self.transport = transport
        self.retries = retries

    def validate(self, document):
        if not self.schema:
            raise Exception(
                "Cannot validate locally the document, you need to pass a schema."
            )
        validation_errors = validate(self.schema, document)
        if validation_errors:
            raise validation_errors[0]

    def execute(self, document, *args, **kwargs):
        if self.schema:
            self.validate(document)

        result = self._get_result(document, *args, **kwargs)
        if result.errors:
            raise OperationException(
                str(result.errors[0]), result.extensions.get("trace_id", ""))

        return result.data

    def _get_result(self, document, *args, **kwargs):
        if not self.retries:
            return self.transport.execute(document, *args, **kwargs)

        last_exception = None
        retries_count = 0
        while retries_count < self.retries:
            try:
                result = self.transport.execute(document, *args, **kwargs)
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

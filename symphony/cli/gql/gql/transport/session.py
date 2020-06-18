#!/usr/bin/env python3

from __future__ import absolute_import

import json
from datetime import datetime, timedelta, tzinfo
from enum import Enum
from http import HTTPStatus
from typing import Any, Dict, Optional, Union

from gql.gql.transport.http import HTTPTransport
from graphql.language.ast import DocumentNode
from graphql.language.printer import print_ast
from requests.auth import AuthBase
from requests.sessions import Session

from .transport import ExtendedExecutionResult


class MissingEnumException(Exception):
    def __init__(self, variable: Enum) -> None:
        self.enum_type = str(type(variable))

    def __str__(self) -> str:
        return f"Try to encode missing value of enum {self.enum_type}"


class UserDeactivatedException(Exception):
    pass


class simple_utc(tzinfo):
    def tzname(self, dt: Optional[datetime]) -> Optional[str]:
        return "UTC"

    def utcoffset(self, dt: Optional[datetime]) -> Optional[timedelta]:
        return timedelta(0)


def encode_variable(
    variable: Union[datetime, Enum, object]
) -> Union[str, Dict[str, Any]]:
    if isinstance(variable, datetime):
        return datetime.isoformat(variable.replace(tzinfo=simple_utc()))
    elif isinstance(variable, Enum):
        if variable.value == "":
            raise MissingEnumException(variable)
        return variable.value
    else:
        return variable.__dict__


class RequestsHTTPSessionTransport(HTTPTransport):
    def __init__(
        self,
        session: Session,
        url: str,
        headers: Optional[Dict[str, str]] = None,
        auth: Optional[AuthBase] = None,
    ) -> None:
        """
        :param session: The session
        """
        super(RequestsHTTPSessionTransport, self).__init__(url, headers)
        self.session: Session = session
        self.auth = auth

    def execute(
        self, document: DocumentNode, variable_values: Dict[str, Any] = {}  # noqa: B006
    ) -> ExtendedExecutionResult:
        query_str = print_ast(document)
        payload = {"query": query_str, "variables": variable_values}

        response = self.session.post(
            self.url,
            data=json.dumps(payload, default=encode_variable).encode("utf-8"),
            headers=self.headers,
            auth=self.auth,
        )

        if (
            response.status_code == HTTPStatus.FORBIDDEN
            and response.text == "user is deactivated\n"
        ):
            raise UserDeactivatedException()

        if response.status_code == HTTPStatus.GATEWAY_TIMEOUT:
            raise TimeoutError()

        result = response.json()

        extensions = {}
        if "x-correlation-id" in response.headers:
            extensions["trace_id"] = response.headers["x-correlation-id"]

        assert (
            "errors" in result or "data" in result
        ), 'Received non-compatible response "{}"'.format(result)

        return ExtendedExecutionResult(
            response=response.text,
            errors=result.get("errors"),
            data=result.get("data"),
            extensions=extensions,
        )

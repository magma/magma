#!/usr/bin/env python3

from __future__ import absolute_import

import json
from datetime import datetime, timedelta, tzinfo
from enum import Enum
from typing import Optional

from gql.gql.transport.http import HTTPTransport
from graphql.language.printer import print_ast

from .exceptions import MissingEnumException


class simple_utc(tzinfo):
    def tzname(self, dt: Optional[datetime]) -> Optional[str]:
        return "UTC"

    def utcoffset(self, dt: Optional[datetime]) -> Optional[timedelta]:
        return timedelta(0)


def encode_variable(variable):
    if isinstance(variable, datetime):
        return datetime.isoformat(variable.replace(tzinfo=simple_utc()))
    elif isinstance(variable, Enum):
        if variable.value == "":
            raise MissingEnumException(variable)
        return variable.value
    else:
        return variable.__dict__


class ExtendedExecutionResult:
    def __init__(self, errors, data, extensions):
        self.errors = errors
        self.data = data
        self.extensions = extensions


class RequestsHTTPSessionTransport(HTTPTransport):
    def __init__(self, session, url, auth=None, use_json=False, timeout=None, **kwargs):
        """
        :param session: The session
        :param auth: Auth tuple or callable to enable Basic/Digest/Custom HTTP Auth
        :param use_json: Send request body as JSON instead of form-urlencoded
        :param timeout: Specifies a default timeout for requests (Default: None)
        """
        super(RequestsHTTPSessionTransport, self).__init__(url, **kwargs)
        self.session = session
        self.auth = auth
        self.default_timeout = timeout
        self.use_json = use_json

    def execute(self, document, variable_values=None, timeout=None, return_json=True):
        query_str = print_ast(document)
        payload = {"query": query_str, "variables": variable_values or {}}

        request = self.session.post(
            self.url,
            data=json.dumps(payload, default=encode_variable).encode("utf-8"),
            headers=self.headers,
        )

        result = request.json()

        extensions = {}
        if "x-correlation-id" in request.headers:
            extensions["trace_id"] = request.headers["x-correlation-id"]

        assert (
            "errors" in result or "data" in result
        ), 'Received non-compatible response "{}"'.format(result)

        data = result.get("data") if return_json else request.text
        return ExtendedExecutionResult(
            errors=result.get("errors"), data=data, extensions=extensions
        )

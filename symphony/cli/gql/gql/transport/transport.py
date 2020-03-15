#!/usr/bin/env python3

import abc
from typing import Any, Dict

from graphql.language.ast import DocumentNode


class ExtendedExecutionResult:
    def __init__(
        self,
        response: str,
        errors: Dict[str, Any],
        data: Dict[str, Any],
        extensions: Dict[str, Any],
    ) -> None:
        self.response: str = response
        self.errors: Dict[str, Any] = errors
        self.data: Dict[str, Any] = data
        self.extensions: Dict[str, Any] = extensions


class Transport(abc.ABC):
    @abc.abstractmethod
    def execute(
        self, document: DocumentNode, variable_values: Dict[str, Any] = {}  # noqa: B006
    ) -> ExtendedExecutionResult:
        pass

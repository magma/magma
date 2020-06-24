#!/usr/bin/env python3

import abc
from typing import Any, Dict, Optional

from graphql.language.ast import DocumentNode


class ExtendedExecutionResult:
    def __init__(
        self,
        response: str,
        errors: Optional[Dict[str, Any]],
        data: Optional[Dict[str, Any]],
        extensions: Optional[Dict[str, Any]],
    ) -> None:
        self.response: str = response
        self.errors: Dict[str, Any] = errors if errors is not None else {}
        self.data: Dict[str, Any] = data if data is not None else {}
        self.extensions: Dict[str, Any] = extensions if extensions is not None else {}


class Transport(abc.ABC):
    @abc.abstractmethod
    def execute(
        self, document: DocumentNode, variable_values: Dict[str, Any] = {}  # noqa: B006
    ) -> ExtendedExecutionResult:
        pass

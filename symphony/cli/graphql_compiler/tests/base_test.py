#!/usr/bin/env python3

import os
import sys
import unittest
from types import ModuleType

import graphql_compiler.tests.testmodule as testmodule
from graphql_compiler.gql.query_parser import QueryParser
from graphql_compiler.gql.renderer_dataclasses import DataclassesRenderer
from graphql_compiler.gql.utils_schema import load_schema


class BaseTest(unittest.TestCase):
    def setUp(self) -> None:
        if "SWAPI_SCHEMA" not in globals():
            filename = os.path.join(
                os.path.dirname(__file__), "fixtures/swapi-schema.graphql"
            )
            globals()["SWAPI_SCHEMA"] = load_schema(filename)

        self.swapi_schema = globals()["SWAPI_SCHEMA"]
        self.swapi_parser = QueryParser(self.swapi_schema)

        if "GITHUB_SCHEMA" not in globals():
            filename = os.path.join(
                os.path.dirname(__file__), "fixtures/github-schema.graphql"
            )
            globals()["GITHUB_SCHEMA"] = load_schema(filename)

        self.github_schema = globals()["GITHUB_SCHEMA"]
        self.github_parser = QueryParser(self.github_schema)

        self.swapi_dataclass_renderer = DataclassesRenderer(self.swapi_schema)
        self.github_dataclass_renderer = DataclassesRenderer(self.github_schema)

    def load_module(self, code, module_name=None):
        compiled = compile(code, "", "exec")
        module = testmodule
        if module_name is not None:
            module_name = ".".join([testmodule.__name__, module_name])
            module = ModuleType(module_name)
            sys.modules[module_name] = module
        exec(compiled, module.__dict__)
        return testmodule

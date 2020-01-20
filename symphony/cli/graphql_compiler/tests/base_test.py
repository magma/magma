#!/usr/bin/env python3

from fbc.symphony.cli.graphql_compiler.gql.utils_schema import load_schema
from fbc.symphony.cli.graphql_compiler.gql.query_parser import QueryParser
from fbc.symphony.cli.graphql_compiler.gql.renderer_dataclasses import DataclassesRenderer
import fbc.symphony.cli.graphql_compiler.tests.testmodule as testmodule
import os
import unittest


class BaseTest(unittest.TestCase):

    def setUp(self) -> None:
        if 'SWAPI_SCHEMA' not in globals():
            filename = os.path.join(
                os.path.dirname(__file__), 'fixtures/swapi-schema.graphql')
            globals()['SWAPI_SCHEMA'] = load_schema(filename)
        self.swapi_schema = globals()['SWAPI_SCHEMA']
        self.swapi_parser = QueryParser(self.swapi_schema)

        if 'GITHUB_SCHEMA' not in globals():
            filename = os.path.join(
                os.path.dirname(__file__), 'fixtures/github-schema.graphql')
            globals()['GITHUB_SCHEMA'] = load_schema(filename)

        self.github_schema = globals()['GITHUB_SCHEMA']
        self.github_parser = QueryParser(self.github_schema)

        self.swapi_dataclass_renderer = DataclassesRenderer(self.swapi_schema)
        self.github_dataclass_renderer = DataclassesRenderer(self.github_schema)

    def load_module(self, code, module_name=None):
        compiled = compile(code, '', 'exec')
        exec(compiled, testmodule.__dict__)
        return testmodule

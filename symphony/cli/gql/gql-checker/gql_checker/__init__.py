import ast
import json

import pycodestyle

from gql_checker.__about__ import (
    __author__, __copyright__, __email__, __license__, __summary__, __title__,
    __uri__, __version__
)
from gql_checker.stdlib_list import STDLIB_NAMES
from graphql import Source, validate, parse, build_client_schema


__all__ = [
    "__title__", "__summary__", "__uri__", "__version__", "__author__",
    "__email__", "__license__", "__copyright__",
]

GQL_SYNTAX_ERROR = 'GQL100'
GQL_VALIDATION_ERROR = 'GQL101'

class ImportVisitor(ast.NodeVisitor):
    """
    This class visits all the gql calls.
    """

    def __init__(self, filename, options):
        self.filename = filename
        self.options = options or {}
        self.calls = []

    def visit_Call(self, node):  # noqa
        if node.func.id == 'gql':
            self.calls.append(node)

    def node_query(self, node):
        """
        Return the query for the gql call node
        """

        if isinstance(node, ast.Call):
            assert node.args
            arg = node.args[0]
            if not isinstance(arg, ast.Str):
                return
        else:
            raise TypeError(type(node))

        return arg.s


class ImportOrderChecker(object):
    visitor_class = ImportVisitor
    options = None

    def __init__(self, filename, tree):
        self.tree = tree
        self.filename = filename
        self.lines = None

    def load_file(self):
        if self.filename in ("stdin", "-", None):
            self.filename = "stdin"
            self.lines = pycodestyle.stdin_get_value().splitlines(True)
        else:
            self.lines = pycodestyle.readlines(self.filename)

        if not self.tree:
            self.tree = ast.parse("".join(self.lines))

    def get_schema(self):
        gql_introspection_schema = self.options.get('gql_introspection_schema')
        if gql_introspection_schema:
            try:
                with open(gql_introspection_schema) as data_file:
                    introspection_schema = json.load(data_file)
                    return build_client_schema(introspection_schema)
            except IOError as e:
                raise Exception("Cannot find the provided introspection schema. {}".format(str(e)))

        schema = self.options.get('schema')
        assert schema, 'Need to provide schema'

    def validation_errors(self, ast):
        return validate(self.get_schema(), ast)

    def error(self, node, code, message):
        raise NotImplemented()

    def check_gql(self):
        if not self.tree or not self.lines:
            self.load_file()

        visitor = self.visitor_class(self.filename, self.options)
        visitor.visit(self.tree)

        for node in visitor.calls:
            # Lines with the noqa flag are ignored entirely
            if pycodestyle.noqa(self.lines[node.lineno - 1]):
                continue

            query = visitor.node_query(node)
            if not query:
                continue

            try:
                source = Source(query, 'gql query')
                ast = parse(source)
            except Exception as e:
                message = str(e)
                yield self.error(node, GQL_SYNTAX_ERROR, message)
                continue

            validation_errors = self.validation_errors(ast)
            if validation_errors:
                for error in validation_errors:
                    message = str(error)
                    yield self.error(node, GQL_VALIDATION_ERROR, message)

from __future__ import absolute_import

import gql_checker
from gql_checker import ImportOrderChecker


class Linter(ImportOrderChecker):
    name = "gql"
    version = gql_checker.__version__

    def __init__(self, tree, filename):
        super(Linter, self).__init__(filename, tree)

    @classmethod
    def add_options(cls, parser):
        # List of application import names. They go last.
        parser.add_option(
            "--gql-introspection-schema",
            metavar="FILE",
            help="Import names to consider as application specific"
        )
        parser.add_option(
            "--gql-typedef-schema",
            default='',
            action="store",
            type="string",
            help=("Style to follow. Available: "
                  "cryptography, google, smarkets, pep8")
        )
        parser.config_options.append("gql-introspection-schema")
        parser.config_options.append("gql-typedef-schema")

    @classmethod
    def parse_options(cls, options):
        optdict = {}

        optdict = dict(
            gql_introspection_schema=options.gql_introspection_schema,
            gql_typedef_schema=options.gql_typedef_schema,
        )

        cls.options = optdict

    def error(self, node, code, message):
        lineno, col_offset = node.lineno, node.col_offset
        return (lineno, col_offset, '{0} {1}'.format(code, message), Linter)

    def run(self):
        for error in self.check_gql():
            yield error

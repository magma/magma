#!/usr/bin/env python3

from __future__ import absolute_import

import gql_checker
from gql_checker import ImportOrderChecker
from pylama.lint import Linter as BaseLinter


class Linter(ImportOrderChecker, BaseLinter):
    name = "gql"
    version = gql_checker.__version__

    def __init__(self):
        super(Linter, self).__init__(None, None)

    def allow(self, path):
        return path.endswith(".py")

    def error(self, node, code, message):
        lineno, col_offset = node.lineno, node.col_offset
        return {"lnum": lineno, "col": col_offset, "text": message, "type": code}

    def run(self, path, **meta):
        self.filename = path
        self.tree = None
        self.options = dict({"schema": ""}, **meta)

        for error in self.check_gql():
            yield error

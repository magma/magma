#!/usr/bin/env python3

import os
from distutils.version import LooseVersion

import mock
import pkg_resources
from graphql import build_ast_schema
from graphql.language.parser import parse
from graphql.utilities import find_breaking_changes
from graphql_compiler.gql.utils_schema import compile_schema_library
from libfb.py import testutil

from .utils import extract_zip


ERR_MSG_FORMAT = "Diff breaks schema of version {} with changes: {}"


class PyinventoryGraphqlBreakingChangesTest(testutil.BaseFacebookTestCase):
    @mock.patch("os.path.exists", return_value=True)
    @mock.patch("subprocess.Popen")
    def test_graphql_breaking_changes(self, mock_popen, mock_path_exists) -> None:
        schema_versions_library = pkg_resources.resource_filename(
            __name__, "schema_versions"
        )
        schema_library = pkg_resources.resource_filename(__name__, "schema")

        zip_filepath = os.path.join(schema_versions_library, "old_schemas.zip")
        minimal_version_filepath = os.path.join(
            schema_versions_library, "minimal_supported_version"
        )
        schema = compile_schema_library(schema_library)
        schemas = extract_zip(zip_filepath)
        with open(minimal_version_filepath) as minimal_version_file:
            minimal_version = minimal_version_file.read().strip()
            for ver, schema_str in schemas.items():
                if LooseVersion(ver) < LooseVersion(minimal_version):
                    continue
                schema_str = schema_str.decode("utf-8")
                old_schema = build_ast_schema(parse(schema_str))
                breaking_changes = find_breaking_changes(old_schema, schema)
                assert len(breaking_changes) == 0, ERR_MSG_FORMAT.format(
                    ver, breaking_changes
                )

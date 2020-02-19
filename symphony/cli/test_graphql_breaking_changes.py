#!/usr/bin/env python3

import os
from graphql import build_ast_schema
from graphql.language.parser import parse
from graphql.utilities import find_breaking_changes
from libfb.py import testutil
import mock
import pkg_resources
from distutils.version import LooseVersion
from .utils import extract_zip

ERR_MSG_FORMAT = 'Diff breaks schema of version {} with changes: {}'


class PyinventoryGraphqlBreakingChangesTest(testutil.BaseFacebookTestCase):

    @mock.patch("os.path.exists", return_value=True)
    @mock.patch("subprocess.Popen")
    def test_graphql_breaking_changes(self, mock_popen, mock_path_exists):
        schema_versions_library = pkg_resources.resource_filename(
            __name__, "schema_versions")
        schema_filepath = pkg_resources.resource_filename(
            __name__, "schema/symphony.graphql")

        zip_filepath = os.path.join(schema_versions_library, "old_schemas.zip")
        minimal_version_filepath = os.path.join(
            schema_versions_library, "minimal_supported_version")
        with open(schema_filepath) as schema_file:
            schema = build_ast_schema(parse(schema_file.read()))
            schemas = extract_zip(zip_filepath)
            with open(minimal_version_filepath) as minimal_version_file:
                minimal_version = minimal_version_file.read().strip()
                for ver, schema_str in schemas.items():
                    if LooseVersion(ver) < LooseVersion(minimal_version):
                        continue
                    schema_str = schema_str.decode("utf-8")
                    old_schema = build_ast_schema(parse(schema_str))
                    breaking_changes = find_breaking_changes(old_schema, schema)
                    assert len(breaking_changes) == 0, \
                        ERR_MSG_FORMAT.format(ver, breaking_changes)

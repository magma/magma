#!/usr/bin/env python3

import os

from utils import archive_zip, extract_zip


def get_minimal_supported_version(schema_versions_library):
    minimal_version_filepath = os.path.join(
        schema_versions_library, "minimal_supported_version"
    )
    with open(minimal_version_filepath) as minimal_version_file:
        return minimal_version_file.read().strip()


def set_minimal_supported_version(schema_versions_library, version):
    minimal_version_filepath = os.path.join(
        schema_versions_library, "minimal_supported_version"
    )
    with open(minimal_version_filepath, "w") as minimal_version_file:
        minimal_version_file.write(version)


def add_current_schema_with_version(
    schema_versions_library, current_schema_path, version
):
    old_schemas_archive = os.path.join(schema_versions_library, "old_schemas.zip")
    schemas = extract_zip(old_schemas_archive)
    schemas[version] = open(current_schema_path).read()
    archive_zip(old_schemas_archive, schemas)

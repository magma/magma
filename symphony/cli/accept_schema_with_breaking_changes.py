#!/usr/bin/env python3

from distutils.version import LooseVersion

from pyinventory.consts import __version__ as pyinventory_version
from schema_versioning_utils import (
    add_current_schema_with_version,
    get_minimal_supported_version,
    set_minimal_supported_version,
)


if __name__ == "__main__":
    minimal_version = get_minimal_supported_version("graphql_schema_versions")
    latest_version = max(
        LooseVersion(minimal_version), LooseVersion(pyinventory_version)
    )
    version_elements = latest_version.version
    version_elements = version_elements + [0] * (4 - len(version_elements))
    version_elements[-1] = version_elements[-1] + 1
    new_version = ".".join(map(str, version_elements))
    add_current_schema_with_version(
        "graphql_schema_versions",
        "../graph/graphql/schema/symphony.graphql",
        new_version,
    )
    set_minimal_supported_version("graphql_schema_versions", new_version)

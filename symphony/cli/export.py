#!/usr/bin/env python3

import argparse
import json
import os
import re
import sys
import time
from datetime import datetime
from distutils.version import LooseVersion

from export_doc import export_doc
from pyinventory import InventoryClient
from utils import archive_zip, extract_zip


GRAPHQL_PYINVENORY_CONTENT = (
    """// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@"""
    + """generated) by pyinventory, DO NOT EDIT.

package resolver

// PyinventoryConsts consists of metadata on python packages that were exported to store
const PyinventoryConsts = `
{}
`
"""
)

GRAPHQL_PYINVENORY_PATH = "../graph/graphql/resolver/pyinventory.go"


def export(email, password, useLocally, replaceLatestVersion, hasBreakingChange):
    files = os.listdir("./dist")
    whlFiles = [file for file in files if os.path.splitext(file)[1] == ".whl"]
    assert len(whlFiles), "More than one whl files"

    match = re.search("pyinventory-([.0-9]+)-py3-none-any.whl", whlFiles[0])
    version = match.group(1)

    res = "\n".join(open(GRAPHQL_PYINVENORY_PATH, "r").read().splitlines()[10:-1])

    packages = json.loads(res)

    client = InventoryClient(email, password, is_local_host=useLocally)

    if len(packages) != 0:
        latestVersion = packages[0]["version"]
        if LooseVersion(version) == LooseVersion(latestVersion):
            print("version {} is already exported".format(version))
            if replaceLatestVersion:
                print("Replace version {} with new version".format(version))
                latestPackage = packages.pop(0)
                try:
                    client.delete_file(latestPackage["whlFileKey"], True)
                except Exception:
                    print(
                        f'whlFileKey {latestPackage["whlFileKey"]} cannot ' "be deleted"
                    )
            else:
                return
        elif LooseVersion(version) < LooseVersion(latestVersion):
            print(
                "newer version {} is already exported than current version {}".format(
                    latestVersion, version
                )
            )
            return

    whlFileKey = client.store_file(
        os.path.join("./dist", whlFiles[0]), "application/zip", True
    )

    newPackage = {
        "version": version,
        "whlFileKey": whlFileKey,
        "uploadTime": datetime.isoformat(datetime.fromtimestamp(int(time.time())))
        + "+00:00",
        "hasBreakingChange": hasBreakingChange,
    }
    packages.insert(0, newPackage)

    newContent = json.dumps(packages)
    open(GRAPHQL_PYINVENORY_PATH, "w").write(
        GRAPHQL_PYINVENORY_CONTENT.format(newContent)
    )

    export_doc()
    schemas = extract_zip("graphql_schema_versions/old_schemas.zip")
    schemas[version] = open("../graph/graphql/schema/symphony.graphql").read()
    archive_zip("graphql_schema_versions/old_schemas.zip", schemas)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("email", help="email to connect to inventory with", type=str)
    parser.add_argument(
        "password", help="password to connect to inventory with", type=str
    )
    parser.add_argument(
        "-l",
        "--use-locally",
        help="upload new python package to local inventory",
        action="store_true",
    )
    parser.add_argument(
        "-r",
        "--replace-latest-version",
        help="replace the latest package version with this version",
        action="store_true",
    )
    parser.add_argument(
        "-b",
        "--has-breaking-change",
        help="forbid users to use older versions than this version",
        action="store_true",
    )
    args = parser.parse_args()
    export(
        args.email,
        args.password,
        args.use_locally,
        args.replace_latest_version,
        args.has_breaking_change,
    )
    sys.exit(0)

#!/usr/bin/env python3

import argparse
import json
import os
import pdoc
import re
import sys
import time
from datetime import datetime
from distutils.version import LooseVersion

from pyinventory import InventoryClient
from pyinventory._image import delete_file, store_file

GRAPHQL_PYINVENORY_CONTENT = \
    """// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@""" + """generated) by pyinventory, DO NOT EDIT.

package resolver

// PyinventoryConsts consists of metadata on python packages that were exported to store
const PyinventoryConsts = `
{}
`
"""

GRAPHQL_PYINVENORY_PATH = "../graph/graphql/resolver/pyinventory.go"


def module_path(m: pdoc.Module, output_dir: str, ext: str):
    return os.path.join(output_dir, *re.sub(r'\.html$', ext, m.url()).split('/'))


def write_files(m: pdoc.Module, output_dir: str, **kwargs):
    f = module_path(m, output_dir, ".html")

    dirpath = os.path.dirname(f)
    if not os.access(dirpath, os.R_OK):
        os.makedirs(dirpath)

    try:
        with open(f, 'w+', encoding='utf-8') as w:
            w.write(m.html(**kwargs))
    except Exception:
        try:
            os.unlink(f)
        except Exception:
            pass
        raise

    for submodule in m.submodules():
        write_files(submodule, output_dir, **kwargs)


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
                    delete_file(client, latestPackage["whlFileKey"], True)
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

    whlFileKey = store_file(
        client, os.path.join("./dist", whlFiles[0]), "application/zip", True
    )

    newPackage = {
        "version": version,
        "whlFileKey": whlFileKey,
        "uploadTime": datetime.isoformat(
            datetime.fromtimestamp(int(time.time()))) + "+00:00",
        "hasBreakingChange": hasBreakingChange
    }
    packages.insert(0, newPackage)

    newContent = json.dumps(packages)
    open(GRAPHQL_PYINVENORY_PATH, "w").write(
        GRAPHQL_PYINVENORY_CONTENT.format(newContent)
    )

    pyinventory_module = pdoc.Module("pyinventory")
    write_files(pyinventory_module, '../docs/website/static', show_source_code=False)


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
        args.has_breaking_change)
    sys.exit(0)

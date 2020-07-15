#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
import subprocess


BASEPATH = os.path.dirname(os.path.abspath(__file__))
SPHINX_PATH = os.path.join(BASEPATH, "sphinx/")


def export_doc():
    subprocess.run(
        ["sphinx-build", "-M", "html", ".", "../../docs/website/static/pyinventory"],
        cwd=SPHINX_PATH,
    )


if __name__ == "__main__":
    export_doc()

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import time


def main():
    """
    Dummy run-forever loop, to test process stop/start
    """
    while True:
        time.sleep(1)
    return


if __name__ == "__main__":
    main()

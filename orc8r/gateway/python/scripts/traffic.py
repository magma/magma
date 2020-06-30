#!/usr/bin/env python3

"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import os
import sys


def scan_wifi(interface):
    out = os.system('nmcli device wifi list | grep {}'.format(interface))
    if out != 0:
        raise Exception()

def send_traffic():
    pass


def main(argv):
    try:
        scan_wifi(argv[1])
    except Exception:
        print('Error interface {} not found'.format(argv[1]))
    send_traffic()


if __name__ == '__main__':
    main(sys.argv)

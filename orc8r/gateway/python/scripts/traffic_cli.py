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


def connect_to_wifi(interface, password):
    out = os.system('nmcli device wifi connect {} password {}'.format(interface, password))
    if out != 0:
        raise Exception()


def send_traffic():
    pass


def main(argv):
    iface = argv[1]
    pw = argv[2]
    try:
        scan_wifi(iface)
    except Exception:
        print('Error interface {} not found'.format(iface))
    try:
        connect_to_wifi(iface, pw)
    except Exception:
        print('Error could not connect to {}'.format(iface))
    send_traffic()


if __name__ == '__main__':
    main(sys.argv)

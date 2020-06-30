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
        raise OSError()


def connect_to_wifi(interface, password):
    out = os.system('nmcli device wifi connect {} password {}'.format(interface, password))
    if out != 0:
        raise OSError()


def send_traffic(endpt, num_pkts):
    out = os.system('ping -c {} {}'.format(num_pkts, endpt))
    if out != 0:
        raise OSError()


def main():
    argv = sys.argv
    if len(argv) < 5:
        print('Usage: ./traffic_cli.py <iface> <password> <endpoint> <num_packets_to_send>')
        sys.exit(1)
    iface = argv[1]
    pw = argv[2]
    endpt = argv[3]
    num_pkts = argv[4]
    try:
        scan_wifi(iface)
    except OSError:
        print('Error interface {} not found'.format(iface))
    try:
        connect_to_wifi(iface, pw)
    except OSError:
        print('Error could not connect to {}'.format(iface))
    try:
        send_traffic(endpt, num_pkts)
    except OSError:
        print('Error pinging {}'.format(endpt))


if __name__ == '__main__':
    main()

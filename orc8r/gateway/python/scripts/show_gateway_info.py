#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import argparse

import snowflake
from magma.common.cert_utils import load_public_key_to_base64der


def main():
    parser = argparse.ArgumentParser(
        description='Show the UUID and base64 encoded DER public key')

    parser.add_argument("--pub_key", type=str,
                        default="/var/opt/magma/certs/gw_challenge.key")
    opts = parser.parse_args()

    public_key = load_public_key_to_base64der(opts.pub_key)
    print("Hardware ID:\n------------\n%s\n" % snowflake.snowflake())
    print("Challenge Key:\n-----------\n%s" % public_key.decode('utf-8'))


if __name__ == "__main__":
    main()

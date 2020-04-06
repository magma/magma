#!/usr/bin/env python
#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import json
import os.path
import sys
from typing import Dict

import boto3

ORC8R_CERTS = [
    'rootCA.pem',
    'controller.key',
    'controller.crt',
    'certifier.key',
    'certifier.pem',
    'bootstrapper.key',
]

FLUENTD_CERTS = [
    'fluentd.key',
    'fluentd.pem',
]

ADMIN_CERTS = [
    'admin_operator.pem',
    'admin_operator.key.pem',
]

ALL_CERTS = ORC8R_CERTS + FLUENTD_CERTS + ADMIN_CERTS


def main(secret_name: str, aws_region: str, certs_dir: str):
    sec = create_orc8r_secrets(certs_dir)
    set_orc8r_secretsmanager(secret_name, aws_region, sec)


def create_orc8r_secrets(certs_dir: str) -> Dict[str, str]:
    certs_dir_abs = os.path.abspath(
        os.path.expandvars(os.path.expanduser(certs_dir)),
    )

    ret = {}
    for fname in ALL_CERTS:
        full_fpath = os.path.join(certs_dir_abs, fname)
        # admin_operator.* certs are optional
        if not os.path.isfile(full_fpath):
            if fname not in ADMIN_CERTS:
                raise ValueError(f'No cert {fname} found in certs directory')
            else:
                continue
        with open(full_fpath, 'r') as f:
            # readlines elements already have \n at the end
            ret[fname] = ''.join(f.readlines())
    return ret


def set_orc8r_secretsmanager(secret_name: str,
                             region: str,
                             secret_contents: Dict[str, str]):
    secret_string = json.dumps(secret_contents)

    session = boto3.session.Session()
    client = session.client('secretsmanager', region)
    resp = client.update_secret(
        SecretId=secret_name,
        SecretString=secret_string,
    )
    if resp['ResponseMetadata']['HTTPStatusCode'] != 200:
        raise Exception(f'Secretsmanager request failed. '
                        f'AWS Response: \n{json.dumps(resp, indent=2)}')


if __name__ == '__main__':
    # 0: script name
    # 1: secret name
    # 2: AWS region
    # 3: certs dir
    if len(sys.argv) < 4:
        print(f'Expected 3 CLI arguments, got {len(sys.argv) - 1}',
              file=sys.stderr)
        sys.exit(1)
    main(*sys.argv[1:])

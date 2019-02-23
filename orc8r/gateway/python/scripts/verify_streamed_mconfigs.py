"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import json
import sys

import os

legacy_mconfig_file = '/var/opt/magma/configs/gateway.mconfig'
new_mconfig_file = '/var/opt/magma/configs/gateway.streamed.mconfig'


def main():
    """
    Checks computed streamed mconfig against the built on pulled from leagcy
    mconfig builders and makes sure the two are equal.
    """
    if not os.path.isfile(legacy_mconfig_file):
        print('Legacy gateway.mconfig is not a file')
        return 1
    if not os.path.isfile(new_mconfig_file):
        print('New gateway.streamed.mconfig is not a file')
        return 1

    with open(legacy_mconfig_file, 'r') as old_f, open(new_mconfig_file, 'r') as new_f:
        gw_configs_json_deser = json.loads(old_f.read())
        streamed_configs_json_deser = json.loads(new_f.read())
        streamed_gw_configs = streamed_configs_json_deser.get('configs', {})

        if gw_configs_json_deser == streamed_gw_configs:
            print('Mconfigs are equal, all good!')
            return 0
        else:
            key_difference = gw_configs_json_deser['configsByKey'].keys() ^\
                             streamed_gw_configs['configsByKey'].keys()
            if key_difference:
                print('Symmetric difference of mconfig keys: {}'.format(
                    key_difference))

            for k in gw_configs_json_deser['configsByKey'].keys():
                old_v = gw_configs_json_deser['configsByKey'][k]
                new_v = streamed_gw_configs['configsByKey'][k]
                if old_v != new_v:
                    print('Values for key {} differ'.format(k))

            return 1


if __name__ == '__main__':
    sys.exit(main())

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import json
import os
import sys

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
                print(
                    'Symmetric difference of mconfig keys: {}'.format(
                    key_difference,
                    ),
                )

            for k in gw_configs_json_deser['configsByKey'].keys():
                old_v = gw_configs_json_deser['configsByKey'][k]
                new_v = streamed_gw_configs['configsByKey'][k]
                if old_v != new_v:
                    print('Values for key {} differ'.format(k))

            return 1


if __name__ == '__main__':
    sys.exit(main())

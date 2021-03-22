#!/usr/bin/python3.5

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

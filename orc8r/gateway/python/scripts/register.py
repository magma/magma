#!/usr/bin/env python3
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
import os
import sys
import textwrap
import time
from typing import List

import checkin_cli
import snowflake
from magma.common.cert_utils import load_public_key_to_base64der
from magma.common.service_registry import ServiceRegistry
from magma.magmad.bootstrap_manager import BootstrapManager
from orc8r.protos.bootstrapper_pb2 import (
    ChallengeKey,
    RegisterRequest,
    RegisterResponse,
)
from orc8r.protos.bootstrapper_pb2_grpc import RegistrationStub
from orc8r.protos.identity_pb2 import AccessGatewayID

CONFIG_OVERRIDE_DIR = '/var/opt/magma/configs'
CONTROL_PROXY = 'control_proxy.yml'


def register_handler(client: RegistrationStub, args: argparse.Namespace) -> RegisterResponse:
    """
    Register a device and retrieves its control proxy
    Args:
        client: Registration stub
        args: command line arguments
    Returns:
        RegisterRequest: register request, used for printing after function returns
        RegisterResponse: response from gRPC call, either error or the control_proxy
    """
    req = RegisterRequest(
        token=args.token,
        hwid=AccessGatewayID(
            id=snowflake.snowflake(),
        ),
        challenge_key=ChallengeKey(
            key=load_public_key_to_base64der("/var/opt/magma/certs/gw_challenge.key"),
            key_type=ChallengeKey.KeyType.SOFTWARE_ECDSA_SHA256,
        ),
    )

    res = client.Register(req)
    if res.HasField("error"):
        raise Exception(res.error)

    return req, res


def save_control_proxy(control_proxy: str):
    os.makedirs(CONFIG_OVERRIDE_DIR, exist_ok=True)
    control_proxy_path = os.path.join(CONFIG_OVERRIDE_DIR, CONTROL_PROXY)
    with open(control_proxy_path, 'w', encoding='utf-8') as f:
        f.write(control_proxy)


def main():
    """Register a gateway"""
    parser = argparse.ArgumentParser(description="Register a gateway.")
    parser.add_argument(
        "domain",
        metavar="DOMAIN_NAME",
        type=str,
        help="orc8r's domain name",
    )
    parser.add_argument(
        "token",
        metavar="REGISTRATION_TOKEN",
        type=str,
        help="registration token after API call",
    )
    parser.add_argument(
        "--ca-file",
        type=str,
        help="orc8r's root CA file",
    )
    parser.add_argument(
        "--cloud-port",
        type=str,
        help="orc8r's port",
    )
    parser.add_argument(
        "--no-control-proxy",
        action="store_true",
        help="disables writing the control proxy file",
    )
    args = parser.parse_args()

    chan = ServiceRegistry.get_bare_bootstrap_rpc_channel(
        args.domain,
        '8444' if not args.cloud_port else args.cloud_port,
        args.ca_file,
    )
    client = RegistrationStub(chan)

    try:
        req, res = register_handler(client, args)
        msg = textwrap.dedent(
            """
            > Registered gateway
            Hardware ID
            -----------
            {}
            Challenge Key
            -----------
            {}
            """,
        )
        print(msg.format(req.hwid, req.challenge_key, res.control_proxy))
        if not args.no_control_proxy:
            save_control_proxy("control_proxy", res.control_proxy)
            msg = textwrap.dedent(
                """
                Control Proxy
                -----------
                {}
                """,
            )
            print(msg.format(res.control_proxy))
    except Exception as e:
        msg = textwrap.dedent(" > Error: {} ")
        print(msg.format(e))
        sys.exit(1)

    print(
        "> Waiting {} seconds for next bootstrap".format(
            BootstrapManager.LONG_BOOTSTRAP_RETRY_INTERVAL.total_seconds(),
        ),
    )
    time.sleep(BootstrapManager.LONG_BOOTSTRAP_RETRY_INTERVAL.total_seconds())

    print("> Running checkin_cli")
    checkin_cli.main()


if __name__ == "__main__":
    main()

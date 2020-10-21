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

"""
Script to config AGW in stateful and stateless mode
"""

import argparse
import os
import subprocess
import sys
import time

from enum import Enum

from magma.configuration.service_configs import (
    load_override_config,
    load_service_config,
    save_override_config,
)

return_codes = Enum(
    "return_codes", "STATELESS STATEFUL CORRUPT INVALID", start=0
)
SERVICE_CONFIG_NAMES = [
    ("mme", "use_stateless"),
    ("mobilityd", "persist_to_redis"),
    ("pipelined", "clean_restart"),
    ("sessiond", "support_stateless"),
]


def _check_stateless_service_config(service, config_name):
    config_value = True
    if service == "pipelined":
        config_value = False  # pipelined uses inverse logic

    service_config = load_service_config(service)
    if service_config.get(config_name) == config_value:
        return return_codes.STATELESS

    print(service, "is stateful")
    return return_codes.STATEFUL


def _check_stateless_services():
    num_stateful = 0
    for service, config in SERVICE_CONFIG_NAMES:
        if (
            _check_stateless_service_config(service, config)
            == return_codes.STATEFUL
        ):
            num_stateful += 1

    if num_stateful == 0:
        print("Check returning", return_codes.STATELESS)
        return return_codes.STATELESS
    elif num_stateful == len(SERVICE_CONFIG_NAMES):
        print("Check returning", return_codes.STATEFUL)
        return return_codes.STATEFUL

    print("Check returning", return_codes.CORRUPT)
    return return_codes.CORRUPT


def check_stateless_agw():
    sys.exit(_check_stateless_services().value)


def _clear_redis_state():
    if os.getuid() != 0:
        print("Need to run as root to clear Redis state.")
        sys.exit(return_codes.INVALID)
    subprocess.call("service magma@* stop".split())
    subprocess.call("service magma@redis start".split())
    subprocess.call("redis-cli -p 6380 FLUSHALL".split())
    subprocess.call("service magma@redis stop".split())


def _start_magmad():
    if os.getuid() != 0:
        print("Need to run as root to start magmad.")
        sys.exit(return_codes.INVALID)
    subprocess.call("service magma@magmad start".split())


def _restart_sctpd():
    if os.getuid() != 0:
        print("Need to run as root to restart sctpd.")
        sys.exit(return_codes.INVALID)
    print("Restarting sctpd")
    subprocess.call("service sctpd restart".split())


def enable_stateless_agw():
    if _check_stateless_services() == return_codes.STATELESS:
        print("Nothing to enable, AGW is stateless")
        sys.exit(return_codes.STATELESS.value)
    for service, config in SERVICE_CONFIG_NAMES:
        cfg = load_override_config(service) or {}
        if service == "pipelined":
            cfg[config] = False
        else:
            cfg[config] = True

        save_override_config(service, cfg)

    # restart Sctpd so that eNB connections are reset and local state cleared
    _restart_sctpd()
    sys.exit(_check_stateless_services().value)


def disable_stateless_agw():
    if _check_stateless_services() == return_codes.STATEFUL:
        print("Nothing to disable, AGW is stateful")
        sys.exit(return_codes.STATEFUL.value)
    for service, config in SERVICE_CONFIG_NAMES:
        cfg = load_override_config(service)
        if cfg is None:
            cfg = {}

        # remove the stateless override
        cfg.pop(config, None)

        save_override_config(service, cfg)

    # restart Sctpd so that eNB connections are reset and local state cleared
    _restart_sctpd()
    sys.exit(_check_stateless_services().value)


def sctpd_pre_start():
    if _check_stateless_services() == return_codes.STATEFUL:
        # switching from stateless to stateful
        print("AGW is stateful, nothing to be done")
    else:
        _clear_redis_state()
    sys.exit(0)


def sctpd_post_start():
    _start_magmad()
    time.sleep(15)  # sleep for a bit to ensure OVS and Magma services are up
    sys.exit(0)


def clear_redis_and_restart():
    _clear_redis_state()
    _start_magmad()
    sys.exit(0)


STATELESS_FUNC_DICT = {
    "check": check_stateless_agw,
    "enable": enable_stateless_agw,
    "disable": disable_stateless_agw,
    "sctpd_pre": sctpd_pre_start,
    "sctpd_post": sctpd_post_start,
    "clear_redis": clear_redis_and_restart,
}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=STATELESS_FUNC_DICT.keys())

    args = parser.parse_args()

    func = STATELESS_FUNC_DICT[args.command]
    func()


if __name__ == "__main__":
    main()

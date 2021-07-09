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

Script to trigger pre and post start commands for the Sctpd systemd unit
"""

import argparse
import os
import subprocess
import sys
import time
from enum import Enum

from magma.common.redis.client import get_default_client
from magma.configuration.service_configs import (
    load_override_config,
    load_service_config,
    save_override_config,
)

return_codes = Enum(
    "return_codes", "STATELESS STATEFUL CORRUPT INVALID", start=0
)
STATELESS_SERVICE_CONFIGS = [
    ("mme", "use_stateless", True),
    ("mobilityd", "persist_to_redis", True),
    ("pipelined", "clean_restart", False),
    ("pipelined", "redis_enabled", True),
    ("sessiond", "support_stateless", True),
]


def check_stateless_service_config(service, config_name, config_value):
    service_config = load_service_config(service)
    if service_config.get(config_name) == config_value:
        print("STATELESS\t%s -> %s" % (service, config_name))
        return return_codes.STATELESS

    print("STATEFUL\t%s -> %s" % (service, config_name))
    return return_codes.STATEFUL


def check_stateless_services():
    num_stateful = 0
    for service, config, value in STATELESS_SERVICE_CONFIGS:
        if (
                check_stateless_service_config(service, config, value)
                == return_codes.STATEFUL
        ):
            num_stateful += 1

    if num_stateful == 0:
        res = return_codes.STATELESS
    elif num_stateful == len(STATELESS_SERVICE_CONFIGS):
        res = return_codes.STATEFUL
    else:
        res = return_codes.CORRUPT

    print("Check returning", res)
    return res


def check_stateless_agw():
    sys.exit(check_stateless_services().value)


def clear_redis_state():
    if os.getuid() != 0:
        print("Need to run as root to clear Redis state.")
        sys.exit(return_codes.INVALID)
    # stop MME, which in turn stops mobilityd, pipelined and sessiond
    subprocess.call("service magma@mme stop".split())
    # delete all keys from Redis which capture service state
    redis_client = get_default_client()
    for key_regex in [
        "*_state",
        "IMSI*MME",
        "IMSI*S1AP",
        "IMSI*SPGW",
        "IMSI*mobilityd*",
        "mobilityd:assigned_ip_blocks",
        "mobilityd:ip_states:*",
        "NO_VLAN:mobilityd_gw_info",
        "QosManager",
        "s1ap_imsi_map",
        "sessiond:sessions",
        "*pipelined:rule_ids",
        "*pipelined:rule_versions",
        "*pipelined:rule_names",
    ]:
        for key in redis_client.scan_iter(key_regex):
            redis_client.delete(key)


def flushall_redis():
    if os.getuid() != 0:
        print("Need to run as root to clear Redis state.")
        sys.exit(return_codes.INVALID)
    print("Flushing all content in Redis")
    subprocess.call("service magma@* stop".split())
    subprocess.call("service magma@redis start".split())
    subprocess.call("redis-cli -p 6380 flushall".split())
    subprocess.call("service magma@redis stop".split())


def start_magmad():
    if os.getuid() != 0:
        print("Need to run as root to start magmad.")
        sys.exit(return_codes.INVALID)
    subprocess.call("service magma@magmad start".split())


def restart_sctpd():
    if os.getuid() != 0:
        print("Need to run as root to restart sctpd.")
        sys.exit(return_codes.INVALID)
    print("Restarting sctpd")
    subprocess.call("service sctpd restart".split())
    # delay return after restarting so that Magma and OVS services come up
    time.sleep(30)


def enable_stateless_agw():
    if check_stateless_services() == return_codes.STATELESS:
        print("Nothing to enable, AGW is stateless")
        sys.exit(return_codes.STATELESS.value)
    for service, config, value in STATELESS_SERVICE_CONFIGS:
        cfg = load_override_config(service) or {}
        cfg[config] = value
        save_override_config(service, cfg)

    # restart Sctpd so that eNB connections are reset and local state cleared
    restart_sctpd()
    sys.exit(check_stateless_services().value)


def disable_stateless_agw():
    if check_stateless_services() == return_codes.STATEFUL:
        print("Nothing to disable, AGW is stateful")
        sys.exit(return_codes.STATEFUL.value)
    for service, config, value in STATELESS_SERVICE_CONFIGS:
        cfg = load_override_config(service) or {}
        cfg[config] = not value
        save_override_config(service, cfg)

    # restart Sctpd so that eNB connections are reset and local state cleared
    restart_sctpd()
    sys.exit(check_stateless_services().value)


def ovs_reset_bridges():
    subprocess.call(
        "ovs-vsctl --all destroy Flow_Sample_Collector_Set".split())
    subprocess.call("ifdown uplink_br0".split())
    subprocess.call("ifdown gtp_br0".split())
    subprocess.call("ifdown patch-up".split())
    subprocess.call("service openvswitch-switch restart".split())
    subprocess.call("ifup uplink_br0".split())
    subprocess.call("ifup gtp_br0".split())
    subprocess.call("ifup patch-up".split())


def sctpd_pre_start():
    subprocess.Popen("service procps restart".split())

    if check_stateless_services() == return_codes.STATEFUL:
        # switching from stateless to stateful
        print("AGW is stateful, nothing to be done")
    else:
        # Clean up all mobilityd, MME, pipelined and sessiond Redis keys
        clear_redis_state()
        # Clean up OVS flows
        ovs_reset_bridges()
    sys.exit(0)


def sctpd_post_start():
    subprocess.Popen("/bin/systemctl start magma@mme".split())
    subprocess.Popen("/bin/systemctl start magma@pipelined".split())
    subprocess.Popen("/bin/systemctl start magma@envoy_controller".split())
    subprocess.Popen("/bin/systemctl start magma@sessiond".split())
    subprocess.Popen("/bin/systemctl start magma@mobilityd".split())
    sys.exit(0)


def clear_redis_and_restart():
    clear_redis_state()
    sctpd_post_start()
    sys.exit(0)


def flushall_redis_and_restart():
    flushall_redis()
    start_magmad()
    restart_sctpd()
    sys.exit(0)


def reset_sctpd_for_stateful():
    if check_stateless_services() == return_codes.STATELESS:
        print("AGW is stateless, no need to restart Sctpd")
        sys.exit(0)
    restart_sctpd()


STATELESS_FUNC_DICT = {
    "check": check_stateless_agw,
    "enable": enable_stateless_agw,
    "disable": disable_stateless_agw,
    "sctpd_pre": sctpd_pre_start,
    "sctpd_post": sctpd_post_start,
    "clear_redis": clear_redis_and_restart,
    "flushall_redis": flushall_redis_and_restart,
    "reset_sctpd_for_stateful": reset_sctpd_for_stateful
}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=STATELESS_FUNC_DICT.keys())

    args = parser.parse_args()

    func = STATELESS_FUNC_DICT[args.command]
    func()


if __name__ == "__main__":
    main()

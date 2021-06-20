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
# pylint: skip-file
# pylint does not play well with aioeventlet, as it uses asyncio.async which
# produces a parse error

import asyncio
import logging
import threading

import aioeventlet
from lte.protos.mconfig import mconfigs_pb2
from magma.common.misc_utils import call_process, get_ip_from_if
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.configuration import environment
from magma.pipelined.app import of_rest_server
from magma.pipelined.app.he import PROXY_PORT_NAME
from magma.pipelined.app.uplink_bridge import UPLINK_OVS_BRIDGE_NAME
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.check_quota_server import run_flask
from magma.pipelined.gtp_stats_collector import (
    MIN_OVSDB_DUMP_POLLING_INTERVAL,
    GTPStatsCollector,
)
from magma.pipelined.ifaces import monitor_ifaces
from magma.pipelined.rpc_servicer import PipelinedRpcServicer
from magma.pipelined.service_manager import ServiceManager
from ryu import cfg
from ryu.base.app_manager import AppManager
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL
from scapy.arch import get_if_hwaddr
from magma.pipelined.datapath_setup import tune_datapath


def main():
    """
    Loads the Ryu apps we want to run from the config file.
    This should exit on keyboard interrupt.
    """

    # Run asyncio loop in a greenthread so we can evaluate other eventlets
    # TODO: Remove once Ryu migrates to asyncio
    asyncio.set_event_loop_policy(aioeventlet.EventLoopPolicy())

    service = MagmaService('pipelined', mconfigs_pb2.PipelineD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name)

    service_config = service.config

    if environment.is_dev_mode():
        of_rest_server.configure(service_config)

    # Set Ryu config params
    cfg.CONF.ofp_listen_host = "127.0.0.1"

    # override mconfig using local config.
    # TODO: move config compilation to separate module.
    enable_nat = service.config.get('enable_nat', service.mconfig.nat_enabled)
    service.config['enable_nat'] = enable_nat
    logging.info("Nat: %s", enable_nat)
    vlan_tag = service.config.get(
        'sgi_management_iface_vlan',
        service.mconfig.sgi_management_iface_vlan,
    )
    service.config['sgi_management_iface_vlan'] = vlan_tag

    sgi_ip = service.config.get(
        'sgi_management_iface_ip_addr',
        service.mconfig.sgi_management_iface_ip_addr,
    )
    service.config['sgi_management_iface_ip_addr'] = sgi_ip

    sgi_gateway_ip = service.config.get(
        'sgi_management_iface_gw',
        service.mconfig.sgi_management_iface_gw,
    )
    service.config['sgi_management_iface_gw'] = sgi_gateway_ip

    # Keep router mode off for smooth upgrade path
    service.config['dp_router_enabled'] = service.config.get(
        'dp_router_enabled',
        False,
    )
    if 'virtual_mac' not in service.config:
        if service.config['dp_router_enabled']:
            up_bridge_name = service.config.get(
                'uplink_bridge', UPLINK_OVS_BRIDGE_NAME,
            )
            mac_addr = get_if_hwaddr(up_bridge_name)
        else:
            mac_addr = get_if_hwaddr(service.config.get('bridge_name'))

        service.config['virtual_mac'] = mac_addr

    # this is not read from yml file.
    service.config['uplink_port'] = OFPP_LOCAL
    uplink_port_name = service.config.get('ovs_uplink_port_name', None)
    if enable_nat is False and uplink_port_name is not None:
        service.config['uplink_port'] = BridgeTools.get_ofport(
            uplink_port_name,
        )

    # header enrichment related configuration.
    service.config['proxy_port_name'] = PROXY_PORT_NAME
    he_enabled_flag = False
    if service.mconfig.he_config:
        he_enabled_flag = service.mconfig.he_config.enable_header_enrichment
    he_enabled = service.config.get('he_enabled', he_enabled_flag)
    service.config['he_enabled'] = he_enabled

    # tune datapath according to config
    tune_datapath(service.config)

    # monitoring related configuration
    mtr_interface = service.config.get('mtr_interface', None)
    if mtr_interface:
        mtr_ip = get_ip_from_if(mtr_interface)
        service.config['mtr_ip'] = mtr_ip

    # Load the ryu apps
    service_manager = ServiceManager(service)
    service_manager.load()

    def callback(returncode):
        if returncode != 0:
            logging.error(
                "Failed to set MASQUERADE: %d", returncode,
            )

    # TODO fix this hack for XWF
    if enable_nat is True or service.config.get('setup_type') == 'XWF':
        ip_table_rule = 'POSTROUTING -o %s -j MASQUERADE' % service.config['nat_iface']
        check_and_add = 'iptables -t nat -C %s || iptables -t nat -A %s' % \
                (ip_table_rule, ip_table_rule)
        logging.debug("check_and_add: %s", check_and_add)
        call_process(check_and_add,
            callback,
            service.loop,
        )

    service.loop.create_task(
        monitor_ifaces(
            service.config['monitored_ifaces'],
        ),
    )

    manager = AppManager.get_instance()
    # Add pipelined rpc servicer
    pipelined_srv = PipelinedRpcServicer(
        service.loop,
        manager.applications.get('GYController', None),
        manager.applications.get('EnforcementController', None),
        manager.applications.get('EnforcementStatsController', None),
        manager.applications.get('DPIController', None),
        manager.applications.get('UEMacAddressController', None),
        manager.applications.get('CheckQuotaController', None),
        manager.applications.get('IPFIXController', None),
        manager.applications.get('VlanLearnController', None),
        manager.applications.get('TunnelLearnController', None),
        manager.applications.get('Classifier', None),
        manager.applications.get('InOutController', None),
        manager.applications.get('NGServiceController', None),
        service.config,
        service_manager,
    )
    pipelined_srv.add_to_server(service.rpc_server)

    if service.config['setup_type'] == 'CWF':
        bridge_ip = service.config['bridge_ip_address']
        has_quota_port = service.config['has_quota_port']
        no_quota_port = service.config['no_quota_port']

        def on_exit_server_thread():
            service.StopService(None, None)

        # For CWF start quota check servers
        start_check_quota_server(
            run_flask, bridge_ip, has_quota_port, True,
            on_exit_server_thread,
        )
        start_check_quota_server(
            run_flask, bridge_ip, no_quota_port, False,
            on_exit_server_thread,
        )

    if service.config['setup_type'] == 'LTE':
        polling_interval = service.config.get(
            'ovs_gtp_stats_polling_interval',
            MIN_OVSDB_DUMP_POLLING_INTERVAL,
        )
        collector = GTPStatsCollector(
            polling_interval,
            service.loop,
        )
        collector.start()

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def start_check_quota_server(target, ip, port, response, exit_callback):
    """ Starts service server threads """
    thread = threading.Thread(
        target=target,
        args=(ip, port, response, exit_callback),
    )
    thread.daemon = True
    thread.start()


if __name__ == "__main__":
    main()

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
import subprocess
from typing import Dict, NamedTuple

from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from ryu.controller.controller import Datapath


class IPFIXController(MagmaController):
    """
    IPFIXController

    The IPFIX controller installs sample flows for exporting IPFIX statistics
    to the controller. Each imsi gets a sample flow. After sampling traffic is
    forwarded to the next table.
    """

    APP_NAME = "ipfix"
    APP_TYPE = ControllerType.LOGICAL

    IPFIXConfig = NamedTuple(
        'IPFIXConfig',
        [('enabled', bool), ('collector_ip', str), ('collector_port', int),
         ('probability', int), ('collector_set_id', int),
         ('obs_domain_id', int), ('obs_point_id', int), ('cache_timeout', int),
         ('sampling_port', int)],
    )

    def __init__(self, *args, **kwargs):
        super(IPFIXController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self._app_set_tbl_num = self._service_manager.INTERNAL_APP_SET_TABLE_NUM
        self._imsi_set_tbl_num = \
            self._service_manager.INTERNAL_IMSI_SET_TABLE_NUM
        self._dpi_enabled = kwargs['config']['dpi']['enabled']
        self._bridge_name = kwargs['config']['bridge_name']
        self._conntrackd_enabled = kwargs['config']['conntrackd']['enabled']
        self.ipfix_config = self._get_ipfix_config(kwargs['config'],
                                                   kwargs['mconfig'])
        # If DPI enabled don't sample normal traffic, sample only internal pkts
        if self._dpi_enabled or self._conntrackd_enabled:
            self._ipfix_sample_tbl_num = \
                self._service_manager.INTERNAL_IPFIX_SAMPLE_TABLE_NUM
        else:
            self._ipfix_sample_tbl_num = self.tbl_num
        self._datapath = None

    def _get_ipfix_config(self, config_dict: Dict,
                          mconfig) -> NamedTuple:
        if 'ipfix' not in config_dict or not config_dict['ipfix']['enabled']:
            return self.IPFIXConfig(enabled=False, probability=0,
                collector_ip='', collector_port=0, collector_set_id=0,
                obs_domain_id=0, obs_point_id=0, cache_timeout=0,
                sampling_port=0)
        collector_ip = mconfig.ipdr_export_dst.ip
        collector_port = mconfig.ipdr_export_dst.port
        if not mconfig.ipdr_export_dst.ip:
            if 'collector_ip' in config_dict['ipfix']:
                self.logger.error("Missing IPDR dest IP, using val from .yml")
                collector_ip = config_dict['ipfix']['collector_ip']
                collector_port = config_dict['ipfix']['collector_port']
            else:
                self.logger.error("Missing mconfig IPDR dest IP")
                return self.IPFIXConfig(enabled=False, probability=0,
                    collector_ip='', collector_port=0, collector_set_id=0,
                    obs_domain_id=0, obs_point_id=0, cache_timeout=0,
                    sampling_port=0)

        if collector_port == 0:
            self.logger.error("Missing mconfig IPDR dest port")
            return self.IPFIXConfig(enabled=False, probability=0,
                collector_ip='', collector_port=0, collector_set_id=0,
                obs_domain_id=0, obs_point_id=0, cache_timeout=0,
                sampling_port=0)

        if self._dpi_enabled or self._conntrackd_enabled:
            probability = 65535
        else:
            probability = config_dict['ipfix']['probability']

        return self.IPFIXConfig(
            enabled=config_dict['ipfix']['enabled'],
            collector_ip=collector_ip,
            collector_port=collector_port,
            probability=probability,
            collector_set_id=config_dict['ipfix']['collector_set_id'],
            obs_domain_id=config_dict['ipfix']['obs_domain_id'],
            obs_point_id=config_dict['ipfix']['obs_point_id'],
            cache_timeout=config_dict['ipfix']['cache_timeout'],
            sampling_port=config_dict['ovs_gtp_port_number']
        )

    def initialize_on_connect(self, datapath: Datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self._delete_all_flows(datapath)
        self._install_default_flows(datapath)

        if not self.ipfix_config.enabled:
            return

        action_str = (
            'ovs-vsctl -- --id=@{} get Bridge {} -- --id=@cs create '
            'Flow_Sample_Collector_Set id={} bridge=@{} ipfix=@i -- --id=@i '
            'create IPFIX targets=\"{}\\:{}\" obs_domain_id={} obs_point_id={} '
            'cache_active_timeout={}, other_config:enable-tunnel-sampling=false'
        ).format(
            self._bridge_name, self._bridge_name,
            self.ipfix_config.collector_set_id, self._bridge_name,
            self.ipfix_config.collector_ip, self.ipfix_config.collector_port,
            self.ipfix_config.obs_domain_id, self.ipfix_config.obs_point_id,
            self.ipfix_config.cache_timeout
        )
        try:
            p = subprocess.Popen(action_str, shell=True,
                                 stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            _, err = p.communicate()
            err_str = err.decode('utf-8')
            if err_str:
                self.logger.error("Failed setting up ipfix sampling %s",
                                  err_str)
        except subprocess.CalledProcessError as e:
            raise Exception('Error: {} failed with: {}'.format(action_str, e))

    def cleanup_on_disconnect(self, datapath: Datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self._delete_all_flows(datapath)

    def _delete_all_flows(self, datapath: Datapath) -> None:
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._app_set_tbl_num,
                                          cookie=self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._imsi_set_tbl_num,
                                          cookie=self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._ipfix_sample_tbl_num)

    def _install_default_flows(self, datapath: Datapath) -> None:
        """
        For each direction set the default flows to just forward to next app.

        Args:
            datapath: ryu datapath struct
        """
        parser = self._datapath.ofproto_parser
        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        if not self._service_manager.is_app_enabled(DPIController.APP_NAME):
            flows.add_resubmit_next_service_flow(
                self._datapath, self._app_set_tbl_num, MagmaMatch(),
                priority=flows.MINIMUM_PRIORITY, cookie=self.tbl_num,
                resubmit_table=self._imsi_set_tbl_num)

        flows.add_resubmit_next_service_flow(
            self._datapath, self._imsi_set_tbl_num, MagmaMatch(),
            priority=flows.MINIMUM_PRIORITY, cookie=self.tbl_num,
            resubmit_table=self._ipfix_sample_tbl_num)

        if self.ipfix_config.enabled and (self._dpi_enabled or self._conntrackd_enabled):
            pdp = 1
            actions = [parser.NXActionSample2(
                probability=self.ipfix_config.probability,
                collector_set_id=self.ipfix_config.collector_set_id,
                obs_domain_id=self.ipfix_config.obs_domain_id,
                obs_point_id=self.ipfix_config.obs_point_id,
                apn_mac_addr=[0, 0, 0, 0, 0, 0],
                msisdn="defau".encode('ascii'),
                apn_name="default".encode('ascii'),
                pdp_start_epoch=pdp.to_bytes(8, byteorder='little'),
                sampling_port=self.ipfix_config.sampling_port)]
            flows.add_drop_flow(
                self._datapath, self._ipfix_sample_tbl_num, match, actions,
                priority=flows.DEFAULT_PRIORITY)

    def add_ue_sample_flow(self, imsi: str, msisdn: str,
                           apn_mac_addr: str, apn_name: str,
                           pdp_start_time: int) -> None:
        """
        Install a flow to sample packets for IPFIX for specific imsi

        Args:
            imsi (string): subscriber to install rule for
            msisdn (string): msisdn string
            apn_mac_addr (string): AP mac address string
            apn_name (string): AP name
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized for adding flows')
            return

        if not self.ipfix_config.enabled:
            #TODO logging higher than debug here will provide too much noise
            # possible fix is making ipfix a dynamic service enabled from orc8r
            self.logger.debug('IPFIX export dst not setup for adding flows')
            return

        parser = self._datapath.ofproto_parser
        if not apn_mac_addr or '-' not in apn_mac_addr:
            apn_mac_bytes = [0, 0, 0, 0, 0, 0]
        else:
            apn_mac_bytes = [int(a, 16) for a in apn_mac_addr.split('-')]

        if not msisdn:
            msisdn='no_msisdn'

        actions = [parser.NXActionSample2(
            probability=self.ipfix_config.probability,
            collector_set_id=self.ipfix_config.collector_set_id,
            obs_domain_id=self.ipfix_config.obs_domain_id,
            obs_point_id=self.ipfix_config.obs_point_id,
            apn_mac_addr=apn_mac_bytes,
            msisdn=msisdn.encode('ascii'),
            apn_name=apn_name.encode('ascii'),
            pdp_start_epoch=pdp_start_time.to_bytes(8, byteorder='little'),
            sampling_port=self.ipfix_config.sampling_port)]

        match = MagmaMatch(imsi=encode_imsi(imsi))
        if self._dpi_enabled or self._conntrackd_enabled:
            flows.add_drop_flow(
                self._datapath, self._ipfix_sample_tbl_num, match, actions,
                priority=flows.UE_FLOW_PRIORITY)
        else:
            flows.add_resubmit_next_service_flow(
                self._datapath, self._ipfix_sample_tbl_num, match, actions,
                priority=flows.UE_FLOW_PRIORITY,
                resubmit_table=self.next_main_table)

    def delete_ue_sample_flow(self, imsi: str) -> None:
        """
        Delete a flow to sample packets for IPFIX for specific imsi

        Args:
            imsi (string): subscriber to install rule for
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized')
            return

        if not imsi:
            self.logger.error('No subscriber specified')
            return

        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self._ipfix_sample_tbl_num, match)

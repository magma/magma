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
import logging
import shlex
import subprocess

from lte.protos.pipelined_pb2 import FlowRequest
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import DPI_REG
from magma.pipelined.policy_converters import (
    FlowMatchError,
    flip_flow_match,
    flow_match_to_magma_match,
)

# TODO might move to config file
# Current classification will finalize if found in APP_PROTOS, if found in
# PARENT_PROTOS we will also add the SERVICE_IDS id to the final classification
PARENT_PROTOS = {"facebook": 10, "google_gen": 20, "viber": 30, "imo": 40}
APP_PROTOS = {"facebook_messenger": 1, "instagram": 2, "youtube": 3,
              "gmail": 4, "google_docs": 5, "netflix": 6,
              "apple": 7, "microsoft": 8, 'reddit': 9, 'whatsapp': 101,
              "google_play": 102, "appstore": 103, "amazon": 104, "wechat": 105,
              "tiktok": 106, "twitter": 107, "wikipedia": 108, "yahoo": 109}
SERVICE_IDS = {"other": 0, "chat": 1, "audio": 2, "video": 3}
DEFAULT_DPI_ID = 0
# Max register value
UNCLASSIFIED_PROTO_ID = 0xFFFFFFFF

LOG = logging.getLogger('pipelined.app.dpi')


class DPIController(MagmaController):
    """
    DPI controller.

    The DPI controller is responsible for marking a flow with an App ID derived
    from DPI. The APP ID should be stored in register 3
    """

    APP_NAME = "dpi"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(DPIController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.setup_type = kwargs['config']['setup_type']
        self._datapath = None
        self._dpi_enabled = kwargs['config']['dpi']['enabled']
        self._mon_port = kwargs['config']['dpi']['mon_port']
        self._mon_port_number = kwargs['config']['dpi']['mon_port_number']
        self._idle_timeout = kwargs['config']['dpi']['idle_timeout']
        self._bridge_name = kwargs['config']['bridge_name']
        self._classify_app_tbl_num = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self._app_set_tbl_num = self._service_manager.INTERNAL_APP_SET_TABLE_NUM
        self._imsi_set_tbl_num = \
            self._service_manager.INTERNAL_IMSI_SET_TABLE_NUM
        if self._dpi_enabled:
            self._create_monitor_port()

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._app_set_tbl_num,
                                          cookie=self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._classify_app_tbl_num)

    def add_classify_flow(self, flow_match, flow_state, app: str,
                          service_type: str):
        """
        Parse DPI output and set the register for future packets matching this
        flow. APP is split into tokens as the top level app is not supported,
        but the parent protocol might be.
        Example we care about google traffic, but don't neccessarily want to
        classify every specific google service.
        """
        # TODO add error return
        if self._datapath is None:
            return
        parser = self._datapath.ofproto_parser

        app_id = get_app_id(app, service_type)

        try:
            ul_match = flow_match_to_magma_match(flow_match)
            ul_match.direction = None
            dl_match = flow_match_to_magma_match(flip_flow_match(flow_match))
            dl_match.direction = None
        except FlowMatchError as e:
            self.logger.error(e)
            return

        actions = [parser.NXActionRegLoad2(dst=DPI_REG, value=app_id)]
        # No reason to create a flow here
        if flow_state != FlowRequest.FLOW_CREATED:
            flows.add_flow(self._datapath, self._classify_app_tbl_num,
                ul_match, actions, priority=flows.DEFAULT_PRIORITY,
                idle_timeout=self._idle_timeout)
            flows.add_flow(self._datapath, self._classify_app_tbl_num,
                dl_match, actions, priority=flows.DEFAULT_PRIORITY,
                idle_timeout=self._idle_timeout)

    def remove_classify_flow(self, flow_match):
        try:
            ul_match = flow_match_to_magma_match(flow_match)
            ul_match.direction = None
            dl_match = flow_match_to_magma_match(flip_flow_match(flow_match))
            dl_match.direction = None
        except FlowMatchError as e:
            self.logger.error(e)
            return False

        flows.delete_flow(self._datapath, self._classify_app_tbl_num, ul_match)

        return True

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next table.
        The policies for each subscriber would be added when the IP session is
        created, by reaching out to the controller/PCRF.

        Args:
            datapath: ryu datapath struct
        """
        parser = self._datapath.ofproto_parser

        # Setup flows to classify & mirror to sampling port
        match = MagmaMatch()
        actions = [
            parser.NXActionResubmitTable(table_id=self._classify_app_tbl_num)]

        if self._dpi_enabled:
            actions.append(parser.OFPActionOutput(self._mon_port_number))

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

        # Setup flows for internal IPFIX sampling
        actions = [
            parser.NXActionResubmitTable(table_id=self._classify_app_tbl_num)]

        flows.add_resubmit_next_service_flow(
            self._datapath, self._app_set_tbl_num, MagmaMatch(), actions,
            priority=flows.MINIMUM_PRIORITY, cookie=self.tbl_num,
            resubmit_table=self._imsi_set_tbl_num)

        # Setup flows for the application reg classifier tbl
        actions = [parser.NXActionRegLoad2(dst=DPI_REG,
                                           value=UNCLASSIFIED_PROTO_ID)]
        flows.add_flow(datapath, self._classify_app_tbl_num, MagmaMatch(),
                       actions, priority=flows.MINIMUM_PRIORITY)

    def _create_monitor_port(self):
        """
        For cwf we set this up when running docker compose as we can't modify
        interfaces from inside the container

        For lte just add the port.
        """
        if self.setup_type == 'CWF':
            self._mon_port_number = BridgeTools.get_ofport(self._mon_port)
            return

        add_cmd = "sudo ovs-vsctl add-port {} mon1 -- set interface {} \
            ofport_request={} type=internal" \
            .format(self._bridge_name, self._mon_port, self._mon_port_number)

        args = shlex.split(add_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Created monitor port ret %d", ret)

        enable_cmd = "sudo ifconfig {} up".format(self._mon_port)
        args = shlex.split(enable_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Enabled monitor port ret %d", ret)


def get_app_id(app: str, service_type: str) -> int:
    """
    Classify the app/service_type to a numeric identifier to export
    """
    if not app or not service_type:
        return DEFAULT_DPI_ID

    app = app.lower()
    service_type = service_type.lower()
    tokens = app.split('.')
    app_match = [app for app in tokens if app in APP_PROTOS]
    if len(app_match) > 1:
        LOG.warning("Found more than 1 app match in %s", app)
        return DEFAULT_DPI_ID

    if (len(app_match) == 1):
        app_id = APP_PROTOS[app_match[0]]
        LOG.debug("Classified %s-%s as %d", app, service_type,
                            app_id)
        return app_id
    parent_match = [app for app in tokens if app in PARENT_PROTOS]

    # This shoudn't happen as we confirmed the match exists
    if len(parent_match) == 0:
        LOG.debug("Didn't find a match for app name %s", app)
        return DEFAULT_DPI_ID
    if len(parent_match) > 1:
        LOG.debug("Found more than 1 parent app match in %s", app)
        return DEFAULT_DPI_ID
    app_id = PARENT_PROTOS[parent_match[0]]

    service_id = SERVICE_IDS['other']
    for serv in SERVICE_IDS:
        if serv in service_type:
            service_id = SERVICE_IDS[serv]
            break
    app_id += service_id
    LOG.debug("Classified %s-%s as %d", app, service_type, app_id)
    return app_id

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
import time
from enum import Enum

from lte.protos.pipelined_pb2 import SetupFlowsResult
from magma.pipelined.bridge_util import BridgeTools, DatapathLookupError
from magma.pipelined.metrics import OPENFLOW_ERROR_MSG
from magma.pipelined.openflow.exceptions import MagmaOFError
from ryu import utils
from ryu.base import app_manager
from ryu.controller import dpset, ofp_event
from ryu.controller.handler import (
    CONFIG_DISPATCHER,
    HANDSHAKE_DISPATCHER,
    MAIN_DISPATCHER,
    set_ev_cls,
)
from ryu.ofproto import ofproto_v1_4

global_epoch = int(time.time())


class ControllerType(Enum):
    PHYSICAL = 1
    LOGICAL = 2
    SPECIAL = 3


class ControllerNotReadyException(Exception):
    pass


class MagmaController(app_manager.RyuApp):
    """
    The base class for all MagmaControllers. Does not itself manage any tables,
    but instead handles shared state for subclass controllers.

    Applications should subclass this and can own some number of tables to
    implement their own logic.
    """
    # Inherited from RyuApp base class
    OFP_VERSIONS = [ofproto_v1_4.OFP_VERSION]

    # App name that should be overridden by the controller implementation
    APP_NAME = ""

    def __init__(self, service_manager, *args, **kwargs):
        """ Try to lookup the datapath_id of the bridge to run the app on """
        super(MagmaController, self).__init__(*args, **kwargs)
        self._app_futures = kwargs['app_futures']
        try:
            self._datapath_id = BridgeTools.get_datapath_id(
                kwargs['config']['bridge_name']
            )
        except DatapathLookupError as e:
            self.logger.error(
                'Exception in %s contoller: %s', self.APP_NAME, e)
            raise
        if 'controller_port' in kwargs['config']:
            self.CONF.ofp_tcp_listen_port = kwargs['config']['controller_port']
        self._service_manager = service_manager
        self._startup_flow_controller = None
        self._startup_flows_fut = kwargs['app_futures']['startup_flows']
        self.init_finished = False

    @set_ev_cls(ofp_event.EventOFPErrorMsg,
                [HANDSHAKE_DISPATCHER, CONFIG_DISPATCHER, MAIN_DISPATCHER])
    def record_of_errors(self, ev):
        msg = ev.msg
        self.logger.error("OF Error: type=0x%02x code=0x%02x "
                          "message=%s",
                          msg.type, msg.code, utils.hex_array(msg.data))
        OPENFLOW_ERROR_MSG.labels(
            error_type="0x%02x" % msg.type,
            error_code="0x%02x" % msg.code).inc()

    @set_ev_cls(dpset.EventDP, MAIN_DISPATCHER)
    def datapath_event_handler(self, ev):
        """
        This event handler is called on datapath connect and disconnect
        Check datapath_id in case of multiple bridges

        Args:
            ev (dpset.EventDP):  ryu event for connect/disconnect
        """
        datapath = ev.dp

        if self._datapath_id != datapath.id:
            return

        try:
            if ev.enter:
                self.initialize_on_connect(datapath)
                # set a barrier to ensure things are applied
                if self.APP_NAME in self._app_futures:
                    self._app_futures[self.APP_NAME].set_result(self)
            else:
                self.cleanup_on_disconnect(datapath)
        except MagmaOFError as e:
            act = 'initializing' if ev.enter else 'cleaning'
            self.logger.error(
                'Error %s %s flow rules: %s', act, self.APP_NAME, e)

    def check_setup_request_epoch(self, epoch):
        """
        Check if the controller is ready to be initialized after restart,
        returns:    status code if epoch is invalid/controller is initialized
                    None if controller can be initialized
        """
        self.logger.info("Received Setup request with epoch - %d, current "
                         "epoch  is - %d", epoch, global_epoch)
        if epoch != global_epoch:
            self.logger.warning(
                "Received SetupFlowsRequest has outdated epoch - %d, current "
                "epoch is - %d.", epoch, global_epoch)
            return SetupFlowsResult.OUTDATED_EPOCH

        if self._datapath is None:
            self.logger.warning("Datapath not initilized, setup failed")
            return SetupFlowsResult.FAILURE

        if self.init_finished:
            self.logger.warning('Controller already initialized, ignoring')
            return SetupFlowsResult.SUCCESS

        return None

    def is_controller_ready(self):
        """
        Check if the controller is setup & ready to process requests
        """
        return self._datapath and self.init_finished

    def initialize_on_connect(self, datapath):
        """
        Initialize the app on the datapath connect event.
        Subclasses can override this method to init default flows for
        the table that they handle.
        """
        pass

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup the app on the datapath disconnect event.
        Subclasses can override this method to cleanup flows for
        the table that they handle.
        """
        pass

    def delete_all_flows(self, datapath):
        """
        Delete all flows in tables that the controller is responsible for.
        """
        pass

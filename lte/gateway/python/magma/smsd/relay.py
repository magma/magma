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
import asyncio
import logging
from typing import List

import grpc
import lte.protos.sms_orc8r_pb2 as sms_orc8r_pb2
import lte.protos.sms_orc8r_pb2_grpc as sms_orc8r_pb2_grpc
from lte.protos.mconfig.mconfigs_pb2 import MME
from magma.common.job import Job
from magma.common.rpc_utils import grpc_async_wrapper, return_void
from magma.configuration.mconfig_managers import load_service_mconfig
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

POLL_INTERVAL = 5   # seconds. Implies minimum SMS delivery latency.
TIMEOUT_SECS = 15
SMS_TIMEOUT_SECS = 5


class SmsRelay(Job):
    def __init__(
            self,
            loop: asyncio.AbstractEventLoop,
            directoryd: GatewayDirectoryServiceStub,
            sms_orc8r_gw_mme: sms_orc8r_pb2_grpc.SMSOrc8rGatewayServiceStub,
            smsd: sms_orc8r_pb2_grpc.SmsDStub,
    ) -> None:
        super().__init__(interval=POLL_INTERVAL, loop=loop)
        self._directoryd = directoryd
        self._mme_sms = sms_orc8r_gw_mme
        self._smsd = smsd

    def add_to_server(self, server):
        """ Add ourselves to the gRPC servicer """
        sms_orc8r_pb2_grpc.add_SMSOrc8rServiceServicer_to_server(self, server)

    async def _run(self) -> None:
        if not self._is_enabled():
            # sleep, and don't poll for messages
            logging.info(
                "mme non_eps_service_config is not SMS_ORC8R, sleeping.",
            )
            return

        imsis = await self._get_attached_imsis()
        if len(imsis) == 0:
            logging.debug("No active subs")
            return

        logging.info("Checking SMS for %d IMSIs", len(imsis))
        try:
            smsd_resp = await grpc_async_wrapper(
                self._smsd.GetMessages.future(
                    sms_orc8r_pb2.GetMessagesRequest(imsis=imsis),
                    TIMEOUT_SECS,
                ),
            )
        except grpc.RpcError as err:
            logging.error("GRPC call failed while fetching messages: %s", err)
            return

        for msg in smsd_resp.messages:
            logging.error('%s', msg)
            await self._send_sms(msg)

    async def _get_attached_imsis(self) -> List[str]:
        try:
            smsd_resp = await grpc_async_wrapper(
                self._directoryd.GetAllDirectoryRecords.future(
                    Void(), TIMEOUT_SECS,
                ),
            )
            return [r.id for r in smsd_resp.records]
        except grpc.RpcError as err:
            logging.error("Error fetching IMSIs from directoryd: %s", err)
            return []

    async def _send_sms(self, dl: sms_orc8r_pb2.SMODownlinkUnitdata):
        try:
            await grpc_async_wrapper(
                self._mme_sms.SMODownlink.future(dl, SMS_TIMEOUT_SECS),
            )
        except grpc.RpcError as err:
            logging.error("RPC call to MME failed: %s", err)
            return

    @return_void
    def SMOUplink(self, request: sms_orc8r_pb2.SMOUplinkUnitdata, context):
        logging.debug(
            "got an uplink: %s: %s",
            request.imsi, request.nas_message_container.hex(),
        )

        if not self._is_enabled():
            # sleep, and don't poll for messages
            logging.info(
                "mme non_eps_service_config is not SMS_ORC8R, ignoring uplink message.",
            )
            return

        try:
            self._smsd.ReportDelivery(
                sms_orc8r_pb2.ReportDeliveryRequest(
                    report=sms_orc8r_pb2.SMOUplinkUnitdata(
                        imsi="IMSI" + request.imsi,
                        nas_message_container=request.nas_message_container,
                    ),
                ),
            )
        except grpc.RpcError as err:
            context.set_details('SMS delivery report to smsd failed: %s' % err)
            context.set_code(grpc.StatusCode.INTERNAL)
            return

    def _is_enabled(self) -> bool:
        """Return whether SMS should act as a relay

        SMS_ORC8R has value 3
        Returns:
        bool: True if MME's NON_EPS_SERVICE_CONFIG is set to SMS_ORC8R
        False otherwise
        """
        mme_service_config = load_service_mconfig("mme", MME())
        non_eps_service_control = mme_service_config.non_eps_service_control
        return non_eps_service_control and non_eps_service_control == 3

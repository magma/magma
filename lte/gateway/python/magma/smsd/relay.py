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
from typing import Any, Dict

import grpc

from magma.common.rpc_utils import return_void
from magma.common.service import MagmaService
from magma.common.sdwatchdog import SDWatchdogTask
import lte.protos.sms_orc8r_pb2_grpc as sms_orc8r_pb2_grpc
import lte.protos.sms_orc8r_pb2 as sms_orc8r_pb2
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

POLL_INTERVAL = 5 # seconds. Implies minimum SMS delivery latency.

class SmsRelay(SDWatchdogTask):
    def __init__(self,
                 service: MagmaService,
                 directoryd: GatewayDirectoryServiceStub,
                 sms_orc8r_gw_mme: sms_orc8r_pb2_grpc.SMSOrc8rGatewayServiceStub,
                 smsd: sms_orc8r_pb2_grpc.SmsDStub
                ) -> None:
        super().__init__(POLL_INTERVAL, service.loop)
        self._service = service
        self._directoryd = directoryd
        self._sms_orc8r_gw_mme = sms_orc8r_gw_mme
        self._smsd = smsd
        self._add_to_server(self._service.rpc_server)
        

    def _add_to_server(self, server):
        """ Add ourselves to the gRPC servicer """
        sms_orc8r_pb2_grpc.add_SMSOrc8rServiceServicer_to_server(self, server)

    async def _run(self):
        try:
            res = self._directoryd.GetAllDirectoryRecords(Void())
        except grpc.RpcError as err:
            logging.error("Error while fetching active IMSIs from directoryd: %s", err)
            return
        if len(res.records) > 0:
            imsis = [_.id for _ in res.records]
            gmr = sms_orc8r_pb2.GetMessagesRequest(imsis=imsis)
            logging.debug("Checking SMS for %d IMSIs" % len(imsis))
            try:
                res = await self._smsd.GetMessages(gmr)
            except grpc.RpcError as err:
                logging.error("GRPC call failed while fetching messages: %s", err)
                return

            try:
                await self._send_sms_to_mme(res)
            except grpc.RpcError as err:
                logging.error("GRPC call failed while sending messages to MME: %s", err)
                return
        else:
            logging.debug("No active subs")

    async def _send_sms_to_mme(self,
                               resp: sms_orc8r_pb2.GetMessagesResponse):
        for sms in resp:
            self._sms_orc8r_gw_mme.SMODownlink(sms)


    @return_void
    def SMOUplink(self, request, context):
        logging.debug("got an uplink: %s: %s" % (request.imsi, request.nas_message_container.hex()))
        report = sms_orc8r_pb2.SMOUplinkUnitdata(imsi=request.imsi, nas_message_container=request.nas_message_container)
        rdr = sms_orc8r_pb2.ReportDeliveryRequest(report=report)

        try:
            self._smsd.ReportDelivery(rdr)
        except grpc.RpcError as err:
            logging.error("GRPC call failed while sending report to smsd: %s", err)
            return


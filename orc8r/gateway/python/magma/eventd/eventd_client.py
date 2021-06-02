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

import grpc
from google.protobuf.json_format import MessageToDict
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.eventd_pb2 import Event
from orc8r.protos.eventd_pb2_grpc import EventServiceStub

EVENTD_SERVICE_NAME = "eventd"
DEFAULT_GRPC_TIMEOUT = 10


def log_event(event: Event) -> None:
    """
    Make RPC call to 'LogEvent' method of local eventD service
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(
            EVENTD_SERVICE_NAME, ServiceRegistry.LOCAL,
        )
    except ValueError:
        logging.error("Cant get RPC channel to %s", EVENTD_SERVICE_NAME)
        return
    client = EventServiceStub(chan)
    try:
        # Location will be filled in by directory service
        client.LogEvent(event, DEFAULT_GRPC_TIMEOUT)
    except grpc.RpcError as err:
        if err.code() == grpc.StatusCode.UNAVAILABLE:
            logging.debug(
                "LogEvent will not occur unless eventd configuration "
                "is set up.",
            )
        else:
            logging.error(
                "LogEvent error for event: %s, [%s] %s",
                MessageToDict(event),
                err.code(),
                err.details(),
            )

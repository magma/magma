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

import json
import logging
import socket
from contextlib import closing
from typing import Any, Dict
import grpc
import jsonschema

from magma.common.rpc_utils import return_void
from orc8r.protos import eventd_pb2_grpc, eventd_pb2
from .event_validator import EventValidator


class EventDRpcServicer(eventd_pb2_grpc.EventServiceServicer):
    """
    gRPC based server for EventD.
    """

    def __init__(self, config: Dict[str, Any], validator: EventValidator):
        self.fluent_bit_port = config['fluent_bit_port']
        self.tcp_timeout = config['tcp_timeout']
        self._validator = validator

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        eventd_pb2_grpc.add_EventServiceServicer_to_server(self, server)

    @return_void
    def LogEvent(self, request: eventd_pb2.Event, context):
        """
        Logs an event.
        """
        logging.debug("Logging event: %s", request)

        try:
            self._validator.validate_event(request.value, request.event_type)
        except (KeyError, jsonschema.ValidationError) as e:
            logging.error("KeyError for log: %s. Error: %s", request, e)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(
                'Event validation failed, Details: {}'.format(e))
            return

        try:
            with closing(socket.create_connection(
                    ('localhost', self.fluent_bit_port),
                    timeout=self.tcp_timeout)) as sock:
                logging.debug('Sending log to FluentBit')
                value = {
                    'stream_name': request.stream_name,
                    'event_type': request.event_type,
                    # We use event_tag as FluentD uses the "tag" field
                    'event_tag': request.tag,
                    'value': request.value
                }
                sock.sendall(json.dumps(value).encode('utf-8'))
        except socket.error as e:
            logging.error('Connection to FluentBit failed: %s', e)
            logging.info('FluentBit (td-agent-bit) may not be enabled '
                         'or configured correctly')
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details(
                'Could not connect to FluentBit locally, Details: {}'
                .format(e))
            return

        logging.debug("Successfully logged event: %s", request)

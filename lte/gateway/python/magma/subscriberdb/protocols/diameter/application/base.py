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
from enum import IntEnum, unique

from magma.subscriberdb.metrics import (
    DIAMETER_CEX_TOTAL,
    DIAMETER_DISCONECT_TOTAL,
    DIAMETER_WATCHDOG_TOTAL,
)
from magma.subscriberdb.protocols.diameter import avp, message

from . import abc


@unique
class BaseApplicationCommands(IntEnum):
    # Command codes defined in this application used in msg header
    CAPABILITIES_EXCHANGE = 257
    DEVICE_WATCHDOG = 280
    DISCONNECT_PEER = 282


class BaseApplication(abc.Application):
    """
    This is where we implement the Diameter Base Protocol Application which
    is defined in RFC6733. All diameters servers should implement this spec.
    The base Diameter protocol concerns itself with establishing connections
    to peers, capabilities negotiation, how messages are sent and routed
    through peers, and how the connections are eventually torn down.

    This implements a subset of the diameter Base Protcol that OAI MME requires
    This includes handling connection initiation and transport failure
    detection via watchdog requests
    """
    # The ID used for common messages addressed to the base application
    APP_ID = 0
    # The Vendor-Specific-Application-Id and VendorId AVPs that
    # the Common Application should advertise
    CAPABILITIES_EXCHANGE_AVPS = [avp.AVP('Vendor-Id', 0)]
    # Required fields for requests of each command type
    REQUIRED_FIELDS = {
        BaseApplicationCommands.CAPABILITIES_EXCHANGE:
            [
                'Host-IP-Address',
                'Inband-Security-Id',
                'Vendor-Id',
                'Supported-Vendor-Id',
                'Vendor-Specific-Application-Id',
            ],
        BaseApplicationCommands.DEVICE_WATCHDOG:
            [],
        BaseApplicationCommands.DISCONNECT_PEER:
            [],
    }

    def __init__(self, realm, host, host_ip):
        """Each application has access to a write stream and a collection of
        settings, currently limited to realm and host

        Args:
            realm: the realm the application should serve
            host: the host name the application should serve
            host_ip: the IP address of the host
        """
        super(BaseApplication, self).__init__(realm, host, host_ip)
        self.applications = []

    def validate_message(self, state_id, msg):
        """
        Validate a message and send the appropriate error response
        if necessary

        Args:
            msg: the message to validate
        Returns:
            True if the message validated
        """
        # Validate we have all required fields
        required_fields = self.REQUIRED_FIELDS[msg.header.command_code]
        if not msg.has_fields(required_fields):
            logging.error(
                "Missing AVP for diameter command %d",
                msg.header.command_code,
            )
            resp = self._gen_response(
                state_id, msg,
                avp.ResultCode.DIAMETER_MISSING_AVP,
            )
            self.writer.send_msg(resp)
            return False
        return True

    def register(self, application):
        """
        Registers an application for the Diameter Base application to advertise
        in the capabilities exchange

        Args:
            application: an application instance
        Returns:
            None
        Raises:
            TypeError: application is not an Application subtype
        """
        if not issubclass(type(application), abc.Application):
            raise TypeError('Not a valid application')
        self.applications.append(application)

    def handle_msg(self, state_id, msg):
        """
        Handle a message bound for the common application

        Args:
            state_id: the server state id
            msg: the message to handle
        Returns:
            None
        """
        if not msg.header.request:
            logging.warning("Received unsolicited answer")
            return

        if msg.header.command_code == BaseApplicationCommands.CAPABILITIES_EXCHANGE:
            self._send_capabilities(state_id, msg)
        elif msg.header.command_code == BaseApplicationCommands.DEVICE_WATCHDOG:
            self._send_device_watchdog(state_id, msg)
        elif msg.header.command_code == BaseApplicationCommands.DISCONNECT_PEER:
            self._send_disconnect_peer(state_id, msg)
        else:
            logging.error('Unsupported command: %d', msg.header.command_code)

    def _gen_response(self, state_id, msg, result_code, body_avps=None):
        """
        Generate a response message with all the fields that are found
        in base application responses.

        Args:
            state_id: the server state identifer
            msg: the message to respond to
            result_code: the result code to send
            body_avps: the AVPs to include in the response body
        Returns:
            Message instance containing the response
        """
        # Generate an empty message with the response headers to the msg
        if body_avps is None:
            body_avps = []
        resp_msg = message.Message.create_response_msg(msg)

        resp_msg.append_avp(avp.AVP('Result-Code', result_code))
        resp_msg.append_avp(avp.AVP('Origin-Host', self.host))
        resp_msg.append_avp(avp.AVP('Origin-Realm', self.realm))
        resp_msg.append_avp(avp.AVP('Origin-State-Id', state_id))

        # Add body AVPs
        for body_avp in body_avps:
            resp_msg.append_avp(body_avp)
        return resp_msg

    def _send_capabilities(self, state_id, msg):
        """
        Responds to a Capability-Exchange request by sending Vendor specific
        AVPs from the base and registered applications

        Args:
            state_id: the server state identifier
            msg: the message to respond to
        Returns:
            None
        """
        if self.validate_message(state_id, msg):
            # Generate the capabilities body
            body_avps = [avp.AVP('Host-IP-Address', self.host_ip)]
            body_avps.extend(self.CAPABILITIES_EXCHANGE_AVPS)
            for application in self.applications:
                body_avps.extend(application.CAPABILITIES_EXCHANGE_AVPS)
            body_avps.append(avp.AVP('Product-Name', avp.PRODUCT_NAME))
            DIAMETER_CEX_TOTAL.inc()
            resp = self._gen_response(
                state_id, msg,
                avp.ResultCode.DIAMETER_SUCCESS,
                body_avps,
            )
            self.writer.send_msg(resp)

    def _send_device_watchdog(self, state_id, msg):
        """
        Responds to a Device-Watchdog requests which are pings to detection
        transport failures

        Args:
            state_id: the server state identifier
            msg: the message to respond to
        Returns:
            None
        """
        if self.validate_message(state_id, msg):
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_SUCCESS,
            )
            DIAMETER_WATCHDOG_TOTAL.inc()
            self.writer.send_msg(resp)

    def _send_disconnect_peer(self, state_id, msg):
        """
        Responds to a Disconnect-Peer-Request. Upon receipt of this
        message, the transport connection is shut down.

        Args:
            state_id: the server state identifier
            msg: the message to respond to
        Returns:
            None
        """
        logging.info('Received disconnect request, state id: %d', state_id)
        if self.validate_message(state_id, msg):
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_SUCCESS,
            )
            DIAMETER_DISCONECT_TOTAL.inc()
            self.writer.send_msg(resp)

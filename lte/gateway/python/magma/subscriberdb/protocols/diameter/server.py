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
import random

from magma.subscriberdb.protocols.diameter.application import base, s6a

from . import exception, message


class S6aServer(asyncio.Protocol):
    """
    This is a Diameter 3GPP S6A Server. It sits between the MME and
    subscriberdb to exchange auth information. This class handles TCP
    connection initialization and handling incoming data from the network
    """

    def __init__(self, base_manager, s6a_manager, realm, host, loop=None):
        self.realm = realm
        self.host = host
        self.state_id = random.randint(0, 100000000)
        self._s6a_manager = s6a_manager
        self._readbuf = None
        self._base_manager = base_manager
        self.writer = None
        self.loop = loop

    def connection_made(self, transport):
        """
        Handle a new Diameter connection.

        Args:
            transport (asynio.Transport): the transport for the new connection
        Returns:
            None
        """
        logging.info("Connection received, state id: %d", self.state_id)
        # bytesarray is more efficient to append fragments of reads
        self._readbuf = bytearray()
        self.writer = Writer(
            self.realm, self.host,
            self.state_id, transport,
        )
        self._base_manager.set_writer(self.writer)
        self._s6a_manager.set_writer(self.writer)

    def data_received(self, data):
        """
        Append the 'data' bytes to the readbuf, and parse the message
        if the entire message has been received.
        Unparsed bytes will be left in readbuf and will be parsed when
        more data is received in the future.

        Args:
            data (bytes): new data read from the transport
        Returns:
            None
        """
        logging.debug("Bytes read: %s", data)
        self._readbuf.extend(data)

        # Use memoryview to prevent copies when slicing
        memview = memoryview(self._readbuf)
        remain = len(memview)
        begin = 0  # beginning of message

        while remain >= message.HEADER_LEN:
            try:
                msg = message.decode(memview[begin:])
                logging.debug("Handling diameter message:\n%s", msg)
                self._handle_msg(msg.header.application_id, msg)
                # Get ready for the next message
                begin += msg.length
                remain -= msg.length
            except exception.TooShortException:
                logging.error("Diameter message too short to decode")
                return
            except Exception as exc:  # pylint: disable=broad-except
                # Handle any exceptions with message handling, without
                # affecting other messages/users
                logging.exception(exc)

                # Clear past garbage
                length = len(memview[begin:])
                begin += length
                remain -= length

        # Get the unparsed bytes
        self._readbuf = bytearray(memview[begin:])

    def connection_lost(self, exc):
        """
        The IPA connection has been lost.

        Args:
            exc: exception object or None if EOF
        Returns:
            None
        """
        logging.warning("Connection lost!")

    def _handle_msg(self, application_id, msg):
        """
        Handles a message bound for an application.

        Args:
            application_id: the application id the message is addressed to
            msg: the actual message
        Returns:
            None
        """
        # TOOD(oramadan) move this distpatch loging out of server
        if application_id == base.BaseApplication.APP_ID:
            self._base_manager.handle_msg(self.state_id, msg)
        elif application_id == s6a.S6AApplication.APP_ID:
            self._s6a_manager.handle_msg(self.state_id, msg)
        else:
            logging.error(
                "Unknown application: %d",
                msg.header.application_id,
            )


class Writer:
    """The writer abstracts away a client connection for an
    application to be able to send messages to.
    """

    def __init__(self, realm, host, state_id, transport):
        self.realm = realm
        self.host = host
        self.state_id = state_id
        self._transport = transport

    def send_msg(self, msg):
        """
        Sends a message. Prepares a writer buffer to send,
        and encodes the message into it, then writes to transport.

        Args:
            msg: the message to send
        Returns:
            None
        """
        buf = self._get_write_buf(msg)
        logging.debug("Sending diameter response:\n%s", msg)
        try:
            msg.encode(buf, 0)
        except exception.CodecException as e:
            logging.fatal("Encoding failed with err: %s", e)
            return
        self._write(buf)

    def _get_write_buf(self, msg):
        """
        Allocates one chunk of memory for the entire message including
        the header which needs to be prepended. This avoids extra allocations
        and copying.

        Args:
            msg: message to be encapsulated by IPA
        Returns:
            memoryview: Allocated buf of header + length bytes.
        """
        buf = memoryview(bytearray(msg.length))
        return buf

    def _write(self, buf):
        """
        Write the buffer to the underlying socket

        Args:
            buf: the message buffer to wrtie to the transport
        Returns:
            None
        """
        self._transport.write(buf)

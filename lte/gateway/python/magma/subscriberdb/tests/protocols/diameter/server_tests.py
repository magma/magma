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

# pylint:disable=protected-access

import unittest
from unittest.mock import Mock

from magma.subscriberdb.protocols.diameter import message, server

from .common import MockTransport


class ServerTests(unittest.TestCase):
    """
    Test class for Diameter Server dispatch to Applications
    """

    def setUp(self):

        self._server = server.S6aServer(
            Mock(),
            Mock(),
            "mai.facebook.com",
            "hss.mai.facebook.com",
        )

        # Mock the message handler
        self._server._handle_msg = Mock()

        # Mock the writes to check responses
        self._writes = Mock()

        def convert_memview_to_bytes(memview):
            """ Deep copy the memoryview for checking later  """
            return self._writes(memview.tobytes())

        self._transport = MockTransport()
        self._transport.write = Mock(side_effect=convert_memview_to_bytes)

        # Here goes nothing..
        self._server.connection_made(self._transport)

    def _check_handler(self, req_bytes, application_id):
        """
        Send data to the protocol in different step lengths to
        verify that we assemble all segments and invoke the correct handler

        Args:
            req_bytes (bytes): request which would be sent
                multiple times with different step sizes
            application_id: the application handler which should be invoked
        Returns:
            None
        """
        for step in range(1, len(req_bytes) + 1):
            offset = 0
            while offset < len(req_bytes):
                self._server.data_received(req_bytes[offset:offset + step])
                offset += step
            self.assertTrue(self._server._handle_msg.called)
            # pylint:disable=unsubscriptable-object
            self.assertEqual(
                self._server._handle_msg.call_args[0][0],
                application_id,
            )
            self._server._handle_msg.reset_mock()

    def test_application_dispatch(self):
        """Check that we can decode an inbound message and call
        the appropriate handler"""
        msg = message.Message()
        msg.header.application_id = 0xfac3b00c
        msg.header.request = True

        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)
        self._check_handler(req_buf, 0xfac3b00c)

    def test_too_short(self):
        """Check that if we didn't receive enough data
        we keep it in the buffer"""
        # Read in less than the header length
        req_buf = bytearray(b'\x01' * 19)
        self._server.data_received(req_buf)
        self.assertEqual(len(self._server._readbuf), 19)

    def test_decode_error(self):
        """Check that we can seek past garbage"""
        # Feed garbage past the header length
        req_buf = bytearray(b'\x01' * 20)
        self._server.data_received(req_buf)
        # We should flush the read buffer
        self.assertEqual(len(self._server._readbuf), 0)


class WriterTests(unittest.TestCase):
    """
    Test the Writer class for the diameter server
    """

    def setUp(self):
        # Mock the writes to check responses
        self._writes = Mock()

        def convert_memview_to_bytes(memview):
            """ Deep copy the memoryview for checking later  """
            return self._writes(memview.tobytes())

        self._transport = MockTransport()
        self._transport.write = Mock(side_effect=convert_memview_to_bytes)

        self.writer = server.Writer(
            "mai.facebook.com",
            "hss.mai.facebook.com",
            "127.0.0.1",
            self._transport,
        )

    def test_send_msg(self):
        """Test that the writer will encode a message and write
        it to the transport"""
        msg = message.Message()
        self.writer.send_msg(msg)
        self._writes.assert_called_once_with(
            b'\x01\x00\x00\x14'
            b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
            b'\x00\x00\x00',
        )
        self._writes.reset_mock()

    def test_gen_buf(self):
        """Test that the writer will generate a buffer of the
        length of the message"""
        msg = message.Message()
        buf = self.writer._get_write_buf(msg)
        self.assertEqual(len(buf), msg.length)

    def test_write(self):
        """Test that the writer will push to the transport"""
        msg = memoryview(b'helloworld')
        self.writer._write(msg)
        self._writes.assert_called_once_with(msg)
        self._writes.reset_mock()


if __name__ == "__main__":
    unittest.main()

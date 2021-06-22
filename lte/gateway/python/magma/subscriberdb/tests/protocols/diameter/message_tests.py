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

# pylint: disable=protected-access

import unittest

from magma.subscriberdb.protocols.diameter import avp, message
from magma.subscriberdb.protocols.diameter.exception import (
    CodecException,
    TooShortException,
)


class MessageHeaderCodecTests(unittest.TestCase):
    """
    Tests for encoding and decoding message headers
    """

    def _decode_check(self, header, header_bytes):
        decoded_header = message.MessageHeader.decode(header_bytes)
        self.assertEqual(header.version, decoded_header.version)
        self.assertEqual(header.command_flags, decoded_header.command_flags)
        self.assertEqual(header.command_code, decoded_header.command_code)
        self.assertEqual(header.application_id, decoded_header.application_id)
        self.assertEqual(header.hop_by_hop_id, decoded_header.hop_by_hop_id)
        self.assertEqual(header.end_to_end_id, decoded_header.end_to_end_id)

    def _encode_check(self, header, header_bytes, length):
        out_buf = bytearray(header.length)
        out_len = header.encode(out_buf, 0, length)
        self.assertGreater(out_len, 0)
        self.assertEqual(out_buf, bytes(header_bytes))

    def _compare_header(self, header, header_bytes, length):
        # Test encoder
        self._encode_check(header, header_bytes, length)

        # Test decoder
        self._decode_check(header, header_bytes)

    def test_header_basic(self):
        """Test that we can encode and decode message header"""
        # The default message header is all zeros except the version number
        header = message.MessageHeader()
        self._compare_header(
            header, b'\x01'
            + b'\x00' * (message.HEADER_LEN - 1), 0,
        )

        # Next three bytes are the length
        self._compare_header(
            header, b'\x01\x00\x00\x02'
            + b'\x00' * (message.HEADER_LEN - 4), 2,
        )

        # Next is our command flags
        header.command_flags = 0x1f
        self._compare_header(
            header, b'\x01\x00\x00\x02\x1f'
            + b'\x00' * (message.HEADER_LEN - 5), 2,
        )

        # Command code is the next three bytes
        header.command_code = 0x1
        self._compare_header(
            header, b'\x01\x00\x00\x02\x1f\x00\x00\x01'
            + b'\x00' * (message.HEADER_LEN - 8), 2,
        )

        # Application id is the next four bytes
        header.application_id = 0x7
        self._compare_header(
            header, b'\x01\x00\x00\x02\x1f\x00\x00\x01'
            + b'\x00\x00\x00\x07'
            + b'\x00' * (message.HEADER_LEN - 12), 2,
        )

        # hop by hop id is the next four bytes
        header.hop_by_hop_id = 0xb33f0000
        self._compare_header(
            header, b'\x01\x00\x00\x02\x1f\x00\x00\x01'
            + b'\x00\x00\x00\x07\xb3\x3f\x00\x00'
            + b'\x00' * (message.HEADER_LEN - 16), 2,
        )

        # end to end id is the last four
        header.end_to_end_id = 0x0000dead
        self._compare_header(
            header, b'\x01\x00\x00\x02\x1f\x00\x00\x01'
            + b'\x00\x00\x00\x07\xb3\x3f\x00\x00'
            + b'\x00\x00\xde\xad', 2,
        )

    def test_header_decode_validate(self):
        """
        Tests that we validate the payload is decodable length
        """

        # Must be at least 20 Bytes to decode
        with self.assertRaises(CodecException):
            message.MessageHeader.decode(
                b'\x00'
                * (message.HEADER_LEN - 1),
            )

    def test_header_validate_length(self):
        """
        Tests that we validate the length is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.encode(out_buf, 0, 0x00FFFFFF + 1)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.encode(out_buf, 0, -1)

    def test_header_validate_version(self):
        """
        Tests that we validate the version is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.version = 0xFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.version = -1
            header.encode(out_buf, 0, 0)

    def test_header_validate_command_flags(self):
        """
        Tests that we validate the command flags are encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.command_flags = 0xFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.command_flags = -1
            header.encode(out_buf, 0, 0)

    def test_header_validate_command_code(self):
        """
        Tests that we validate the command code is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.command_code = 0x00FFFFFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.command_code = -1
            header.encode(out_buf, 0, 0)

    def test_header_validate_app_id(self):
        """
        Tests that we validate the app id is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.application_id = 0xFFFFFFFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.application_id = -1
            header.encode(out_buf, 0, 0)

    def test_header_validate_hbh_id(self):
        """
        Tests that we validate the hbh id is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.application_id = 0xFFFFFFFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.application_id = -1
            header.encode(out_buf, 0, 0)

    def test_header_validate_ete_id(self):
        """
        Tests that we validate the ete id is encodable
        """
        out_buf = bytearray(message.HEADER_LEN)
        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.end_to_end_id = 0xFFFFFFFF + 1
            header.encode(out_buf, 0, 0)

        with self.assertRaises(CodecException):
            header = message.MessageHeader()
            header.end_to_end_id = -1
            header.encode(out_buf, 0, 0)


class MessageHeaderFlagTests(unittest.TestCase):
    """
    Tests for reading and writing the message flags
    """

    def test_read(self):
        """
        Tests that we can use the flag properties to read bits of command_flags
        """
        header = message.MessageHeader()

        # It all starts off as false
        self.assertFalse(header.retransmitted)
        self.assertFalse(header.error)
        self.assertFalse(header.proxiable)
        self.assertFalse(header.request)

        header.command_flags = 0x10
        self.assertTrue(header.retransmitted)

        header.command_flags = 0x20
        self.assertTrue(header.error)

        header.command_flags = 0x40
        self.assertTrue(header.proxiable)

        header.command_flags = 0x80
        self.assertTrue(header.request)

    def test_write(self):
        """
        Test that we can use the message header flag properties to flip
        the bits of the command_flags field
        """
        header = message.MessageHeader()
        # Set reserved bit to validate we are not overriding command_flags
        # and are properly doing the bit flipping
        header.command_flags = 0x01

        header.retransmitted = True
        self.assertEqual(header.command_flags, 0x11)
        header.retransmitted = False
        self.assertEqual(header.command_flags, 0x01)

        header.error = True
        self.assertEqual(header.command_flags, 0x21)
        header.error = False
        self.assertEqual(header.command_flags, 0x01)

        header.proxiable = True
        self.assertEqual(header.command_flags, 0x41)
        header.proxiable = False
        self.assertEqual(header.command_flags, 0x01)

        header.request = True
        self.assertEqual(header.command_flags, 0x81)
        header.request = False
        self.assertEqual(header.command_flags, 0x01)

    def test_clone_and_respond(self):
        """
        Tests that we can use the convenience clone constructor
        and create_response_header method for copying fields from header
        """
        # Set some non-defaults we will clone
        header = message.MessageHeader()
        header.version = 0x2
        header.command_code = 0x3
        header.application_id = 0x4
        header.hop_by_hop_id = 0x5
        header.end_to_end_id = 0x6
        header.request = True
        header.retransmitted = True
        header.error = True
        header.proxiable = True

        # A clone should be identical
        clone = message.MessageHeader.copy(header)
        self.assertEqual(header, clone)

        # A clone should be almost equal
        resp = message.MessageHeader.create_response_header(header)
        # These fields are reset
        self.assertFalse(resp.request)
        self.assertFalse(resp.retransmitted)
        self.assertFalse(resp.error)

        # Now they should be equal
        resp.request = header.request
        resp.retransmitted = header.retransmitted
        resp.error = header.error
        self.assertEqual(header, resp)


class MessageCodecTests(unittest.TestCase):
    """
    Tests for encoding and decoding Messages
    """

    def _decode_check(self, msg, msg_bytes):
        decoded_msg = message.decode(msg_bytes)
        self.assertEqual(msg.header, decoded_msg.header)
        self.assertEqual(msg._avps, decoded_msg._avps)

    def _encode_check(self, msg, msg_bytes):
        out_buf = bytearray(msg.length)
        out_len = msg.encode(out_buf, 0)
        self.assertGreater(out_len, 0)
        self.assertEqual(out_buf, bytes(msg_bytes))

    def _compare_msg(self, msg, msg_bytes):
        # Test encoder
        self._encode_check(msg, msg_bytes)

        # Test decoder
        self._decode_check(msg, msg_bytes)

    def test_empty_message(self):
        """
        Tests that we can encode and decode an empty message
        """
        # A message with no arguments is just an empty header
        msg = message.Message()
        self._compare_msg(
            msg, b'\x01'  # Version
            + b'\x00\x00\x14'  # Length 20
            + b'\x00' * (message.HEADER_LEN - 4),
        )

    def test_message_with_avp(self):
        """
        Adding an AVP to the message is just tacking on the AVP encoded
        output
        """
        msg = message.Message()
        msg.append_avp(avp.AVP('User-Name', 'hello'))
        self._compare_msg(
            msg, b'\x01'  # Version
            + b'\x00\x00\x24'  # Length 36
            + b'\x00' * (message.HEADER_LEN - 4)
            + b'\x00\x00\x00\x01@\x00\x00\rhello\x00\x00\x00',
        )

    def test_decode_validate_length(self):
        """
        If too short to decode we should raise a unique exception that
        Will allow us to retry
        """
        # Shorter than header by 1
        with self.assertRaises(TooShortException):
            message.decode(
                b'\x01'
                + b'\x00' * (message.HEADER_LEN - 2),
            )

        # Header says 24 bytes but we give 23
        with self.assertRaises(TooShortException):
            payload = (
                b'\x01'  # Version
                + b'\x00\x00\x18'  # Length 24
                + b'\x00' * (message.HEADER_LEN - 1)
            )
            self.assertEqual(len(payload), 23)
            message.decode(payload)

    def test_decode_garbage(self):
        """
        If too short to decode we should raise a unique exception that
        Will allow us to retry
        """
        # Gave a length that was not a multiple of 4
        with self.assertRaises(CodecException):
            message.decode(
                b'\x01'  # Version
                + b'\x00\x00\x15'  # Length 21
                + b'\x00' * (message.HEADER_LEN - 3),
            )

        # Garbage AVPs
        with self.assertRaises(CodecException):
            message.decode(
                b'\x01'  # Version
                + b'\x00\x00\x18'  # Length 24
                + b'\x00' * (message.HEADER_LEN),
            )

    def test_respond(self):
        """
        Tests that we can use the convenience clone constructor
        and create_response_header method for copying fields from header
        """
        # Set some non-defaults we will clone
        header = message.MessageHeader()
        header.version = 0x2
        header.command_code = 0x3
        header.application_id = 0x4
        header.hop_by_hop_id = 0x5
        header.end_to_end_id = 0x6
        header.request = True
        header.retransmitted = True
        header.error = True
        header.proxiable = True
        msg = message.Message(header)
        msg.append_avp(avp.AVP('User-Name', ''))

        # A clone should be almost equal
        resp = message.Message.create_response_msg(msg)
        # These fields are reset
        self.assertFalse(resp.header.request)
        self.assertFalse(resp.header.retransmitted)
        self.assertFalse(resp.header.error)
        self.assertEqual(len(resp._avps), 0)

        # Now they should be equal
        resp.header.request = msg.header.request
        resp.header.retransmitted = msg.header.retransmitted
        resp.header.error = msg.header.error
        self.assertEqual(msg.header, resp.header)


class MessageAVPTests(unittest.TestCase):
    """
    Tests for adding and reading AVPs in a message
    """
    @classmethod
    def setUpClass(cls):
        cls.msg = message.Message()
        cls.msg.append_avp(avp.AVP('User-Name', 'hello'))
        cls.msg.append_avp(avp.AVP('User-Name', 'world'))
        cls.msg.append_avp(avp.AVP('Host-IP-Address', '127.0.0.1'))

    def test_message_add_avp(self):
        """We successfully added 3 avps"""
        self.assertEqual(len(self.msg._avps), 3)

    def test_message_filter_avp(self):
        """We can get a list of AVPs that matches a code and vendor"""
        # Filter for code 1 (User-Name AVP)
        filtered_avps = list(self.msg.filter_avps(avp.VendorId.DEFAULT, 1))
        self.assertEqual(len(filtered_avps), 2)
        self.assertEqual(filtered_avps[0].value, 'hello')
        self.assertEqual(filtered_avps[1].value, 'world')

        # Filter for non-existent
        self.assertEqual(
            len(list(self.msg.filter_avps(*avp.resolve('Session-Id')))), 0,
        )

    def test_get_avp(self):
        """We can get the first occurance of an AVP by name or code"""
        self.assertEqual(
            self.msg.find_avp(
                *avp.resolve('User-Name'),
            ).value, 'hello',
        )
        self.assertEqual(self.msg.find_avp(0, 257).value, '127.0.0.1')

        # Doesn't exist so returns None
        self.assertEqual(self.msg.find_avp(0, 1337), None)

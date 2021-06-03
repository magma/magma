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

import unittest

from magma.subscriberdb.protocols.diameter import avp
from magma.subscriberdb.protocols.diameter.exception import CodecException


class AVPHeaderTests(unittest.TestCase):
    """
    Tests for encoding and decoding AVP headers
    """

    def _decode_check(self, avp_val, msg_bytes):
        decoded_val = avp.decode(msg_bytes)
        self.assertEqual(avp_val.code, decoded_val.code)
        self.assertEqual(avp_val.vendor, decoded_val.vendor)
        self.assertEqual(avp_val.value, decoded_val.value)
        self.assertEqual(avp_val.flags, decoded_val.flags)

    def _encode_check(self, avp_val, msg_bytes):
        out_buf = bytearray(avp_val.length)
        out_len = avp_val.encode(out_buf, 0)
        self.assertGreater(out_len, 0)
        self.assertEqual(out_buf, bytes(msg_bytes))

    def _compare_avp(self, avp_val, msg_bytes):
        # Test encoder
        self._encode_check(avp_val, msg_bytes)

        # Test decoder
        self._decode_check(avp_val, msg_bytes)

    def test_avp_flags(self):
        """
        Tests we can encode and decode AVPs with different flags
        """
        self._compare_avp(
            avp.UnknownAVP(0, b''),
            memoryview(b'\x00\x00\x00\x00\x00\x00\x00\x08'),
        )

        avp_val = avp.UnknownAVP(0, b'', flags=avp.FLAG_MANDATORY)
        self._compare_avp(
            avp_val,
            b'\x00\x00\x00\x00@\x00\x00\x08',
        )
        self.assertFalse(avp_val.vendor_specific)
        self.assertTrue(avp_val.mandatory)
        self.assertFalse(avp_val.protected)

        avp_val = avp.UnknownAVP(0, b'', flags=avp.FLAG_PROTECTED)
        self._compare_avp(
            avp_val,
            b'\x00\x00\x00\x00 \x00\x00\x08',
        )
        self.assertFalse(avp_val.vendor_specific)
        self.assertFalse(avp_val.mandatory)
        self.assertTrue(avp_val.protected)

        avp_val = avp.UnknownAVP(
            0, b'', flags=avp.FLAG_VENDOR,
            vendor=avp.VendorId.TGPP,
        )
        self._compare_avp(
            avp_val,
            b'\x00\x00\x00\x00\x80\x00\x00\x0c\x00\x00(\xaf',
        )
        self.assertTrue(avp_val.vendor_specific)
        self.assertFalse(avp_val.mandatory)
        self.assertFalse(avp_val.protected)

        avp_val = avp.UnknownAVP(
            0, b'', flags=avp.FLAG_VENDOR
            | avp.FLAG_MANDATORY,
            vendor=avp.VendorId.TGPP,
        )
        self._compare_avp(
            avp_val,
            b'\x00\x00\x00\x00\xc0\x00\x00\x0c\x00\x00(\xaf',
        )
        self.assertTrue(avp_val.vendor_specific)
        self.assertTrue(avp_val.mandatory)
        self.assertFalse(avp_val.protected)

        avp_val = avp.UnknownAVP(
            0, b'', flags=avp.FLAG_VENDOR
            | avp.FLAG_MANDATORY
            | avp.FLAG_PROTECTED,
            vendor=avp.VendorId.TGPP,
        )
        self._compare_avp(
            avp_val,
            b'\x00\x00\x00\x00\xe0\x00\x00\x0c\x00\x00(\xaf',
        )
        self.assertTrue(avp_val.vendor_specific)
        self.assertTrue(avp_val.mandatory)
        self.assertTrue(avp_val.protected)

    def test_avp_code(self):
        """
        Test the range of AVP codes is 0-0xFFFFFFFF
        """

        avp_val = avp.UnknownAVP(0, b'')
        out_buf = bytearray(avp_val.length)
        avp_val.encode(out_buf, 0)

        avp_val = avp.UnknownAVP(0xFFFFFFFF, b'')
        out_buf = bytearray(avp_val.length)
        avp_val.encode(out_buf, 0)

        with self.assertRaises(CodecException):
            avp_val = avp.UnknownAVP(-1, b'')
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

        with self.assertRaises(CodecException):
            avp_val = avp.UnknownAVP(0xFFFFFFFF + 1, b'')
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

    def test_avp_vendor(self):
        """
        Test the range of AVP vendors is 1-0xFFFFFFFF
        """
        # Vendor specific flags means you need a non default vendor ID
        with self.assertRaises(CodecException):
            avp_val = avp.UnknownAVP(
                0, b'',
                flags=avp.FLAG_VENDOR,
                vendor=avp.VendorId.DEFAULT,
            )
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

        avp_val = avp.UnknownAVP(
            0, b'',
            flags=avp.FLAG_VENDOR,
            vendor=1,
        )
        out_buf = bytearray(avp_val.length)
        avp_val.encode(out_buf, 0)
        self._compare_avp(avp_val, out_buf)

        avp_val = avp.UnknownAVP(
            0, b'',
            flags=avp.FLAG_VENDOR,
            vendor=0x00FFFFFF,
        )
        out_buf = bytearray(avp_val.length)
        avp_val.encode(out_buf, 0)
        self._compare_avp(avp_val, out_buf)

        # Avp vendor in range
        with self.assertRaises(CodecException):
            avp_val = avp.UnknownAVP(
                0, b'',
                flags=avp.FLAG_VENDOR,
                vendor=-1,
            )
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

        # Avp vendor in range
        with self.assertRaises(CodecException):
            avp_val = avp.UnknownAVP(
                0, b'',
                flags=avp.FLAG_VENDOR,
                vendor=0xFFFFFFFF + 1,
            )
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

    def test_avp_length(self):
        """
        Tests we validate AVPs lengths are longer than minumum to decode
        and no longer than the maximum length encodable
        """
        # Avp that has no payload isn't encodable
        with self.assertRaises(CodecException):
            avp_val = avp.AVP(0)
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)

        # Avp shorter than header
        with self.assertRaises(CodecException):
            avp.decode(b'\x00' * (avp.HEADER_LEN - 1))

        # Too short with vendor bit set
        with self.assertRaises(CodecException):
            avp.decode(b'\x00\x00\x00\x00\x80\x00\x00\x00')

        # Max allowable length of payload
        avp_val = avp.UTF8StringAVP(
            1,
            'a' * (0x00FFFFFF - avp.HEADER_LEN),
        )
        out_buf = bytearray(avp_val.length)
        avp_val.encode(out_buf, 0)
        self._compare_avp(avp_val, out_buf)

        # Avp length out of range
        with self.assertRaises(CodecException):
            avp_val = avp.UTF8StringAVP(
                1,
                'a' * (0x00FFFFFF - avp.HEADER_LEN + 1),
            )
            out_buf = bytearray(avp_val.length)
            avp_val.encode(out_buf, 0)


class AVPValueTests(unittest.TestCase):
    """
    Tests for encoding and decoding values
    """

    def _decode_check(self, avp_val, msg_bytes):
        decoded_val = avp.decode(msg_bytes)
        self.assertEqual(avp_val.value, decoded_val.value)
        self.assertEqual(avp_val.payload, decoded_val.payload)

    def _encode_check(self, avp_val, msg_bytes):
        out_buf = bytearray(avp_val.length)
        out_len = avp_val.encode(out_buf, 0)
        self.assertEqual(out_len % 4, 0)  # encoded AVP is end aligned
        self.assertEqual(out_buf, bytes(msg_bytes))

    def _compare_avp(self, avp_val, msg_bytes):
        # Test encoder
        self._encode_check(avp_val, msg_bytes)

        # Test decoder
        self._decode_check(avp_val, msg_bytes)

    def test_octet_strings(self):
        """
        Tests we can encode and decode octet strings
        """
        self._compare_avp(
            avp.UnknownAVP(0, b'hello\x23'),
            memoryview(b'\x00\x00\x00\x00\x00\x00\x00\x0ehello#\x00\x00'),
        )

        # Unicode strings won't load
        with self.assertRaises(CodecException):
            avp.OctetStringAVP(0, u'hello')

    def test_unicode_strings(self):
        """
        Tests we can encode and decode unicode strings
        """
        self._compare_avp(
            avp.UTF8StringAVP(1, u'\u0123\u0490'),
            memoryview(b'\x00\x00\x00\x01\x00\x00\x00\x0c\xc4\xa3\xd2\x90'),
        )

        # Octet strings won't load
        with self.assertRaises(CodecException):
            avp.UTF8StringAVP(1, b'hello')

    def test_unsigned_integers(self):
        """
        Tests we can encode and decode unsigned integers
        """

        self._compare_avp(
            avp.Unsigned32AVP(299, 1234),
            memoryview(b'\x00\x00\x01+\x00\x00\x00\x0c\x00\x00\x04\xd2'),
        )

        with self.assertRaises(CodecException):
            avp.Unsigned32AVP(299, -1234)

    def test_addresses(self):
        """
        Tests we can encode and decode addresses and handle invalid payloads
        """
        # pylint:disable=expression-not-assigned

        self._compare_avp(
            avp.AddressAVP(257, '127.0.0.1'),
            memoryview(
                b'\x00\x00\x01\x01\x00\x00\x00\x0e'
                b'\x00\x01\x7f\x00\x00\x01\x00\x00',
            ),
        )

        self._compare_avp(
            avp.AddressAVP(257, '2001:db8::1'),
            memoryview(
                b'\x00\x00\x01\x01\x00\x00\x00\x1a\x00\x02 \x01\r'
                b'\xb8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
                b'\x01\x00\x00',
            ),
        )

        # Can't read invalid address type \x03
        with self.assertRaises(CodecException):
            avp.decode(
                b'\x00\x00\x01\x01\x00\x00\x00\x0e'
                b'\x00\x03\x7f\x00\x00\x01\x00\x00',
            ).value

        # Can't read too short IPV4
        with self.assertRaises(CodecException):
            avp.decode(
                b'\x00\x00\x01\x01\x00\x00\x00\x0e'
                b'\x00\x01\x7f',
            ).value

        # Can't read too short IPV6
        with self.assertRaises(CodecException):
            avp.decode(
                b'\x00\x00\x01\x01\x00\x00\x00\x0e'
                b'\x00\x02\x7f\x00\x00\x01\x00\x00',
            ).value

        # Cant encode non-ips
        with self.assertRaises(CodecException):
            avp.Unsigned32AVP(257, 'facebook.com')

    def test_result_code(self):
        """
        Test we can encode and decode result codes
        """
        self._compare_avp(
            avp.ResultCodeAVP(268, avp.ResultCode.DIAMETER_SUCCESS),
            memoryview(b'\x00\x00\x01\x0c\x00\x00\x00\x0c\x00\x00\x07\xd1'),
        )

        # Test a value we haven't defined
        self._compare_avp(
            avp.ResultCodeAVP(268, 1337),
            memoryview(b'\x00\x00\x01\x0c\x00\x00\x00\x0c\x00\x00\x059'),
        )

    def test_grouped(self):
        """
        Tests we can encode and decode grouped AVPs
        """

        grouped_avp = avp.GroupedAVP(
            260, [
                avp.UTF8StringAVP(1, 'Hello'),
                avp.UTF8StringAVP(1, 'World'),
            ],
        )
        self._compare_avp(
            grouped_avp,
            (
                b'\x00\x00\x01\x04\x00\x00\x00(\x00\x00'
                b'\x00\x01\x00\x00\x00\rHello\x00\x00\x00'
                b'\x00\x00\x00\x01\x00\x00\x00\rWorld\x00'
                b'\x00\x00'
            ),
        )

        self._compare_avp(
            avp.GroupedAVP(260, [grouped_avp, grouped_avp]),
            (
                b'\x00\x00\x01\x04\x00\x00\x00X\x00\x00\x01\x04'
                b'\x00\x00\x00(\x00\x00\x00\x01\x00\x00\x00\rHello'
                b'\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\rWorld'
                b'\x00\x00\x00\x00\x00\x01\x04\x00\x00\x00(\x00\x00'
                b'\x00\x01\x00\x00\x00\rHello\x00\x00\x00\x00\x00'
                b'\x00\x01\x00\x00\x00\rWorld\x00\x00\x00'
            ),
        )

        # Test filtering
        self.assertEqual(len(list(grouped_avp.filter_avps(0, 1))), 2)
        self.assertEqual(len(list(grouped_avp.filter_avps(0, 2))), 0)

        # Test find returns first
        self.assertEqual(grouped_avp.find_avp(0, 1).value, 'Hello')
        self.assertEqual(grouped_avp.find_avp(0, 2), None)


class AVPConstructorTests(unittest.TestCase):
    """
    Tests for AVP convenience function
    """

    def _decode_check(self, avp_val, msg_bytes):
        decoded_val = avp.decode(msg_bytes)
        self.assertEqual(avp_val.value, decoded_val.value)
        self.assertEqual(avp_val.payload, decoded_val.payload)

    def _encode_check(self, avp_val, msg_bytes):
        out_buf = bytearray(avp_val.length)
        out_len = avp_val.encode(out_buf, 0)
        self.assertEqual(out_len % 4, 0)  # encoded AVP is end aligned
        self.assertEqual(out_buf, bytes(msg_bytes))

    def _compare_avp(self, avp1, avp2):
        self.assertEqual(avp1.code, avp2.code)
        self.assertEqual(avp1.vendor, avp2.vendor)
        self.assertEqual(avp1.flags, avp2.flags)
        self.assertEqual(avp1.payload, avp2.payload)
        self.assertEqual(avp1.name, avp2.name)

    def test_empty_value(self):
        """
        Tests we can initialize an AVP with no value
        """
        avp_val = avp.AVP(0)
        self.assertEqual(avp_val.value, None)
        self.assertEqual(avp_val.payload, None)

        # We can then set its value
        avp_val.value = b''
        self.assertEqual(avp_val.value, b'')
        self.assertEqual(avp_val.payload, b'')

        # And unset it again
        avp_val.value = None
        self.assertEqual(avp_val.value, None)
        self.assertEqual(avp_val.payload, None)

    def test_tuple_identifier(self):
        """
        Tests we can create an AVP with a (vendor, code) tuple
        """

        # This will resolve to the Username AVP
        self._compare_avp(
            avp.AVP((avp.VendorId.DEFAULT, 1), 'a username'),
            avp.UTF8StringAVP(
                1, value='a username', vendor=avp.VendorId.DEFAULT,
                flags=avp.FLAG_MANDATORY,
                name='User-Name',
            ),
        )

        self._compare_avp(
            avp.AVP((avp.VendorId.TGPP, 701), b'msisdn'),
            avp.OctetStringAVP(
                701, value=b'msisdn', vendor=avp.VendorId.TGPP,
                flags=avp.FLAG_MANDATORY | avp.FLAG_VENDOR,
                name='MSISDN',
            ),
        )

        # Unknown AVPs default to unknown AVP
        self._compare_avp(
            avp.AVP((0xfac3b00c, 1), b'wut'),
            avp.UnknownAVP(
                1, value=b'wut', vendor=0xfac3b00c,
                flags=0, name='Unknown-AVP',
            ),
        )

    def test_integer_identifier(self):
        """
        Tests we can create an AVP with a code and it defaults to the defaults
        vendor.
        """
        self._compare_avp(
            avp.AVP(1, 'Hello'),
            avp.UTF8StringAVP(
                1, value='Hello', vendor=avp.VendorId.DEFAULT,
                flags=avp.FLAG_MANDATORY,
                name='User-Name',
            ),
        )

        # Unknown AVPs default to unknown AVP
        self._compare_avp(
            avp.AVP(0xdeadb33f, b'wut'),
            avp.UnknownAVP(
                0xdeadb33f, value=b'wut',
                vendor=avp.VendorId.DEFAULT,
                flags=0, name='Unknown-AVP',
            ),
        )

    def test_string_identifier(self):
        """
        Tests we can create an AVP with a code and it defaults to the defaults
        vendor.
        """
        self._compare_avp(
            avp.AVP('User-Name', 'Hello'),
            avp.UTF8StringAVP(
                1, value='Hello', vendor=avp.VendorId.DEFAULT,
                flags=avp.FLAG_MANDATORY,
                name='User-Name',
            ),
        )

        # Unknown names will cause an error
        with self.assertRaises(ValueError):
            avp.AVP('Wut', 'error')

    def test_result_code(self):
        """
        Tests we can create an AVP with a code and it defaults to the defaults
        vendor.
        """
        result_avp = avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS)
        self.assertEqual(result_avp.value, avp.ResultCode.DIAMETER_SUCCESS)

        self._compare_avp(
            avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS),
            avp.ResultCodeAVP(
                268, 2001, vendor=avp.VendorId.DEFAULT,
                flags=avp.FLAG_MANDATORY,
                name='Result-Code',
            ),
        )

    def test_unknown_identifier(self):
        """
        Tests we get a type error if using an unsupported identifier type
        """
        # Lists are not supported and should cause an error
        with self.assertRaises(TypeError):
            avp.AVP([0, 3])

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

import abc
import socket
import struct
from enum import IntEnum, unique

from . import exception

# Constants for AVPs, also used in other places
# The length of the AVP header excluding the optional vendor-id
HEADER_LEN = 8
# Flag masks for the AVP codes
FLAG_VENDOR = 0x80
FLAG_MANDATORY = 0x40
FLAG_PROTECTED = 0x20
# Value for Product-Name AVP
PRODUCT_NAME = 'magma'


@unique
class VendorId(IntEnum):
    """
    AVP codes are defined under the namespace of a vendor. A vendor id of 0
    is default and indicates that the field should be ignored and not included
    in the message.
    """
    DEFAULT = 0
    # 3GPP
    TGPP = 10415


class BaseAVP(metaclass=abc.ABCMeta):
    """
    This is a base class for storing Attribute-Value-Pair payloads and
    convenience methods for encoding and decoding their fields, flags
    and data. This is to be implemented by the different AVP types.
    The data is encoded as defined in RFC3588 Section 4.1.

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           AVP Code                            |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |V M P r r r r r|                  AVP Length                   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                        Vendor-ID (opt)                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |    Data ...
   +-+-+-+-+-+-+-+-+

    """

    def __init__(
        self, code, value=None, flags=0, name='',
        vendor=VendorId.DEFAULT,
    ):
        self.code = code
        self.vendor = vendor
        self.name = name
        self.flags = flags
        self.payload = None
        self.value = value

    @staticmethod
    @abc.abstractmethod
    def encode_value(value):
        """
        Encode a value. To be implemented by types

        Returns:
            encoded payload bytes

        Raises:
            CodecException: if encode failed
        """
        pass

    @staticmethod
    @abc.abstractmethod
    def decode_payload(payload):
        """
        Decode a value. To be implemented by types

        Returns:
            decoded value
        Raises:
            CodecException: if decode failed
        """
        pass

    @property
    def value(self):
        """
        Decode the payload and return its value. If there is no payload,
        return None

        Returns:
            decoded value or None if no payload is set
        Raises:
            CodecException: if decode failed
        """
        if self.payload is None:
            return None

        return self.decode_payload(self.payload)

    @value.setter
    def value(self, value):
        """
        Encode the value and store it in the payload. If the value is None,
        clear the payload
        """
        if value is None:
            self.payload = None
            return
        self.payload = self.encode_value(value)

    def __repr__(self):
        flags_str = ''
        flags_str += 'M' if self.mandatory else ''
        flags_str += 'P' if self.protected else ''
        flags_str += 'V' if self.vendor_specific else ''
        return (
            "%s: type=%s vendor=%s, code=%s, "
            "flags=%s(0x%x) length=%d payload_len=%d: %s" %
            (
                self.name,
                self.__class__.__name__,
                self.vendor,
                self.code,
                flags_str,
                self.flags,
                self.length,
                self._payload_length(),
                self.value,
            )
        )

    def __eq__(self, other):
        """
        Two AVPs are equal if they respresent the same data.
        """
        return repr(self) == repr(other)

    @property
    def length(self):
        """
        Compute the length of the AVP when encoded which includes padding
        """
        return (self._encoded_length() + 3) & ~3

    def _payload_length(self):
        """
        Compute the length of the payload
        """
        if self.payload is None:
            return 0
        else:
            return len(self.payload)

    def _encoded_length(self):
        """
        Compute the length field of the AVP based which includes the length of
        the header, vendor identifier, and payload
        """
        length = HEADER_LEN
        if self.vendor_specific:
            length += 4
        length += self._payload_length()
        return length

    def validate(self):
        """
        Validates that when encoded, the AVP will be valid

        Raises:
            CodecException: the encoding error we will encounter
        """
        if self.vendor_specific and self.vendor == VendorId.DEFAULT:
            raise exception.CodecException('Vendor must be set')
        if not 0x0 <= self.code <= 0xFFFFFFFF:
            raise exception.CodecException('Code out of range')
        if self.vendor_specific and not 0x0 <= self.vendor <= 0xFFFFFFFF:
            raise exception.CodecException('Vendor out of range')
        if not 0x0 <= self._encoded_length() <= 0x00FFFFFF:
            raise exception.CodecException('AVP too long')
        if self.payload is None:
            raise exception.CodecException('Empty payload')

    def encode(self, buf, begin):
        """
        Write the AVP into a buffer at an offset

            Returns: The number of bytes written
            Raises: exception.CodecException if encoding failed
        """

        # Validate the AVP first
        self.validate()

        # Keep note of where we started
        offset = begin

        # Write the header
        struct.pack_into(
            "!II", buf, offset,
            self.code,
            self.flags << 24 | self._encoded_length(),
        )
        offset += HEADER_LEN

        # Do we need to add a vendor?
        if self.vendor_specific:
            struct.pack_into("!I", buf, offset, self.vendor)
            offset += 4

        # Payload
        payload_length = self._payload_length()
        buf[offset:offset + payload_length] = self.payload
        offset += payload_length

        # Write null bytes at end to keep 32 bit alignment
        padding_length = self.length - (offset - begin)
        struct.pack_into("x" * padding_length, buf, offset)
        offset += padding_length

        return offset - begin

    # pylint:disable=no-self-argument
    def flag_getter(mask):
        """Convenience getter for reading AVP flags"""

        def getter(self):
            return self.flags & mask != 0

        return getter

    # pylint:disable=no-self-argument
    def flag_setter(mask):
        """Convenience setter for reading AVP flags"""

        def setter(self, value):
            self.flags &= ~mask  # pylint:disable=invalid-unary-operand-type
            if value:
                self.flags |= mask

        return setter

    # Convenience flag properties
    # If an AVP is mandatory the responder must handle it to return success
    mandatory = property(
        flag_getter(FLAG_MANDATORY),
        flag_setter(FLAG_MANDATORY),
    )
    # This bit indicated the data has been end to end encrypted
    protected = property(
        flag_getter(FLAG_PROTECTED),
        flag_setter(FLAG_PROTECTED),
    )
    # This bit inidicated that there are four bytes of vendor ID before the
    # payload
    vendor_specific = property(
        flag_getter(FLAG_VENDOR),
        flag_setter(FLAG_VENDOR),
    )


class OctetStringAVP(BaseAVP):
    """Implements an Octet String AVP"""

    @staticmethod
    def decode_payload(payload):
        return bytes(payload)

    @staticmethod
    def encode_value(value):
        try:
            return bytearray(value)
        except TypeError as err:
            raise exception.CodecException(err)


class UTF8StringAVP(BaseAVP):
    """Implements an UTF-8 String AVP"""

    @staticmethod
    def decode_payload(payload):
        return bytes(payload).decode(encoding='UTF-8')

    @staticmethod
    def encode_value(value):
        try:
            return bytearray(value, 'utf-8')
        except TypeError as err:
            raise exception.CodecException(err)


class Unsigned32AVP(BaseAVP):
    """Implements an Unsigned Integer 32 AVP."""

    @staticmethod
    def decode_payload(payload):
        return struct.unpack("!I", payload)[0]

    @staticmethod
    def encode_value(value):
        try:
            return struct.pack("!I", value)
        except struct.error as err:
            raise exception.CodecException(err)


class GroupedAVP(BaseAVP):
    """Implements a Grouped AVP"""

    @staticmethod
    def decode_payload(payload):
        """Returns a list of AVPs from the decoded payload"""
        avps = []
        offset = 0
        while offset < len(payload):
            avp = decode(payload[offset:])
            offset += avp.length
            avps.append(avp)
        return avps

    @staticmethod
    def encode_value(value):
        """Encodes a list of AVPs as the new payload"""
        buf_length = sum([avp.length for avp in value])
        buf = bytearray(buf_length)
        offset = 0
        for avp in value:
            offset += avp.encode(buf, offset)
        return buf

    def filter_avps(self, vendor, code):
        """
        Return an iterator of all message AVPs that match the vendor and code

        Args:
            vendor: the AVP vendor
            code: the AVP code
        Return:
            an iterator on all AVPs that match
        """
        return filter(
            lambda element: element.vendor == vendor
            and element.code == code, self.value,
        )

    def find_avp(self, vendor, code):
        """
        Return the first AVP that matches the code and vendor

        Args:
            vendor: the AVP vendor
            code: the AVP code
        Return:
            the first AVP that matches or None if no match exists
        """
        result = list(self.filter_avps(vendor, code))
        if len(result):
            return result[0]


class AddressAVP(BaseAVP):
    @staticmethod
    def decode_payload(payload):
        """Decode the payload as IPV4 or IPV6 string"""
        if len(payload) == 2 + 4 and payload[:2] == b'\x00\x01':
            return socket.inet_ntop(socket.AF_INET, bytes(payload[2:]))
        elif len(payload) == 2 + 16 and payload[:2] == b'\x00\x02':
            return socket.inet_ntop(socket.AF_INET6, bytes(payload[2:]))
        raise exception.CodecException("Invalid Address payload")

    @staticmethod
    def encode_value(value):
        """Encode the payload given IPV4 or IPV6 string"""
        try:
            raw = socket.inet_pton(socket.AF_INET, value)
            return struct.pack("!h4s", 1, raw)
        except socket.error:
            pass

        try:
            raw = socket.inet_pton(socket.AF_INET6, value)
            return struct.pack("!h16s", 2, raw)
        except socket.error:
            pass

        raise exception.CodecException("Not a valid address")


class EnumAVP(Unsigned32AVP):
    """
    Decode Enum AVPs given a Enum type
    """

    @classmethod
    def decode_payload(cls, payload):
        """Decode as an enum if we can"""
        val = Unsigned32AVP.decode_payload(payload)
        try:
            return cls.enum(val)
        except ValueError:
            return val


@unique
class ResultCode(IntEnum):
    """
    Result-Code AVP values as defined in RFC6733 7.1
    """
    DIAMETER_SUCCESS = 2001
    DIAMETER_COMMAND_UNSUPPORTED = 3001
    DIAMETER_APPLICATION_UNSUPPORTED = 3007
    DIAMETER_ERROR_USER_UNKNOWN = 5001
    DIAMETER_MISSING_AVP = 5005
    DIAMETER_AUTHORIZATION_REJECTED = 5003
    DIAMETER_AUTHENTICATION_REJECTED = 4001
    DIAMETER_UNABLE_TO_COMPLY = 5012


class ResultCodeAVP(EnumAVP):
    """
    Decode Result-Code AVP as ResultCode Enums
    """
    enum = ResultCode


@unique
class DisconnectCause(IntEnum):
    """
    Disconnect-Cause AVP values as defined in RFC6733 5.4.3
    """
    REBOOTING = 0
    BUSY = 1
    DO_NOT_WANT_TO_TALK_TO_YOU = 2


class DisconnectCauseAVP(EnumAVP):
    """
    Decode Disconnect-Cause AVP as DisconnectCause Enums
    """
    enum = DisconnectCause


class UnknownAVP(BaseAVP):
    """
    If we dont know the AVP code, we will play it safe and not try
    to encode or decode the payload
    """

    @staticmethod
    def decode_payload(payload):
        return bytes(payload)

    @staticmethod
    def encode_value(value):
        return bytearray(value)


def AVP(ident, value=None, **kwargs):
    """
    Convenience method for constructing an AVP using an identifier. It will
    also load in default flags for the AVP from the AVPDict, as well as the
    the AVP name and type. If an AVP is not in the dictionary, an UnknownAVP
    will be used.

    Args:
        ident: this can be a integer AVP code for which the default vendor will
            be used. You can also use a tuple to identify a (vendor,code) pair.
            Lastly you can lookup an AVP by name in the AVPDict
        value: the valiue to construct the AVP with
        kwargs: additional args to pass to the AVP Constructor
    Returns:
        an AVP instance
    Raises:
        ValueError: if a string AVP identifier was used and it wasn't found
        TypeError if an identifier of unknown type is used
    """
    avp_unknown = ('Unknown-AVP', UnknownAVP, 0)
    vendor = None
    code = None
    if isinstance(ident, int):
        vendor = VendorId.DEFAULT
        code = ident
    elif isinstance(ident, tuple):
        vendor = ident[0]
        code = ident[1]
    elif isinstance(ident, str):
        vendor, code = resolve(ident)
    else:
        raise TypeError('Invalid key type')
    avp_def = AVPDict.get(vendor, {}).get(code, avp_unknown)
    name = avp_def[0]
    avp_type = avp_def[1]
    flags = kwargs.pop('flags', avp_def[2])
    return avp_type(
        code, value=value,
        vendor=vendor, flags=flags, name=name, **kwargs,
    )


def resolve(name):
    """
    Resolve an AVPs vendor and code from a name in the AVPDict

    Returns:
        a vendor, code tuple
    Raises:
        ValueError if not found
    """
    for vendor in AVPDict.keys():
        for code in AVPDict[vendor].keys():
            if AVPDict[vendor][code][0] == name:
                return vendor, code
    raise ValueError('AVP not found')


def decode(payload):
    """
    Decodes one AVP from the payload

    Args:
        payload: AVP bytestream
    Return:
        an AVP instance
    Raises:
        exception.CodecException if the AVP length was too short to decode
    """

    if len(payload) < HEADER_LEN:
        raise exception.CodecException('AVP shorter than header length')

    offset = 0

    # Load the header
    code, flags_and_length = struct.unpack_from('!II', payload, offset)
    length = flags_and_length & 0x00FFFFFF
    flags = flags_and_length >> 24
    offset += HEADER_LEN

    # Resolve the vendor
    if flags & FLAG_VENDOR != 0:
        if len(payload) - HEADER_LEN < 4:
            raise exception.CodecException('AVP too short to decode vendor')
        vendor = struct.unpack_from('!I', payload, offset)[0]
        offset += 4
    else:
        vendor = VendorId.DEFAULT

    # Lookup the type of the AVP in our dictionary or use Unknown
    avp = AVP((vendor, code), None, flags=flags)

    # Set the payload
    avp.payload = payload[offset:length]
    return avp


# The AVP dictionary the first level is the vendor identifier, and the second
# level is the avp code. The values are tuples containing the name, type and
# default flags for the AVP definition

AVPDict = {
    VendorId.DEFAULT: {
        # Diameter Base Protocol AVPs RFC3588 Section 4.5
        1: ('User-Name', UTF8StringAVP, FLAG_MANDATORY),
        257: ('Host-IP-Address', AddressAVP, FLAG_MANDATORY),
        258: ('Auth-Application-Id', Unsigned32AVP, FLAG_MANDATORY),
        260: ('Vendor-Specific-Application-Id', GroupedAVP, FLAG_MANDATORY),
        263: ('Session-Id', UTF8StringAVP, FLAG_MANDATORY),
        264: ('Origin-Host', UTF8StringAVP, FLAG_MANDATORY),
        265: ('Supported-Vendor-Id', Unsigned32AVP, FLAG_MANDATORY),
        266: ('Vendor-Id', Unsigned32AVP, FLAG_MANDATORY),
        267: ('Firmware-Revision', Unsigned32AVP, 0),
        268: ('Result-Code', ResultCodeAVP, FLAG_MANDATORY),
        269: ('Product-Name', UTF8StringAVP, 0),
        273: ('Disconnect-Cause', DisconnectCauseAVP, 0),
        277: ('Auth-Session-State', Unsigned32AVP, FLAG_MANDATORY),
        278: ('Origin-State-Id', Unsigned32AVP, FLAG_MANDATORY),
        282: ('Route-Record', UTF8StringAVP, FLAG_MANDATORY),
        283: ('Destination-Realm', UTF8StringAVP, FLAG_MANDATORY),
        293: ('Destination-Host', UTF8StringAVP, FLAG_MANDATORY),
        296: ('Origin-Realm', UTF8StringAVP, FLAG_MANDATORY),
        297: ('Experimental-Result', GroupedAVP, FLAG_MANDATORY),
        298: ('Experimental-Result-Code', Unsigned32AVP, FLAG_MANDATORY),
        299: ('Inband-Security-Id', Unsigned32AVP, FLAG_MANDATORY),
    },
    VendorId.TGPP: {
        493: (
            'Service-Selection', UTF8StringAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        # 3GPP 29.214-b80 (11.8.0 2013.03.15) Section 5.3
        515: (
            'Max-Requested-Bandwidth-DL', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        516: (
            'Max-Requested-Bandwidth-UL', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        # 3GPP 29.140-700 (7.0.0 2007.07.05) Section 6.3
        701: ('MSISDN', OctetStringAVP, FLAG_MANDATORY | FLAG_VENDOR),
        # 3GPP 29.212-c00 (12.0.0 2013.03.15) Section 5.3
        1028: (
            'QoS-Class-Identifier', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1032: (
            'RAT-Type', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1034: (
            'Allocation-Retention-Priority', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1046: (
            'Priority-Level', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1047: (
            'Pre-emption-Capability', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1048: (
            'Pre-emption-Vulnerability', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        # 3GPP 29.272-c00 (12.0.0 2013.03.13) Section 7.3
        1400: (
            'Subscription-Data', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1405: (
            'ULR-Flags', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1406: (
            'ULA-Flags', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1407: (
            'Visited-PLMN-Id', OctetStringAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1408: (
            'Requested-EUTRAN-Authentication-Info', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1410: (
            'Number-Of-Requested-Vectors', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1411: (
            'Re-Synchronization-Info', OctetStringAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1412: (
            'Immediate-Response-Preferred', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1413: (
            'Authentication-Info', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1414: (
            'E-UTRAN-Vector', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1417: (
            'Network-Access-Mode', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1423: (
            'Context-Identifier', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1424: (
            'Subscriber-Status', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1426: (
            'Access-Restriction-Data', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1428: (
            'All-APN-Configurations-Included-Indicator', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1429: (
            'APN-Configuration-Profile', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1430: (
            'APN-Configuration', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1431: (
            'EPS-Subscribed-QoS-Profile', GroupedAVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
        1435: ('AMBR', GroupedAVP, FLAG_MANDATORY | FLAG_VENDOR),
        1447: ('RAND', OctetStringAVP, FLAG_MANDATORY | FLAG_VENDOR),
        1448: ('XRES', OctetStringAVP, FLAG_MANDATORY | FLAG_VENDOR),
        1449: ('AUTN', OctetStringAVP, FLAG_MANDATORY | FLAG_VENDOR),
        1450: ('KASME', OctetStringAVP, FLAG_MANDATORY | FLAG_VENDOR),
        1456: (
            'PDN-Type', Unsigned32AVP,
            FLAG_MANDATORY | FLAG_VENDOR,
        ),
    },
}

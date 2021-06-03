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
import struct

from . import avp
from .exception import CodecException, TooShortException

# Constants for Diameter messages
HEADER_LEN = 20
# Masks for command flag bits
FLAG_REQUEST = 0x80
FLAG_PROXIABLE = 0x40
FLAG_ERROR = 0x20
FLAG_RETRANSMITTED = 0x10


def flag_getter(mask):
    """Convenience method for reading command flags"""

    def func(self):
        return self.command_flags & mask != 0
    return func


def flag_setter(mask):
    """Convenience method for setting command flags"""

    def func(self, value):
        self.command_flags &= ~mask
        if value:
            self.command_flags |= mask
    return func


class MessageHeader(object):
    """
    This is a container for a generic Diameter message's header
    as defined in RFC3588 seciton 3 with some utilities for encoding
    and decoding and constructing responses.

        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |    Version    |                 Message Length                |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       | command flags |                  Command-Code                 |
       |R P E T r r r r|                                               |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Application-ID                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      Hop-by-Hop Identifier                    |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      End-to-End Identifier                    |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |  AVPs ...
       +-+-+-+-+-+-+-+-+-+-+-+-+-
    """

    def __init__(self):
        """
        Initialize with the default where all values are zero,
        except the version number which is always 1.
        """
        self.version = 1  # The only version in existence
        self.command_flags = 0x0  # answer, not proxiable or error or retransmit
        self.command_code = 0
        self.application_id = 0
        self.hop_by_hop_id = 0
        self.end_to_end_id = 0

    @classmethod
    def copy(cls, other):
        header = cls()
        header.version = other.version
        header.command_flags = other.command_flags
        header.command_code = other.command_code
        header.application_id = other.application_id
        header.hop_by_hop_id = other.hop_by_hop_id
        header.end_to_end_id = other.end_to_end_id
        return header

    @classmethod
    def create_response_header(cls, header):
        """
        Prepare a response header to a request message header. All fields
        are copied, and the command flag request bit is cleared. This will also
        clear the error, and retransmit flags and copy the proxiable flag
        """
        resp = cls.copy(header)
        resp.command_flags = header.command_flags & FLAG_PROXIABLE
        return resp

    def validate(self, length):
        """
        Validates that when encoded, the header will be valid

        Raises:
            CodecException: the encoding error we will encounter
        """
        if not 0x0 <= length <= 0x00FFFFFF:
            raise CodecException('Length out of range')
        if not 0x0 <= self.version <= 0xFF:
            raise CodecException('Version out of range')
        if not 0x0 <= self.command_flags <= 0xFF:
            raise CodecException('Command flags out of range')
        if not 0x0 <= self.command_code <= 0x00FFFFFF:
            raise CodecException('Command code out of range')
        if not 0x0 <= self.application_id <= 0xFFFFFFFF:
            raise CodecException('Application ID out of range')
        if not 0x0 <= self.hop_by_hop_id <= 0xFFFFFFFF:
            raise CodecException('HbH ID out of range')
        if not 0x0 <= self.end_to_end_id <= 0xFFFFFFFF:
            raise CodecException('EtE ID out of range')

    def encode(self, buf, offset, length):
        """
        Encodes the diameter header into a buffer at offset

        Raises:
            CodecException: if an encoding error is encountered

        Returns:
            the number of bytes written to the buffer
        """
        self.validate(length)
        struct.pack_into(
            '!IIIII', buf, offset,
            self.version << 24 | length,
            self.command_flags << 24 | self.command_code,
            self.application_id,
            self.hop_by_hop_id,
            self.end_to_end_id,
        )
        return HEADER_LEN

    @classmethod
    def decode(cls, payload):
        """
        Decodes the first 20 bytes from a payload bytestream as a Diameter
        message header

        Args:
            payload: a byte stream
        Return:
            MessageHeader instance with respresentative data
        """
        if len(payload) < HEADER_LEN:
            raise CodecException('Payload shorter than header length')

        header_words = struct.unpack_from('!IIIII', payload, 0)

        header = cls()
        header.version = header_words[0] >> 24
        header.command_flags = header_words[1] >> 24
        header.command_code = (header_words[1] & 0x00FFFFFF)
        header.application_id = header_words[2]
        header.hop_by_hop_id = header_words[3]
        header.end_to_end_id = header_words[4]

        return header

    def __repr__(self):
        flags_str = ''
        flags_str += 'R' if self.request else ''
        flags_str += 'P' if self.proxiable else ''
        flags_str += 'E' if self.error else ''
        flags_str += 'T' if self.retransmitted else ''
        return (
            "Header version=%d, application=%d, command=%d, "
            "hbh=0x%x, ete=0x%x, flags=%s(0x%x)" %
            (
                self.version,
                self.application_id,
                self.command_code,
                self.hop_by_hop_id,
                self.end_to_end_id,
                flags_str,
                self.command_flags,
            )
        )

    def __eq__(self, other):
        """Two message headers are equal if they represent the same payload"""
        return repr(self) == repr(other)

    @property
    def length(self):
        """Message header encode length"""
        return HEADER_LEN

    request = property(
        flag_getter(FLAG_REQUEST),
        flag_setter(FLAG_REQUEST),
    )
    proxiable = property(
        flag_getter(FLAG_PROXIABLE),
        flag_setter(FLAG_PROXIABLE),
    )
    error = property(
        flag_getter(FLAG_ERROR),
        flag_setter(FLAG_ERROR),
    )
    retransmitted = property(
        flag_getter(FLAG_RETRANSMITTED),
        flag_setter(FLAG_RETRANSMITTED),
    )


class Message(object):
    """
    This is a container for Diameter messages. A message consists of a header, and
    a list of AVPs. This provides utilities for decoding a message payload into its
    constituint components, and encoding it back. There are also convenience methods
    for adding and retreiving AVPs in the instance.
    """

    def __init__(self, header=None):
        self.header = header if header else MessageHeader()
        self._avps = []

    @classmethod
    def create_response_msg(cls, msg):
        """
        Prepare a response message to a request message. All fields
        are copied, and the command flag request bit is cleared. This will also
        clear the error, and retransmit flags and copy the proxiable flag. AVPs
        will also be cleared.
        """
        resp = cls(MessageHeader.create_response_header(msg.header))
        return resp

    def __repr__(self):
        return (
            "DiameterMessage length=%d:\n\t%s\nAVPs:\n\t%s" %
            (
                self.length, self.header,
                "\n\t".join([str(x) for x in self._avps]),
            )
        )

    @property
    def length(self):
        """
        Compute the length of the message when encoded

        Returns:
            the length of the encoded message in bytes
        """
        length = self.header.length
        for avp_ in self._avps:
            length += avp_.length
        return length

    def encode(self, buf, begin):
        """
        Encodes the diameter message into a buffer at offset

        Returns:
            the number of bytes written to the stream
        Raises:
            CodecException if the encoding failed
        """
        offset = begin
        offset += self.header.encode(buf, offset, self.length)
        for avp_ in self._avps:
            offset += avp_.encode(buf, offset)
        return offset - begin

    def append_avp(self, avp_):
        """
        Append an AVP to the message

        Args:
            avp_: an AVP instance
        """
        self._avps.append(avp_)

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
            and element.code == code, self._avps,
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

    def has_fields(self, fields):
        """
        Given a list of field names validates that they are in the message

        Args:
            fields: a list of field names that are required
        Returns:
            True if the required fields are found in the message
        """
        for name in fields:
            if not self.find_avp(*avp.resolve(name)):
                return False
        return True


def decode(payload):
    """
    Decodes a diameter message from the wire

    Args:
        payload: the byte stream from the wire
    Return:
        DiameterMessage instance if the decode was successful
    Raises:
        CodecException if could not decode
        TooShortException if the payload was not long enough to decode. This is
            uniquely raised so that the we can get more data and try again
    """
    if len(payload) < HEADER_LEN:
        raise TooShortException()

    length = struct.unpack_from('!I', payload, 0)[0] & 0x00FFFFFF

    if length % 4 != 0:
        raise CodecException("Received garbage")

    if len(payload) < length:
        raise TooShortException()

    msg = Message(MessageHeader.decode(payload))

    offset = msg.header.length
    while offset < length:
        avp_entry = avp.decode(payload[offset:])
        offset += avp_entry.length
        msg.append_avp(avp_entry)

    return msg

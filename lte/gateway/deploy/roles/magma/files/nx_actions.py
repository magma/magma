# Copyright (C) 2015 Nippon Telegraph and Telephone Corporation.
# Copyright (C) 2015 YAMAMOTO Takashi <yamamoto at valinux co jp>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import struct

import six
from ryu import utils
from ryu.lib import type_desc
from ryu.lib.pack_utils import msg_pack_into
from ryu.ofproto import nicira_ext, ofproto_common
from ryu.ofproto.ofproto_parser import StringifyMixin


def generate(ofp_name, ofpp_name):
    import sys

    ofp = sys.modules[ofp_name]
    ofpp = sys.modules[ofpp_name]

    class _NXFlowSpec(StringifyMixin):
        _hdr_fmt_str = '!H'  # 2 bit 0s, 1 bit src, 2 bit dst, 11 bit n_bits
        _dst_type = None
        _subclasses = {}
        _TYPE = {
            'nx-flow-spec-field': [
                'src',
                'dst',
            ],
        }

        def __init__(self, src, dst, n_bits):
            self.src = src
            self.dst = dst
            self.n_bits = n_bits

        @classmethod
        def register(cls, subcls):
            assert issubclass(subcls, cls)
            assert subcls._dst_type not in cls._subclasses
            cls._subclasses[subcls._dst_type] = subcls

        @classmethod
        def parse(cls, buf):
            (hdr,) = struct.unpack_from(cls._hdr_fmt_str, buf, 0)
            rest = buf[struct.calcsize(cls._hdr_fmt_str):]
            if hdr == 0:
                return None, rest  # all-0 header is no-op for padding
            src_type = (hdr >> 13) & 0x1
            dst_type = (hdr >> 11) & 0x3
            n_bits = hdr & 0x3ff
            subcls = cls._subclasses[dst_type]
            if src_type == 0:  # subfield
                src = cls._parse_subfield(rest)
                rest = rest[6:]
            elif src_type == 1:  # immediate
                src_len = (n_bits + 15) // 16 * 2
                src_bin = rest[:src_len]
                src = type_desc.IntDescr(size=src_len).to_user(src_bin)
                rest = rest[src_len:]
            if dst_type == 0:  # match
                dst = cls._parse_subfield(rest)
                rest = rest[6:]
            elif dst_type == 1:  # load
                dst = cls._parse_subfield(rest)
                rest = rest[6:]
            elif dst_type == 2:  # output
                dst = ''  # empty
            return subcls(src=src, dst=dst, n_bits=n_bits), rest

        def serialize(self):
            buf = bytearray()
            if isinstance(self.src, tuple):
                src_type = 0  # subfield
            else:
                src_type = 1  # immediate
            # header
            val = (src_type << 13) | (self._dst_type << 11) | self.n_bits
            msg_pack_into(self._hdr_fmt_str, buf, 0, val)
            # src
            if src_type == 0:  # subfield
                buf += self._serialize_subfield(self.src)
            elif src_type == 1:  # immediate
                src_len = (self.n_bits + 15) // 16 * 2
                buf += type_desc.IntDescr(size=src_len).from_user(self.src)
            # dst
            if self._dst_type == 0:  # match
                buf += self._serialize_subfield(self.dst)
            elif self._dst_type == 1:  # load
                buf += self._serialize_subfield(self.dst)
            elif self._dst_type == 2:  # output
                pass  # empty
            return buf

        @staticmethod
        def _parse_subfield(buf):
            (n, len) = ofp.oxm_parse_header(buf, 0)
            assert len == 4  # only 4-bytes NXM/OXM are defined
            field = ofp.oxm_to_user_header(n)
            rest = buf[len:]
            (ofs,) = struct.unpack_from('!H', rest, 0)
            return (field, ofs)

        @staticmethod
        def _serialize_subfield(subfield):
            (field, ofs) = subfield
            buf = bytearray()
            n = ofp.oxm_from_user_header(field)
            ofp.oxm_serialize_header(n, buf, 0)
            assert len(buf) == 4  # only 4-bytes NXM/OXM are defined
            msg_pack_into('!H', buf, 4, ofs)
            return buf

    class NXFlowSpecMatch(_NXFlowSpec):
        """
        Specification for adding match criterion

        This class is used by ``NXActionLearn``.

        For the usage of this class, please refer to ``NXActionLearn``.

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        src              OXM/NXM header and Start bit for source field
        dst              OXM/NXM header and Start bit for destination field
        n_bits           The number of bits from the start bit
        ================ ======================================================
        """
        # Add a match criteria
        # an example of the corresponding ovs-ofctl syntax:
        #    NXM_OF_VLAN_TCI[0..11]
        _dst_type = 0

    class NXFlowSpecLoad(_NXFlowSpec):
        """
        Add NXAST_REG_LOAD actions

        This class is used by ``NXActionLearn``.

        For the usage of this class, please refer to ``NXActionLearn``.

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        src              OXM/NXM header and Start bit for source field
        dst              OXM/NXM header and Start bit for destination field
        n_bits           The number of bits from the start bit
        ================ ======================================================
        """
        # Add NXAST_REG_LOAD actions
        # an example of the corresponding ovs-ofctl syntax:
        #    NXM_OF_ETH_DST[]=NXM_OF_ETH_SRC[]
        _dst_type = 1

    class NXFlowSpecOutput(_NXFlowSpec):
        """
        Add an OFPAT_OUTPUT action

        This class is used by ``NXActionLearn``.

        For the usage of this class, please refer to ``NXActionLearn``.

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        src              OXM/NXM header and Start bit for source field
        dst              Must be ''
        n_bits           The number of bits from the start bit
        ================ ======================================================
        """
        # Add an OFPAT_OUTPUT action
        # an example of the corresponding ovs-ofctl syntax:
        #    output:NXM_OF_IN_PORT[]
        _dst_type = 2

        def __init__(self, src, n_bits, dst=''):
            assert dst == ''
            super(NXFlowSpecOutput, self).__init__(
                src=src, dst=dst,
                n_bits=n_bits,
            )

    class NXAction(ofpp.OFPActionExperimenter):
        _fmt_str = '!H'  # subtype
        _subtypes = {}
        _experimenter = ofproto_common.NX_EXPERIMENTER_ID

        def __init__(self):
            super(NXAction, self).__init__(self._experimenter)
            self.subtype = self._subtype

        @classmethod
        def parse(cls, buf):
            fmt_str = NXAction._fmt_str
            (subtype,) = struct.unpack_from(fmt_str, buf, 0)
            subtype_cls = cls._subtypes.get(subtype)
            rest = buf[struct.calcsize(fmt_str):]
            if subtype_cls is None:
                return NXActionUnknown(subtype, rest)
            return subtype_cls.parser(rest)

        def serialize(self, buf, offset):
            data = self.serialize_body()
            payload_offset = (
                    ofp.OFP_ACTION_EXPERIMENTER_HEADER_SIZE +
                    struct.calcsize(NXAction._fmt_str)
            )
            self.len = utils.round_up(payload_offset + len(data), 8)
            super(NXAction, self).serialize(buf, offset)
            msg_pack_into(
                NXAction._fmt_str,
                buf,
                offset + ofp.OFP_ACTION_EXPERIMENTER_HEADER_SIZE,
                self.subtype,
            )
            buf += data

        @classmethod
        def register(cls, subtype_cls):
            assert subtype_cls._subtype is not cls._subtypes
            cls._subtypes[subtype_cls._subtype] = subtype_cls

    class NXActionUnknown(NXAction):
        def __init__(
            self, subtype, data=None,
            type_=None, len_=None, experimenter=None,
        ):
            self._subtype = subtype
            super(NXActionUnknown, self).__init__()
            self.data = data

        @classmethod
        def parser(cls, buf):
            return cls(data=buf)

        def serialize_body(self):
            # fixup
            return bytearray() if self.data is None else self.data

    # For OpenFlow1.0 only
    class NXActionSetQueue(NXAction):
        r"""
        Set queue action

        This action sets the queue that should be used to queue
        when packets are output.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          set_queue:queue
        ..

        +-------------------------+
        | **set_queue**\:\ *queue*|
        +-------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        queue_id         Queue ID for the packets
        ================ ======================================================

        .. note::
            This actions is supported by
            ``OFPActionSetQueue``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionSetQueue(queue_id=10)]
        """
        _subtype = nicira_ext.NXAST_SET_QUEUE

        # queue_id
        _fmt_str = '!2xI'

        def __init__(
            self, queue_id,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionSetQueue, self).__init__()
            self.queue_id = queue_id

        @classmethod
        def parser(cls, buf):
            (queue_id,) = struct.unpack_from(cls._fmt_str, buf, 0)
            return cls(queue_id)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0, self.queue_id)
            return data

    class NXActionPopQueue(NXAction):
        """
        Pop queue action

        This action restors the queue to the value it was before any
        set_queue actions were applied.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          pop_queue
        ..

        +---------------+
        | **pop_queue** |
        +---------------+

        Example::

            actions += [parser.NXActionPopQueue()]
        """
        _subtype = nicira_ext.NXAST_POP_QUEUE

        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionPopQueue, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    class NXActionRegLoad(NXAction):
        r"""
        Load literal value action

        This action loads a literal value into a field or part of a field.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          load:value->dst[start..end]
        ..

        +-----------------------------------------------------------------+
        | **load**\:\ *value*\->\ *dst*\ **[**\ *start*\..\ *end*\ **]**  |
        +-----------------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ofs_nbits        Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``
        dst              OXM/NXM header for destination field
        value            OXM/NXM value to be loaded
        ================ ======================================================

        Example::

            actions += [parser.NXActionRegLoad(
                            ofs_nbits=nicira_ext.ofs_nbits(4, 31),
                            dst="eth_dst",
                            value=0x112233)]
        """
        _subtype = nicira_ext.NXAST_REG_LOAD
        _fmt_str = '!HIQ'  # ofs_nbits, dst, value
        _TYPE = {
            'ascii': [
                'dst',
            ],
        }

        def __init__(
            self, ofs_nbits, dst, value,
            type_=None, len_=None, experimenter=None,
            subtype=None,
        ):
            super(NXActionRegLoad, self).__init__()
            self.ofs_nbits = ofs_nbits
            self.dst = dst
            self.value = value

        @classmethod
        def parser(cls, buf):
            (ofs_nbits, dst, value) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            # Right-shift instead of using oxm_parse_header for simplicity...
            dst_name = ofp.oxm_to_user_header(dst >> 9)
            return cls(ofs_nbits, dst_name, value)

        def serialize_body(self):
            hdr_data = bytearray()
            n = ofp.oxm_from_user_header(self.dst)
            ofp.oxm_serialize_header(n, hdr_data, 0)
            (dst_num,) = struct.unpack_from('!I', six.binary_type(hdr_data), 0)

            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.ofs_nbits, dst_num, self.value,
            )
            return data

    class NXActionRegLoad2(NXAction):
        r"""
        Load literal value action

        This action loads a literal value into a field or part of a field.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          set_field:value[/mask]->dst
        ..

        +------------------------------------------------------------+
        | **set_field**\:\ *value*\ **[**\/\ *mask*\ **]**\->\ *dst* |
        +------------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        value            OXM/NXM value to be loaded
        mask             Mask for destination field
        dst              OXM/NXM header for destination field
        ================ ======================================================

        Example::

            actions += [parser.NXActionRegLoad2(dst="tun_ipv4_src",
                                                value="192.168.10.0",
                                                mask="255.255.255.0")]
        """
        _subtype = nicira_ext.NXAST_REG_LOAD2
        _TYPE = {
            'ascii': [
                'dst',
                'value',
            ],
        }

        def __init__(
            self, dst, value, mask=None,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionRegLoad2, self).__init__()
            self.dst = dst
            self.value = value
            self.mask = mask

        @classmethod
        def parser(cls, buf):
            (n, uv, mask, _len) = ofp.oxm_parse(buf, 0)
            dst, value = ofp.oxm_to_user(n, uv, mask)

            if isinstance(value, (tuple, list)):
                return cls(dst, value[0], value[1])
            else:
                return cls(dst, value, None)

        def serialize_body(self):
            data = bytearray()
            if self.mask is None:
                value = self.value
            else:
                value = (self.value, self.mask)
                self._TYPE['ascii'].append('mask')

            n, value, mask = ofp.oxm_from_user(self.dst, value)
            len_ = ofp.oxm_serialize(n, value, mask, data, 0)
            msg_pack_into("!%dx" % (14 - len_), data, len_)

            return data

    class NXActionNote(NXAction):
        r"""
        Note action

        This action does nothing at all.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          note:[hh]..
        ..

        +-----------------------------------+
        | **note**\:\ **[**\ *hh*\ **]**\.. |
        +-----------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        note             A list of integer type values
        ================ ======================================================

        Example::

            actions += [parser.NXActionNote(note=[0xaa,0xbb,0xcc,0xdd])]
        """
        _subtype = nicira_ext.NXAST_NOTE

        # note
        _fmt_str = '!%dB'

        # set the integer array in a note
        def __init__(
            self,
            note,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionNote, self).__init__()
            self.note = note

        @classmethod
        def parser(cls, buf):
            note = struct.unpack_from(
                cls._fmt_str % len(buf), buf, 0,
            )
            return cls(list(note))

        def serialize_body(self):
            assert isinstance(self.note, (tuple, list))
            for n in self.note:
                assert isinstance(n, six.integer_types)

            pad = (len(self.note) + nicira_ext.NX_ACTION_HEADER_0_SIZE) % 8
            if pad:
                self.note += [0x0 for i in range(8 - pad)]
            note_len = len(self.note)
            data = bytearray()
            msg_pack_into(
                self._fmt_str % note_len, data, 0,
                *self.note,
            )
            return data

    class _NXActionSetTunnelBase(NXAction):
        # _subtype, _fmt_str must be attributes of subclass.

        def __init__(
            self,
            tun_id,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(_NXActionSetTunnelBase, self).__init__()
            self.tun_id = tun_id

        @classmethod
        def parser(cls, buf):
            (tun_id,) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(tun_id)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.tun_id,
            )
            return data

    class NXActionSetTunnel(_NXActionSetTunnelBase):
        r"""
        Set Tunnel action

        This action sets the identifier (such as GRE) to the specified id.

        And equivalent to the followings action of ovs-ofctl command.

        .. note::
            This actions is supported by
            ``OFPActionSetField``
            in OpenFlow1.2 or later.

        ..
          set_tunnel:id
        ..

        +------------------------+
        | **set_tunnel**\:\ *id* |
        +------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        tun_id           Tunnel ID(32bits)
        ================ ======================================================

        Example::

            actions += [parser.NXActionSetTunnel(tun_id=0xa)]
        """
        _subtype = nicira_ext.NXAST_SET_TUNNEL

        # tun_id
        _fmt_str = '!2xI'

    class NXActionSetTunnel64(_NXActionSetTunnelBase):
        r"""
        Set Tunnel action

        This action outputs to a port that encapsulates
        the packet in a tunnel.

        And equivalent to the followings action of ovs-ofctl command.

        .. note::
            This actions is supported by
            ``OFPActionSetField``
            in OpenFlow1.2 or later.

        ..
          set_tunnel64:id
        ..

        +--------------------------+
        | **set_tunnel64**\:\ *id* |
        +--------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        tun_id           Tunnel ID(64bits)
        ================ ======================================================

        Example::

            actions += [parser.NXActionSetTunnel64(tun_id=0xa)]
        """
        _subtype = nicira_ext.NXAST_SET_TUNNEL64

        # tun_id
        _fmt_str = '!6xQ'

    class NXActionRegMove(NXAction):
        r"""
        Move register action

        This action copies the src to dst.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          move:src[start..end]->dst[start..end]
        ..

        +--------------------------------------------------------+
        | **move**\:\ *src*\ **[**\ *start*\..\ *end*\ **]**\->\ |
        | *dst*\ **[**\ *start*\..\ *end* \ **]**                |
        +--------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        src_field        OXM/NXM header for source field
        dst_field        OXM/NXM header for destination field
        n_bits           Number of bits
        src_ofs          Starting bit offset in source
        dst_ofs          Starting bit offset in destination
        ================ ======================================================

        .. CAUTION::
            **src_start**\  and \ **src_end**\  difference and \ **dst_start**\
             and \ **dst_end**\  difference must be the same.

        Example::

            actions += [parser.NXActionRegMove(src_field="reg0",
                                               dst_field="reg1",
                                               n_bits=5,
                                               src_ofs=0
                                               dst_ofs=10)]
        """
        _subtype = nicira_ext.NXAST_REG_MOVE
        _fmt_str = '!HHH'  # n_bits, src_ofs, dst_ofs
        # Followed by OXM fields (src, dst) and padding to 8 bytes boundary
        _TYPE = {
            'ascii': [
                'src_field',
                'dst_field',
            ],
        }

        def __init__(
            self, src_field, dst_field, n_bits, src_ofs=0, dst_ofs=0,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionRegMove, self).__init__()
            self.n_bits = n_bits
            self.src_ofs = src_ofs
            self.dst_ofs = dst_ofs
            self.src_field = src_field
            self.dst_field = dst_field

        @classmethod
        def parser(cls, buf):
            (n_bits, src_ofs, dst_ofs) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            rest = buf[struct.calcsize(NXActionRegMove._fmt_str):]

            # src field
            (n, len) = ofp.oxm_parse_header(rest, 0)
            src_field = ofp.oxm_to_user_header(n)
            rest = rest[len:]
            # dst field
            (n, len) = ofp.oxm_parse_header(rest, 0)
            dst_field = ofp.oxm_to_user_header(n)
            rest = rest[len:]
            # ignore padding
            return cls(
                src_field, dst_field=dst_field, n_bits=n_bits,
                src_ofs=src_ofs, dst_ofs=dst_ofs,
            )

        def serialize_body(self):
            # fixup
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.n_bits, self.src_ofs, self.dst_ofs,
            )
            # src field
            n = ofp.oxm_from_user_header(self.src_field)
            ofp.oxm_serialize_header(n, data, len(data))
            # dst field
            n = ofp.oxm_from_user_header(self.dst_field)
            ofp.oxm_serialize_header(n, data, len(data))
            return data

    class NXActionResubmit(NXAction):
        r"""
        Resubmit action

        This action searches one of the switch's flow tables.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          resubmit:port
        ..

        +------------------------+
        | **resubmit**\:\ *port* |
        +------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        in_port          New in_port for checking flow table
        ================ ======================================================

        Example::

            actions += [parser.NXActionResubmit(in_port=8080)]
        """
        _subtype = nicira_ext.NXAST_RESUBMIT

        # in_port
        _fmt_str = '!H4x'

        def __init__(
            self,
            in_port=0xfff8,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionResubmit, self).__init__()
            self.in_port = in_port

        @classmethod
        def parser(cls, buf):
            (in_port,) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(in_port)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.in_port,
            )
            return data

    class NXActionResubmitTable(NXAction):
        r"""
        Resubmit action

        This action searches one of the switch's flow tables.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          resubmit([port],[table])
        ..

        +------------------------------------------------+
        | **resubmit(**\[\ *port*\]\,[\ *table*\]\ **)** |
        +------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        in_port          New in_port for checking flow table
        table_id         Checking flow tables
        ================ ======================================================

        Example::

            actions += [parser.NXActionResubmit(in_port=8080,
                                                table_id=10)]
        """
        _subtype = nicira_ext.NXAST_RESUBMIT_TABLE

        # in_port, table_id
        _fmt_str = '!HB3x'

        def __init__(
            self,
            in_port=0xfff8,
            table_id=0xff,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionResubmitTable, self).__init__()
            self.in_port = in_port
            self.table_id = table_id

        @classmethod
        def parser(cls, buf):
            (
                in_port,
                table_id,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(in_port, table_id)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.in_port, self.table_id,
            )
            return data

    class NXActionOutputReg(NXAction):
        r"""
        Add output action

        This action outputs the packet to the OpenFlow port number read from
        src.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          output:src[start...end]
        ..

        +-------------------------------------------------------+
        | **output**\:\ *src*\ **[**\ *start*\...\ *end*\ **]** |
        +-------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ofs_nbits        Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``
        src              OXM/NXM header for source field
        max_len          Max length to send to controller
        ================ ======================================================

        Example::

            actions += [parser.NXActionOutputReg(
                            ofs_nbits=nicira_ext.ofs_nbits(4, 31),
                            src="reg0",
                            max_len=1024)]
        """
        _subtype = nicira_ext.NXAST_OUTPUT_REG

        # ofs_nbits, src, max_len
        _fmt_str = '!H4sH6x'
        _TYPE = {
            'ascii': [
                'src',
            ],
        }

        def __init__(
            self,
            ofs_nbits,
            src,
            max_len,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionOutputReg, self).__init__()
            self.ofs_nbits = ofs_nbits
            self.src = src
            self.max_len = max_len

        @classmethod
        def parser(cls, buf):
            (ofs_nbits, oxm_data, max_len) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            (n, len_) = ofp.oxm_parse_header(oxm_data, 0)
            src = ofp.oxm_to_user_header(n)
            return cls(
                ofs_nbits,
                src,
                max_len,
            )

        def serialize_body(self):
            data = bytearray()
            src = bytearray()
            oxm = ofp.oxm_from_user_header(self.src)
            ofp.oxm_serialize_header(oxm, src, 0),
            msg_pack_into(
                self._fmt_str, data, 0,
                self.ofs_nbits,
                six.binary_type(src),
                self.max_len,
            )
            return data

    class NXActionOutputReg2(NXAction):
        r"""
        Add output action

        This action outputs the packet to the OpenFlow port number read from
        src.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          output:src[start...end]
        ..

        +-------------------------------------------------------+
        | **output**\:\ *src*\ **[**\ *start*\...\ *end*\ **]** |
        +-------------------------------------------------------+

        .. NOTE::
             Like the ``NXActionOutputReg`` but organized so
             that there is room for a 64-bit experimenter OXM as 'src'.

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ofs_nbits        Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``
        src              OXM/NXM header for source field
        max_len          Max length to send to controller
        ================ ======================================================

        Example::

            actions += [parser.NXActionOutputReg2(
                            ofs_nbits=nicira_ext.ofs_nbits(4, 31),
                            src="reg0",
                            max_len=1024)]
        """
        _subtype = nicira_ext.NXAST_OUTPUT_REG2

        # ofs_nbits, src, max_len
        _fmt_str = '!HH4s'
        _TYPE = {
            'ascii': [
                'src',
            ],
        }

        def __init__(
            self,
            ofs_nbits,
            src,
            max_len,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionOutputReg2, self).__init__()
            self.ofs_nbits = ofs_nbits
            self.src = src
            self.max_len = max_len

        @classmethod
        def parser(cls, buf):
            (
                ofs_nbits,
                max_len,
                oxm_data,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            (n, len_) = ofp.oxm_parse_header(oxm_data, 0)
            src = ofp.oxm_to_user_header(n)
            return cls(
                ofs_nbits,
                src,
                max_len,
            )

        def serialize_body(self):
            data = bytearray()
            oxm_data = bytearray()
            oxm = ofp.oxm_from_user_header(self.src)
            ofp.oxm_serialize_header(oxm, oxm_data, 0),
            msg_pack_into(
                self._fmt_str, data, 0,
                self.ofs_nbits,
                self.max_len,
                six.binary_type(oxm_data),
            )
            offset = len(data)
            msg_pack_into("!%dx" % (14 - offset), data, offset)
            return data

    class NXActionLearn(NXAction):
        r"""
        Adds or modifies flow action

        This action adds or modifies a flow in OpenFlow table.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          learn(argument[,argument]...)
        ..

        +---------------------------------------------------+
        | **learn(**\ *argument*\[,\ *argument*\]...\ **)** |
        +---------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        table_id         The table in which the new flow should be inserted
        specs            Adds a match criterion to the new flow

                         Please use the
                         ``NXFlowSpecMatch``
                         in order to set the following format

                         ..
                           field=value
                           field[start..end]=src[start..end]
                           field[start..end]
                         ..

                         | *field*\=\ *value*
                         | *field*\ **[**\ *start*\..\ *end*\ **]**\  =\
                         *src*\ **[**\ *start*\..\ *end*\ **]**
                         | *field*\ **[**\ *start*\..\ *end*\ **]**
                         |

                         Please use the
                         ``NXFlowSpecLoad``
                         in order to set the following format

                         ..
                           load:value->dst[start..end]
                           load:src[start..end]->dst[start..end]
                         ..

                         | **load**\:\ *value*\ **->**\ *dst*\
                         **[**\ *start*\..\ *end*\ **]**
                         | **load**\:\ *src*\ **[**\ *start*\..\ *end*\
                         **] ->**\ *dst*\ **[**\ *start*\..\ *end*\ **]**
                         |

                         Please use the
                         ``NXFlowSpecOutput``
                         in order to set the following format

                         ..
                           output:field[start..end]
                         ..

                         | **output:**\ field\ **[**\ *start*\..\ *end*\ **]**

        idle_timeout     Idle time before discarding(seconds)
        hard_timeout     Max time before discarding(seconds)
        priority         Priority level of flow entry
        cookie           Cookie for new flow
        flags            send_flow_rem
        fin_idle_timeout Idle timeout after FIN(seconds)
        fin_hard_timeout Hard timeout after FIN(seconds)
        ================ ======================================================

        .. CAUTION::
            The arguments specify the flow's match fields, actions,
            and other properties, as follows.
            At least one match criterion and one action argument
            should ordinarily be specified.

        Example::

            actions += [
                parser.NXActionLearn(able_id=10,
                     specs=[parser.NXFlowSpecMatch(src=0x800,
                                                   dst=('eth_type_nxm', 0),
                                                   n_bits=16),
                            parser.NXFlowSpecMatch(src=('reg1', 1),
                                                   dst=('reg2', 3),
                                                   n_bits=5),
                            parser.NXFlowSpecMatch(src=('reg3', 1),
                                                   dst=('reg3', 1),
                                                   n_bits=5),
                            parser.NXFlowSpecLoad(src=0,
                                                  dst=('reg4', 3),
                                                  n_bits=5),
                            parser.NXFlowSpecLoad(src=('reg5', 1),
                                                  dst=('reg6', 3),
                                                  n_bits=5),
                            parser.NXFlowSpecOutput(src=('reg7', 1),
                                                    dst="",
                                                    n_bits=5)],
                     idle_timeout=180,
                     hard_timeout=300,
                     priority=1,
                     cookie=0x64,
                     flags=ofproto.OFPFF_SEND_FLOW_REM,
                     fin_idle_timeout=180,
                     fin_hard_timeout=300)]
        """
        _subtype = nicira_ext.NXAST_LEARN

        # idle_timeout, hard_timeout, priority, cookie, flags,
        # table_id, pad, fin_idle_timeout, fin_hard_timeout
        _fmt_str = '!HHHQHBxHH'
        # Followed by flow_mod_specs

        def __init__(
            self,
            table_id,
            specs,
            idle_timeout=0,
            hard_timeout=0,
            priority=ofp.OFP_DEFAULT_PRIORITY,
            cookie=0,
            flags=0,
            fin_idle_timeout=0,
            fin_hard_timeout=0,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionLearn, self).__init__()
            self.idle_timeout = idle_timeout
            self.hard_timeout = hard_timeout
            self.priority = priority
            self.cookie = cookie
            self.flags = flags
            self.table_id = table_id
            self.fin_idle_timeout = fin_idle_timeout
            self.fin_hard_timeout = fin_hard_timeout
            self.specs = specs

        @classmethod
        def parser(cls, buf):
            (
                idle_timeout,
                hard_timeout,
                priority,
                cookie,
                flags,
                table_id,
                fin_idle_timeout,
                fin_hard_timeout,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            rest = buf[struct.calcsize(cls._fmt_str):]
            # specs
            specs = []
            while len(rest) > 0:
                spec, rest = _NXFlowSpec.parse(rest)
                if spec is None:
                    continue
                specs.append(spec)
            return cls(
                idle_timeout=idle_timeout,
                hard_timeout=hard_timeout,
                priority=priority,
                cookie=cookie,
                flags=flags,
                table_id=table_id,
                fin_idle_timeout=fin_idle_timeout,
                fin_hard_timeout=fin_hard_timeout,
                specs=specs,
            )

        def serialize_body(self):
            # fixup
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.idle_timeout,
                self.hard_timeout,
                self.priority,
                self.cookie,
                self.flags,
                self.table_id,
                self.fin_idle_timeout,
                self.fin_hard_timeout,
            )
            for spec in self.specs:
                data += spec.serialize()
            return data

    class NXActionExit(NXAction):
        """
        Halt action

        This action causes OpenvSwitch to immediately halt
        execution of further actions.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          exit
        ..

        +----------+
        | **exit** |
        +----------+

        Example::

            actions += [parser.NXActionExit()]
        """
        _subtype = nicira_ext.NXAST_EXIT

        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionExit, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    # For OpenFlow1.0 only
    class NXActionDecTtl(NXAction):
        """
        Decrement IP TTL action

        This action decrements TTL of IPv4 packet or
        hop limit of IPv6 packet.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          dec_ttl
        ..

        +-------------+
        | **dec_ttl** |
        +-------------+

        .. NOTE::
            This actions is supported by
            ``OFPActionDecNwTtl``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionDecTtl()]
        """
        _subtype = nicira_ext.NXAST_DEC_TTL

        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionDecTtl, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    class NXActionController(NXAction):
        r"""
        Send packet in message action

        This action sends the packet to the OpenFlow controller as
        a packet in message.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          controller(key=value...)
        ..

        +----------------------------------------------+
        | **controller(**\ *key*\=\ *value*\...\ **)** |
        +----------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        max_len          Max length to send to controller
        controller_id    Controller ID to send packet-in
        reason           Reason for sending the message
        ================ ======================================================

        Example::

            actions += [
                parser.NXActionController(max_len=1024,
                                          controller_id=1,
                                          reason=ofproto.OFPR_INVALID_TTL)]
        """
        _subtype = nicira_ext.NXAST_CONTROLLER

        # max_len, controller_id, reason
        _fmt_str = '!HHBx'

        def __init__(
            self,
            max_len,
            controller_id,
            reason,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionController, self).__init__()
            self.max_len = max_len
            self.controller_id = controller_id
            self.reason = reason

        @classmethod
        def parser(cls, buf):
            (
                max_len,
                controller_id,
                reason,
            ) = struct.unpack_from(
                cls._fmt_str, buf,
            )
            return cls(
                max_len,
                controller_id,
                reason,
            )

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.max_len,
                self.controller_id,
                self.reason,
            )
            return data

    class NXActionController2(NXAction):
        r"""
        Send packet in message action

        This action sends the packet to the OpenFlow controller as
        a packet in message.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          controller(key=value...)
        ..

        +----------------------------------------------+
        | **controller(**\ *key*\=\ *value*\...\ **)** |
        +----------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        max_len          Max length to send to controller
        controller_id    Controller ID to send packet-in
        reason           Reason for sending the message
        userdata         Additional data to the controller in the packet-in
                         message
        pause            Flag to pause pipeline to resume later
        ================ ======================================================

        Example::

            actions += [
                parser.NXActionController(max_len=1024,
                                          controller_id=1,
                                          reason=ofproto.OFPR_INVALID_TTL,
                                          userdata=[0xa,0xb,0xc],
                                          pause=True)]
        """
        _subtype = nicira_ext.NXAST_CONTROLLER2
        _fmt_str = '!6x'
        _PACK_STR = '!HH'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
            **kwargs
        ):
            super(NXActionController2, self).__init__()

            for arg in kwargs:
                if arg in NXActionController2Prop._NAMES:
                    setattr(self, arg, kwargs[arg])

        @classmethod
        def parser(cls, buf):
            cls_data = {}
            offset = 6
            buf_len = len(buf)
            while buf_len > offset:
                (type_, length) = struct.unpack_from(cls._PACK_STR, buf, offset)
                offset += 4
                try:
                    subcls = NXActionController2Prop._TYPES[type_]
                except KeyError:
                    subcls = NXActionController2PropUnknown
                data, size = subcls.parser_prop(buf[offset:], length - 4)
                offset += size
                cls_data[subcls._arg_name] = data
            return cls(**cls_data)

        def serialize_body(self):
            body = bytearray()
            msg_pack_into(self._fmt_str, body, 0)
            prop_list = []
            for arg in self.__dict__:
                if arg in NXActionController2Prop._NAMES:
                    prop_list.append((
                        NXActionController2Prop._NAMES[arg],
                        self.__dict__[arg],
                    ))
            prop_list.sort(key=lambda x: x[0].type)

            for subcls, value in prop_list:
                body += subcls.serialize_prop(value)

            return body

    class NXActionController2Prop(object):
        _TYPES = {}
        _NAMES = {}

        @classmethod
        def register_type(cls, type_):
            def _register_type(subcls):
                subcls.type = type_
                NXActionController2Prop._TYPES[type_] = subcls
                NXActionController2Prop._NAMES[subcls._arg_name] = subcls
                return subcls

            return _register_type

    class NXActionController2PropUnknown(NXActionController2Prop):

        @classmethod
        def parser_prop(cls, buf, length):
            size = 4
            return buf, size

        @classmethod
        def serialize_prop(cls, argment):
            data = bytearray()
            return data

    @NXActionController2Prop.register_type(nicira_ext.NXAC2PT_MAX_LEN)
    class NXActionController2PropMaxLen(NXActionController2Prop):
        # max_len
        _fmt_str = "!H2x"
        _arg_name = "max_len"

        @classmethod
        def parser_prop(cls, buf, length):
            size = 4
            (max_len,) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return max_len, size

        @classmethod
        def serialize_prop(cls, max_len):
            data = bytearray()
            msg_pack_into(
                "!HHH2x", data, 0,
                nicira_ext.NXAC2PT_MAX_LEN,
                8,
                max_len,
            )
            return data

    @NXActionController2Prop.register_type(nicira_ext.NXAC2PT_CONTROLLER_ID)
    class NXActionController2PropControllerId(NXActionController2Prop):
        # controller_id
        _fmt_str = "!H2x"
        _arg_name = "controller_id"

        @classmethod
        def parser_prop(cls, buf, length):
            size = 4
            (controller_id,) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return controller_id, size

        @classmethod
        def serialize_prop(cls, controller_id):
            data = bytearray()
            msg_pack_into(
                "!HHH2x", data, 0,
                nicira_ext.NXAC2PT_CONTROLLER_ID,
                8,
                controller_id,
            )
            return data

    @NXActionController2Prop.register_type(nicira_ext.NXAC2PT_REASON)
    class NXActionController2PropReason(NXActionController2Prop):
        # reason
        _fmt_str = "!B3x"
        _arg_name = "reason"

        @classmethod
        def parser_prop(cls, buf, length):
            size = 4
            (reason,) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return reason, size

        @classmethod
        def serialize_prop(cls, reason):
            data = bytearray()
            msg_pack_into(
                "!HHB3x", data, 0,
                nicira_ext.NXAC2PT_REASON,
                5,
                reason,
            )
            return data

    @NXActionController2Prop.register_type(nicira_ext.NXAC2PT_USERDATA)
    class NXActionController2PropUserData(NXActionController2Prop):
        # userdata
        _fmt_str = "!B"
        _arg_name = "userdata"

        @classmethod
        def parser_prop(cls, buf, length):
            userdata = []
            offset = 0

            while offset < length:
                u = struct.unpack_from(cls._fmt_str, buf, offset)
                userdata.append(u[0])
                offset += 1

            user_size = utils.round_up(length, 4)

            if user_size > 4 and (user_size % 8) == 0:
                size = utils.round_up(length, 4) + 4
            else:
                size = utils.round_up(length, 4)

            return userdata, size

        @classmethod
        def serialize_prop(cls, userdata):
            data = bytearray()
            user_buf = bytearray()
            user_offset = 0
            for user in userdata:
                msg_pack_into(
                    '!B', user_buf, user_offset,
                    user,
                )
                user_offset += 1

            msg_pack_into(
                "!HH", data, 0,
                nicira_ext.NXAC2PT_USERDATA,
                4 + user_offset,
            )
            data += user_buf

            if user_offset > 4:
                user_len = utils.round_up(user_offset, 4)
                brank_size = 0
                if (user_len % 8) == 0:
                    brank_size = 4
                msg_pack_into(
                    "!%dx" % (user_len - user_offset + brank_size),
                    data, 4 + user_offset,
                )
            else:
                user_len = utils.round_up(user_offset, 4)

                msg_pack_into(
                    "!%dx" % (user_len - user_offset),
                    data, 4 + user_offset,
                )
            return data

    @NXActionController2Prop.register_type(nicira_ext.NXAC2PT_PAUSE)
    class NXActionController2PropPause(NXActionController2Prop):
        _arg_name = "pause"

        @classmethod
        def parser_prop(cls, buf, length):
            pause = True
            size = 4
            return pause, size

        @classmethod
        def serialize_prop(cls, pause):
            data = bytearray()
            msg_pack_into(
                "!HH4x", data, 0,
                nicira_ext.NXAC2PT_PAUSE,
                4,
            )
            return data

    class NXActionDecTtlCntIds(NXAction):
        r"""
        Decrement TTL action

        This action decrements TTL of IPv4 packet or
        hop limits of IPv6 packet.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          dec_ttl(id1[,id2]...)
        ..

        +-------------------------------------------+
        | **dec_ttl(**\ *id1*\[,\ *id2*\]...\ **)** |
        +-------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        cnt_ids          Controller ids
        ================ ======================================================

        Example::

            actions += [parser.NXActionDecTtlCntIds(cnt_ids=[1,2,3])]

        .. NOTE::
            If you want to set the following ovs-ofctl command.
            Please use ``OFPActionDecNwTtl``.

        +-------------+
        | **dec_ttl** |
        +-------------+
        """
        _subtype = nicira_ext.NXAST_DEC_TTL_CNT_IDS

        # controllers
        _fmt_str = '!H4x'
        _fmt_len = 6

        def __init__(
            self,
            cnt_ids,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionDecTtlCntIds, self).__init__()

            self.cnt_ids = cnt_ids

        @classmethod
        def parser(cls, buf):
            (controllers,) = struct.unpack_from(
                cls._fmt_str, buf,
            )

            offset = cls._fmt_len
            cnt_ids = []

            for i in range(0, controllers):
                id_ = struct.unpack_from('!H', buf, offset)
                cnt_ids.append(id_[0])
                offset += 2

            return cls(cnt_ids)

        def serialize_body(self):
            assert isinstance(self.cnt_ids, (tuple, list))
            for i in self.cnt_ids:
                assert isinstance(i, six.integer_types)

            controllers = len(self.cnt_ids)

            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                controllers,
            )
            offset = self._fmt_len

            for id_ in self.cnt_ids:
                msg_pack_into('!H', data, offset, id_)
                offset += 2

            id_len = (
                utils.round_up(controllers, 4) -
                controllers
            )

            if id_len != 0:
                msg_pack_into('%dx' % id_len * 2, data, offset)

            return data

    # Use in only OpenFlow1.0
    class NXActionMplsBase(NXAction):
        # ethertype
        _fmt_str = '!H4x'

        def __init__(
            self,
            ethertype,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionMplsBase, self).__init__()
            self.ethertype = ethertype

        @classmethod
        def parser(cls, buf):
            (ethertype,) = struct.unpack_from(
                cls._fmt_str, buf,
            )
            return cls(ethertype)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.ethertype,
            )
            return data

    # For OpenFlow1.0 only
    class NXActionPushMpls(NXActionMplsBase):
        r"""
        Push MPLS action

        This action pushes a new MPLS header to the packet.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          push_mpls:ethertype
        ..

        +-------------------------------+
        | **push_mpls**\:\ *ethertype*  |
        +-------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ethertype        Ether type(The value must be either 0x8847 or 0x8848)
        ================ ======================================================

        .. NOTE::
            This actions is supported by
            ``OFPActionPushMpls``
            in OpenFlow1.2 or later.

        Example::

            match = parser.OFPMatch(dl_type=0x0800)
            actions += [parser.NXActionPushMpls(ethertype=0x8847)]
        """
        _subtype = nicira_ext.NXAST_PUSH_MPLS

    # For OpenFlow1.0 only
    class NXActionPopMpls(NXActionMplsBase):
        r"""
        Pop MPLS action

        This action pops the MPLS header from the packet.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          pop_mpls:ethertype
        ..

        +------------------------------+
        | **pop_mpls**\:\ *ethertype*  |
        +------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ethertype        Ether type
        ================ ======================================================

        .. NOTE::
            This actions is supported by
            ``OFPActionPopMpls``
            in OpenFlow1.2 or later.

        Example::

            match = parser.OFPMatch(dl_type=0x8847)
            actions += [parser.NXActionPushMpls(ethertype=0x0800)]
        """
        _subtype = nicira_ext.NXAST_POP_MPLS

    # For OpenFlow1.0 only
    class NXActionSetMplsTtl(NXAction):
        r"""
        Set MPLS TTL action

        This action sets the MPLS TTL.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          set_mpls_ttl:ttl
        ..

        +---------------------------+
        | **set_mpls_ttl**\:\ *ttl* |
        +---------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        ttl              MPLS TTL
        ================ ======================================================

        .. NOTE::
            This actions is supported by
            ``OFPActionSetMplsTtl``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionSetMplsTil(ttl=128)]
        """
        _subtype = nicira_ext.NXAST_SET_MPLS_TTL

        # ethertype
        _fmt_str = '!B5x'

        def __init__(
            self,
            ttl,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionSetMplsTtl, self).__init__()
            self.ttl = ttl

        @classmethod
        def parser(cls, buf):
            (ttl,) = struct.unpack_from(
                cls._fmt_str, buf,
            )
            return cls(ttl)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.ttl,
            )
            return data

    # For OpenFlow1.0 only
    class NXActionDecMplsTtl(NXAction):
        """
        Decrement MPLS TTL action

        This action decrements the MPLS TTL.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          dec_mpls_ttl
        ..

        +------------------+
        | **dec_mpls_ttl** |
        +------------------+

        .. NOTE::
            This actions is supported by
            ``OFPActionDecMplsTtl``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionDecMplsTil()]
        """
        _subtype = nicira_ext.NXAST_DEC_MPLS_TTL

        # ethertype
        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionDecMplsTtl, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    # For OpenFlow1.0 only
    class NXActionSetMplsLabel(NXAction):
        r"""
        Set MPLS Lavel action

        This action sets the MPLS Label.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          set_mpls_label:label
        ..

        +-------------------------------+
        | **set_mpls_label**\:\ *label* |
        +-------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        label            MPLS Label
        ================ ======================================================

        .. NOTE::
            This actions is supported by
            ``OFPActionSetField(mpls_label=label)``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionSetMplsLabel(label=0x10)]
        """
        _subtype = nicira_ext.NXAST_SET_MPLS_LABEL

        # ethertype
        _fmt_str = '!2xI'

        def __init__(
            self,
            label,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionSetMplsLabel, self).__init__()
            self.label = label

        @classmethod
        def parser(cls, buf):
            (label,) = struct.unpack_from(
                cls._fmt_str, buf,
            )
            return cls(label)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.label,
            )
            return data

    # For OpenFlow1.0 only
    class NXActionSetMplsTc(NXAction):
        r"""
        Set MPLS Tc action

        This action sets the MPLS Tc.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          set_mpls_tc:tc
        ..

        +-------------------------+
        | **set_mpls_tc**\:\ *tc* |
        +-------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        tc               MPLS Tc
        ================ ======================================================

        .. NOTE::
            This actions is supported by
            ``OFPActionSetField(mpls_label=tc)``
            in OpenFlow1.2 or later.

        Example::

            actions += [parser.NXActionSetMplsLabel(tc=0x10)]
        """
        _subtype = nicira_ext.NXAST_SET_MPLS_TC

        # ethertype
        _fmt_str = '!B5x'

        def __init__(
            self,
            tc,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionSetMplsTc, self).__init__()
            self.tc = tc

        @classmethod
        def parser(cls, buf):
            (tc,) = struct.unpack_from(
                cls._fmt_str, buf,
            )
            return cls(tc)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.tc,
            )
            return data

    class NXActionStackBase(NXAction):
        # start, field, end
        _fmt_str = '!H4sH'
        _TYPE = {
            'ascii': [
                'field',
            ],
        }

        def __init__(
            self,
            field,
            start,
            end,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionStackBase, self).__init__()
            self.field = field
            self.start = start
            self.end = end

        @classmethod
        def parser(cls, buf):
            (start, oxm_data, end) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            (n, len_) = ofp.oxm_parse_header(oxm_data, 0)
            field = ofp.oxm_to_user_header(n)
            return cls(field, start, end)

        def serialize_body(self):
            data = bytearray()
            oxm_data = bytearray()
            oxm = ofp.oxm_from_user_header(self.field)
            ofp.oxm_serialize_header(oxm, oxm_data, 0)
            msg_pack_into(
                self._fmt_str, data, 0,
                self.start,
                six.binary_type(oxm_data),
                self.end,
            )
            offset = len(data)
            msg_pack_into("!%dx" % (12 - offset), data, offset)
            return data

    class NXActionStackPush(NXActionStackBase):
        r"""
        Push field action

        This action pushes field to top of the stack.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          pop:dst[start...end]
        ..

        +----------------------------------------------------+
        | **pop**\:\ *dst*\ **[**\ *start*\...\ *end*\ **]** |
        +----------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        field            OXM/NXM header for source field
        start            Start bit for source field
        end              End bit for source field
        ================ ======================================================

        Example::

            actions += [parser.NXActionStackPush(field="reg2",
                                                 start=0,
                                                 end=5)]
        """
        _subtype = nicira_ext.NXAST_STACK_PUSH

    class NXActionStackPop(NXActionStackBase):
        r"""
        Pop field action

        This action pops field from top of the stack.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          pop:src[start...end]
        ..

        +----------------------------------------------------+
        | **pop**\:\ *src*\ **[**\ *start*\...\ *end*\ **]** |
        +----------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        field            OXM/NXM header for destination field
        start            Start bit for destination field
        end              End bit for destination field
        ================ ======================================================

        Example::

            actions += [parser.NXActionStackPop(field="reg2",
                                                start=0,
                                                end=5)]
        """
        _subtype = nicira_ext.NXAST_STACK_POP

    class NXActionSample(NXAction):
        r"""
        Sample packets action

        This action samples packets and sends one sample for
        every sampled packet.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          sample(argument[,argument]...)
        ..

        +----------------------------------------------------+
        | **sample(**\ *argument*\[,\ *argument*\]...\ **)** |
        +----------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        probability      The number of sampled packets
        collector_set_id The unsigned 32-bit integer identifier of
                         the set of sample collectors to send sampled packets
                         to
        obs_domain_id    The Unsigned 32-bit integer Observation Domain ID
        obs_point_id     The unsigned 32-bit integer Observation Point ID
        ================ ======================================================

        Example::

            actions += [parser.NXActionSample(probability=3,
                                              collector_set_id=1,
                                              obs_domain_id=2,
                                              obs_point_id=3,)]
        """
        _subtype = nicira_ext.NXAST_SAMPLE

        # probability, collector_set_id, obs_domain_id, obs_point_id
        _fmt_str = '!HIII'

        def __init__(
            self,
            probability,
            collector_set_id=0,
            obs_domain_id=0,
            obs_point_id=0,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionSample, self).__init__()
            self.probability = probability
            self.collector_set_id = collector_set_id
            self.obs_domain_id = obs_domain_id
            self.obs_point_id = obs_point_id

        @classmethod
        def parser(cls, buf):
            (
                probability,
                collector_set_id,
                obs_domain_id,
                obs_point_id,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(
                probability,
                collector_set_id,
                obs_domain_id,
                obs_point_id,
            )

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.probability,
                self.collector_set_id,
                self.obs_domain_id,
                self.obs_point_id,
            )
            return data

    class NXActionSample2(NXAction):
        r"""
        Sample packets action

        This action samples packets and sends one sample for
        every sampled packet.
        'sampling_port' can be equal to ingress port or one of egress ports.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          sample(argument[,argument]...)
        ..

        +----------------------------------------------------+
        | **sample(**\ *argument*\[,\ *argument*\]...\ **)** |
        +----------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        probability      The number of sampled packets
        collector_set_id The unsigned 32-bit integer identifier of
                         the set of sample collectors to send sampled packets to
        obs_domain_id    The Unsigned 32-bit integer Observation Domain ID
        obs_point_id     The unsigned 32-bit integer Observation Point ID
        sampling_port    Sampling port number
        ================ ======================================================

        Example::

            actions += [parser.NXActionSample2(probability=3,
                                               collector_set_id=1,
                                               obs_domain_id=2,
                                               obs_point_id=3,
                                               apn_mac_addr=[10,0,2,0,0,5],
                                               msisdn=b'magmaIsTheBest',
                                               apn_name=b'big_tower123',
                                               pdp_start_epoch=b'90\x00\x00\x00\x00\x00\x00',
                                               sampling_port=8080)]
        """
        _subtype = nicira_ext.NXAST_SAMPLE2

        # probability, collector_set_id, obs_domain_id,
        # obs_point_id, msisdn, apn_mac_addr, apn_name, sampling_port
        _fmt_str = '!HIIIH16s6B24s8s6x'

        def __init__(
            self,
            probability,
            msisdn,
            apn_mac_addr,
            apn_name,
            pdp_start_epoch,
            collector_set_id=0,
            obs_domain_id=0,
            obs_point_id=0,
            sampling_port=0,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionSample2, self).__init__()
            self.probability = probability
            self.collector_set_id = collector_set_id
            self.obs_domain_id = obs_domain_id
            self.obs_point_id = obs_point_id
            self.sampling_port = sampling_port

            self.msisdn = msisdn
            self.apn_mac_addr = apn_mac_addr
            self.apn_name = apn_name
            self.pdp_start_epoch = pdp_start_epoch

        @classmethod
        def parser(cls, buf):
            (
                probability,
                collector_set_id,
                obs_domain_id,
                obs_point_id,
                sampling_port,
                msisdn,
                apn_mac_addr_0,
                apn_mac_addr_1,
                apn_mac_addr_2,
                apn_mac_addr_3,
                apn_mac_addr_4,
                apn_mac_addr_5,
                apn_name,
                pdp_start_epoch,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )

            apn_mac_addr = [apn_mac_addr_0, apn_mac_addr_1, apn_mac_addr_2, apn_mac_addr_3, apn_mac_addr_4, apn_mac_addr_5]
            return cls(
                probability,
                msisdn,
                apn_mac_addr,
                apn_name,
                pdp_start_epoch,
                collector_set_id,
                obs_domain_id,
                obs_point_id,
                sampling_port,
            )

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.probability,
                self.collector_set_id,
                self.obs_domain_id,
                self.obs_point_id,
                self.sampling_port,
                self.msisdn,
                *self.apn_mac_addr,
                self.apn_name,
                self.pdp_start_epoch,
            )

            return data

    class NXActionFinTimeout(NXAction):
        r"""
        Change TCP timeout action

        This action changes the idle timeout or hard timeout or
        both, of this OpenFlow rule when the rule matches a TCP
        packet with the FIN or RST flag.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          fin_timeout(argument[,argument]...)
        ..

        +---------------------------------------------------------+
        | **fin_timeout(**\ *argument*\[,\ *argument*\]...\ **)** |
        +---------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        fin_idle_timeout Causes the flow to expire after the given number
                         of seconds of inactivity
        fin_idle_timeout Causes the flow to expire after the given number
                         of second, regardless of activity
        ================ ======================================================

        Example::

            match = parser.OFPMatch(ip_proto=6, eth_type=0x0800)
            actions += [parser.NXActionFinTimeout(fin_idle_timeout=30,
                                                  fin_hard_timeout=60)]
        """
        _subtype = nicira_ext.NXAST_FIN_TIMEOUT

        # fin_idle_timeout, fin_hard_timeout
        _fmt_str = '!HH2x'

        def __init__(
            self,
            fin_idle_timeout,
            fin_hard_timeout,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionFinTimeout, self).__init__()
            self.fin_idle_timeout = fin_idle_timeout
            self.fin_hard_timeout = fin_hard_timeout

        @classmethod
        def parser(cls, buf):
            (
                fin_idle_timeout,
                fin_hard_timeout,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(
                fin_idle_timeout,
                fin_hard_timeout,
            )

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.fin_idle_timeout,
                self.fin_hard_timeout,
            )
            return data

    class NXActionConjunction(NXAction):
        r"""
        Conjunctive matches action

        This action ties groups of individual OpenFlow flows into
        higher-level conjunctive flows.
        Please refer to the ovs-ofctl command manual for details.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          conjunction(id,k/n)
        ..

        +--------------------------------------------------+
        | **conjunction(**\ *id*\,\ *k*\ **/**\ *n*\ **)** |
        +--------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        clause           Number assigned to the flow's dimension
        n_clauses        Specify the conjunctive flow's match condition
        id\_             Conjunction ID
        ================ ======================================================

        Example::

            actions += [parser.NXActionConjunction(clause=1,
                                                   n_clauses=2,
                                                   id_=10)]
        """
        _subtype = nicira_ext.NXAST_CONJUNCTION

        # clause, n_clauses, id
        _fmt_str = '!BBI'

        def __init__(
            self,
            clause,
            n_clauses,
            id_,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionConjunction, self).__init__()
            self.clause = clause
            self.n_clauses = n_clauses
            self.id = id_

        @classmethod
        def parser(cls, buf):
            (
                clause,
                n_clauses,
                id_,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(clause, n_clauses, id_)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.clause,
                self.n_clauses,
                self.id,
            )
            return data

    class NXActionMultipath(NXAction):
        r"""
        Select multipath link action

        This action selects multipath link based on the specified parameters.
        Please refer to the ovs-ofctl command manual for details.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          multipath(fields, basis, algorithm, n_links, arg, dst[start..end])
        ..

        +-------------------------------------------------------------+
        | **multipath(**\ *fields*\, \ *basis*\, \ *algorithm*\,      |
        | *n_links*\, \ *arg*\, \ *dst*\[\ *start*\..\ *end*\]\ **)** |
        +-------------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        fields           One of NX_HASH_FIELDS_*
        basis            Universal hash parameter
        algorithm        One of NX_MP_ALG_*.
        max_link         Number of output links
        arg              Algorithm-specific argument
        ofs_nbits        Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``
        dst              OXM/NXM header for source field
        ================ ======================================================

        Example::

            actions += [parser.NXActionMultipath(
                            fields=nicira_ext.NX_HASH_FIELDS_SYMMETRIC_L4,
                            basis=1024,
                            algorithm=nicira_ext.NX_MP_ALG_HRW,
                            max_link=5,
                            arg=0,
                            ofs_nbits=nicira_ext.ofs_nbits(4, 31),
                            dst="reg2")]
        """
        _subtype = nicira_ext.NXAST_MULTIPATH

        # fields, basis, algorithm, max_link,
        # arg, ofs_nbits, dst
        _fmt_str = '!HH2xHHI2xH4s'
        _TYPE = {
            'ascii': [
                'dst',
            ],
        }

        def __init__(
            self,
            fields,
            basis,
            algorithm,
            max_link,
            arg,
            ofs_nbits,
            dst,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionMultipath, self).__init__()
            self.fields = fields
            self.basis = basis
            self.algorithm = algorithm
            self.max_link = max_link
            self.arg = arg
            self.ofs_nbits = ofs_nbits
            self.dst = dst

        @classmethod
        def parser(cls, buf):
            (
                fields,
                basis,
                algorithm,
                max_link,
                arg,
                ofs_nbits,
                oxm_data,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            (n, len_) = ofp.oxm_parse_header(oxm_data, 0)
            dst = ofp.oxm_to_user_header(n)
            return cls(
                fields,
                basis,
                algorithm,
                max_link,
                arg,
                ofs_nbits,
                dst,
            )

        def serialize_body(self):
            data = bytearray()
            dst = bytearray()
            oxm = ofp.oxm_from_user_header(self.dst)
            ofp.oxm_serialize_header(oxm, dst, 0),
            msg_pack_into(
                self._fmt_str, data, 0,
                self.fields,
                self.basis,
                self.algorithm,
                self.max_link,
                self.arg,
                self.ofs_nbits,
                six.binary_type(dst),
            )

            return data

    class _NXActionBundleBase(NXAction):
        # algorithm, fields, basis, slave_type, n_slaves
        # ofs_nbits
        _fmt_str = '!HHHIHH'

        def __init__(
            self, algorithm, fields, basis, slave_type, n_slaves,
            ofs_nbits, dst, slaves,
        ):
            super(_NXActionBundleBase, self).__init__()
            self.len = utils.round_up(
                nicira_ext.NX_ACTION_BUNDLE_0_SIZE + len(slaves) * 2, 8,
            )

            self.algorithm = algorithm
            self.fields = fields
            self.basis = basis
            self.slave_type = slave_type
            self.n_slaves = n_slaves
            self.ofs_nbits = ofs_nbits
            self.dst = dst

            assert isinstance(slaves, (list, tuple))
            for s in slaves:
                assert isinstance(s, six.integer_types)

            self.slaves = slaves

        @classmethod
        def parser(cls, buf):
            # Add dst ('I') to _fmt_str
            (
                algorithm, fields, basis,
                slave_type, n_slaves, ofs_nbits, dst,
            ) = struct.unpack_from(
                cls._fmt_str + 'I', buf, 0,
            )

            offset = (
                nicira_ext.NX_ACTION_BUNDLE_0_SIZE -
                nicira_ext.NX_ACTION_HEADER_0_SIZE - 8
            )

            if dst != 0:
                (n, len_) = ofp.oxm_parse_header(buf, offset)
                dst = ofp.oxm_to_user_header(n)

            slave_offset = (
                nicira_ext.NX_ACTION_BUNDLE_0_SIZE -
                nicira_ext.NX_ACTION_HEADER_0_SIZE
            )

            slaves = []
            for i in range(0, n_slaves):
                s = struct.unpack_from('!H', buf, slave_offset)
                slaves.append(s[0])
                slave_offset += 2

            return cls(
                algorithm, fields, basis, slave_type,
                n_slaves, ofs_nbits, dst, slaves,
            )

        def serialize_body(self):
            data = bytearray()
            slave_offset = (
                nicira_ext.NX_ACTION_BUNDLE_0_SIZE -
                nicira_ext.NX_ACTION_HEADER_0_SIZE
            )
            self.n_slaves = len(self.slaves)
            for s in self.slaves:
                msg_pack_into('!H', data, slave_offset, s)
                slave_offset += 2
            pad_len = (
                utils.round_up(self.n_slaves, 4) -
                self.n_slaves
            )

            if pad_len != 0:
                msg_pack_into('%dx' % pad_len * 2, data, slave_offset)

            msg_pack_into(
                self._fmt_str, data, 0,
                self.algorithm, self.fields, self.basis,
                self.slave_type, self.n_slaves,
                self.ofs_nbits,
            )
            offset = (
                nicira_ext.NX_ACTION_BUNDLE_0_SIZE -
                nicira_ext.NX_ACTION_HEADER_0_SIZE - 8
            )

            if self.dst == 0:
                msg_pack_into('I', data, offset, self.dst)
            else:
                oxm_data = ofp.oxm_from_user_header(self.dst)
                ofp.oxm_serialize_header(oxm_data, data, offset)
            return data

    class NXActionBundle(_NXActionBundleBase):
        r"""
        Select bundle link action

        This action selects bundle link based on the specified parameters.
        Please refer to the ovs-ofctl command manual for details.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          bundle(fields, basis, algorithm, slave_type, slaves:[ s1, s2,...])
        ..

        +-----------------------------------------------------------+
        | **bundle(**\ *fields*\, \ *basis*\, \ *algorithm*\,       |
        | *slave_type*\, \ *slaves*\:[ \ *s1*\, \ *s2*\,...]\ **)** |
        +-----------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        algorithm        One of NX_MP_ALG_*.
        fields           One of NX_HASH_FIELDS_*
        basis            Universal hash parameter
        slave_type       Type of slaves(must be NXM_OF_IN_PORT)
        n_slaves         Number of slaves
        ofs_nbits        Start and End for the OXM/NXM field. (must be zero)
        dst              OXM/NXM header for source field(must be zero)
        slaves           List of slaves
        ================ ======================================================


        Example::

            actions += [parser.NXActionBundle(
                            algorithm=nicira_ext.NX_MP_ALG_HRW,
                            fields=nicira_ext.NX_HASH_FIELDS_ETH_SRC,
                            basis=0,
                            slave_type=nicira_ext.NXM_OF_IN_PORT,
                            n_slaves=2,
                            ofs_nbits=0,
                            dst=0,
                            slaves=[2, 3])]
        """
        _subtype = nicira_ext.NXAST_BUNDLE

        def __init__(
            self, algorithm, fields, basis, slave_type, n_slaves,
            ofs_nbits, dst, slaves,
        ):
            # NXAST_BUNDLE actions should have 'sofs_nbits' and 'dst' zeroed.
            super(NXActionBundle, self).__init__(
                algorithm, fields, basis, slave_type, n_slaves,
                ofs_nbits=0, dst=0, slaves=slaves,
            )

    class NXActionBundleLoad(_NXActionBundleBase):
        r"""
        Select bundle link action

        This action has the same behavior as the bundle action,
        with one exception.
        Please refer to the ovs-ofctl command manual for details.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          bundle_load(fields, basis, algorithm, slave_type,
                      dst[start..end], slaves:[ s1, s2,...])
        ..

        +-----------------------------------------------------------+
        | **bundle_load(**\ *fields*\, \ *basis*\, \ *algorithm*\,  |
        | *slave_type*\, \ *dst*\[\ *start*\... \*emd*\],           |
        | \ *slaves*\:[ \ *s1*\, \ *s2*\,...]\ **)** |              |
        +-----------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        algorithm        One of NX_MP_ALG_*.
        fields           One of NX_HASH_FIELDS_*
        basis            Universal hash parameter
        slave_type       Type of slaves(must be NXM_OF_IN_PORT)
        n_slaves         Number of slaves
        ofs_nbits        Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``
        dst              OXM/NXM header for source field
        slaves           List of slaves
        ================ ======================================================


        Example::

            actions += [parser.NXActionBundleLoad(
                            algorithm=nicira_ext.NX_MP_ALG_HRW,
                            fields=nicira_ext.NX_HASH_FIELDS_ETH_SRC,
                            basis=0,
                            slave_type=nicira_ext.NXM_OF_IN_PORT,
                            n_slaves=2,
                            ofs_nbits=nicira_ext.ofs_nbits(4, 31),
                            dst="reg0",
                            slaves=[2, 3])]
        """
        _subtype = nicira_ext.NXAST_BUNDLE_LOAD
        _TYPE = {
            'ascii': [
                'dst',
            ],
        }

        def __init__(
            self, algorithm, fields, basis, slave_type, n_slaves,
            ofs_nbits, dst, slaves,
        ):
            super(NXActionBundleLoad, self).__init__(
                algorithm, fields, basis, slave_type, n_slaves,
                ofs_nbits, dst, slaves,
            )

    class NXActionCT(NXAction):
        r"""
        Pass traffic to the connection tracker action

        This action sends the packet through the connection tracker.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          ct(argument[,argument]...)
        ..

        +------------------------------------------------+
        | **ct(**\ *argument*\[,\ *argument*\]...\ **)** |
        +------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        flags            Zero or more(Unspecified flag bits must be zero.)
        zone_src         OXM/NXM header for source field
        zone_ofs_nbits   Start and End for the OXM/NXM field.
                         Setting method refer to the ``nicira_ext.ofs_nbits``.
                         If you need set the Immediate value for zone,
                         zone_src must be set to None or empty character string.
        recirc_table     Recirculate to a specific table
        alg              Well-known port number for the protocol
        actions          Zero or more actions may immediately follow this
                         action
        ================ ======================================================

        .. NOTE::

            If you set number to zone_src,
            Traceback occurs when you run the to_jsondict.

        Example::

            match = parser.OFPMatch(eth_type=0x0800, ct_state=(0,32))
            actions += [parser.NXActionCT(
                            flags = 1,
                            zone_src = "reg0",
                            zone_ofs_nbits = nicira_ext.ofs_nbits(4, 31),
                            recirc_table = 4,
                            alg = 0,
                            actions = [])]
        """
        _subtype = nicira_ext.NXAST_CT

        # flags, zone_src, zone_ofs_nbits, recirc_table,
        # pad, alg
        _fmt_str = '!H4sHB3xH'
        _TYPE = {
            'ascii': [
                'zone_src',
            ],
        }

        # Followed by actions

        def __init__(
            self,
            flags,
            zone_src,
            zone_ofs_nbits,
            recirc_table,
            alg,
            actions,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionCT, self).__init__()
            self.flags = flags
            self.zone_src = zone_src
            self.zone_ofs_nbits = zone_ofs_nbits
            self.recirc_table = recirc_table
            self.alg = alg
            self.actions = actions

        @classmethod
        def parser(cls, buf):
            (
                flags,
                oxm_data,
                zone_ofs_nbits,
                recirc_table,
                alg,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            rest = buf[struct.calcsize(cls._fmt_str):]

            # OXM/NXM field
            if oxm_data == b'\x00' * 4:
                zone_src = ""
            else:
                (n, len_) = ofp.oxm_parse_header(oxm_data, 0)
                zone_src = ofp.oxm_to_user_header(n)

            # actions
            actions = []
            while len(rest) > 0:
                action = ofpp.OFPAction.parser(rest, 0)
                actions.append(action)
                rest = rest[action.len:]

            return cls(
                flags, zone_src, zone_ofs_nbits, recirc_table,
                alg, actions,
            )

        def serialize_body(self):
            data = bytearray()
            # If zone_src is zero, zone_ofs_nbits is zone_imm
            if not self.zone_src:
                zone_src = b'\x00' * 4
            elif isinstance(self.zone_src, six.integer_types):
                zone_src = struct.pack("!I", self.zone_src)
            else:
                zone_src = bytearray()
                oxm = ofp.oxm_from_user_header(self.zone_src)
                ofp.oxm_serialize_header(oxm, zone_src, 0)

            msg_pack_into(
                self._fmt_str, data, 0,
                self.flags,
                six.binary_type(zone_src),
                self.zone_ofs_nbits,
                self.recirc_table,
                self.alg,
            )
            for a in self.actions:
                a.serialize(data, len(data))
            return data

    class NXActionCTClear(NXAction):
        """
        Clear connection tracking state action

        This action clears connection tracking state from packets.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          ct_clear
        ..

        +--------------+
        | **ct_clear** |
        +--------------+

        Example::

            actions += [parser.NXActionCTClear()]
        """
        _subtype = nicira_ext.NXAST_CT_CLEAR

        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionCTClear, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    class NXActionNAT(NXAction):
        r"""
        Network address translation action

        This action sends the packet through the connection tracker.

        And equivalent to the followings action of ovs-ofctl command.

        .. NOTE::
            The following command image does not exist in ovs-ofctl command
            manual and has been created from the command response.

        ..
          nat(src=ip_min-ip_max : proto_min-proto-max)
        ..

        +--------------------------------------------------+
        | **nat(src**\=\ *ip_min*\ **-**\ *ip_max*\  **:** |
        | *proto_min*\ **-**\ *proto-max*\ **)**           |
        +--------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        flags            Zero or more(Unspecified flag bits must be zero.)
        range_ipv4_min   Range ipv4 address minimun
        range_ipv4_max   Range ipv4 address maximun
        range_ipv6_min   Range ipv6 address minimun
        range_ipv6_max   Range ipv6 address maximun
        range_proto_min  Range protocol minimum
        range_proto_max  Range protocol maximun
        ================ ======================================================

        .. CAUTION::
            ``NXActionNAT`` must be defined in the actions in the
            ``NXActionCT``.

        Example::

            match = parser.OFPMatch(eth_type=0x0800)
            actions += [
                parser.NXActionCT(
                    flags = 1,
                    zone_src = "reg0",
                    zone_ofs_nbits = nicira_ext.ofs_nbits(4, 31),
                    recirc_table = 255,
                    alg = 0,
                    actions = [
                        parser.NXActionNAT(
                            flags = 1,
                            range_ipv4_min = "10.1.12.0",
                            range_ipv4_max = "10.1.13.255",
                            range_ipv6_min = "",
                            range_ipv6_max = "",
                            range_proto_min = 1,
                            range_proto_max = 1023
                        )
                    ]
                )
            ]
        """
        _subtype = nicira_ext.NXAST_NAT

        # pad, flags, range_present
        _fmt_str = '!2xHH'
        # Followed by optional parameters

        _TYPE = {
            'ascii': [
                'range_ipv4_max',
                'range_ipv4_min',
                'range_ipv6_max',
                'range_ipv6_min',
            ],
        }

        def __init__(
            self,
            flags,
            range_ipv4_min='',
            range_ipv4_max='',
            range_ipv6_min='',
            range_ipv6_max='',
            range_proto_min=None,
            range_proto_max=None,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionNAT, self).__init__()
            self.flags = flags
            self.range_ipv4_min = range_ipv4_min
            self.range_ipv4_max = range_ipv4_max
            self.range_ipv6_min = range_ipv6_min
            self.range_ipv6_max = range_ipv6_max
            self.range_proto_min = range_proto_min
            self.range_proto_max = range_proto_max

        @classmethod
        def parser(cls, buf):
            (
                flags,
                range_present,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            rest = buf[struct.calcsize(cls._fmt_str):]
            # optional parameters
            kwargs = dict()
            if range_present & nicira_ext.NX_NAT_RANGE_IPV4_MIN:
                kwargs['range_ipv4_min'] = type_desc.IPv4Addr.to_user(rest[:4])
                rest = rest[4:]
            if range_present & nicira_ext.NX_NAT_RANGE_IPV4_MAX:
                kwargs['range_ipv4_max'] = type_desc.IPv4Addr.to_user(rest[:4])
                rest = rest[4:]
            if range_present & nicira_ext.NX_NAT_RANGE_IPV6_MIN:
                kwargs['range_ipv6_min'] = (
                    type_desc.IPv6Addr.to_user(rest[:16])
                )
                rest = rest[16:]
            if range_present & nicira_ext.NX_NAT_RANGE_IPV6_MAX:
                kwargs['range_ipv6_max'] = (
                    type_desc.IPv6Addr.to_user(rest[:16])
                )
                rest = rest[16:]
            if range_present & nicira_ext.NX_NAT_RANGE_PROTO_MIN:
                kwargs['range_proto_min'] = type_desc.Int2.to_user(rest[:2])
                rest = rest[2:]
            if range_present & nicira_ext.NX_NAT_RANGE_PROTO_MAX:
                kwargs['range_proto_max'] = type_desc.Int2.to_user(rest[:2])

            return cls(flags, **kwargs)

        def serialize_body(self):
            # Pack optional parameters first, as range_present needs
            # to be calculated.
            optional_data = b''
            range_present = 0
            if self.range_ipv4_min != '':
                range_present |= nicira_ext.NX_NAT_RANGE_IPV4_MIN
                optional_data += type_desc.IPv4Addr.from_user(
                    self.range_ipv4_min,
                )
            if self.range_ipv4_max != '':
                range_present |= nicira_ext.NX_NAT_RANGE_IPV4_MAX
                optional_data += type_desc.IPv4Addr.from_user(
                    self.range_ipv4_max,
                )
            if self.range_ipv6_min != '':
                range_present |= nicira_ext.NX_NAT_RANGE_IPV6_MIN
                optional_data += type_desc.IPv6Addr.from_user(
                    self.range_ipv6_min,
                )
            if self.range_ipv6_max != '':
                range_present |= nicira_ext.NX_NAT_RANGE_IPV6_MAX
                optional_data += type_desc.IPv6Addr.from_user(
                    self.range_ipv6_max,
                )
            if self.range_proto_min is not None:
                range_present |= nicira_ext.NX_NAT_RANGE_PROTO_MIN
                optional_data += type_desc.Int2.from_user(
                    self.range_proto_min,
                )
            if self.range_proto_max is not None:
                range_present |= nicira_ext.NX_NAT_RANGE_PROTO_MAX
                optional_data += type_desc.Int2.from_user(
                    self.range_proto_max,
                )

            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.flags,
                range_present,
            )
            msg_pack_into(
                '!%ds' % len(optional_data), data, len(data),
                optional_data,
            )

            return data

    class NXActionOutputTrunc(NXAction):
        r"""
        Truncate output action

        This action truncate a packet into the specified size and outputs it.

        And equivalent to the followings action of ovs-ofctl command.

        ..
          output(port=port,max_len=max_len)
        ..

        +--------------------------------------------------------------+
        | **output(port**\=\ *port*\,\ **max_len**\=\ *max_len*\ **)** |
        +--------------------------------------------------------------+

        ================ ======================================================
        Attribute        Description
        ================ ======================================================
        port             Output port
        max_len          Max bytes to send
        ================ ======================================================

        Example::

            actions += [parser.NXActionOutputTrunc(port=8080,
                                                   max_len=1024)]
        """
        _subtype = nicira_ext.NXAST_OUTPUT_TRUNC

        # port, max_len
        _fmt_str = '!HI'

        def __init__(
            self,
            port,
            max_len,
            type_=None, len_=None, experimenter=None, subtype=None,
        ):
            super(NXActionOutputTrunc, self).__init__()
            self.port = port
            self.max_len = max_len

        @classmethod
        def parser(cls, buf):
            (
                port,
                max_len,
            ) = struct.unpack_from(
                cls._fmt_str, buf, 0,
            )
            return cls(port, max_len)

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(
                self._fmt_str, data, 0,
                self.port,
                self.max_len,
            )
            return data

    class NXActionEncapEther(NXAction):
        """
        Encap Ether

        This action encaps package with ethernet

        And equivalent to the followings action of ovs-ofctl command.

        ::

            encap(ethernet)

        Example::

            actions += [parser.NXActionEncapEther()]
        """
        _subtype = nicira_ext.NXAST_RAW_ENCAP

        _fmt_str = '!HI'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionEncapEther, self).__init__()
            self.hdr_size = 0
            self.new_pkt_type = 0x00000000

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0, self.hdr_size, self.new_pkt_type)
            return data

    class NXActionEncapNsh(NXAction):
        """
        Encap nsh

        This action encaps package with nsh

        And equivalent to the followings action of ovs-ofctl command.

        ::

            encap(nsh(md_type=1))

        Example::

            actions += [parser.NXActionEncapNsh()]
        """
        _subtype = nicira_ext.NXAST_RAW_ENCAP

        _fmt_str = '!HI'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionEncapNsh, self).__init__()
            self.hdr_size = hdr_size
            self.new_pkt_type = 0x0001894F

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0, self.hdr_size, self.new_pkt_type)
            return data

    class NXActionDecNshTtl(NXAction):
        """
        Decrement NSH TTL action

        This action decrements the TTL in the Network Service Header(NSH).

        This action was added in OVS v2.9.

        And equivalent to the followings action of ovs-ofctl command.

        ::

            dec_nsh_ttl

        Example::

            actions += [parser.NXActionDecNshTtl()]
        """
        _subtype = nicira_ext.NXAST_DEC_NSH_TTL

        _fmt_str = '!6x'

        def __init__(
            self,
            type_=None, len_=None, vendor=None, subtype=None,
        ):
            super(NXActionDecNshTtl, self).__init__()

        @classmethod
        def parser(cls, buf):
            return cls()

        def serialize_body(self):
            data = bytearray()
            msg_pack_into(self._fmt_str, data, 0)
            return data

    def add_attr(k, v):
        v.__module__ = ofpp.__name__  # Necessary for stringify stuff
        setattr(ofpp, k, v)

    add_attr('NXAction', NXAction)
    add_attr('NXActionUnknown', NXActionUnknown)

    classes = [
        'NXActionSetQueue',
        'NXActionPopQueue',
        'NXActionRegLoad',
        'NXActionRegLoad2',
        'NXActionNote',
        'NXActionSetTunnel',
        'NXActionSetTunnel64',
        'NXActionRegMove',
        'NXActionResubmit',
        'NXActionResubmitTable',
        'NXActionOutputReg',
        'NXActionOutputReg2',
        'NXActionLearn',
        'NXActionExit',
        'NXActionDecTtl',
        'NXActionController',
        'NXActionController2',
        'NXActionDecTtlCntIds',
        'NXActionPushMpls',
        'NXActionPopMpls',
        'NXActionSetMplsTtl',
        'NXActionDecMplsTtl',
        'NXActionSetMplsLabel',
        'NXActionSetMplsTc',
        'NXActionStackPush',
        'NXActionStackPop',
        'NXActionSample',
        'NXActionSample2',
        'NXActionFinTimeout',
        'NXActionConjunction',
        'NXActionMultipath',
        'NXActionBundle',
        'NXActionBundleLoad',
        'NXActionCT',
        'NXActionCTClear',
        'NXActionNAT',
        'NXActionOutputTrunc',
        '_NXFlowSpec',  # exported for testing
        'NXFlowSpecMatch',
        'NXFlowSpecLoad',
        'NXFlowSpecOutput',
        'NXActionEncapNsh',
        'NXActionEncapEther',
        'NXActionDecNshTtl',
    ]
    vars = locals()
    for name in classes:
        cls = vars[name]
        add_attr(name, cls)
        if issubclass(cls, NXAction):
            NXAction.register(cls)
        if issubclass(cls, _NXFlowSpec):
            _NXFlowSpec.register(cls)

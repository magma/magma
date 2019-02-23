"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from lte.protos.subscriberdb_pb2 import SubscriberID


class SIDUtils:
    """
    Utility functions for SubscriberIDs (a.k.a SIDs)
    """

    @staticmethod
    def to_str(sid_pb):
        """
        Standard way of converting a SubscriberID to a string representation

        Args:
            sid_pb - SubscriberID protobuf message
        Returns:
            string representation of the SubscriberID
        Raises:
            ValueError
        """
        if sid_pb.type == SubscriberID.IMSI:
            return 'IMSI' + sid_pb.id
        raise ValueError('Invalid sid! type:%s id:%s' %
                         (sid_pb.type, sid_pb.id))

    @staticmethod
    def to_pb(sid_str):
        """
        Converts a string to SubscriberID protobuf message

        Args:
            sid_str: string representation of SubscriberID
        Returns:
            SubscriberID protobuf message
        Raises:
            ValueError
        """
        if sid_str.startswith('IMSI'):
            imsi = sid_str[4:]
            if not imsi.isdigit():
                raise ValueError('Invalid imsi: %s' % imsi)
            return SubscriberID(id=imsi, type=SubscriberID.IMSI)
        raise ValueError('Invalid sid: %s' % sid_str)

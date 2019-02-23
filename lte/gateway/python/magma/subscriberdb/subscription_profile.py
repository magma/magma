"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from lte.protos.mconfig.mconfigs_pb2 import SubscriberDB


def get_default_sub_profile(service):
    """
    Returns the default subscription profile to be used when the
    subcribers don't have a profile associated with them.
    """
    if 'default' in service.mconfig.sub_profiles:
        return service.mconfig.sub_profiles['default']
    # No default profile configured for the network. Use the default defined
    # in the code.
    return SubscriberDB.SubscriptionProfile(
        max_ul_bit_rate=service.config['default_max_ul_bit_rate'],
        max_dl_bit_rate=service.config['default_max_dl_bit_rate']
    )

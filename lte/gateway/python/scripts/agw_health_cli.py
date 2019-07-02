#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sys

import fire
from magma.configuration.mconfig_managers import load_service_mconfig_as_json


class AGWHealthSummary:
    def __init__(self, relay_enabled):
        self.relay_enabled = relay_enabled

    def __str__(self):
        return """
{}""".format(
            'Using Feg' if self.relay_enabled else 'Using subscriberdb',
        )


def gateway_health_status():
    config = load_service_mconfig_as_json('mme')

    health_summary = AGWHealthSummary(relay_enabled=config['relayEnabled'])
    return str(health_summary)


if __name__ == '__main__':
    print('Access Gateway health summary')
    if len(sys.argv) == 1:
        fire.Fire(gateway_health_status)
    else:
        fire.Fire({
            'status': gateway_health_status,
        })

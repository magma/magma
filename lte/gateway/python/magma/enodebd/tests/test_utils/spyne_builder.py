"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from unittest import mock
from spyne.server.wsgi import WsgiMethodContext


def get_spyne_context_with_ip(
    req_ip: str = "192.168.60.145",
) -> WsgiMethodContext:
    with mock.patch('spyne.server.wsgi.WsgiApplication') as MockTransport:
        MockTransport.req_env = {"REMOTE_ADDR": req_ip}
        with mock.patch('spyne.server.wsgi.WsgiMethodContext') as MockContext:
            MockContext.transport = MockTransport
            return MockContext

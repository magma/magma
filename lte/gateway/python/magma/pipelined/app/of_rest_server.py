"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
from ryu import cfg
from ryu.app import wsgi
from ryu.lib import hub


def configure(pipelined_config):
    CONF = cfg.CONF
    CONF.wsapi_port = pipelined_config['of_server_port']


def start(app_manager):
    webapp = wsgi.start_service(app_manager)
    return hub.spawn(webapp)

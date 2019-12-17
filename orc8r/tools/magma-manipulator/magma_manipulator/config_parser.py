"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import os
from jsonschema import validate
import yaml

schema = """
    type: object
    properties:
      k8s:
        type: object
        properties:
          kubeconfig_path:
            type: string
          namespace:
            type: string
      orc8r_api_url:
        type: string
      magma_certs_path:
        type: array
      gateways:
        type: object
        properties:
          configs_dir:
            type: string
          username:
            type: string
          rsa_private_key_path:
            type: string
"""


def parse_config(cfg_rel_path):
    dirname = os.path.dirname(__file__)
    cfg_path = os.path.join(dirname, cfg_rel_path)
    with open(cfg_path, 'r') as ymlfile:
        yml_cfg = yaml.load(ymlfile, Loader=yaml.FullLoader)
    validate(yml_cfg, yaml.load(schema))
    return yml_cfg


class Config(object):
    def __init__(self, yml_cfg):
        for k, v in yml_cfg.items():
            if isinstance(v, dict):
                self.__dict__[k] = Config(v)
            else:
                self.__dict__[k] = v


cfg = Config(parse_config('config.yml'))

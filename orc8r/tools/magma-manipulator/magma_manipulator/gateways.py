"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma_manipulator.config_parser import cfg as CONF
from magma_manipulator import k8s_tools
from magma_manipulator import magma_api
from magma_manipulator import utils


class GatewaysManager(object):
    _gateways = {}

    def __init__(self):
        self._get_magma_gateways()

    def _get_magma_gateways(self):
        networks = magma_api.get_networks(
            CONF.orc8r_api_url, CONF.magma_certs_path)
        for net in networks:
            net_type = magma_api.get_network_type(
                CONF.orc8r_api_url, net, CONF.magma_certs_path)
            gws = magma_api.get_gateways(
                CONF.orc8r_api_url, net, net_type, CONF.magma_certs_path)
            for gw_id, gw_desc in gws.items():
                gw_config = magma_api.get_gateway_config(
                    CONF.orc8r_api_url, net, net_type,
                    gw_id, CONF.magma_certs_path)
                config_path = utils.save_gateway_config(
                        gw_id, CONF.gateways.configs_dir, gw_config)

                self._gateways[gw_desc['name']] = Gateway(
                        gw_id, gw_desc['name'], net, net_type, config_path)

    def get_gateway(self, gw_pod_name):
        gw_name = gw_pod_name.split('-')[0]
        return self._gateways[gw_name]

    def get_gateways(self):
        return self._gateways

    def delete_gateway(self, gw_name):
        del self._gateways[gw_name]

    def get_gateway_names(self):
        return list(self._gateways.keys())


class Gateway(object):
    def __init__(self, gw_id, gw_name, gw_network,
                 gw_network_type, gw_config_path):
        self.id = gw_id
        self.name = gw_name

        self.network = gw_network
        self.network_type = gw_network_type

        self.config_path = gw_config_path

        self._ip = None
        self._uuid = None
        self._key = None
        self._pod_name = None

    def get_ip(self, pod_name):
        if not self._ip:
            self._ip = k8s_tools.get_gw_ip(CONF.k8s.kubeconfig_path,
                                           CONF.k8s.namespace,
                                           pod_name)
            self._pod_name = pod_name
        return self._ip

    def get_uuid_and_key(self):
        if not self._uuid and not self._key:
            self._uuid, self._key = utils.get_gw_uuid_and_key(
                self._ip, CONF.gateways.username,
                CONF.gateways.rsa_private_key_path)
        return (self._uuid, self._key)

    def get_config(self):
        return utils.load_gateway_config(self.name, self.config_path)

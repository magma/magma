"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
import logging
import requests
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from urllib.parse import urljoin

from magma_manipulator import exceptions

requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
LOG = logging.getLogger(__name__)


def is_network_exist(orc8r_api_url, gw_net, certs):
    LOG.info('Check if network {gw_net} exists'.format(gw_net=gw_net))
    magma_net_url = urljoin(orc8r_api_url,
                            'magma/v1/networks/{gw_net}'.format(gw_net=gw_net))
    LOG.debug('Make get request to {url}'.format(url=magma_net_url))
    resp = requests.get(magma_net_url, verify=False, cert=certs)
    str_result = resp.content.decode('ascii')
    json_result = json.loads(str_result)
    LOG.debug('Received result {result}'.format(result=json_result))
    if 'id' in json_result:
        return json_result['id'] == gw_net
    return False


def create_network(orc8r_api_url, gw_net, certs):
    LOG.info('Start to create network {gw_net}'.format(gw_net=gw_net))
    magma_net_url = urljoin(
        orc8r_api_url,
        'magma/v1/networks')
    data = {
        'description': 'This network created from automation tool',
        'dns': {
          'enable_caching': False,
          'local_ttl': 0,
        },
        'id': gw_net,
        'name': gw_net
      }

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}
    resp = requests.post(magma_net_url,
                         data=json.dumps(data),
                         headers=headers,
                         verify=False,
                         cert=certs)
    msg = 'Receive response {text} with status code '\
          '{status_code} afte network {gw_net} creation.'.format(
                  text=resp.text,
                  status_code=resp.status_code,
                  gw_net=gw_net)
    LOG.info(msg)
    if resp.status_code not in [200, 201, 204]:
        raise exceptions.MagmaRequestException(msg)


def _get_register_gateway_url(gw_net_type, gw_net):
    if gw_net_type == 'carrier_wifi_network':
        return 'magma/v1/cwf/{gw_net}/gateways'.format(gw_net=gw_net)
    elif gw_net_type == 'feg':
        return 'magma/v1/feg/{gw_net}/gateways'.format(gw_net=gw_net)


def _get_register_gateway_data(gw_net_type, gw_id, gw_uuid, gw_key, gw_name):
    if gw_net_type == 'carrier_wifi_network':
        data = {
          "carrier_wifi": {
            "allowed_gre_peers": [
              {
                "ip": "192.168.127.1",
                "key": 1
              }
            ]
          },
          "description": "Gateway was created from magma-manipulator",
          "device": {
            "hardware_id": gw_uuid,
            "key": {
              "key": gw_key,
              "key_type": "SOFTWARE_ECDSA_SHA256"
            }
          },
          "id": gw_id,
          "magmad": {
            "autoupgrade_enabled": True,
            "autoupgrade_poll_interval": 300,
            "checkin_interval": 60,
            "checkin_timeout": 10,
            "dynamic_services": [],
            "feature_flags": {
              "newfeature1": True,
              "newfeature2": False
            }
          },
          "name": gw_name,
          "tier": "default"
        }
        return data

    elif gw_net_type == 'feg':
        data = {
            "description": "Sample Gateway description",
            "device": {
              "hardware_id": gw_uuid,
              "key": {
                "key": gw_key,
                "key_type": "SOFTWARE_ECDSA_SHA256"
              }
            },
            "federation": {
              "aaa_server": {
                "accounting_enabled": True,
                "create_session_on_auth": True,
                "idle_session_timeout_ms": 21600000
              },
              "eap_aka": {
                "plmn_ids": [
                  "123456"
                ],
                "timeout": {
                  "challenge_ms": 20000,
                  "error_notification_ms": 10000,
                  "session_authenticated_ms": 5000,
                  "session_ms": 43200000
                }
              },
              "gx": {
                "server": {
                  "address": "foo.bar.com:5555",
                  "dest_host": "magma-fedgw.magma.com",
                  "dest_realm": "magma.com",
                  "disable_dest_host": False,
                  "host": "string",
                  "local_address": ":56789",
                  "product_name": "string",
                  "protocol": "tcp",
                  "realm": "string",
                  "retransmits": 0,
                  "retry_count": 0,
                  "watchdog_interval": 0
                }
              },
              "gy": {
                "init_method": 2,
                "server": {
                  "address": "foo.bar.com:5555",
                  "dest_host": "magma-fedgw.magma.com",
                  "dest_realm": "magma.com",
                  "disable_dest_host": False,
                  "host": "string",
                  "local_address": ":56789",
                  "product_name": "string",
                  "protocol": "tcp",
                  "realm": "string",
                  "retransmits": 0,
                  "retry_count": 0,
                  "watchdog_interval": 0
                }
              },
              "health": {
                "cloud_disable_period_secs": 10,
                "cpu_utilization_threshold": 0.9,
                "health_services": [
                  "SESSION_PROXY",
                  "SWX_PROXY"
                ],
                "local_disable_period_secs": 1,
                "memory_available_threshold": 0.75,
                "minimum_request_threshold": 1,
                "request_failure_threshold": 0.5,
                "update_failure_threshold": 3,
                "update_interval_secs": 10
              },
              "hss": {
                "default_sub_profile": {
                  "max_dl_bit_rate": 200000000,
                  "max_ul_bit_rate": 100000000
                },
                "lte_auth_amf": "gAA=",
                "lte_auth_op": "EREREREREREREREREREREQ==",
                "server": {
                  "address": "foo.bar.com:5555",
                  "dest_host": "magma-fedgw.magma.com",
                  "dest_realm": "magma.com",
                  "local_address": ":56789",
                  "protocol": "tcp"
                },
                "stream_subscribers": False,
                "sub_profiles": {
                  "additionalProp1": {
                    "max_dl_bit_rate": 200000000,
                    "max_ul_bit_rate": 100000000
                  },
                  "additionalProp2": {
                    "max_dl_bit_rate": 200000000,
                    "max_ul_bit_rate": 100000000
                  },
                  "additionalProp3": {
                    "max_dl_bit_rate": 200000000,
                    "max_ul_bit_rate": 100000000
                  }
                }
              },
              "s6a": {
                "server": {
                  "address": "foo.bar.com:5555",
                  "dest_host": "magma-fedgw.magma.com",
                  "dest_realm": "magma.com",
                  "disable_dest_host": False,
                  "host": "string",
                  "local_address": ":56789",
                  "product_name": "string",
                  "protocol": "tcp",
                  "realm": "string",
                  "retransmits": 0,
                  "retry_count": 0,
                  "watchdog_interval": 0
                }
              },
              "served_network_ids": [
                "string"
              ],
              "swx": {
                "cache_TTL_seconds": 10800,
                "derive_unregister_realm": False,
                "register_on_auth": False,
                "server": {
                  "address": "foo.bar.com:5555",
                  "dest_host": "magma-fedgw.magma.com",
                  "dest_realm": "magma.com",
                  "disable_dest_host": False,
                  "host": "string",
                  "local_address": ":56789",
                  "product_name": "string",
                  "protocol": "tcp",
                  "realm": "string",
                  "retransmits": 0,
                  "retry_count": 0,
                  "watchdog_interval": 0
                },
                "verify_authorization": False
              }
            },
            "id": gw_id,
            "magmad": {
              "autoupgrade_enabled": True,
              "autoupgrade_poll_interval": 300,
              "checkin_interval": 60,
              "checkin_timeout": 10,
              "dynamic_services": [],
              "feature_flags": {
                "newfeature1": True,
                "newfeature2": False
              }
            },
            "name": gw_name,
            "tier": "default"
        }
        return data


def register_gateway(orc8r_api_url, gw_net, gw_net_type,
                     gw_id, gw_uuid, gw_key, gw_name, certs):
    msg = 'Register gateway {gw_name} with {gw_id} in network {gw_net_type} '\
          '{gw_net} with hardware_id {gw_uuid} and key {gw_key}'.format(
                  gw_name=gw_name,
                  gw_id=gw_id,
                  gw_net_type=gw_net_type,
                  gw_net=gw_net,
                  gw_uuid=gw_uuid,
                  gw_key=gw_key)
    LOG.info(msg)

    register_gw_url = _get_register_gateway_url(gw_net_type, gw_net)
    magma_gw_url = urljoin(orc8r_api_url, register_gw_url)

    data = _get_register_gateway_data(gw_net_type,
                                      gw_id, gw_uuid,
                                      gw_key, gw_name)

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}
    resp = requests.post(magma_gw_url,
                         data=json.dumps(data),
                         headers=headers,
                         verify=False,
                         cert=certs)
    msg = 'Receive response {text} with status code {status_code} '\
          'after {gw_name} creation'\
          .format(text=resp.text,
                  status_code=resp.status_code,
                  gw_name=gw_name)
    LOG.info(msg)
    if resp.status_code not in [200, 201, 204]:
        raise exceptions.MagmaRequestException(msg)


def is_gateway_in_network(orc8r_api_url, gw_net, gw_id, certs):
    LOG.info('Check if gateway {gw_id} exists in network {gw_net}'.format(
        gw_id=gw_id, gw_net=gw_net))
    magma_gw_url = urljoin(
        orc8r_api_url,
        'magma/v1/networks/{gw_net}/gateways'.format(
            gw_net=gw_net))
    headers = {'content-type': 'application/json',
               'accept': 'application/json'}
    resp = requests.get(magma_gw_url,
                        headers=headers,
                        verify=False,
                        cert=certs)
    data = json.loads(resp.content.decode('ascii'))
    LOG.info('Gateways {gws} presented in network {gw_net}'.format(
        gws=data, gw_net=gw_net))
    return gw_id in data


def _get_delete_gateway_url(gw_net_type, gw_net, gw_id):
    if gw_net_type == 'carrier_wifi_network':
        url = 'magma/v1/cwf/{gw_net}/gateways/{gw_id}'.format(
                gw_net=gw_net, gw_id=gw_id)
    elif gw_net_type == 'feg':
        url = 'magma/v1/feg/{gw_net}/gateways/{gw_id}'.format(
            gw_net=gw_net, gw_id=gw_id)
    return url


def delete_gateway(orc8r_api_url, gw_net, gw_net_type, gw_id, certs):
    LOG.info('Delete gateway {gw_id} in network {gw_net}'.format(
        gw_id=gw_id,
        gw_net=gw_net))

    delete_gw_url = _get_delete_gateway_url(gw_net_type, gw_net, gw_id)
    magma_gw_url = urljoin(orc8r_api_url, delete_gw_url)

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}
    resp = requests.delete(magma_gw_url,
                           headers=headers,
                           verify=False,
                           cert=certs)
    msg = 'Received response {text} with status code {status_code} '\
          'after gateway {gw_id} deletion'\
          .format(text=resp.text,
                  status_code=resp.status_code,
                  gw_id=gw_id)
    LOG.info(msg)
    if resp.status_code not in [200, 201, 204]:
        raise exceptions.MagmaRequestException(msg)


def get_networks(orc8r_api_url, certs):
    LOG.info('Get all networks from Magma')
    magma_nets_url = urljoin(
        orc8r_api_url,
        'magma/v1/networks')

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}

    resp = requests.get(magma_nets_url,
                        headers=headers,
                        verify=False,
                        cert=certs)
    data = json.loads(resp.content.decode('ascii'))
    LOG.info('Received networks {nets} from Magma'.format(
        nets=data))
    return data


def get_network_type(orc8r_api_url, net_id, certs):
    LOG.info('Get type for network {net_id}'.format(net_id=net_id))
    magma_net_type_url = urljoin(
        orc8r_api_url,
        'magma/v1/networks/{net_id}/type'.format(net_id=net_id))

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}

    resp = requests.get(magma_net_type_url,
                        headers=headers,
                        verify=False,
                        cert=certs)

    data = json.loads(resp.content.decode('ascii'))
    LOG.info('Type of network {net_id} is {net_type}'.format(
       net_id=net_id, net_type=data))
    return data


def _get_gws_url(net_id, net_type):
    if net_type == 'carrier_wifi_network':
        url = 'magma/v1/cwf/{net_id}/gateways'.format(net_id=net_id)
    elif net_type == 'feg':
        url = 'magma/v1/feg/{net_id}/gateways'.format(net_id=net_id)
    return url


def get_gateways(orc8r_api_url, net_id, net_type, certs):
    LOG.info('Get all gateways from {net_id} {net_type}'.format(
        net_id=net_id, net_type=net_type))
    gws_url = _get_gws_url(net_id, net_type)
    magma_gws_url = urljoin(
        orc8r_api_url,
        gws_url)

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}

    resp = requests.get(magma_gws_url,
                        headers=headers,
                        verify=False,
                        cert=certs)

    data = json.loads(resp.content.decode('ascii'))
    LOG.info('Received gateways {gws} from network {net_id}'.format(
        gws=list(data.keys()), net_id=net_id))
    LOG.debug('Gateways in network {net_id} {gws}'.format(
        net_id=net_id, gws=data))
    return data


def _get_gw_config_url(net_id, net_type, gw_id):
    if net_type == 'carrier_wifi_network':
        url = 'magma/v1/cwf/{net_id}/gateways/{gw_id}/carrier_wifi'.format(
            net_id=net_id, gw_id=gw_id)
    elif net_type == 'feg':
        url = 'magma/v1/feg/{net_id}/gateways/{gw_id}/federation'.format(
            net_id=net_id, gw_id=gw_id)
    return url


def get_gateway_config(orc8r_api_url, net_id, net_type, gw_id, certs):
    LOG.info('Get config for gateway {gw_id} in {net_type} {net_id}'.format(
        gw_id=gw_id, net_type=net_type, net_id=net_id))
    gw_cfg_url = _get_gw_config_url(net_id, net_type, gw_id)
    magma_gw_cfg_url = urljoin(orc8r_api_url, gw_cfg_url)

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}

    resp = requests.get(magma_gw_cfg_url,
                        headers=headers,
                        verify=False,
                        cert=certs)

    data = json.loads(resp.content.decode('ascii'))
    LOG.info('Received config for gateway {gw_id} '
             'from network {net_id} {net_type}'.format(
                gw_id=gw_id, net_id=net_id, net_type=net_type))
    LOG.debug('Config for cateway {gw_id} {cfg}'.format(
        gw_id=gw_id, cfg=data))
    return data


def apply_gateway_config(orc8r_api_url, net_id, net_type, gw_id, cfg, certs):
    LOG.info('Apply config to gateway {gw_id} in {net_type} {net_id}'.format(
        gw_id=gw_id, net_type=net_type, net_id=net_id))
    gw_cfg_url = _get_gw_config_url(net_id, net_type, gw_id)
    magma_gw_cfg_url = urljoin(orc8r_api_url, gw_cfg_url)

    headers = {'content-type': 'application/json',
               'accept': 'application/json'}
    resp = requests.put(magma_gw_cfg_url,
                        data=json.dumps(cfg),
                        headers=headers,
                        verify=False,
                        cert=certs)
    msg = 'Received response {text} with status code {status_code} after '\
          'applying the configuration to gateway {gw_id}'.format(
              text=resp.text,
              status_code=resp.status_code,
              gw_id=gw_id)
    LOG.info(msg)
    if resp.status_code not in [200, 201, 204]:
        raise exceptions.MagmaRequestException(msg)

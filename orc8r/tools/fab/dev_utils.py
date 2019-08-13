"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
import os
import requests
import urllib3

from fabric.api import run, hide

import fab.vagrant as vagrant


PORTAL_URL = 'https://127.0.0.1:9443/magma'

def register_vm(vm_type="magma", admin_cert=(
              './../../.cache/test_certs/admin_operator.pem',
              './../../.cache/test_certs/admin_operator.key.pem')):
    """
    Provisions the gateway vm with the cloud vm
    """
    print('Please ensure that you did "make run" in both VMs! '
          'Linking gateway and cloud VMs...')
    with hide('output', 'running', 'warnings'):
        vagrant.setup_env_vagrant(vm_type)
        hardware_id = run('cat /etc/snowflake')
    print('Found Hardware ID for gateway: %s' % hardware_id)

    # Validate if we have the right admin certs
    _validate_certs(admin_cert)
    # Create the test network
    network_id = 'test'
    networks = _cloud_get('/networks', admin_cert)
    if network_id in networks:
        print('Test network already exists!')
    else:
        print('Creating a test network...')
        _cloud_post('/networks', data={'name': 'TestNetwork'},
                    params={'requested_id': network_id}, admin_cert=admin_cert)

    # Provision the gateway
    gateways = _cloud_get('/networks/%s/gateways' % network_id, admin_cert)
    gateway_id = 'gw' + str(len(gateways) + 1)
    print('Provisioning gateway as %s...' % gateway_id)
    data = {'hardware_id': hardware_id, 'name': 'TestGateway',
            'key': {'key_type': 'ECHO'}}
    _cloud_post('/networks/%s/gateways' % network_id,
                data=data, params={'requested_id': gateway_id}, admin_cert=admin_cert)
    print('Gateway successfully provisioned as: %s' % gateway_id)


def connect_gateway_to_cloud(control_proxy_setting_path, cert_path):
    """
    Setup the gateway VM to connect to the cloud
    Path to control_proxy.yml and rootCA.pem could be specified to use
    non-default control proxy setting and certificates
    """
    # Add the override for the production endpoints
    run("sudo rm -rf /var/opt/magma/configs")
    run("sudo mkdir /var/opt/magma/configs")
    if control_proxy_setting_path is not None:
        run("sudo cp " + control_proxy_setting_path
            + " /var/opt/magma/configs/control_proxy.yml")

    # Copy certs which will be used by the bootstrapper
    run("sudo rm -rf /var/opt/magma/certs")
    run("sudo mkdir /var/opt/magma/certs")
    run("sudo cp " + cert_path + " /var/opt/magma/certs/")

    # Restart the bootstrapper in the gateway to use the new certs
    run("sudo systemctl stop magma@*")
    run("sudo systemctl restart magma@magmad")


def _validate_certs(admin_cert):
    if not os.path.isfile(admin_cert[0]) or not os.path.isfile(admin_cert[1]):
        raise Exception('Admin cert or key missing %s. Please provision VMs!' %
                        str(admin_cert))
    # Disable warnings about SSL verification since its a local VM
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def _cloud_get(resource, admin_cert):
    resp = requests.get(PORTAL_URL + resource, verify=False, cert=admin_cert)
    if resp.status_code != 200:
        raise Exception('Received a %d response: %s' %
                        (resp.status_code, resp.text))
    return resp.json()


def _cloud_post(resource, data, admin_cert, params=None):
    resp = requests.post(PORTAL_URL + resource,
                         data=json.dumps(data),
                         params=params,
                         headers={'content-type': 'application/json'},
                         verify=False,
                         cert=admin_cert)
    if resp.status_code not in [200, 201]:
        raise Exception('Received a %d response: %s' %
                        (resp.status_code, resp.text))

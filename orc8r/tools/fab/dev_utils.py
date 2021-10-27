"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import json
import os
import subprocess
from typing import Any, Dict

import jsonpickle
import requests
from fabric.api import hide, lcd, local, run, settings
from fabric.context_managers import cd
from fabric.operations import sudo
from tools.fab import types, vagrant


def register_generic_gateway(
        network_id: str,
        vm_name: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> None:
    """
    Register a generic magmad gateway.

    Args:
        network_id: Network to register inside
        vm_name: Vagrant VM name to pull HWID from
        admin_cert: Cert for API access
    """
    if not does_network_exist(network_id, admin_cert=admin_cert):
        network_payload = types.GenericNetwork(
            id=network_id, name='Test Network', description='Test Network',
            dns=types.NetworkDNSConfig(enable_caching=True, local_ttl=60),
        )
        cloud_post('networks', network_payload, admin_cert=admin_cert)

    create_tier_if_not_exists(network_id, 'default')
    hw_id = get_gateway_hardware_id_from_vagrant(vm_name=vm_name)
    already_registered, registered_as = is_hw_id_registered(network_id, hw_id)
    if already_registered:
        print(f'VM is already registered as {registered_as}')
        return

    gw_id = get_next_available_gateway_id(network_id)
    payload = construct_magmad_gateway_payload(gw_id, hw_id)
    cloud_post(f'networks/{network_id}/gateways', payload)

    print(f'Gateway {gw_id} successfully provisioned')


def construct_magmad_gateway_payload(
    gateway_id: str,
    hardware_id: str,
) -> types.Gateway:
    """
    Returns a default development magmad gateway entity given a desired gateway
    ID and a hardware ID pulled from the hardware secrets.

    Args:
        gateway_id: Desired gateway ID
        hardware_id: Hardware ID pulled from the VM

    Returns:
        Gateway object with fields filled in with reasonable default values
    """
    return types.Gateway(
        name='TestGateway',
        description='Test Gateway',
        tier='default',
        id=gateway_id,
        device=types.GatewayDevice(
            hardware_id=hardware_id,
            key=types.ChallengeKey(
                key_type='ECHO',
            ),
        ),
        magmad=types.MagmadGatewayConfigs(
            autoupgrade_enabled=True,
            autoupgrade_poll_interval=60,
            checkin_interval=60,
            checkin_timeout=30,
            dynamic_services=[],
        ),
    )


PORTAL_URL = 'https://127.0.0.1:9443/magma/v1/'


def get_next_available_gateway_id(
        network_id: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> str:
    """
    Returns the next available gateway ID in the sequence gwN for the given
    network.

    Args:
        network_id: Network to check for available gateways
        admin_cert: Client cert to use with the API

    Returns:
        Next available gateway ID in the form gwN
    """
    # gateways is a dict mapping gw ID to full resource
    gateways = cloud_get(
        f'networks/{network_id}/gateways',
        admin_cert=admin_cert,
    )

    n = len(gateways) + 1
    candidate = f'gw{n}'
    while candidate in gateways:
        n += 1
        candidate = f'gw{n}'
    return candidate


def does_network_exist(
        network_id: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> bool:
    """
    Check for the existence of a network ID
    Args:
        network_id: Network to check
        admin_cert: Cert for API access

    Returns:
        True if the network exists, False otherwise
    """
    networks = cloud_get('/networks', admin_cert)
    return network_id in networks


def create_tier_if_not_exists(
        network_id: str, tier_id: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> None:
    """
    Create a placeholder tier on Orchestrator if the specified one doesn't
    already exist.

    Args:
        network_id: Network the tier belongs to
        tier_id: ID for the tier
        admin_cert: Cert for API access
    """
    tiers = cloud_get(f'networks/{network_id}/tiers', admin_cert=admin_cert)
    if tier_id in tiers:
        return

    tier_payload = types.Tier(
        id=tier_id, version='0.0.0-0', images=[],
        gateways=[],
    )
    cloud_post(
        f'networks/{network_id}/tiers', tier_payload,
        admin_cert=admin_cert,
    )


def get_gateway_hardware_id_from_vagrant(vm_name: str) -> str:
    """
    Get the hardware ID of a gateway running on Vagrant VM

    Args:
        vm_name: Name of the vagrant machine to use

    Returns:
        Hardware snowflake from the VM
    """
    with hide('output', 'running', 'warnings'):
        vagrant.setup_env_vagrant(vm_name)
        hardware_id = run('cat /etc/snowflake')
    return str(hardware_id)


def get_gateway_hardware_id_from_docker(location_docker_compose: str) -> str:
    """
    Get the hardware ID of a gateway running on Docker

    Args:
        location_docker_compose: location of docker compose used to run FEG
        by default feg/gateway/docker
    Returns:
        Hardware snowflake from the VM
    """
    with lcd('docker'), hide('output', 'running', 'warnings'), \
            cd(location_docker_compose):
        hardware_id = local(
            'docker-compose exec magmad bash -c "cat /etc/snowflake"',
        capture=True,
        )
    return str(hardware_id)


def delete_gateway_certs_from_vagrant(vm_name: str):
    """
    Delete certificates and gw_challenge of a gateway running on Vagrant VM

    Args:
        vm_name: Name of the vagrant machine to use
    """
    with settings(warn_only=True), hide('output', 'running', 'warnings'), \
            cd('/var/opt/magma/certs'):
        vagrant.setup_env_vagrant(vm_name)
        sudo('rm gateway.*')
        sudo('rm gw_challenge.key')


def delete_gateway_certs_from_docker(location_docker_compose: str):
    """
        Delete certificates and gw_challenge of a gateway running on Docker

    Args:
        location_docker_compose: location of docker compose used to run FEG
    """
    print("delete_feg_certs is running on directory %s" % os.getcwd())

    subprocess.check_call(
        [
            'docker-compose exec magmad bash -c '
            '"rm -f /var/opt/magma/certs/gateway.*"',
        ],
        shell=True,
        cwd=location_docker_compose,
    )

    subprocess.check_call(
        [
            'docker-compose exec magmad bash -c '
            '"rm -f /var/opt/magma/certs/gw_challenge.key    "',
        ],
        shell=True,
        cwd=location_docker_compose,
    )


def is_hw_id_registered(
        network_id: str, hw_id: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> (bool, str):
    """
    Check if a hardware ID is already registered for a given network. Note that
    this is not a true guarantee that a VM is not already registered, as the
    HW ID could be taken on another network.

    Args:
        network_id: Network to check
        hw_id: HW ID to check
        admin_cert: Cert for API access

    Returns:
        (True, gw_id) if the HWID is already registered, (False, '') otherwise
    """
    # gateways is a dict mapping gw ID to full resource
    paginated_gateways = cloud_get(
        f'networks/{network_id}/gateways',
        admin_cert=admin_cert,
    )
    gateways = paginated_gateways['gateways']
    for gw in gateways.values():
        if gw['device']['hardware_id'] == hw_id:
            return True, gw['id']
    return False, ''


def connect_gateway_to_cloud(control_proxy_setting_path, cert_path):
    """
    Setup the gateway Vagrant VM to connect to the cloud
    Path to control_proxy.yml and rootCA.pem could be specified to use
    non-default control proxy setting and certificates
    """
    # Add the override for the production endpoints
    run("sudo rm -rf /var/opt/magma/configs")
    run("sudo mkdir /var/opt/magma/configs")
    if control_proxy_setting_path is not None:
        run(
            "sudo cp " + control_proxy_setting_path
            + " /var/opt/magma/configs/control_proxy.yml",
        )

    # Copy certs which will be used by the bootstrapper
    run("sudo rm -rf /var/opt/magma/certs")
    run("sudo mkdir /var/opt/magma/certs")
    run("sudo cp " + cert_path + " /var/opt/magma/certs/")

    # Restart the bootstrapper in the gateway to use the new certs
    run("sudo systemctl stop magma@*")
    run("sudo systemctl restart magma@magmad")


def cloud_get(
        resource: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> Any:
    """
    Send a GET request to an API URI
    Args:
        resource: URI to request
        admin_cert: API client certificate

    Returns:
        JSON-encoded response content
    """
    if resource.startswith("/"):
        resource = resource[1:]
    resp = requests.get(PORTAL_URL + resource, verify=False, cert=admin_cert)
    if resp.status_code != 200:
        raise Exception(
            'Received a %d response: %s' %
            (resp.status_code, resp.text),
        )
    return resp.json()


def cloud_post(
        resource: str,
        data: Any,
        params: Dict[str, str] = None,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
):
    """
    Send a POST request to an API URI

    Args:
        resource: URI to request
        data: JSON-serializable payload
        params: Params to include with the request
        admin_cert: API client certificate
    """
    resp = requests.post(
        PORTAL_URL + resource,
        data=jsonpickle.pickler.encode(data),
        params=params,
        headers={'content-type': 'application/json'},
        verify=False,
        cert=admin_cert,
    )
    if resp.status_code not in [200, 201, 204]:
        parsed = json.loads(jsonpickle.pickler.encode(data))
        raise Exception(
            'Post Request failed: \n%s\n%s \nReceived a %d response: %s\nFAILED!' %
            (resp.url, json.dumps(parsed, indent=4, sort_keys=False), resp.status_code, resp.text),
        )


def cloud_delete(
        resource: str,
        admin_cert: types.ClientCert = types.ClientCert(
            cert='./../../.cache/test_certs/admin_operator.pem',
            key='./../../.cache/test_certs/admin_operator.key.pem',
        ),
) -> Any:
    """
    Send a delete request to an API URI

    Args:
        resource: URI to request
        admin_cert: API client certificate

    Returns:
        JSON-encoded response content
    """
    if resource.startswith("/"):
        resource = resource[1:]
    resp = requests.delete(PORTAL_URL + resource, verify=False, cert=admin_cert)
    if resp.status_code not in [200, 201, 204]:
        raise Exception(
            'Delete Request failed: \n%s \nReceived a %d response: %s\nFAILED!' %
            (resp.url, resp.status_code, resp.text),
        )

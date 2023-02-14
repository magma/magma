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
import sys
import time
from typing import Any, Dict, Optional, Tuple

import jsonpickle
import requests
from fabric import Connection
from tools.fab import types
from tools.fab.hosts import vagrant_connection

AGW_ROOT = "$MAGMA_ROOT/lte/gateway"
FEG_INTEG_TEST_ROOT = "$MAGMA_ROOT/lte/gateway/python/integ_tests/federated_tests/"


def register_generic_gateway(
        c: Connection,
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
        c: Fabric connection
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
    hw_id = get_gateway_hardware_id_from_vagrant(c, vm_name=vm_name)
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
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> str:
    """
    Returns the next available gateway ID in the sequence gwN for the given
    network.

    Args:
        network_id: Network to check for available gateways
        url: API base URL
        admin_cert: Client cert to use with the API

    Returns:
        Next available gateway ID in the form gwN
    """
    # res is a dict mapping gw ID to full resource
    res = cloud_get(f'networks/{network_id}/gateways', url, admin_cert)

    # res could also return paginated gateways, so need to unwrap the mapping
    gateways = res.get('gateways', res)

    n = len(gateways) + 1
    candidate = f'gw{n}'
    while candidate in gateways:
        n += 1
        candidate = f'gw{n}'
    return candidate


def does_network_exist(
        network_id: str,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> bool:
    """
    Check for the existence of a network ID
    Args:
        network_id: Network to check
        url: API base URL
        admin_cert: Cert for API access

    Returns:
        True if the network exists, False otherwise
    """
    networks = cloud_get('/networks', url, admin_cert)
    return network_id in networks


def create_tier_if_not_exists(
        network_id: str,
        tier_id: str,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> None:
    """
    Create a placeholder tier on Orchestrator if the specified one doesn't
    already exist.

    Args:
        network_id: Network the tier belongs to
        tier_id: ID for the tier
        url: API base URL
        admin_cert: Cert for API access
    """
    tiers = cloud_get(f'networks/{network_id}/tiers', url, admin_cert)

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


def get_gateway_hardware_id_from_vagrant(c: Connection, vm_name: str) -> str:
    """
    Get the hardware ID of a gateway running on Vagrant VM

    Args:
        c: Fabric connection
        vm_name: Name of the vagrant machine to use

    Returns:
        Hardware snowflake from the VM
    """
    with vagrant_connection(c, vm_name) as c_gw:
        hardware_id = c_gw.run('cat /etc/snowflake', hide=True).stdout
    return str(hardware_id).strip()


def get_gateway_hardware_id_from_docker(c: Connection, location_docker_compose: str) -> str:
    """
    Get the hardware ID of a gateway running on Docker

    Args:
        c: Fabric connection
        location_docker_compose: location of docker compose used to run FEG
        by default feg/gateway/docker
    Returns:
        Hardware snowflake from the VM
    """
    with vagrant_connection(c, 'magma') as c_gw:
        with c_gw.cd(location_docker_compose):
            hardware_id = c_gw.run(
                'docker compose exec magmad bash -c "cat /etc/snowflake"',
                hide=True,
            ).stdout
    return str(hardware_id).strip()


def delete_gateway_certs_from_vagrant(c: Connection, vm_name: str):
    """
    Delete certificates and gw_challenge of a gateway running on Vagrant VM

    Args:
        c: Fabric connection
        vm_name: Name of the vagrant machine to use
    """
    with c.cd(AGW_ROOT):
        with vagrant_connection(c, vm_name) as c_gw:
            with c_gw.cd('/var/opt/magma/certs'):
                c_gw.run('sudo rm gateway.*', hide=True, warn=True)
                c_gw.run('sudo rm gw_challenge.key', hide=True, warn=True)


def delete_gateway_certs_from_docker(c: Connection, location_docker_compose: str):
    """
        Delete certificates and gw_challenge of a gateway running on Docker

    Args:
        c: Fabric connection
        location_docker_compose: location of docker compose used to run FEG
    """
    with c.cd(AGW_ROOT):
        with vagrant_connection(c, 'magma') as c_gw:
            with c_gw.cd(FEG_INTEG_TEST_ROOT + location_docker_compose):
                c_gw.run('echo "delete_feg_certs is running on directory $PWD"')
                c_gw.run(
                    'docker compose exec -t magmad bash -c '
                    '"rm -f /var/opt/magma/certs/gateway.*"',
                    warn=True,
                )
                c_gw.run(
                    'docker compose exec -t magmad bash -c '
                    '"rm -f /var/opt/magma/certs/gw_challenge.key"',
                    warn=True,
                )


def is_hw_id_registered(
        network_id: str,
        hw_id: str,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> Tuple[bool, str]:
    """
    Check if a hardware ID is already registered for a given network. Note that
    this is not a true guarantee that a VM is not already registered, as the
    HW ID could be taken on another network.

    Args:
        network_id: Network to check
        hw_id: HW ID to check
        url: API base URL
        admin_cert: Cert for API access

    Returns:
        (True, gw_id) if the HWID is already registered, (False, '') otherwise
    """
    # gateways is a dict mapping gw ID to full resource
    paginated_gateways = cloud_get(
        f'networks/{network_id}/gateways',
        url,
        admin_cert,
    )
    # Handle paginated or flat response
    gateways = paginated_gateways.get('gateways', paginated_gateways)
    for gw in gateways.values():
        if gw['device']['hardware_id'] == hw_id:
            return True, gw['id']
    return False, ''


def connect_gateway_to_cloud(
        c_gw: Connection,
        control_proxy_setting_path: str,
        cert_path: str,
):
    """
    Setup the gateway Vagrant VM to connect to the cloud
    Path to control_proxy.yml and rootCA.pem could be specified to use
    non-default control proxy setting and certificates
    """
    # Add the override for the production endpoints
    c_gw.run("sudo rm -rf /var/opt/magma/configs")
    c_gw.run("sudo mkdir /var/opt/magma/configs")
    if control_proxy_setting_path is not None:
        c_gw.run(
            "sudo cp " + control_proxy_setting_path
            + " /var/opt/magma/configs/control_proxy.yml",
        )

    # Copy certs which will be used by the bootstrapper
    c_gw.run("sudo rm -rf /var/opt/magma/certs")
    c_gw.run("sudo mkdir /var/opt/magma/certs")
    c_gw.run("sudo cp " + cert_path + " /var/opt/magma/certs/")

    # Restart the bootstrapper in the gateway to use the new certs
    c_gw.run("sudo systemctl stop magma@*")
    c_gw.run("sudo systemctl restart magma@magmad")


def cloud_get(
        resource: str,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> Any:
    """
    Send a GET request to an API URI
    Args:
        resource: URI to request
        url: API base URL
        admin_cert: API client certificate

    Returns:
        JSON-encoded response content
    """
    url = url or PORTAL_URL
    admin_cert = admin_cert or types.ClientCert(
        cert='./../../.cache/test_certs/admin_operator.pem',
        key='./../../.cache/test_certs/admin_operator.key.pem',
    )

    if resource.startswith("/"):
        resource = resource[1:]
    resp = requests.get(url + resource, verify=False, cert=admin_cert)
    if resp.status_code != 200:
        raise Exception(
            'Received a %d response: %s' %
            (resp.status_code, resp.text),
        )
    return resp.json()


def cloud_post(
        resource: str,
        data: Any,
        params: Optional[Dict[str, str]] = None,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
):
    """
    Send a POST request to an API URI

    Args:
        resource: URI to request
        data: JSON-serializable payload
        params: Params to include with the request
        url: API base URL
        admin_cert: API client certificate
    """
    url = url or PORTAL_URL
    admin_cert = admin_cert or types.ClientCert(
        cert='./../../.cache/test_certs/admin_operator.pem',
        key='./../../.cache/test_certs/admin_operator.key.pem',
    )
    resp = requests.post(
        url + resource,
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
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
) -> Any:
    """
    Send a delete request to an API URI

    Args:
        resource: URI to request
        url: API base URL
        admin_cert: API client certificate

    Returns:
        JSON-encoded response content
    """
    url = url or PORTAL_URL
    admin_cert = admin_cert or types.ClientCert(
        cert='./../../.cache/test_certs/admin_operator.pem',
        key='./../../.cache/test_certs/admin_operator.key.pem',
    )

    if resource.startswith("/"):
        resource = resource[1:]
    resp = requests.delete(url + resource, verify=False, cert=admin_cert)
    if resp.status_code not in [200, 201, 204]:
        raise Exception(
            'Delete Request failed: \n%s \nReceived a %d response: %s\nFAILED!' %
            (resp.url, resp.status_code, resp.text),
        )


def run_local_command_with_repetition(c, command, timeout=5):
    """
    Run command on local machine using fabric.connection.Connection.local.
    Repeats on error
    Args:
        c: fabric context
        command: command to issue
        timeout: time to execute the command while it fails
    """
    _command_with_repetition(c.local, command, timeout)


def run_remote_command_with_repetition(c, command, timeout=5):
    """
    Run command on remote machine using fabric.connection.Connection.run.
    Repeats on error
    Args:
        c: fabric connection to remote machine
        command: command to run
        timeout: time to execute the command while it fails
    """
    _command_with_repetition(c.run, command, timeout)


def _command_with_repetition(func, command, timeout=5):
    timeout = int(timeout)  # in seconds
    start_time = time.time()
    while time.time() - start_time <= timeout:
        # func will be either 'run' or 'local'
        result = func(command, hide='out', warn=True)
        if result.return_code == 0:
            # GOOD execution
            print(result)
            return
        print(
            " ┗ command failed. Trying again",
        )
        time.sleep(1)
    print(f"\nERROR on {command}\nError message:\n{result}")
    sys.exit(1)

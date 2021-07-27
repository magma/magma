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

import logging
import os

import grpc
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.service_configs import load_service_config

GRPC_KEEPALIVE_MS = 30 * 1000


class ServiceRegistry:
    """
    ServiceRegistry provides the framework to discover services.

    ServiceRegistry takes care of service naming, and sets the connection
    params like ip/port, TLS, certs, etc based on service level configuration.
    """

    _REGISTRY = {}
    _PROXY_CONFIG = {}
    _CHANNELS_CACHE = {}

    LOCAL = 'local'
    CLOUD = 'cloud'

    @staticmethod
    def get_service_address(service):
        """
        Returns the (host, port) tuple for the service.

        Args:
            service (string): Name of the service
        Returns:
            (host, port) tuple
        Raises:
            ValueError if the service is unknown
        """
        registry = ServiceRegistry.get_registry()
        if service not in registry["services"]:
            raise ValueError("Invalid service name: %s" % service)
        service_conf = registry["services"][service]
        return service_conf["ip_address"], service_conf["port"]

    @staticmethod
    def add_service(name, ip_address, port):
        """
        Adds a service to the registry.

        Args:
            name (string): Service name
            ip_address (string): ip address string
            port (int): service port
        """
        registry = ServiceRegistry.get_registry()
        service = {"ip_address": ip_address, "port": port}
        registry["services"][name] = service

    @staticmethod
    def list_services():
        """
        Returns the list of services in the registry.

        Returns:
            list of services
        """
        return ServiceRegistry.get_registry()["services"]

    @staticmethod
    def reset():
        """
        Removes all the entries in the registry
        """
        ServiceRegistry.get_registry()["services"] = {}

    @staticmethod
    def get_bootstrap_rpc_channel():
        """
        Returns a RPC channel to the bootstrap service in CLOUD.
        Returns:
            grpc channel
        """
        proxy_config = ServiceRegistry.get_proxy_config()
        (ip, port) = (
            proxy_config['bootstrap_address'],
            proxy_config['bootstrap_port'],
        )
        authority = proxy_config['bootstrap_address']

        try:
            rootca = open(proxy_config['rootca_cert'], 'rb').read()
        except FileNotFoundError as exp:
            raise ValueError("SSL cert not found: %s" % exp)

        ssl_creds = grpc.ssl_channel_credentials(rootca)
        return create_grpc_channel(ip, port, authority, ssl_creds)

    @staticmethod
    def get_rpc_channel(
        service, destination, proxy_cloud_connections=True,
        grpc_options=None,
    ):
        """
        Returns a RPC channel to the service. The connection params
        are obtained from the service registry and used.
        TBD: pool connections to a service and reuse them. Right
        now each call creates a new TCP/SSL/HTTP2 connection.

        Args:
            service (string): Name of the service
            destination (string): ServiceRegistry.LOCAL or ServiceRegistry.CLOUD
            proxy_cloud_connections (bool): Override to connect direct to cloud
            grpc_options (list): list of gRPC options params for the channel
        Returns:
            grpc channel
        Raises:
            ValueError if the service is unknown
        """
        proxy_config = ServiceRegistry.get_proxy_config()

        # Control proxy uses the :authority: HTTP header to route to services.
        if destination == ServiceRegistry.LOCAL:
            authority = '%s.local' % (service)
        else:
            authority = '%s-%s' % (service, proxy_config['cloud_address'])

        should_use_proxy = proxy_config['proxy_cloud_connections'] and \
            proxy_cloud_connections

        # If speaking to a local service or to the proxy, the grpc channel
        # can be reused. If speaking to the cloud directly, the client cert
        # could become stale after the next bootstrapper run.
        should_reuse_channel = should_use_proxy or \
            (destination == ServiceRegistry.LOCAL)
        if should_reuse_channel:
            channel = ServiceRegistry._CHANNELS_CACHE.get(authority, None)
            if channel is not None:
                return channel

        if grpc_options is None:
            grpc_options = [
                ("grpc.keepalive_time_ms", GRPC_KEEPALIVE_MS),
            ]
        # We need to figure out the ip and port to connnect, if we need to use
        # SSL and the authority to use.
        if destination == ServiceRegistry.LOCAL:
            # Connect to the local service directly
            (ip, port) = ServiceRegistry.get_service_address(service)
            channel = create_grpc_channel(
                ip, port, authority,
                options=grpc_options,
            )
        elif should_use_proxy:
            # Connect to the cloud via local control proxy
            try:
                (ip, unused_port) = ServiceRegistry.get_service_address(
                    "control_proxy",
                )
                port = proxy_config['local_port']
            except ValueError as err:
                logging.error(err)
                (ip, port) = ('127.0.0.1', proxy_config['local_port'])
            channel = create_grpc_channel(
                ip, port, authority,
                options=grpc_options,
            )
        else:
            # Connect to the cloud directly
            ip = proxy_config['cloud_address']
            port = proxy_config['cloud_port']
            ssl_creds = get_ssl_creds()
            channel = create_grpc_channel(
                ip, port, authority, ssl_creds,
                options=grpc_options,
            )
        if should_reuse_channel:
            ServiceRegistry._CHANNELS_CACHE[authority] = channel
        return channel

    @staticmethod
    def get_registry():
        """
        Returns _REGISTRY which holds the contents from the
        config/service/service_registry.yml file. Its a static member and the
        .yml file is loaded only once.
        """
        if not ServiceRegistry._REGISTRY:
            try:
                ServiceRegistry._REGISTRY = load_service_config(
                    "service_registry",
                )
            except LoadConfigError as err:
                logging.error(err)
                ServiceRegistry._REGISTRY = {"services": {}}
        return ServiceRegistry._REGISTRY

    @staticmethod
    def get_proxy_config():
        """
        Returns the control proxy config. The config file is loaded only
        once and cached.
        """
        if not ServiceRegistry._PROXY_CONFIG:
            try:
                ServiceRegistry._PROXY_CONFIG = load_service_config(
                    'control_proxy',
                )
            except LoadConfigError as err:
                logging.error(err)
                ServiceRegistry._PROXY_CONFIG = {
                    'proxy_cloud_connections': True,
                }
        return ServiceRegistry._PROXY_CONFIG


def set_grpc_cipher_suites():
    """
    Set the cipher suites to be used for the gRPC TLS connection.
    TODO (praveenr) t19265877: Update nghttpx in the cloud to recent version
        and delete this. The current nghttpx version doesn't support the
        ciphers needed by default for gRPC.
    """
    os.environ["GRPC_SSL_CIPHER_SUITES"] = "ECDHE-ECDSA-AES256-GCM-SHA384:"\
        "ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:"\
        "ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:"\
        "ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-SHA384:"\
        "ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA256:"\
        "ECDHE-RSA-AES128-SHA256"


def get_ssl_creds():
    """
    Get the SSL credentials to use to communicate securely.
    We use client side TLS auth, with the cert and keys
    obtained during bootstrapping of the gateway.

    Returns:
        gRPC ssl creds
    Raises:
        ValueError if the cert or key filename in the
            control proxy config is incorrect.
    """
    proxy_config = ServiceRegistry.get_proxy_config()
    try:
        rootca = open(proxy_config['rootca_cert'], 'rb').read()
        cert = open(proxy_config['gateway_cert']).read().encode()
        key = open(proxy_config['gateway_key']).read().encode()
        ssl_creds = grpc.ssl_channel_credentials(
            root_certificates=rootca,
            certificate_chain=cert,
            private_key=key,
        )
    except FileNotFoundError as exp:
        raise ValueError("SSL cert not found: %s" % exp)
    return ssl_creds


def create_grpc_channel(ip, port, authority, ssl_creds=None, options=None):
    """
    Helper function to create a grpc channel.

    Args:
       ip: IP address of the remote endpoint
       port: port of the remote endpoint
       authority: HTTP header that control proxy uses for routing
       ssl_creds: Enables SSL
       options: configuration options for gRPC channel
    Returns:
        grpc channel
    """
    grpc_options = [('grpc.default_authority', authority)]
    if options is not None:
        grpc_options.extend(options)
    if ssl_creds is not None:
        set_grpc_cipher_suites()
        channel = grpc.secure_channel(
            target='%s:%s' % (ip, port),
            credentials=ssl_creds,
            options=grpc_options,
        )
    else:
        channel = grpc.insecure_channel(
            target='%s:%s' % (ip, port),
            options=grpc_options,
        )
    return channel

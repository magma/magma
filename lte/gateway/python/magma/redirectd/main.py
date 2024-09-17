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
import threading

from lte.protos.mconfig import mconfigs_pb2
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.configuration.service_configs import get_service_config_value
from magma.redirectd.redirect_server import run_flask


def main():
    """
    main() for redirectd. Starts the server threads.
    """
    service = MagmaService('redirectd', mconfigs_pb2.RedirectD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    redirect_ip = get_service_config_value(
        'pipelined',
        'bridge_ip_address', None,
    )
    if redirect_ip is None:
        logging.error("ERROR bridge_ip_address not found in pipelined config")
        service.close()
        return

    http_port = service.config['http_port']
    exit_callback = get_exit_server_thread_callback(service)
    run_server_thread(run_flask, redirect_ip, http_port, exit_callback)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def get_exit_server_thread_callback(service):
    def on_exit_server_thread():
        service.StopService(None, None)

    return on_exit_server_thread


def run_server_thread(target, ip, port, exit_callback):
    """ Start redirectd service server thread """
    thread = threading.Thread(
        target=target,
        args=(ip, port, exit_callback),
    )
    thread.daemon = True
    thread.start()


if __name__ == "__main__":
    main()

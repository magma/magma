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
import asyncio
import logging

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBCloudStub
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.subscriberdb.client import SubscriberDBCloudClient
from magma.subscriberdb.processor import Processor
from magma.subscriberdb.protocols.diameter.application import base, s6a
from magma.subscriberdb.protocols.diameter.server import S6aServer
from magma.subscriberdb.protocols.m5g_auth_servicer import M5GAuthRpcServicer
from magma.subscriberdb.protocols.s6a_proxy_servicer import S6aProxyRpcServicer
from magma.subscriberdb.rpc_servicer import SubscriberDBRpcServicer
from magma.subscriberdb.store.sqlite import SqliteStore
from magma.subscriberdb.subscription_profile import get_default_sub_profile


def main():
    """Main routine for subscriberdb service."""  # noqa: D401
    service = MagmaService('subscriberdb', mconfigs_pb2.SubscriberDB())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name)

    # Initialize a store to keep all subscriber data.
    store = SqliteStore(
        service.config.get('db_path'), loop=service.loop,
        sid_digits=service.config.get('sid_last_n'),
    )

    # Initialize the processor
    processor = Processor(
        store,
        get_default_sub_profile(service),
        service.mconfig.sub_profiles,
        service.mconfig.lte_auth_op,
        service.mconfig.lte_auth_amf,
    )

    # Add all servicers to the server
    subscriberdb_servicer = SubscriberDBRpcServicer(
        store,
        service.config.get('print_grpc_payload', False),
    )
    subscriberdb_servicer.add_to_server(service.rpc_server)

    # Start a background thread to stream updates from the cloud
    if service.config.get('enable_streaming'):
        grpc_client_manager = GRPCClientManager(
            service_name="subscriberdb",
            service_stub=SubscriberDBCloudStub,
            max_client_reuse=60,
        )
        sync_interval = service.mconfig.sync_interval
        subscriber_page_size = service.config.get('subscriber_page_size')
        subscriberdb_cloud_client = SubscriberDBCloudClient(
            service.loop,
            store,
            subscriber_page_size,
            sync_interval,
            grpc_client_manager,
        )

        subscriberdb_cloud_client.start()
    else:
        logging.info(
            'enable_streaming set to False. Subscriber streaming '
            'disabled!',
        )

    # Wait until the datastore is populated by addition or resync before
    # listening for clients.
    async def serve():  # noqa: WPS430
        if not store.list_subscribers():
            # Waiting for subscribers to be added to store
            await store.on_ready()

        if service.config.get('m5g_auth_proc'):
            logging.info('Cater to 5G Authentication')
            m5g_subs_auth_servicer = M5GAuthRpcServicer(
                processor,
                service.config.get('print_grpc_payload', False),
            )
            m5g_subs_auth_servicer.add_to_server(service.rpc_server)

        if service.config.get('s6a_over_grpc'):
            logging.info('Running s6a over grpc')
            s6a_proxy_servicer = S6aProxyRpcServicer(
                processor,
                service.config.get('print_grpc_payload', False),
            )
            s6a_proxy_servicer.add_to_server(service.rpc_server)
        else:
            logging.info('Running s6a over DIAMETER')
            base_manager = base.BaseApplication(
                service.config.get('mme_realm'),
                service.config.get('mme_host_name'),
                service.config.get('mme_host_address'),
            )
            s6a_manager = _get_s6a_manager(service, processor)
            base_manager.register(s6a_manager)

            # Setup the Diameter/s6a MME
            s6a_server = service.loop.create_server(
                lambda: S6aServer(
                    base_manager,
                    s6a_manager,
                    service.config.get('mme_realm'),
                    service.config.get('mme_host_name'),
                    loop=service.loop,
                ),
                service.config.get('host_address'), service.config.get('mme_port'),
            )
            asyncio.ensure_future(s6a_server, loop=service.loop)
    asyncio.ensure_future(serve(), loop=service.loop)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def _get_s6a_manager(service, processor):
    return s6a.S6AApplication(
        processor,
        service.config.get('mme_realm'),
        service.config.get('mme_host_name'),
        service.config.get('mme_host_address'),
        service.loop,
    )


if __name__ == "__main__":
    main()


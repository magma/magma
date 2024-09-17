"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import logging

import grpc
import prometheus_client
from magma.common.metrics_export import get_metrics
from magma.metricsd_client.config import get_config
from orc8r.protos.metricsd_pb2 import RawMetricsContainer
from orc8r.protos.metricsd_pb2_grpc import CloudMetricsControllerStub

logging.basicConfig(
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S',
    format='%(asctime)s %(levelname)-8s %(message)s',
)
logger = logging.getLogger("metricsd_client.client")


config = get_config()


def get_metricsd_client() -> CloudMetricsControllerStub:
    """
    Get metricsd gRPC client

    """
    logger.info("getting metricsd GRPC channel")
    grpc_channel = grpc.insecure_channel(f"{config.ORC8R_METRICSD_GRPC_SERVICE}:{config.ORC8R_METRICSD_GRPC_PORT}")
    logger.info(f"{grpc_channel=}")
    return CloudMetricsControllerStub(grpc_channel)


def process_metrics(client: CloudMetricsControllerStub, host_name: str, service_name: str):
    """
    Get service metrics from the registry and push them to metricsd.

    Args:
        client (CloudMetricsControllerStub): metricsd client instance
        host_name (str): source host name
        service_name (str): source service name

    """
    logger.info(f"Processing Metrics for {service_name}")
    try:
        metric_families = list(get_metrics(prometheus_client.REGISTRY))
        if not metric_families:
            return
        container = RawMetricsContainer(
            hostName=host_name, families=metric_families, service=service_name,
        )
        client.PushRaw(
            request=container, metadata=(
                (config.MAGMA_CLIENT_CERT_SERIAL_KEY, config.MAGMA_CLIENT_CERT_SERIAL_VALUE),
            ),
        )
    except Exception as e:
        logger.error(f"Failed processing metrics for {service_name}: {e}")

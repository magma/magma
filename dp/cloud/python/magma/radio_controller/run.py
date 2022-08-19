"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging
from concurrent import futures
from signal import SIGTERM, signal

import grpc
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.interval import IntervalTrigger
from dp.protos.active_mode_pb2_grpc import (
    add_ActiveModeControllerServicer_to_server,
)
from dp.protos.requests_pb2_grpc import add_RadioControllerServicer_to_server
from magma.db_service.models import DBCbsdState, DBRequestType
from magma.db_service.session_manager import SessionManager
from magma.metricsd_client.client import get_metricsd_client, process_metrics
from magma.radio_controller.config import get_config
from magma.radio_controller.services.active_mode_controller.service import (
    ActiveModeControllerService,
)
from magma.radio_controller.services.radio_controller.service import (
    RadioControllerService,
)
from sqlalchemy import create_engine

logging.basicConfig(
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S',
    format='%(asctime)s %(levelname)-8s %(message)s',
)
logger = logging.getLogger("radio_controller.run")


def run():
    """
    Top-level function for radio controller
    """
    logger.info("Starting grpc server")
    config = get_config()
    scheduler = BackgroundScheduler()
    metricsd_client = get_metricsd_client()
    scheduler.add_job(
        process_metrics,
        args=[metricsd_client, config.SERVICE_HOSTNAME, "radio_controller"],
        trigger=IntervalTrigger(
            seconds=config.METRICS_PROCESSING_INTERVAL_SEC,
        ),
        max_instances=1,
        name="metrics_processing_job",
    )
    scheduler.start()

    logger.info(f"grpc port is: {config.GRPC_PORT}")
    db_engine = create_engine(
        url=config.SQLALCHEMY_DB_URI,
        encoding=config.SQLALCHEMY_DB_ENCODING,
        echo=config.SQLALCHEMY_ECHO,
        future=config.SQLALCHEMY_FUTURE,
        pool_size=config.SQLALCHEMY_ENGINE_POOL_SIZE,
        max_overflow=config.SQLALCHEMY_ENGINE_MAX_OVERFLOW,
    )
    session_manager = SessionManager(db_engine)
    with session_manager.session_scope() as session:
        cbsd_states = {state.name: state.id for state in session.query(DBCbsdState).all()}
        request_types = {req_type.name: req_type.id for req_type in session.query(DBRequestType).all()}
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_RadioControllerServicer_to_server(
        RadioControllerService(
            session_manager=session_manager, cbsd_states_map=cbsd_states, request_types_map=request_types,
        ), server,
    )
    add_ActiveModeControllerServicer_to_server(
        ActiveModeControllerService(session_manager=session_manager), server,
    )
    server.add_insecure_port(f"[::]:{config.GRPC_PORT}")
    server.start()
    logger.info(f"GRPC Server started on port {config.GRPC_PORT}")

    def handle_sigterm(*_):
        logger.info("Received shutdown signal")
        all_rpcs_done_event = server.stop(30)
        all_rpcs_done_event.wait(30)
        logger.info("Shut down gracefully")

    signal(SIGTERM, handle_sigterm)
    server.wait_for_termination()


if __name__ == "__main__":
    run()

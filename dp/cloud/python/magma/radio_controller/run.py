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

import importlib
import logging
import os
from concurrent import futures
from datetime import datetime
from signal import SIGTERM, signal

import grpc
from dp.protos.active_mode_pb2_grpc import (
    add_ActiveModeControllerServicer_to_server,
)
from dp.protos.enodebd_dp_pb2_grpc import add_DPServiceServicer_to_server
from dp.protos.requests_pb2_grpc import add_RadioControllerServicer_to_server
from magma.db_service.session_manager import SessionManager
from magma.radio_controller.config import Config
from magma.radio_controller.services.active_mode_controller.service import (
    ActiveModeControllerService,
)
from magma.radio_controller.services.dp.service import DPService
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
    logger.info(f"grpc port is: {config.GRPC_PORT}")
    db_engine = create_engine(
        url=config.SQLALCHEMY_DB_URI,
        encoding=config.SQLALCHEMY_DB_ENCODING,
        echo=config.SQLALCHEMY_ECHO,
        future=config.SQLALCHEMY_FUTURE,
    )
    session_manager = SessionManager(db_engine)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_RadioControllerServicer_to_server(
        RadioControllerService(session_manager=session_manager), server,
    )
    add_ActiveModeControllerServicer_to_server(
        ActiveModeControllerService(session_manager=session_manager), server,
    )
    add_DPServiceServicer_to_server(
        DPService(
            session_manager=session_manager,
            now_func=datetime.now,
        ), server,
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


def get_config() -> Config:
    """
    Get Configuration object for radio controller
    """
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module(
        '.'.join(
            f"magma.radio_controller.config.{app_config}".split('.')[:-1],
        ),
    )
    config_class = getattr(config_module, app_config.split('.')[-1])
    logger.info(str(config_class))

    return config_class()


if __name__ == "__main__":
    run()

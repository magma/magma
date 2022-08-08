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
import time
from typing import Optional

import click
import requests
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.interval import IntervalTrigger
from magma.configuration_controller.config import get_config
from magma.configuration_controller.crl_validator.crl_validator import (
    CRLValidator,
)
from magma.configuration_controller.request_consumer.request_db_consumer import (
    RequestDBConsumer,
)
from magma.configuration_controller.request_formatting.merger import (
    merge_requests,
)
from magma.configuration_controller.request_router.exceptions import (
    RequestRouterError,
)
from magma.configuration_controller.request_router.request_router import (
    RequestRouter,
)
from magma.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from magma.configuration_controller.response_processor.strategies.strategies_mapping import (
    processor_strategies,
)
from magma.db_service.session_manager import SessionManager
from magma.fluentd_client.client import FluentdClient, FluentdClientException
from magma.fluentd_client.dp_logs import make_dp_log
from magma.mappings.request_mapping import request_mapping
from magma.mappings.request_response_mapping import request_response
from magma.mappings.types import RequestTypes
from magma.metricsd_client.client import get_metricsd_client, process_metrics
from sqlalchemy import create_engine

logging.basicConfig(
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S',
    format='%(asctime)s %(levelname)-8s %(message)s',
)
logger = logging.getLogger("configuration_controller.run")


@click.group()
def cli():
    """
    CLI function for click module
    """

    pass


@cli.command()
def run():
    """
    Top-level function for configuration controller
    """
    config = get_config()
    scheduler = BackgroundScheduler()
    db_engine = create_engine(
        url=config.SQLALCHEMY_DB_URI,
        encoding=config.SQLALCHEMY_DB_ENCODING,
        echo=config.SQLALCHEMY_ECHO,
        future=config.SQLALCHEMY_FUTURE,
        pool_size=config.SQLALCHEMY_ENGINE_POOL_SIZE,
        max_overflow=config.SQLALCHEMY_ENGINE_MAX_OVERFLOW,
    )
    session_manager = SessionManager(db_engine=db_engine)
    ssl_validator = CRLValidator(urls=[config.SAS_URL])
    router = RequestRouter(
        sas_url=config.SAS_URL,
        rc_ingest_url=config.RC_INGEST_URL,
        cert_path=config.CC_CERT_PATH,
        ssl_key_path=config.CC_SSL_KEY_PATH,
        request_mapping=request_mapping,
        ssl_verify=config.SAS_CERT_PATH,
        crl_validator=ssl_validator,
    )
    fluentd_client = FluentdClient()
    metricsd_client = get_metricsd_client()

    for request_type in RequestTypes:
        req_type = request_type.value
        response_type = request_response[req_type]
        consumer = RequestDBConsumer(
            request_type=req_type,
            request_processing_limit=config.REQUEST_PROCESSING_LIMIT,
        )
        processor = ResponseDBProcessor(
            response_type=response_type,
            process_responses_func=processor_strategies[req_type]["process_responses"],
            fluentd_client=fluentd_client,
        )

        scheduler.add_job(
            process_requests,
            args=[consumer, processor, router, session_manager, fluentd_client],
            trigger=IntervalTrigger(
                seconds=config.REQUEST_PROCESSING_INTERVAL_SEC,
            ),
            max_instances=1,
            name=f"{req_type}_job",
        )

    scheduler.add_job(
        process_metrics,
        args=[metricsd_client, config.SERVICE_HOSTNAME, "configuration_controller"],
        trigger=IntervalTrigger(
            seconds=config.METRICS_PROCESSING_INTERVAL_SEC,
        ),
        max_instances=1,
        name="metrics_processing_job",
    )
    scheduler.add_job(
        ssl_validator.update_certificates,
        trigger=IntervalTrigger(
            seconds=config.CRL_CACHE_TIME,
        ),
        max_instances=1,
        name="crl_validator_certs_update_job",
    )
    scheduler.start()

    while True:
        time.sleep(1)


def process_requests(
        consumer: RequestDBConsumer,
        processor: ResponseDBProcessor,
        router: RequestRouter,
        session_manager: SessionManager,
        fluentd_client: FluentdClient,
) -> Optional[requests.Response]:
    """
    Process SAS requests
    """

    with session_manager.session_scope() as session:
        requests_map = consumer.get_pending_requests(session)
        requests_type = next(iter(requests_map))
        requests_list = requests_map[requests_type]

        if not requests_list:
            logger.debug(f"Received no {requests_type} requests.")
            return None

        no_of_requests = len(requests_list)
        logger.info(
            f'Processing {no_of_requests} {requests_type} requests',
        )
        bulked_sas_requests = merge_requests(requests_map)

        _log_requests_map(requests_map, fluentd_client)
        try:
            sas_response = router.post_to_sas(bulked_sas_requests)
            logger.info(
                f"Sent {bulked_sas_requests} to SAS and got the following response: {sas_response.content}",
            )
        except RequestRouterError as e:
            logging.error(f"Error posting request to SAS: {e}")
            return None

        logger.info(f"About to process responses {sas_response=}")
        processor.process_response(requests_list, sas_response, session)

        session.commit()

        return sas_response


def _log_requests_map(requests_map: dict, fluentd_client: FluentdClient):
    requests_type = next(iter(requests_map))
    for request in requests_map[requests_type]:
        try:
            log = make_dp_log(request)
            fluentd_client.send_dp_log(log)
        except (FluentdClientException, TypeError) as err:
            logging.error(f"Failed to log {requests_type} request. {err}")


if __name__ == '__main__':
    run()

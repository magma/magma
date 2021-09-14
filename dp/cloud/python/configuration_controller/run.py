import importlib
import logging
import os
import time
from typing import Optional

import click
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.interval import IntervalTrigger
from dp.cloud.python.configuration_controller.config import Config
from dp.cloud.python.configuration_controller.request_consumer.request_db_consumer import (
    RequestDBConsumer,
)
from dp.cloud.python.configuration_controller.request_formatting.merger import (
    merge_requests,
)
from dp.cloud.python.configuration_controller.request_router.exceptions import (
    RequestRouterException,
)
from dp.cloud.python.configuration_controller.request_router.request_router import (
    RequestRouter,
)
from dp.cloud.python.configuration_controller.response_processor.response_db_processor import (
    ResponseDBProcessor,
)
from dp.cloud.python.configuration_controller.response_processor.strategies.strategies_mapping import (
    processor_strategies,
)
from dp.cloud.python.db_service.session_manager import SessionManager
from dp.cloud.python.mappings.request_mapping import request_mapping
from dp.cloud.python.mappings.request_response_mapping import request_response
from dp.cloud.python.mappings.types import RequestTypes
from requests import Response
from sqlalchemy import create_engine

logging.basicConfig(level=logging.DEBUG,
                    datefmt='%Y-%m-%d %H:%M:%S',
                    format='%(asctime)s %(levelname)-8s %(message)s')
logger = logging.getLogger("configuration_controller.run")


@click.group()
def cli():
    pass


@cli.command()
def run():
    config = get_config()
    scheduler = BackgroundScheduler()
    db_engine = create_engine(
        url=config.SQLALCHEMY_DB_URI,
        encoding=config.SQLALCHEMY_DB_ENCODING,
        echo=config.SQLALCHEMY_ECHO,
        future=config.SQLALCHEMY_FUTURE
    )
    session_manager = SessionManager(db_engine=db_engine)
    router = RequestRouter(
        sas_url=config.SAS_URL,
        rc_ingest_url=config.RC_INGEST_URL,
        cert_path=config.CC_CERT_PATH,
        ssl_key_path=config.CC_SSL_KEY_PATH,
        request_mapping=request_mapping,
        ssl_verify=config.SAS_CERT_PATH
    )
    for request_type in RequestTypes:
        req_type = request_type.value
        response_type = request_response[req_type]
        consumer = RequestDBConsumer(
            request_type=req_type,
            request_processing_limit=config.REQUEST_PROCESSING_LIMIT,
        )
        processor = ResponseDBProcessor(
            response_type=response_type,
            request_map_key_func=processor_strategies[req_type]["request_map_key"],
            response_map_key_func=processor_strategies[req_type]["response_map_key"],
            process_responses_func=processor_strategies[req_type]["process_responses"],
        )
        scheduler.add_job(
            process_requests,
            args=[consumer, processor, router, session_manager],
            trigger=IntervalTrigger(
                seconds=config.REQUEST_PROCESSING_INTERVAL),
            max_instances=1,
            name=f"{req_type}_job"
        )
    scheduler.start()

    while True:
        time.sleep(1)


def get_config() -> Config:
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module('.'.join(f"dp.cloud.python.configuration_controller.config.{app_config}".split('.')[:-1]))
    config_class = getattr(config_module, app_config.split('.')[-1])
    return config_class()


def process_requests(
        consumer: RequestDBConsumer,
        processor: ResponseDBProcessor,
        router: RequestRouter,
        session_manager: SessionManager) -> Optional[Response]:

    with session_manager.session_scope() as session:
        requests_map = consumer.get_pending_requests(session)
        requests_type = next(iter(requests_map))
        requests_list = requests_map[requests_type]

        if not requests_list:
            logger.debug(f"Received no {requests_type} requests.")
            return

        logger.info(f'Processing {len(requests_list)} {requests_type} requests')
        bulked_sas_requests = merge_requests(requests_map)

        try:
            sas_response = router.post_to_sas(bulked_sas_requests)
            logger.info(f"Sent {bulked_sas_requests} to SAS and got the following response: {sas_response.json()}")
        except RequestRouterException as e:
            logging.error(f"Error posting request to SAS: {e}")
            return

        logger.info(f"About to process responses {sas_response=}")
        processor.process_response(requests_list, sas_response, session)

        session.commit()

        return sas_response


if __name__ == '__main__':
    run()

from time import sleep
from unittest import TestCase

import grpc
import requests
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from magma.test_runner.config import TestConfig
from retrying import retry

config = TestConfig()


class DomainProxyIntegrationTestCase(TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()

        grpc_channel = grpc.insecure_channel(
            f"{config.GRPC_SERVICE}:{config.GRPC_PORT}",
        )
        cls.dp_client = DPServiceStub(grpc_channel)
        wait_for_elastic_to_start()

    @classmethod
    def tearDownClass(cls) -> None:
        _delete_dp_elasticsearch_indices()


@retry(stop_max_attempt_number=30, wait_fixed=1000)
def wait_for_elastic_to_start() -> None:
    requests.get(f'{config.ELASTICSEARCH_URL}/_status')


def when_elastic_indexes_data():
    # TODO use retrying instead
    sleep(15)


def _delete_dp_elasticsearch_indices() -> None:
    requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")

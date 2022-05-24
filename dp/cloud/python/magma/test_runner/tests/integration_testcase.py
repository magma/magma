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

    @classmethod
    def tearDownClass(cls) -> None:
        _delete_dp_elasticsearch_indices()


def _delete_dp_elasticsearch_indices() -> None:
    requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")

from unittest import TestCase

import grpc
import requests
from dp.protos.cbsd_pb2_grpc import CbsdManagementStub
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from magma.test_runner.config import TestConfig

config = TestConfig()


class DomainProxyIntegrationTestCase(TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.maxDiff = None
        grpc_channel = grpc.insecure_channel(
            f"{config.GRPC_SERVICE}:{config.GRPC_PORT}",
        )
        cls.dp_client = DPServiceStub(grpc_channel)

    @classmethod
    def tearDownClass(cls) -> None:
        _delete_dp_elasticsearch_indices()


class Orc8rIntegrationTestCase(TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.maxDiff = None
        grpc_channel = grpc.insecure_channel(
            f"{config.ORC8R_DP_GRPC_SERVICE}:{config.ORC8R_DP_GRPC_PORT}",
        )
        cls.orc8r_dp_client = CbsdManagementStub(grpc_channel)


def _delete_dp_elasticsearch_indices() -> None:
    requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")

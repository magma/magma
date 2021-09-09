import os


class TestConfig:
    # General
    CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX = os.environ.get('CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX',
                                                             "http://domain-proxy-protocol-controller:8080/sas/v1")
    GRPC_SERVICE = os.environ.get('GRPC_SERVICE', 'domain-proxy-radio-controller')
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))

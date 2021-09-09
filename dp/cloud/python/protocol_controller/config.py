import os


class Config:
    # General
    TESTING = False
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'DEBUG')
    API_PREFIX = os.environ.get('API_PREFIX', '/sas/v1')
    RC_RESPONSE_WAIT_TIMEOUT = int(os.environ.get('RC_RESPONSE_WAIT_TIMEOUT', 60))
    RC_RESPONSE_WAIT_INTERVAL = int(os.environ.get('RC_RESPONSE_WAIT_INTERVAL', 1))

    PROTOCOL_PLUGIN = os.environ.get('PROTOCOL_PLUGIN', 'dp.cloud.python.protocol_controller.plugins.cbsd_sas.CBSDSASProtocolPlugin')

    # gRPC
    GRPC_SERVICE = os.environ.get('GRPC_SERVICE', 'domain-proxy-radio-controller')
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))

    JSON_ADD_STATUS = True
    JSON_STATUS_FIELD_NAME = '__status'
    JSON_JSONIFY_HTTP_ERRORS = True
    JSON_USE_ENCODE_METHODS = True


class DevelopmentConfig(Config):
    pass


class TestConfig(Config):
    pass


class ProductionConfig(Config):
    TESTING = False

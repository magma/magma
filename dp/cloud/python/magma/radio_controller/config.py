import os

import magma.db_service.config as conf


class Config:
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

    # gRPC
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))

    # SQLAlchemy DB URI (scheme + url)
    SQLALCHEMY_DB_URI = conf.Config().SQLALCHEMY_DB_URI
    SQLALCHEMY_DB_ENCODING = conf.Config().SQLALCHEMY_DB_ENCODING
    SQLALCHEMY_ECHO = conf.Config().SQLALCHEMY_ECHO
    SQLALCHEMY_FUTURE = conf.Config().SQLALCHEMY_FUTURE


class DevelopmentConfig(Config):
    pass


class TestConfig(Config):
    SQLALCHEMY_DB_URI = conf.TestConfig().SQLALCHEMY_DB_URI


class ProductionConfig(Config):
    pass

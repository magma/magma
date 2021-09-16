import os


class Config:
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

    # SQLAlchemy DB URI (scheme + url)
    SQLALCHEMY_DB_URI = os.environ.get('SQLALCHEMY_DB_URI', 'postgresql+psycopg2://postgres:postgres@db:5432/dp')
    SQLALCHEMY_DB_ENCODING = os.environ.get('SQLALCHEMY_DB_ENCODING', 'utf8')
    SQLALCHEMY_ECHO = False
    SQLALCHEMY_FUTURE = False


class DevelopmentConfig(Config):
    pass


class TestConfig(Config):
    SQLALCHEMY_DB_URI = os.environ.get('SQLALCHEMY_DB_URI', 'postgresql+psycopg2://postgres:postgres@db:5433/dp_test')


class ProductionConfig(Config):
    pass

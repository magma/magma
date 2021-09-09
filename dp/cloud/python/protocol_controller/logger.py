import logging
from typing import Type

from dp.cloud.python.protocol_controller.config import Config


def configure_logger(config: Type[Config]):
    logging.basicConfig(level=config.LOG_LEVEL)

from typing import Type

from dp.cloud.python.protocol_controller import config
from dp.cloud.python.protocol_controller.config import Config
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.app import create_app
from flask import Flask
from flask_testing import TestCase


class AppTestCase(TestCase):
    conf: Type[Config] = config.TestConfig()

    def create_app(self) -> Flask:
        app = create_app(self.conf)
        return app

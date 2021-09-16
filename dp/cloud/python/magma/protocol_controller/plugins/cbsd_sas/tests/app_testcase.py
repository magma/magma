from typing import Type

from magma.protocol_controller import config
from magma.protocol_controller.config import Config
from magma.protocol_controller.plugins.cbsd_sas.app import create_app
from flask import Flask
from flask_testing import TestCase


class AppTestCase(TestCase):
    conf: Type[Config] = config.TestConfig()

    def create_app(self) -> Flask:
        app = create_app(self.conf)
        return app

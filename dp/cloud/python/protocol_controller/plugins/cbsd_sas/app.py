from typing import Type

import grpc
from dp.cloud.python.protocol_controller.config import Config
from dp.cloud.python.protocol_controller.grpc_client.grpc_client import (
    GrpcClient,
)
from dp.cloud.python.protocol_controller.logger import configure_logger
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.deregistration import (
    deregistration_page,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.grant import (
    grant_page,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.heartbeat import (
    heartbeat_page,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.registration import (
    registration_page,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.relinquishment import (
    relinquishment_page,
)
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.views.spectrumInquiry import (
    spectrum_inquiry_page,
)
from flask import Flask


def create_app(conf: Type[Config]):
    app = Flask(__name__)
    app.config.from_object(conf)
    configure_logger(conf)
    register_pc_blueprints(app)
    register_extensions(app)
    return app


def register_pc_blueprints(app):
    blueprints = [
        registration_page,
        spectrum_inquiry_page,
        grant_page,
        heartbeat_page,
        relinquishment_page,
        deregistration_page,
    ]
    register_blueprints(app, blueprints, app.config['API_PREFIX'])


def register_blueprints(app, blueprints, url_prefix):
    for blueprint in blueprints:
        app.register_blueprint(blueprint, url_prefix=url_prefix)


def register_extensions(app):
    grpc_channel = grpc.insecure_channel(f"{app.config['GRPC_SERVICE']}:{app.config['GRPC_PORT']}")
    grpc_client = GrpcClient(grpc_channel)
    grpc_client.init_app(app)

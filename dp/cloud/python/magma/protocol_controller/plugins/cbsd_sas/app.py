"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from typing import Type

import grpc
from flask import Flask
from magma.protocol_controller.config import Config
from magma.protocol_controller.grpc_client.grpc_client import GrpcClient
from magma.protocol_controller.logger import configure_logger
from magma.protocol_controller.plugins.cbsd_sas.views.deregistration import (
    deregistration_page,
)
from magma.protocol_controller.plugins.cbsd_sas.views.grant import grant_page
from magma.protocol_controller.plugins.cbsd_sas.views.heartbeat import (
    heartbeat_page,
)
from magma.protocol_controller.plugins.cbsd_sas.views.registration import (
    registration_page,
)
from magma.protocol_controller.plugins.cbsd_sas.views.relinquishment import (
    relinquishment_page,
)
from magma.protocol_controller.plugins.cbsd_sas.views.spectrumInquiry import (
    spectrum_inquiry_page,
)


def create_app(conf: Type[Config]):
    """
    Create Flask application from configuration

    Parameters:
        conf: protocol controller configuration object

    Returns:
        Flask: a flask application
    """
    app = Flask(__name__)
    app.config.from_object(conf)
    configure_logger(conf)
    register_pc_blueprints(app)
    register_extensions(app)
    return app


def register_pc_blueprints(app):
    """
    Register protocol controller blueprints

    Parameters:
        app: Flask application
    """
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
    """
    Register list of blueprints in flask application

    Parameters:
        app: Flask application
        blueprints: flask blueprints
        url_prefix: URL prefix for the blueprints
    """
    for blueprint in blueprints:
        app.register_blueprint(blueprint, url_prefix=url_prefix)


def register_extensions(app):
    """
    Register protocol controller extensions

    Parameters:
        app: Flask application
    """
    grpc_channel = grpc.insecure_channel(
        f"{app.config['GRPC_SERVICE']}:{app.config['GRPC_PORT']}",
    )
    grpc_client = GrpcClient(grpc_channel)
    grpc_client.init_app(app)

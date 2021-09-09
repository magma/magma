from dp.cloud.python.protocol_controller.config import Config
from dp.cloud.python.protocol_controller.plugins.cbsd_sas.app import create_app

application = create_app(Config)

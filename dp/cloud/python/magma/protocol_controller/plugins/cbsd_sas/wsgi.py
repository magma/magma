from magma.protocol_controller.config import Config
from magma.protocol_controller.plugins.cbsd_sas.app import create_app

application = create_app(Config)

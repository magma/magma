from magma.protocol_controller.plugin import ProtocolPlugin
from magma.protocol_controller.plugins.cbsd_sas.wsgi import (
    application,
)


class CBSDSASProtocolPlugin(ProtocolPlugin):
    def initialize(self):
        application.run(host='0.0.0.0', port=8080)

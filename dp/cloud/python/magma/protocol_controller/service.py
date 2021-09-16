import importlib

from magma.protocol_controller.config import Config
from magma.protocol_controller.plugin import ProtocolPlugin


def get_plugin(config: Config) -> ProtocolPlugin:
    plugin_module = importlib.import_module('.'.join(config.PROTOCOL_PLUGIN.split('.')[:-1]))
    plugin_class = getattr(plugin_module, config.PROTOCOL_PLUGIN.split('.')[-1])
    return plugin_class()


if __name__ == "__main__":
    pc_plugin = get_plugin(Config)
    pc_plugin.initialize()

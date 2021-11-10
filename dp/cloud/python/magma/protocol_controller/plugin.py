from abc import abstractmethod


class ProtocolPlugin:
    '''
    A very simple plugin class.

    Since protocol controllers do not really have any generic shareable API or design,
    we can only assume that each plugin will have its own 'initialize' method, which
    will be responsibe for running the plugin.

    For instance, in the case of the `cbsd_sas` plugin, the plugin will start a flask server.

    The interface with other services (like Radio Controller) is based on gRPC so any other plugin
    is free to implement their own gRPC client and interface with RC.
    '''
    @abstractmethod
    def initialize(self):
        pass

import logging

from dp.protos.requests_pb2_grpc import RadioControllerStub


class GrpcClient(RadioControllerStub):
    def init_app(self, app):
        logging.info("Initializing GRPC Client")
        app.extensions = getattr(app, "extensions", {})
        app.extensions[self.__class__.__name__] = self

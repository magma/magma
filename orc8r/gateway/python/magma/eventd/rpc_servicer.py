"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
import logging
import socket
from contextlib import closing
from typing import Any, Dict

import grpc
import jsonschema
import pkg_resources
import yaml

from bravado_core.spec import Spec
from bravado_core.validate import validate_object as bravado_validate
from magma.common.rpc_utils import return_void
from orc8r.protos import eventd_pb2_grpc, eventd_pb2


class EventDRpcServicer(eventd_pb2_grpc.EventServiceServicer):
    """
    gRPC based server for EventD.
    """

    def __init__(self, config: Dict[str, Any]):
        self.fluent_bit_port = config['fluent_bit_port']
        self.tcp_timeout = config['tcp_timeout']
        self.event_registry = config['event_registry']
        # To be initialized in load_specs_from_registry
        self.event_type_to_spec = {}

    def load_specs_from_registry(self):
        """
        Loads all swagger definitions from the files specified in the
        event registry.
        """
        for event_type, info in self.event_registry.items():
            module = '{}.swagger.specs'.format(info['module'])
            filename = info['filename']
            if not pkg_resources.resource_exists(module, filename):
                raise LookupError(
                    'File {} not found under {}/swagger, please ensure that '
                    'it exists'.format(filename, info['module']))

            stream = pkg_resources.resource_stream(module, filename)
            with closing(stream) as spec_file:
                spec = yaml.safe_load(spec_file)
                if event_type not in spec['definitions']:
                    raise KeyError(
                        'Event type {} is not defined in {}, '
                        'please add the definition and re-generate '
                        'swagger specifications'.format(event_type, filename))
                self.event_type_to_spec[event_type] = spec

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        eventd_pb2_grpc.add_EventServiceServicer_to_server(self, server)

    def _validate_event(self, raw_event: str, event_type: str) -> None:
        """
        Checks if an event is registered and validates it based on
        a registered schema.
        Args:
            raw_event: The event to be validated, as a JSON-encoded string
            event_type: The type of an event, which corresponds
            to a generated model
        Returns:
            Does not return, but throws exceptions if validation fails.
        """
        event = json.loads(raw_event)

        # Event not in registry
        if event_type not in self.event_registry:
            logging.debug(
                'Event type %s not among registered event types (%s)',
                event_type, self.event_registry)
            raise KeyError(
                'Event type {} not registered, '
                'please add it to the eventd config'.format(event_type))

        # swagger_spec exists because we load it up for every event_type
        # in load_specs_from_registry()
        swagger_spec = self.event_type_to_spec[event_type]

        # Field and type checking
        bravado_spec = Spec.from_dict(swagger_spec,
                                      config={'validate_swagger_spec': False})
        bravado_validate(
            bravado_spec,
            swagger_spec['definitions'][event_type],
            event)

    @return_void
    def LogEvent(self, request: eventd_pb2.Event, context):
        """
        Logs an event.
        """
        logging.debug("Logging event: %s", request)

        try:
            self._validate_event(request.value, request.event_type)
        except (KeyError, jsonschema.ValidationError) as e:
            logging.error("KeyError for log: %s. Error: %s", request, e)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(
                'Event validation failed, Details: {}'.format(e))
            return

        try:
            with closing(socket.create_connection(
                    ('localhost', self.fluent_bit_port),
                    timeout=self.tcp_timeout)) as sock:
                logging.debug('Sending log to FluentBit')
                sock.sendall(json.dumps({
                    'stream_name': request.stream_name,
                    'event_type': request.event_type,
                    # We use event_tag as fluentd uses the "tag" field
                    'event_tag': request.tag,
                    'value': request.value
                }).encode('utf-8'))
        except socket.error as e:
            logging.error('Connection to FluentBit failed: %s', e)
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details(
                'Could not connect to FluentBit locally, Details: {}'
                .format(e))
            return

        logging.debug("Successfully logged event: %s", request)

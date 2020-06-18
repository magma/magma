"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
from unittest import TestCase

from jsonschema import ValidationError

from magma.eventd.rpc_servicer import EventDRpcServicer
from magma.common.service import MagmaService


class EventValidationTests(TestCase):

    def setUp(self):
        # A test event registry that specifies the test events
        test_events_location = {
            'module': 'orc8r',
            'filename': 'test_event_definitions.yml'
        }
        config = {
            'fluent_bit_port': '',
            'tcp_timeout': '',
            'event_registry': {
                'simple_event': test_events_location,
                'array_and_object_event': test_events_location,
                'null_event': test_events_location,
            },
        }
        servicer = EventDRpcServicer(config)
        servicer.load_specs_from_registry()
        self.validate_event = servicer._validate_event

    def test_event_registration(self):
        # Errors when event is not registered
        with self.assertRaises(Exception):
            self.validate_event(json.dumps({'foo': 'asdf', 'bar': 123}),
                                'non_existent_event')

        # Does not error when event is registered
        self.validate_event(json.dumps({'foo': 'asdf', 'bar': 123}),
                            'simple_event')

    def test_field_consistency(self):
        # Errors when there are missing fields (required fields)
        with self.assertRaises(ValidationError):
            # foo is missing
            self.validate_event(json.dumps({'bar': 123}), 'simple_event')

        # Errors on excess fields (additionalProperties set to false)
        with self.assertRaises(ValidationError):
            self.validate_event(
                json.dumps({'extra_field': 12, 'foo': 'asdf', 'bar': 123}),
                'simple_event')

        # Errors when there are missing AND excess fields
        with self.assertRaises(ValidationError):
            # foo is missing
            self.validate_event(json.dumps({'extra_field': 12, 'bar': 123}),
                                'simple_event')

        # Does not error when the fields are equivalent
        self.validate_event(json.dumps({'foo': 'asdf', 'bar': 123}),
                            'simple_event')

        # Does not error when event has no fields
        self.validate_event(json.dumps({}), 'null_event')

    def test_type_checking(self):
        # Does not error when the types match
        self.validate_event(
            json.dumps({
                'an_array': ["a", "b"],
                'an_object': {
                    "a_key": 1,
                    "b_key": 1
                }
            }),
            'array_and_object_event')

        # Errors when the type is wrong for primitive fields
        with self.assertRaises(ValidationError):
            self.validate_event(json.dumps({'foo': 123, 'bar': 'asdf'}),
                                'simple_event')

        # Errors when the type is wrong for array
        with self.assertRaises(ValidationError):
            self.validate_event(
                json.dumps({
                    'an_array': [1, 2, 3],
                    'an_object': {}
                }),
                'array_and_object_event')

        # Errors when the value type is wrong for object
        with self.assertRaises(ValidationError):
            self.validate_event(
                json.dumps({
                    'an_array': ["a", "b"],
                    'an_object': {
                        "a_key": "wrong_value"
                    }
                }),
                'array_and_object_event')

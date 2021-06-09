"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import json
from unittest import TestCase

from jsonschema import ValidationError
from magma.eventd.event_validator import EventValidator


class EventValidationTests(TestCase):

    def setUp(self):
        # A test event registry that specifies the test events
        test_events_location = {
            'module': 'orc8r',
            'filename': 'test_event_definitions.yml',
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
        self.validator = EventValidator(config)

    def test_event_registration(self):
        data = json.dumps({
            'foo': 'magma',  # required
            'bar': 123,
        })
        # Errors when event is not registered
        with self.assertRaises(Exception):
            self.validator.validate_event(data, 'non_existent_event')

        # Does not error when event is registered
        self.validator.validate_event(data, 'simple_event')

    def test_field_consistency(self):
        # Errors when there are missing fields (required fields)
        with self.assertRaises(ValidationError):
            # foo is missing
            data = json.dumps({
                'bar': 123,
            })
            self.validator.validate_event(data, 'simple_event')

        # Errors on excess fields (additionalProperties set to false)
        with self.assertRaises(ValidationError):
            data = json.dumps({
                'extra_field': 12,
                'foo': 'asdf',
                'bar': 123,
            })
            self.validator.validate_event(data, 'simple_event')

        # Errors when there are missing AND excess fields
        with self.assertRaises(ValidationError):
            data = json.dumps({
                'extra_field': 12,
                'bar': 123,
            })
            # foo is missing
            self.validator.validate_event(data, 'simple_event')

        # Does not error when the fields are equivalent
        data = json.dumps({
            'foo': 'magma',  # required
            'bar': 123,
        })
        self.validator.validate_event(data, 'simple_event')

        # Does not error when event has no fields
        self.validator.validate_event(json.dumps({}), 'null_event')

    def test_type_checking(self):
        data = json.dumps({
            'an_array': ["a", "b"],
            'an_object': {
                "a_key": 1,
                "b_key": 1,
            },
        })
        # Does not error when the types match
        self.validator.validate_event(data, 'array_and_object_event')

        # Errors when the type is wrong for primitive fields
        with self.assertRaises(ValidationError):
            data = json.dumps({
                'foo': 123,
                'bar': 'asdf',
            })
            self.validator.validate_event(data, 'simple_event')

        # Errors when the type is wrong for array
        with self.assertRaises(ValidationError):
            data = json.dumps({
                'an_array': [1, 2, 3],
                'an_object': {},
            })
            self.validator.validate_event(data, 'array_and_object_event')

        # Errors when the value type is wrong for object
        with self.assertRaises(ValidationError):
            data = json.dumps({
                'an_array': ["a", "b"],
                'an_object': {
                    "a_key": "wrong_value",
                },
            })
            self.validator.validate_event(data, 'array_and_object_event')

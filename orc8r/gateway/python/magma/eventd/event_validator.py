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
import logging
from contextlib import closing
from typing import Any, Dict
import pkg_resources
import yaml

from bravado_core.spec import Spec
from bravado_core.validate import validate_object as bravado_validate

class EventValidator(object):
    """
    gRPC based server for EventD.
    """
    def __init__(self, config: Dict[str, Any]):
        self.event_registry = config['event_registry']
        self.event_type_to_spec = self._load_specs_from_registry()

    def validate_event(self, raw_event: str, event_type: str) -> None:
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
                'please add it to the EventD config'.format(event_type))

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


    def _load_specs_from_registry(self) -> Dict[str, Any]:
        """
        Loads all swagger definitions from the files specified in the
        event registry.
        """
        event_type_to_spec = {}
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
                event_type_to_spec[event_type] = spec
        return event_type_to_spec

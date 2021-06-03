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

EVENT_REGISTRY = 'event_registry'
SWAGGER_SPEC = 'swagger_spec'
BRAVADO_SPEC = 'bravado_spec'
MODULE = 'module'
FILENAME = 'filename'
DEFINITIONS = 'definitions'


class EventValidator(object):
    """
    gRPC based server for EventD.
    """

    def __init__(self, config: Dict[str, Any]):
        self.event_registry = config[EVENT_REGISTRY]
        self.specs_by_filename = self._load_specs_from_registry()

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
                event_type, self.event_registry,
            )
            raise KeyError(
                'Event type {} not registered, '
                'please add it to the EventD config'.format(event_type),
            )
        filename = self.event_registry[event_type][FILENAME]
        bravado_validate(
            self.specs_by_filename[filename][BRAVADO_SPEC],
            self.specs_by_filename[filename][SWAGGER_SPEC][event_type],
            event,
        )

    def _load_specs_from_registry(self) -> Dict[str, Any]:
        """
        Loads all swagger definitions from the files specified in the
        event registry.
        """
        specs_by_filename = {}
        for event_type, info in self.event_registry.items():
            filename = info[FILENAME]
            if filename in specs_by_filename:
                # Spec for this file is already registered
                self._check_event_exists_in_spec(
                    specs_by_filename[filename][SWAGGER_SPEC],
                    filename,
                    event_type,
                )
                continue

            module = '{}.swagger.specs'.format(info[MODULE])
            if not pkg_resources.resource_exists(module, filename):
                raise LookupError(
                    'File {} not found under {}/swagger, please ensure that '
                    'it exists'.format(filename, info[MODULE]),
                )

            stream = pkg_resources.resource_stream(module, filename)
            with closing(stream) as spec_file:
                swagger_spec = yaml.safe_load(spec_file)
                self._check_event_exists_in_spec(
                    swagger_spec[DEFINITIONS], filename, event_type,
                )

                config = {'validate_swagger_spec': False}
                bravado_spec = Spec.from_dict(swagger_spec, config=config)
                specs_by_filename[filename] = {
                    SWAGGER_SPEC: swagger_spec[DEFINITIONS],
                    BRAVADO_SPEC: bravado_spec,
                }

        return specs_by_filename

    @staticmethod
    def _check_event_exists_in_spec(
            swagger_definitions: Dict[str, Any],
            filename: str,
            event_type: str,
    ):
        """
        Throw a KeyError if the event_type does not exist in swagger_definitions
        """
        if event_type not in swagger_definitions:
            raise KeyError(
                'Event type {} is not defined in {}, '
                'please add the definition and re-generate '
                'swagger specifications'.format(event_type, filename),
            )

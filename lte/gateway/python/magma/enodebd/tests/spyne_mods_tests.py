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
from datetime import datetime
from unittest import TestCase

from magma.enodebd.tr069 import models
from magma.enodebd.tr069.spyne_mods import as_dict


class SpineModsTests(TestCase):
    def test_as_dict(self):
        inp = models.Inform(
            DeviceId=models.DeviceIdStruct(
                Manufacturer='some_manufacturer',
                OUI='123456',
                ProductClass='some_product_class',
                SerialNumber='some_serial_number',
            ),
            Event=models.EventList(
                EventStruct=[
                    models.EventStruct(
                        EventCode='some_event_code',
                        CommandKey='some_command_key',
                    ),
                    models.EventStruct(
                        EventCode='other_event_code',
                        CommandKey='other_command_key',
                    ),
                ],
            ),
            MaxEnvelopes=1234,
            CurrentTime=datetime.fromisoformat('2021-09-15 16:15:43.351680'),
        )
        out = as_dict(inp)
        expected = {
            'DeviceId': {
                'Manufacturer': 'some_manufacturer',
                'OUI': '123456',
                'ProductClass': 'some_product_class',
                'SerialNumber': 'some_serial_number',
            },
            'Event': {
                'EventStruct': [
                    {
                        'EventCode': 'some_event_code',
                        'CommandKey': 'some_command_key',
                    },
                    {
                        'EventCode': 'other_event_code',
                        'CommandKey': 'other_command_key',
                    },
                ],
            },
            'MaxEnvelopes': '1234',
            'CurrentTime': '2021-09-15 16:15:43.351680',
        }
        self.assertEqual(out, expected)

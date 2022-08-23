"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import testing.postgresql
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBChannel,
    DBRequest,
    DBRequestType,
)
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.mappings.cbsd_states import CbsdStates
from magma.radio_controller.services.radio_controller.service import (
    RadioControllerService,
)
from parameterized import parameterized

Postgresql = testing.postgresql.PostgresqlFactory(cache_initialized_db=True)


class RadioControllerTestCase(LocalDBTestCase):

    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()

        self.cbsd_states = {state.name: state.id for state in self.session.query(DBCbsdState).all()}
        self.request_types = {req_type.name: req_type.id for req_type in self.session.query(DBRequestType).all()}

        self.rc_service = RadioControllerService(
            SessionManager(self.engine), cbsd_states_map=self.cbsd_states, request_types_map=self.request_types,
        )

    @parameterized.expand([
        (
            {
                "registrationRequest":
                [
                    {"fccId": "foo1", "cbsdSerialNumber": "foo2"},
                    {"fccId": "foo1", "cbsdSerialNumber": "foo2"},
                ],
            }, [1, 2],
        ),
        (
            {
                "deregistrationRequest":
                [
                    {"cbsdId": "foo1"},
                    {"cbsdId": "foo1"},
                ],
            }, [1, 2],
        ),
        (
            {
                "relinquishmentRequest":
                [
                    {"cbsdId": "foo1"},
                    {"cbsdId": "foo1"},
                ],
            }, [1, 2],
        ),
        (
            {
                "heartbeatRequest":
                [
                    {"cbsdId": "foo1"},
                    {"cbsdId": "foo1"},
                ],
            }, [1, 2],
        ),
        (
            {
                "grantRequest":
                [
                    {"cbsdId": "foo1"},
                    {"cbsdId": "foo1"},
                ],
            }, [1, 2],
        ),
        (
            {
                "spectrumInquiryRequest":
                [
                    {"cbsdId": "foo1"},
                    {"cbsdId": "foo1"},
                ],
            }, [1, 2],
        ),
    ])
    def test_store_requests_from_map_stores_requests_in_db(self, request_map, expected_list):
        # Given

        # When
        self.rc_service._store_requests_from_map_in_db(request_map)
        db_request_ids = self.session.query(DBRequest.id).all()
        db_request_ids = [_id for (_id,) in db_request_ids]

        # Then
        self.assertListEqual(db_request_ids, expected_list)

    def test_get_or_create_cbsd_doesnt_create_already_existing_entities(self):
        # Given
        payload = {"fccId": "foo1", "cbsdSerialNumber": "foo2"}
        # No cbsds in the db
        # When
        self.rc_service._get_or_create_cbsd(
            self.session,
            "registrationRequest",
            payload,
        )
        self.session.commit()

        cbsd1 = self.session.query(DBCbsd).first()

        self.rc_service._get_or_create_cbsd(
            self.session,
            "registrationRequest",
            payload,
        )
        self.session.commit()
        cbsd2 = self.session.query(DBCbsd).first()

        # Then
        self.assertEqual(cbsd1.id, cbsd2.id)

    @parameterized.expand([
        (0,),
        (1,),
        (2,),
    ])
    def test_channels_not_deleted_when_new_spectrum_inquiry_request_arrives(self, number_of_channels):
        # Given
        unregistered = self.cbsd_states[CbsdStates.UNREGISTERED.value]
        cbsd = DBCbsd(id=1, cbsd_id="foo1", state_id=unregistered, desired_state_id=unregistered)

        self._create_channels_for_cbsd(cbsd, number_of_channels)

        cbsd_channels_count_pre_request = len(cbsd.channels)

        self.assertEqual(number_of_channels, cbsd_channels_count_pre_request)

        request_map = {"spectrumInquiryRequest": [{"cbsdId": "foo1"}]}

        # When
        self.rc_service._store_requests_from_map_in_db(request_map)
        self.session.commit()

        cbsd_channels_count_post_request = len(cbsd.channels)

        # Then
        self.assertEqual(number_of_channels, cbsd_channels_count_post_request)

    def _create_channels_for_cbsd(self, cbsd: DBCbsd, number: int):
        channels = [
            DBChannel(
                cbsd=cbsd,
                low_frequency=number,
                high_frequency=number + 1,
                channel_type=f"test_type{number}",
                rule_applied=f"test_rule{number}",
                max_eirp=0.1 + number,
            ) for _ in range(0, number)
        ]
        self.session.add_all(channels)
        self.session.commit()
        return channels

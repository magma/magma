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
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.models import DBCbsdState, DBGrantState, DBRequestType
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from parameterized import parameterized


class DBInitializationTestCase(LocalDBTestCase):

    def setUp(self):
        super().setUp()
        self.initializer = DBInitializer(SessionManager(db_engine=self.engine))

    @parameterized.expand([
        (DBRequestType, 6),
        (DBGrantState, 3),
        (DBCbsdState, 2),
    ])
    def test_db_is_initialized_with_db_states_and_types(self, model, expected_post_init_count):
        # Given
        model_entities_pre_init = self.session.query(model).all()

        # When
        self.initializer.initialize()

        model_entities_post_init = self.session.query(model).all()

        # Then
        self.assertEqual(0, len(model_entities_pre_init))
        self.assertEqual(
            expected_post_init_count,
            len(model_entities_post_init),
        )

    @parameterized.expand([
        (DBRequestType,),
        (DBGrantState,),
        (DBCbsdState,),
    ])
    def test_db_is_initialized_only_once(self, model):
        # Given / When
        self.initializer.initialize()
        model_entities_post_init_1 = self.session.query(model).all()

        self.initializer.initialize()
        model_entities_post_init_2 = self.session.query(model).all()

        # Then
        self.assertListEqual(
            model_entities_post_init_1,
            model_entities_post_init_2,
        )

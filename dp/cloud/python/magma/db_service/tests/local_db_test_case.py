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
from magma.db_service.models import Base
from magma.db_service.tests.db_testcase import BaseDBTestCase


class LocalDBTestCase(BaseDBTestCase):
    postgresql: testing.postgresql.Postgresql

    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.postgresql = testing.postgresql.PostgresqlFactory(cache_initialized_db=True)()

    @classmethod
    def tearDownClass(cls) -> None:
        cls.postgresql.stop()

    def setUp(self):
        self.set_up_db_test_case(SQLALCHEMY_DB_URI=self.postgresql.url())
        self.create_all()

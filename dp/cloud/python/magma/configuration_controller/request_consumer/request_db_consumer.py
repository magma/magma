"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging

from magma.configuration_controller.custom_types.custom_types import RequestsMap
from magma.configuration_controller.metrics import PENDING_REQUESTS_FETCH_TIME
from magma.db_service.models import DBRequest, DBRequestType

logger = logging.getLogger(__name__)


class RequestDBConsumer(object):
    """Class to consume requests from the db"""

    def __init__(self, request_type: str, request_processing_limit: int):
        self.request_type = request_type
        self.request_processing_limit = request_processing_limit

    @PENDING_REQUESTS_FETCH_TIME.time()
    def get_pending_requests(self, session) -> RequestsMap:
        """
        Get requests in pending state and acquiring lock on them if they weren't previously locked
        by another process. If they have a lock on them, select the ones that don't.

        Parameters:
            session: Database session

        Returns:
            RequestsMap: Requests map object
        """
        db_requests_query = session.query(DBRequest).join(DBRequestType).filter(
            DBRequestType.name == self.request_type,
        ).with_for_update(skip_locked=True, of=DBRequest)

        if self.request_processing_limit > 0:
            db_requests_query = db_requests_query.limit(
                self.request_processing_limit,
            )

        all_db_requests = db_requests_query.all()
        db_requests_num = len(all_db_requests)

        if db_requests_num:
            logger.info(
                f"[{self.request_type}] Fetched {db_requests_num} pending <{self.request_type}> requests.",
            )

        return {self.request_type: all_db_requests}

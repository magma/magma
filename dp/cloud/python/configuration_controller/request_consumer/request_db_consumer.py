import logging

from dp.cloud.python.configuration_controller.custom_types.custom_types import (
    RequestsMap,
)
from dp.cloud.python.db_service.models import (
    DBRequest,
    DBRequestState,
    DBRequestType,
)
from dp.cloud.python.mappings.types import RequestStates

logger = logging.getLogger(__name__)


class RequestDBConsumer:
    """Class to consume requests from the db"""

    def __init__(self, request_type: str, request_processing_limit: int):
        self.request_type = request_type
        self.request_processing_limit = request_processing_limit

    def get_pending_requests(self, session) -> RequestsMap:
        """
        Getting requests in pending state and acquiring lock on them if they weren't previously locked
        by another process. If they have a lock on them, select the ones that don't.
        """
        db_requests_query = session.query(DBRequest) \
            .join(DBRequestType, DBRequestState) \
            .filter(DBRequestType.name == self.request_type,
                    DBRequestState.name == RequestStates.PENDING.value).with_for_update(skip_locked=True, of=DBRequest)

        if self.request_processing_limit > 0:
            db_requests_query = db_requests_query.limit(self.request_processing_limit)

        all_db_requests = db_requests_query.all()
        db_requests_num = len(all_db_requests)

        if db_requests_num:
            logger.info(f"[{self.request_type}] Fetched {db_requests_num} pending <{self.request_type}> requests.")

        return {self.request_type: all_db_requests}

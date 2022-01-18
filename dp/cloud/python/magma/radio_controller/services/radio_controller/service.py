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

import json
import logging
from typing import Dict, List, Optional

from dp.protos.requests_pb2 import (
    RequestDbId,
    RequestDbIds,
    RequestPayload,
    ResponsePayload,
)
from dp.protos.requests_pb2_grpc import RadioControllerServicer
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBRequest,
    DBRequestState,
    DBRequestType,
    DBResponse,
)
from magma.db_service.session_manager import Session, SessionManager
from magma.mappings.types import CbsdStates, RequestStates
from magma.radio_controller.services.radio_controller.strategies.strategies_mapping import (
    get_cbsd_filter_strategies,
)

logger = logging.getLogger(__name__)

CBSD_ID = "cbsdId"


class RadioControllerService(RadioControllerServicer):
    """
    Radio Controller gRPC Service class
    """

    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def UploadRequests(self, request_payload: RequestPayload, context) -> RequestDbIds:
        """
        Insert uploaded requests to the database

        Parameters:
            request_payload: gRPC RequestPayload message
            context: gRPC context

        Returns:
            RequestDbIds: a list of IDs of inserted database records
        """
        logger.info("Storing requests in DB.")
        requests_map = json.loads(request_payload.payload)
        db_request_ids = self._store_requests_from_map_in_db(requests_map)
        return RequestDbIds(ids=db_request_ids)

    def GetResponse(self, pb2_message: RequestDbId, context) -> ResponsePayload:
        """
        Return response for given request id

        Parameters:
            pb2_message: gRPC RequestDbId message
            context: gRPC context

        Returns:
            ResponsePayload: a gRPC ResponsePayload message
        """
        logger.info(f"Getting SAS response for request {pb2_message.id}")
        response = self._get_request_response(pb2_message.id)
        if response.payload:
            logger.info(
                f"Returning response {response.payload} for request id: {pb2_message.id}.",
            )

        return response

    def _store_requests_from_map_in_db(self, request_map: Dict[str, List[Dict]]) -> List[int]:
        request_db_ids = []
        request_type = next(iter(request_map))
        with self.session_manager.session_scope() as session:
            request_pending_state = session.query(DBRequestState).filter(
                DBRequestState.name == RequestStates.PENDING.value,
            ).scalar()
            req_type = session.query(DBRequestType).filter(
                DBRequestType.name == request_type,
            ).first()
            for request_json in request_map[request_type]:
                cbsd = self._get_or_create_cbsd(
                    session, request_type, request_json,
                )
                if not cbsd:
                    logger.error(
                        f"Could not obtain cbsd to bind to the request: {request_json}",
                    )
                    continue
                db_request = DBRequest(
                    type=req_type,
                    state=request_pending_state,
                    cbsd=cbsd,
                    payload=request_json,
                )
                if db_request:
                    logger.info(f"Adding request {db_request}.")
                    session.add(db_request)
                    session.flush()
                    request_db_ids.append(db_request.id)
            session.commit()
        return request_db_ids

    def _get_request_response(self, request_db_id: int) -> ResponsePayload:
        with self.session_manager.session_scope() as session:
            logger.info(
                f"Trying to fetch DB response for request id: {request_db_id}",
            )
            response = session.query(DBResponse).filter(
                DBResponse.request_id == request_db_id,
            ).first()
            session.commit()
            if not response:
                payload = str({})
                return ResponsePayload(payload=payload)
            return ResponsePayload(payload=json.dumps(response.payload))

    def _get_or_create_cbsd(self, session: Session, request_type: str, request_json: Dict) -> Optional[DBCbsd]:
        filters = self._get_cbsd_filters(request_type, request_json)
        cbsd = session.query(DBCbsd).filter(*filters).first()
        cbsd_id = request_json.get(CBSD_ID)
        return cbsd if cbsd else self._create_cbsd(session, request_json, cbsd_id)

    def _get_cbsd_filters(self, request_name: str, request_payload: Dict) -> List:
        return get_cbsd_filter_strategies[request_name](request_payload)

    # TODO extract this so it can be used by other grpc services
    def _create_cbsd(self, session: Session, request_payload: Dict, cbsd_id: Optional[str]):
        cbsd_state = session.query(DBCbsdState).filter(
            DBCbsdState.name == CbsdStates.UNREGISTERED.value,
        ).scalar()
        user_id = request_payload.get("userId", None)
        fcc_id = request_payload.get("fccId", None)
        cbsd_serial_number = request_payload.get("cbsdSerialNumber", None)
        cbsd = DBCbsd(
            cbsd_id=cbsd_id,
            state=cbsd_state,
            user_id=user_id,
            fcc_id=fcc_id,
            cbsd_serial_number=cbsd_serial_number,
        )
        session.add(cbsd)
        logging.info(f"New CBSD {cbsd=} created.")
        return cbsd

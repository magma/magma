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
from json.decoder import JSONDecodeError
from typing import Callable, List

from magma.db_service.models import (
    DBGrantState,
    DBLog,
    DBRequest,
    DBRequestState,
    DBResponse,
)
from magma.db_service.session_manager import Session
from magma.mappings.request_response_mapping import request_response
from magma.mappings.types import GrantStates, RequestStates
from requests import Response

logger = logging.getLogger(__name__)


class ResponseDBProcessor(object):
    """
    Class responsible for processing Database requests
    """

    def __init__(
        self,
        response_type: str,
        process_responses_func: Callable,
    ):
        self.response_type = response_type
        self.process_responses_func = process_responses_func
        self.grant_states_map = {}
        self.request_states_map = {}

    def process_response(self, requests: List[DBRequest], response: Response, session: Session) -> None:
        """
        Process SAS response.

        Parameters:
            requests: A list of Database request objects
            response: a HTTP response from SAS
            session: A database session

        Returns:
            None
        """
        try:
            logger.debug(
                f"[{self.response_type}] Processing requests: {requests} using response {response.json()}",
            )
            self._populate_grant_states_map(session)
            self._populate_request_states_map(session)
            self._process_responses(requests, response, session)
        except JSONDecodeError:
            logger.warning(
                f"[{self.response_type}] Cannot update requests from SAS reply: {response.content}",
            )
            return

    def _populate_grant_states_map(self, session):
        self.grant_states_map = {
            GrantStates.IDLE.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.IDLE.value,
            ).scalar(),
            GrantStates.GRANTED.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.GRANTED.value,
            ).scalar(),
            GrantStates.AUTHORIZED.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.AUTHORIZED.value,
            ).scalar(),
        }

    def _populate_request_states_map(self, session):
        self.request_states_map = {
            RequestStates.PROCESSED.value: session.query(DBRequestState).filter(
                DBRequestState.name == RequestStates.PROCESSED.value,
            ).scalar(),
        }

    def _process_responses(
            self,
            requests: List[DBRequest],
            sas_response: Response,
            session: Session,
    ) -> None:

        response_json_list = sas_response.json().get(self.response_type, [])
        logger.debug(
            f"[{self.response_type}] requests json list: {response_json_list}",
        )

        no_of_requests = len(requests)
        no_of_responses = len(response_json_list)
        if no_of_responses != no_of_requests:
            logger.warning(
                f"[{self.response_type}] Got {no_of_requests=} and {no_of_responses=}",
            )
        for response_json, db_request in zip(response_json_list, requests):
            db_response = DBResponse(
                response_code=int(response_json["response"]["responseCode"]),
                payload=response_json,
                request=db_request,
            )
            logger.info(
                f"[{self.response_type}] Adding Response: {db_response} for Request {db_request}",
            )
            session.add(db_response)
            self._log_response(session, db_response)
            self._process_request(db_request)
            logger.debug(
                f'[{self.response_type}] About to process Response: {db_response}',
            )
            self._process_response(db_response, session)

    def _process_response(self, response: DBResponse, session: Session) -> None:
        self.process_responses_func(self, response, session)

    def _process_request(self, request: DBRequest) -> None:
        logger.info(
            f"[{self.response_type}] Marking request {request} as processed.",
        )
        request.state = self.request_states_map[RequestStates.PROCESSED.value]

    def _log_response(self, session: Session, response: DBResponse):
        network_id = response.request.cbsd.network_id or ''
        fcc_id = response.request.cbsd.fcc_id or ''
        serial_number = response.request.cbsd.cbsd_serial_number
        log_name = request_response[response.request.type.name]
        response_code = response.payload.get(
            'response', {},
        ).get('responseCode', None)
        log = DBLog(
            log_from='SAS',
            log_to='DP',
            log_name=f'{log_name}',
            log_message=f'{response.payload}',
            cbsd_serial_number=f'{serial_number}',
            network_id=f'{network_id}',
            fcc_id=f'{fcc_id}',
            response_code=response_code,
        )
        session.add(log)

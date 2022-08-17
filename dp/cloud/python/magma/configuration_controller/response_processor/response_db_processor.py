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
import logging
from json.decoder import JSONDecodeError
from typing import Callable, List

import requests
from magma.configuration_controller.config import get_config
from magma.configuration_controller.custom_types.custom_types import DBResponse
from magma.configuration_controller.metrics import SAS_RESPONSE_PROCESSING_TIME
from magma.db_service.models import DBGrantState, DBRequest
from magma.db_service.session_manager import Session
from magma.fluentd_client.client import FluentdClient, FluentdClientException
from magma.fluentd_client.dp_logs import make_dp_log

logger = logging.getLogger(__name__)

conf = get_config()


class ResponseDBProcessor(object):
    """
    Class responsible for processing Database requests
    """

    def __init__(
        self,
        response_type: str,
        process_responses_func: Callable,
        fluentd_client: FluentdClient,
    ):
        self.response_type = response_type
        self.process_responses_func = process_responses_func
        self.grant_states_map = {}
        self.request_states_map = {}
        self.fluentd_client = fluentd_client

    @SAS_RESPONSE_PROCESSING_TIME.time()
    def process_response(self, db_requests: List[DBRequest], response: requests.Response, session: Session) -> None:
        """
        Process SAS response.

        Parameters:
            db_requests: A list of Database request objects
            response: a HTTP response from SAS
            session: A database session

        Returns:
            None
        """
        try:
            logger.debug(
                f"[{self.response_type}] Processing requests: {db_requests} using response {response.json()}",
            )
            self._populate_grant_states_map(session)
            self._process_responses(db_requests, response, session)
        except JSONDecodeError:
            logger.warning(
                f"[{self.response_type}] Cannot update requests from SAS reply: {response.content}",
            )
            return

    def _populate_grant_states_map(self, session):
        grant_states = session.query(DBGrantState).all()
        self.grant_states_map = {gs.name: gs for gs in grant_states}

    def _process_responses(
            self,
            db_requests: List[DBRequest],
            sas_response: requests.Response,
            session: Session,
    ) -> None:

        response_json_list = sas_response.json().get(self.response_type, [])
        logger.debug(
            f"[{self.response_type}] requests json list: {response_json_list}",
        )

        no_of_requests = len(db_requests)
        no_of_responses = len(response_json_list)
        if no_of_responses != no_of_requests:
            logger.warning(
                f"[{self.response_type}] Got {no_of_requests=} and {no_of_responses=}",
            )
        for response_json, db_request in zip(response_json_list, db_requests):
            response_type = next(iter(response_json))
            db_response = DBResponse(
                response_code=int(response_json["response"]["responseCode"]),
                payload=response_json,
                request=db_request,
            )
            logger.info(
                f"[{self.response_type}] Adding Response: {db_response} for Request {db_request}",
            )
            self._log_response(response_type, db_response)
            logger.debug(
                f'[{self.response_type}] About to process Response: {db_response}',
            )
            self._process_response(db_response, session)
            self._process_request(session, db_request)

    def _log_response(self, response_type: str, db_response: DBResponse):
        try:
            log = make_dp_log(db_response)
            self.fluentd_client.send_dp_log(log)
        except (FluentdClientException, TypeError) as err:
            logging.error(f"Failed to log {response_type} response. {err}")

    def _process_response(self, response: DBResponse, session: Session) -> None:
        self.process_responses_func(self, response, session)

    def _process_request(self, session: Session, request: DBRequest) -> None:
        logger.info(
            f"[{self.response_type}] Deleting processed request {request}",
        )
        session.delete(request)

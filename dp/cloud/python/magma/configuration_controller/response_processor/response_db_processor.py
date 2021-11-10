import logging
from typing import Callable, List

from magma.db_service.models import (
    DBGrantState,
    DBRequest,
    DBRequestState,
    DBResponse,
)
from magma.db_service.session_manager import Session
from magma.mappings.types import GrantStates, RequestStates
from requests import Response

logger = logging.getLogger(__name__)


class ResponseDBProcessor:
    def __init__(self,
                 response_type: str,
                 process_responses_func: Callable):
        self.response_type = response_type
        self.process_responses_func = process_responses_func
        self.grant_states_map = dict()
        self.request_states_map = dict()

    def process_response(self, requests: List[DBRequest], response: Response, session: Session) -> None:
        if not response.json():
            logger.warning(f"[{self.response_type}] Cannot update requests from SAS reply: {response.json()}")
            return

        logger.debug(f"[{self.response_type}] Processing requests: {requests} using response {response.json()}")
        self._populate_grant_states_map(session)
        self._populate_request_states_map(session)
        self._process_responses(requests, response, session)

    def _populate_grant_states_map(self, session):
        self.grant_states_map = {
            GrantStates.IDLE.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.IDLE.value).scalar(),
            GrantStates.GRANTED.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.GRANTED.value).scalar(),
            GrantStates.AUTHORIZED.value: session.query(DBGrantState).filter(
                DBGrantState.name == GrantStates.AUTHORIZED.value).scalar(),
        }

    def _populate_request_states_map(self, session):
        self.request_states_map = {
            RequestStates.PROCESSED.value: session.query(DBRequestState).filter(
                DBRequestState.name == RequestStates.PROCESSED.value).scalar()
        }

    def _process_responses(
            self,
            requests: List[DBRequest],
            sas_response: Response,
            session: Session) -> None:

        response_json_list = sas_response.json().get(self.response_type, [])
        logger.debug(f"[{self.response_type}] requests json list: {response_json_list}")

        if len(response_json_list) != len(requests):
            logger.warning(f"[{self.response_type}] not equal number of"
                           f" requests ({len(requests)}) and responses ({len(response_json_list)})")
        for response_json, db_request in zip(response_json_list, requests):
            db_response = DBResponse(
                response_code=int(response_json["response"]["responseCode"]),
                payload=response_json,
                request=db_request,
            )
            logger.info(f"[{self.response_type}] Adding Response: {db_response} for Request {db_request}")
            session.add(db_response)
            self._process_request(db_request)
            self._process_response(db_response, session)

    def _process_response(self, response: DBResponse, session: Session) -> None:
        self.process_responses_func(self, response, session)

    def _process_request(self, request: DBRequest) -> None:
        logger.info(f"[{self.response_type}] Marking request {request} as processed.")
        request.state = self.request_states_map[RequestStates.PROCESSED.value]

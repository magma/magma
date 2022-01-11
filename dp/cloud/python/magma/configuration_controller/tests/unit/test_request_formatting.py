import json
from unittest import TestCase

from magma.configuration_controller.request_formatting.merger import (
    merge_requests,
)
from magma.db_service.models import DBRequest, DBRequestState, DBRequestType
from magma.fixtures.fake_requests.registration_requests import (
    registration_requests,
)


class RequestMergingTestCase(TestCase):

    def test_request_merging_returns_empty_dict_for_empty_request_list(self):
        # Given / When
        merged_requests = merge_requests({})

        # Then
        self.assertEqual({}, merged_requests)

    def test_request_merging_merges_multiple_requests_into_one(self):
        # Given / When
        request_type = "registrationRequest"
        req_state = DBRequestState(name="pending")
        req_type = DBRequestType(name=request_type)
        reqs = [
            DBRequest(
                cbsd_id=1, state=req_state, type=req_type,
                payload=json.dumps(r[request_type]),
            )
            for r in registration_requests
        ]
        merged_requests = merge_requests({request_type: reqs})

        # Then
        self.assertIsInstance(merged_requests, dict)
        self.assertEqual(1, len(merged_requests.keys()))
        self.assertIsInstance(list(merged_requests.values())[0], list)
        self.assertEqual(2, len(list(merged_requests.values())[0]))

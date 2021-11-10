import logging
from typing import List, Optional

from magma.protocol_controller.grpc_client.grpc_client import (
    GrpcClient,
)
from dp.protos.requests_pb2 import RequestPayload


def upload_requests(client: GrpcClient, payload: str) -> Optional[List[int]]:
    grpc_req = RequestPayload(payload=payload)
    grpc_response = client.UploadRequests(grpc_req)
    req_db_ids = grpc_response.ids
    logging.info(f'Executed upload_requests action. GRPC response: {req_db_ids}')
    return req_db_ids

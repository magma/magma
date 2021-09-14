from collections import defaultdict
from typing import Dict, List

from dp.cloud.python.configuration_controller.custom_types.custom_types import (
    DBRequest,
    MergedRequests,
)


def merge_requests(request_map: Dict[str, List[DBRequest]]) -> MergedRequests:
    """
    This function receives an map of Request objects and merges them
    into one JSON object with request names as keys
    and sub-arrays of payload values

    :param request_map: Map containing list of request objects at request type key
    :return: JSON serialized object with a subarray of payloads
    """

    merged_requests = defaultdict(list)

    for request_type in request_map:
        payloads = [req.payload for req in request_map[request_type]]
        merged_requests[request_type].extend(payloads)

    return merged_requests

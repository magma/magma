from typing import Dict, List

import requests
from dp.cloud.python.configuration_controller.request_router.exceptions import (
    RequestRouterException,
)
from requests import Response


class RequestRouter:
    """
    This class is responsible for sending requests to SAS and forwarding SAS responses to Radio Controller.
    """

    def __init__(self,
                 sas_url: str,
                 rc_ingest_url: str,
                 cert_path: str,
                 ssl_key_path: str,
                 request_mapping: Dict,
                 ssl_verify: str):
        self.sas_url = sas_url
        self.rc_ingest_url = rc_ingest_url
        self.cert_path = cert_path
        self.ssl_key_path = ssl_key_path
        self.ssl_verify = ssl_verify
        self.request_mapping = request_mapping

    def post_to_sas(self, request_dict: Dict[str, List[Dict]]) -> Response:
        """
        This method parses a JSON request and sends it to the appropriate SAS endpoint.
        It will only look at the first key of the parsed JSON dict, so if there are multiple request types chunked in
        one dictionary it will send them to a SAS endpoint pertaining to the first key of the dictionary only.
        Therefore it is important to pass the requests grouped under one request name
        :param request_dict: Dictionary with a request name as key and an array of objects as value
        :return: Response object with SAS response as json payload
        """
        try:
            request_name = next(iter(request_dict))
        except StopIteration:
            raise RequestRouterException("Received an empty requests dictionary")

        try:
            sas_method = self.request_mapping[request_name]
        except KeyError:
            raise RequestRouterException(f'Unable to find SAS method matching {request_name}')

        try:
            sas_response = requests.post(
                f'{self.sas_url}/{sas_method}',
                json=request_dict,
                cert=(self.cert_path, self.ssl_key_path),
                verify=self.ssl_verify
            )
        except Exception as e:
            raise RequestRouterException(str(e))

        return sas_response

    def redirect_sas_response_to_radio_controller(self, sas_response: Response):
        """
        The method takes the Response object and passes its payload on to Radio Controller's ingest endpoint
        :param sas_response: Response object
        :return: Response object
        """
        payload = sas_response.json()
        try:
            return requests.post(self.rc_ingest_url, json=payload)
        except Exception as e:
            raise RequestRouterException(str(e))

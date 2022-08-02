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

from typing import Dict, List, Union

import requests
from magma.configuration_controller.crl_validator.crl_validator import (
    CRLValidator,
)
from magma.configuration_controller.metrics import SAS_REQUEST_PROCESSING_TIME
from magma.configuration_controller.request_router.exceptions import (
    RequestRouterError,
)


class RequestRouter(object):
    """
    This class is responsible for sending requests to SAS and forwarding SAS responses to Radio Controller.
    """

    def __init__(
            self,
            sas_url: str,
            rc_ingest_url: str,
            cert_path: str,
            ssl_key_path: str,
            request_mapping: Dict,
            ssl_verify: Union[str, bool],
            crl_validator: CRLValidator = None,
    ) -> None:
        self.sas_url = sas_url
        self.rc_ingest_url = rc_ingest_url
        self.cert_path = cert_path
        self.ssl_key_path = ssl_key_path
        self.ssl_verify = ssl_verify
        self.request_mapping = request_mapping
        self.crl_validator = crl_validator

    @SAS_REQUEST_PROCESSING_TIME.time()
    def post_to_sas(self, request_dict: Dict[str, List[Dict]]) -> requests.Response:
        """
        Parse JSON request and send it to the appropriate SAS endpoint.
        It will only look at the first key of the parsed JSON dict, so if there are multiple request types chunked in
        one dictionary it will send them to a SAS endpoint pertaining to the first key of the dictionary only.
        Therefore it is important to pass the requests grouped under one request name

        Parameters:
            request_dict: Dictionary with a request name as key and an array of objects as value

        Returns:
            requests.Response: Response object with SAS response as json payload

        Raises:
            RequestRouterError: General Request Router error
        """
        try:
            request_name = next(iter(request_dict))
        except StopIteration:
            raise RequestRouterError(
                "Received an empty requests dictionary",
            )

        try:
            sas_method = self.request_mapping[request_name]
        except KeyError:
            raise RequestRouterError(
                f'Unable to find SAS method matching {request_name}',
            )

        try:
            if self.crl_validator:
                self.crl_validator.is_valid(url=self.sas_url)

            sas_response = requests.post(
                f'{self.sas_url}/{sas_method}',
                json=request_dict,
                cert=(self.cert_path, self.ssl_key_path),
                verify=self.ssl_verify,
            )
        except Exception as e:
            raise RequestRouterError(str(e))

        return sas_response

    def redirect_sas_response_to_radio_controller(self, sas_response: requests.Response) -> requests.Response:
        """
        Send Response object to Radio Controller's ingest endpoint

        Parameters:
            sas_response: SAS Response object

        Returns:
            requests.Response: Radio Controller Response object

        Raises:
            RequestRouterError: General Request Router error
        """
        payload = sas_response.json()
        try:
            return requests.post(self.rc_ingest_url, json=payload)
        except Exception as e:
            raise RequestRouterError(str(e))

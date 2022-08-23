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

import enum


class RequestTypes(enum.Enum):
    """
    Request type class
    """
    REGISTRATION = "registrationRequest"
    SPECTRUM_INQUIRY = "spectrumInquiryRequest"
    GRANT = "grantRequest"
    HEARTBEAT = "heartbeatRequest"
    RELINQUISHMENT = "relinquishmentRequest"
    DEREGISTRATION = "deregistrationRequest"


class ResponseTypes(enum.Enum):
    """
    Response type class
    """
    REGISTRATION = "registrationResponse"
    SPECTRUM_INQUIRY = "spectrumInquiryResponse"
    GRANT = "grantResponse"
    HEARTBEAT = "heartbeatResponse"
    RELINQUISHMENT = "relinquishmentResponse"
    DEREGISTRATION = "deregistrationResponse"


class CbsdStates(enum.Enum):
    """
    CBSD SAS registration state class
    """
    UNREGISTERED = "unregistered"
    REGISTERED = "registered"


class GrantStates(enum.Enum):
    """
    Grant SAS state class
    """
    GRANTED = "granted"
    AUTHORIZED = "authorized"
    UNSYNC = "unsync"


class ResponseCodes(enum.Enum):
    """
    SAS response code class
    """
    # Success
    SUCCESS = 0

    # 100 – 199: general errors related to the SAS-CBSD protocol
    VERSION = 100
    BLACKLISTED = 101
    MISSING_PARAM = 102
    INVALID_VALUE = 103
    CERT_ERROR = 104
    DEREGISTER = 105

    # 200 – 299: error events related to the CBSD Registration procedure
    REG_PENDING = 200
    GROUP_ERROR = 201

    # 300 – 399: error events related to the Spectrum Inquiry procedure
    UNSUPPORTED_SPECTRUM = 300

    # 400 – 499: error events related to the Grant procedure
    INTERFERENCE = 400
    GRANT_CONFLICT = 401

    # 500 – 599: error events related to the Heartbeat procedure
    TERMINATED_GRANT = 500
    SUSPENDED_GRANT = 501
    UNSYNC_OP_PARAM = 502

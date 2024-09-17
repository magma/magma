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

from magma.db_service.models import DBCbsd

CBSD_SERIAL_NR = "cbsdSerialNumber"
CBSD_ID = "cbsdId"


def registration_get_cbsd_filters(request_payload):
    """
    Return sqlalchemy filter for CBSD query during registration

    Parameters:
        request_payload: SAS request JSON payload

    Returns:
        list: SQLAlchemy filter list
    """

    return [
        DBCbsd.cbsd_serial_number == request_payload.get(CBSD_SERIAL_NR),
    ]


def simple_get_cbsd_filters(request_payload):
    """
    Return sqlalchemy filter for CBSD query

    Parameters:
        request_payload: SAS request JSON payload

    Returns:
        list: SQLAlchemy filter list
    """

    return [
        DBCbsd.cbsd_id == request_payload.get(CBSD_ID),
    ]

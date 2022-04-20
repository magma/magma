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

from magma.db_service.models import DBCbsd


def get_cbsd_basic_params(cbsd: DBCbsd) -> (str, str, str):
    """

    Args:
        cbsd (DBCbsd): database CBSD entity

    Returns:
        Set of basic cbsd parameters: fcc_id, network_id, cbsd_serial_number

    """
    network_id = ''
    fcc_id = ''
    cbsd_serial_number = ''
    if cbsd:
        network_id = cbsd.network_id or ''
        fcc_id = cbsd.fcc_id or ''
        cbsd_serial_number = cbsd.cbsd_serial_number or ''
    return fcc_id, network_id, cbsd_serial_number

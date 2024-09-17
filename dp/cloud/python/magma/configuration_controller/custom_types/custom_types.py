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

from typing import Any, Dict, List, NamedTuple, Optional

from magma.db_service.models import DBRequest

MergedRequests = Dict[str, List[Dict]]
RequestsMap = Dict[str, List[DBRequest]]


# TODO why does this class have a "DB" prefix. The name is misleading
class DBResponse(NamedTuple):
    """ Class representing single response from SAS. """
    response_code: int
    payload: Dict[str, Any]
    request: DBRequest

    @property
    def cbsd_id(self) -> Optional[str]:
        """ Get cbsd_id that request refers to.

        Returns:
            string
        """
        return self.payload.get('cbsdId') or self.request.payload.get('cbsdId')

    @property
    def grant_id(self) -> Optional[str]:
        """ Get grant_id that request refers to.

        Returns:
            string
        """
        return self.payload.get('grantId') or self.request.payload.get('grantId')

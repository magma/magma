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

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

from abc import ABC, abstractmethod


class TcOpsBase(ABC):
    """
        Implements TC lower level operations to create scheduler and filters.
        Each function has all argments required to create TC object.
        There is minimal state maintained in object.
    """

    @abstractmethod
    def create_htb(self, iface: str, qid: str, max_bw: str, rate:str,
                    parent_qid: str = None) -> int:
        """
        Create HTB scheduler
        """
        ...

    @abstractmethod
    def del_htb(self, iface: str, qid: str) -> int:
        """
        Delete HTB scheduler
        """
        ...

    @abstractmethod
    def create_filter(self, iface: str, mark: str, qid: str, proto: int = 3) -> int:
        """
        Create FW filter
        """
        ...

    @abstractmethod
    def del_filter(self, iface: str, mark: str, qid: str, proto: int = 3) -> int:
        """
        Delete FW filter
        """
        ...

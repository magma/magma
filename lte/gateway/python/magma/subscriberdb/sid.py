"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from lte.protos.subscriberdb_pb2 import SubscriberID


class SIDUtils:
    """
    Utility functions for SubscriberIDs (a.k.a SIDs)
    """

    @staticmethod
    def to_str(sid_pb):
        """
        Standard way of converting a SubscriberID to a string representation

        Args:
            sid_pb - SubscriberID protobuf message
        Returns:
            string representation of the SubscriberID
        Raises:
            ValueError
        """
        if sid_pb.type == SubscriberID.IMSI:
            return 'IMSI' + sid_pb.id
        raise ValueError(
            'Invalid sid! type:%s id:%s' %
            (sid_pb.type, sid_pb.id),
        )

    @staticmethod
    def to_pb(sid_str):
        """
        Converts a string to SubscriberID protobuf message

        Args:
            sid_str: string representation of SubscriberID
        Returns:
            SubscriberID protobuf message
        Raises:
            ValueError
        """
        if sid_str.startswith('IMSI'):
            imsi = sid_str[4:]
            if not imsi.isdigit():
                raise ValueError('Invalid imsi: %s' % imsi)
            return SubscriberID(id=imsi, type=SubscriberID.IMSI)
        raise ValueError('Invalid sid: %s' % sid_str)

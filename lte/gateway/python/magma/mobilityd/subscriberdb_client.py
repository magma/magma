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

import ipaddress
import grpc
import logging

from magma.subscriberdb.sid import SIDUtils


class SubscriberDbClient:
    def __init__(self, subscriberdb_rpc_stub):
        self.subscriber_client = subscriberdb_rpc_stub

    def get_subscriber_ip(self, sid: str) -> ipaddress.ip_address:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        """
        if self.subscriber_client is None:
            return None

        ip_addr = None
        try:
            if '.' in sid:
                imsi, apn_name = sid.split('.', maxsplit=1)
            else:
                imsi, apn_name = sid, ''

            data = self.subscriber_client.GetSubscriberData(SIDUtils.to_pb(imsi))
            if data and data.non_3gpp and data.non_3gpp.apn_config:
                for apn_config in data.non_3gpp.apn_config:
                    if apn_config.service_selection == '*':
                        ip_addr = apn_config.assigned_static_ip
                    if apn_config.service_selection == apn_name:
                        ip_addr = apn_config.assigned_static_ip
                        break

                if ip_addr is not None:
                    return ipaddress.ip_address(ip_addr)
            return None

        except ValueError:
            logging.warning("Invalid sid or ip %s: [%s]", sid, ip_addr)
            return None

        except grpc.RpcError as err:
            logging.error(
                "GetSubscriberData error[%s] %s",
                err.code(),
                err.details())
            return None

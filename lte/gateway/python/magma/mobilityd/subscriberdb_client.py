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
from typing import Optional
from lte.protos.subscriberdb_pb2 import APNConfiguration

import grpc
import logging

from magma.subscriberdb.sid import SIDUtils


class SubscriberDbClient:
    def __init__(self, subscriberdb_rpc_stub):
        self.subscriber_client = subscriberdb_rpc_stub

    def get_subscriber_ip(self, sid: str) -> Optional[ipaddress.ip_address]:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        """
        if self.subscriber_client is None:
            return None

        try:
            apn_config = self._find_ip_and_apn_config(sid)
            logging.debug("ip: Got APN: %s", apn_config)
            if apn_config:
                return ipaddress.ip_address(apn_config.assigned_static_ip)

        except ValueError:
            logging.warning("Invalid data for sid %s: ", sid)
            return None

        except grpc.RpcError as err:
            logging.error(
                "GetSubscriberData while reading static ip, error[%s] %s",
                err.code(),
                err.details())
            return None

    def get_subscriber_apn_vlan(self, sid: str) -> int:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        TODO: Move this API to separate APN configuration service.
        """
        if self.subscriber_client is None:
            return 0

        try:
            apn_config = self._find_ip_and_apn_config(sid)
            logging.debug("vlan: Got APN: %s", apn_config)
            if apn_config:
                return apn_config.resource.vlan_id

        except ValueError:
            logging.warning("Invalid data for sid %s: ", sid)
            return 0

        except grpc.RpcError as err:
            logging.error(
                "GetSubscriberData while reading vlan-id error[%s] %s",
                err.code(),
                err.details())
        return 0

    # use same API to retrieve IP address and related config.
    def _find_ip_and_apn_config(self, sid: str) -> (Optional[APNConfiguration]):
        if '.' in sid:
            imsi, apn_name = sid.split('.', maxsplit=1)
        else:
            imsi, apn_name = sid, ''

        logging.debug("Find APN config for: %s", sid)
        data = self.subscriber_client.GetSubscriberData(SIDUtils.to_pb(imsi))
        if data and data.non_3gpp and data.non_3gpp.apn_config:
            selected_apn_conf = None
            for apn_config in data.non_3gpp.apn_config:
                logging.debug("APN config: %s", apn_config)
                if apn_config.assigned_static_ip is None or \
                   apn_config.assigned_static_ip == "":
                    continue
                try:
                    ipaddress.ip_address(apn_config.assigned_static_ip)
                except ValueError:
                    continue
                if apn_config.service_selection == '*':
                    selected_apn_conf = apn_config
                if apn_config.service_selection == apn_name:
                    selected_apn_conf = apn_config
                    break

            if selected_apn_conf and selected_apn_conf.assigned_static_ip:
                return selected_apn_conf

        return None

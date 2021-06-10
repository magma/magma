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
import logging
from typing import Optional

import grpc
from lte.protos.apn_pb2 import APNConfiguration
from magma.subscriberdb.sid import SIDUtils


class NetworkInfo:
    def __init__(
        self, gw_ip: Optional[str] = None, gw_mac: Optional[str] = None,
        vlan: int = 0,
    ):
        gw_ip_parsed = None
        try:
            gw_ip_parsed = ipaddress.ip_address(gw_ip)
        except ValueError:
            logging.debug("invalid internet gw ip: %s", gw_ip)

        self.gw_ip = gw_ip_parsed
        self.gw_mac = gw_mac
        self.vlan = vlan

    def __str__(self):
        return "GW-IP: {} GW-MAC: {} VLAN: {}".format(
            self.gw_ip,
            self.gw_mac,
            self.vlan,
        )


class StaticIPInfo:
    """
    Operator can configure Static GW IP and MAC.
    This would be used by AGW services to generate networking
    configuration.
    """

    def __init__(
        self, ip: Optional[str],
        gw_ip: Optional[str],
        gw_mac: Optional[str],
        vlan: int,
    ):
        if ip:
            self.ip = ipaddress.ip_address(ip)
        else:
            self.ip = None
        self.net_info = NetworkInfo(gw_ip, gw_mac, vlan)

    def __str__(self):
        return "IP: {} NETWORK: {}".format(self.ip, self.net_info)


class SubscriberDbClient:
    def __init__(self, subscriberdb_rpc_stub):
        self.subscriber_client = subscriberdb_rpc_stub

    def get_subscriber_ip(self, sid: str) -> Optional[StaticIPInfo]:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        """
        if self.subscriber_client is None:
            return None

        try:
            apn_config = self._find_ip_and_apn_config(sid)
            logging.debug("ip: Got APN: %s", apn_config)
            if apn_config and apn_config.assigned_static_ip:
                return StaticIPInfo(
                    ip=apn_config.assigned_static_ip,
                    gw_ip=apn_config.resource.gateway_ip,
                    gw_mac=apn_config.resource.gateway_mac,
                    vlan=apn_config.resource.vlan_id,
                )

        except ValueError as ex:
            logging.warning(
                "static Ip: Invalid or missing data for sid %s: ", sid,
            )
            logging.debug(ex)
            raise SubscriberDBStaticIPValueError(sid)

        except grpc.RpcError as err:
            msg = "GetSubscriberData: while reading vlan-id error[%s] %s" % \
                  (err.code(), err.details())
            logging.error(msg)
            raise SubscriberDBConnectionError(msg)
        return None

    def get_subscriber_apn_network_info(self, sid: str) -> NetworkInfo:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        TODO: Move this API to separate APN configuration service.
        """
        if self.subscriber_client:
            try:
                apn_config = self._find_ip_and_apn_config(sid)
                logging.debug("vlan: Got APN: %s", apn_config)
                if apn_config and apn_config.resource.vlan_id:
                    return NetworkInfo(
                        gw_ip=apn_config.resource.gateway_ip,
                        gw_mac=apn_config.resource.gateway_mac,
                        vlan=apn_config.resource.vlan_id,
                    )

            except ValueError as ex:
                logging.warning(
                    "vlan: Invalid or missing data for sid %s", sid,
                )
                logging.debug(ex)
                raise SubscriberDBMultiAPNValueError(sid)

            except grpc.RpcError as err:
                msg = "GetSubscriberData while reading vlan-id error[%s] %s" % \
                    (err.code(), err.details())
                logging.error(msg)
                raise SubscriberDBConnectionError(msg)

        return NetworkInfo()

    # use same API to retrieve IP address and related config.
    def _find_ip_and_apn_config(
            self, sid: str,
    ) -> (Optional[APNConfiguration]):
        if '.' in sid:
            imsi, apn_name_part = sid.split('.', maxsplit=1)
            apn_name, _ = apn_name_part.split(',', maxsplit=1)
        else:
            imsi, _ = sid.split(',', maxsplit=1)
            apn_name = ''

        logging.debug("Find APN config for: %s", sid)
        data = self.subscriber_client.GetSubscriberData(SIDUtils.to_pb(imsi))
        if data and data.non_3gpp and data.non_3gpp.apn_config:
            selected_apn_conf = None
            for apn_config in data.non_3gpp.apn_config:
                logging.debug("APN config: %s", apn_config)
                try:
                    if apn_config.assigned_static_ip:
                        ipaddress.ip_address(apn_config.assigned_static_ip)
                except ValueError:
                    continue
                if apn_config.service_selection == '*':
                    selected_apn_conf = apn_config
                elif apn_config.service_selection == apn_name:
                    selected_apn_conf = apn_config
                    break

            return selected_apn_conf

        return None


class SubscriberDBConnectionError(Exception):
    """ Exception thrown subscriber DB is not available
    """
    pass


class SubscriberDBStaticIPValueError(Exception):
    """ Exception thrown when subscriber DB has invalid IP value for the subscriber.
    """
    pass


class SubscriberDBMultiAPNValueError(Exception):
    """ Exception thrown when subscriber DB has invalid MultiAPN vlan value
    for the subscriber.
    """
    pass

#!/usr/bin/env python3
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

"""
Pre-run script for services to generate a nghttpx config from a jinja template
and the config/mconfig for the service.
"""

import logging
import os
import socket

from create_oai_certs import generate_mme_certs
from generate_service_config import generate_template_config
from lte.protos.mconfig.mconfigs_pb2 import MME
from magma.common.misc_utils import (
    IpPreference,
    get_ip_from_if,
    get_ip_from_if_cidr,
)
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.configuration.service_configs import get_service_config_value

CONFIG_OVERRIDE_DIR = "/var/opt/magma/tmp"
DEFAULT_DNS_IP_PRIMARY_ADDR = "8.8.8.8"
DEFAULT_DNS_IP_SECONDARY_ADDR = "8.8.4.4"
DEFAULT_DNS_IPV6_ADDR = "2001:4860:4860:0:0:0:0:8888"
DEFAULT_P_CSCF_IPV4_ADDR = "172.27.23.150"
DEFAULT_P_CSCF_IPV6_ADDR = "2a12:577:9941:f99c:0002:0001:c731:f114"


def _get_iface_ip(service, iface_config):
    """
    Get the interface IP given its name.
    """
    iface_name = get_service_config_value(service, iface_config, "")
    return get_ip_from_if_cidr(iface_name)


def _get_primary_dns_ip(service_mconfig, iface_config):
    """
    Get dnsd interface IP without netmask.
    If caching is enabled, use the ip of interface that dnsd listens over.
    Otherwise, use dns server from service mconfig.
    """
    if service_mconfig.enable_dns_caching:
        iface_name = get_service_config_value("dnsd", iface_config, "")
        return get_ip_from_if(iface_name)
    else:
        return service_mconfig.dns_primary or DEFAULT_DNS_IP_PRIMARY_ADDR


def _get_secondary_dns_ip(service_mconfig):
    """
    Get the secondary dns ip from the service mconfig.
    """
    return service_mconfig.dns_secondary or DEFAULT_DNS_IP_SECONDARY_ADDR


def _get_ipv4_pcscf_ip(service_mconfig):
    """
    Get IPv4 P_CSCF IP address value from service mconfig
    """
    return service_mconfig.ipv4_p_cscf_address or DEFAULT_P_CSCF_IPV4_ADDR


def _get_ipv6_pcscf_ip(service_mconfig):
    """
    Get IPv6 P_CSCF IP address value from service mconfig
    """
    return service_mconfig.ipv6_p_cscf_address or DEFAULT_P_CSCF_IPV6_ADDR


def _get_ipv6_dns_ip(service_mconfig):
    """
    Get IPV6 DNS server IP address from service mconfig
    """
    return service_mconfig.ipv6_dns_address or DEFAULT_DNS_IPV6_ADDR


def _get_oai_log_level():
    """
    Convert the logLevel in config into the level which OAI code
    uses. We use OAI's 'TRACE' as the debugging log level and 'CRITICAL'
    as the fatal log level.
    """
    oai_log_level = get_service_config_value("mme", "log_level", "INFO")
    # Translate common log levels to OAI levels
    if oai_log_level == "DEBUG":
        oai_log_level = "TRACE"
    if oai_log_level == "FATAL":
        oai_log_level = "CRITICAL"
    return oai_log_level


def _get_relay_enabled(service_mconfig):
    if service_mconfig.relay_enabled:
        return "yes"
    return "no"


def _get_non_eps_service_control(service_mconfig):
    non_eps_service_control = service_mconfig.non_eps_service_control
    if non_eps_service_control:
        if non_eps_service_control == 0:
            return "OFF"
        elif non_eps_service_control == 1:
            return "CSFB_SMS"
        elif non_eps_service_control == 2:
            return "SMS"
        elif non_eps_service_control == 3:
            return "SMS_ORC8R"
    return "OFF"


def _get_lac(service_mconfig):
    lac = service_mconfig.lac
    if lac:
        return lac
    return 0


def _get_csfb_mcc(service_mconfig):
    csfb_mcc = service_mconfig.csfb_mcc
    if csfb_mcc:
        return csfb_mcc
    return ""


def _get_csfb_mnc(service_mconfig):
    csfb_mnc = service_mconfig.csfb_mnc
    if csfb_mnc:
        return csfb_mnc
    return ""


def _get_identity():
    realm = get_service_config_value("mme", "realm", "")
    return "{}.{}".format(socket.gethostname(), realm)


def _get_enable_nat(service_mconfig):
    """
    Retrieves enable_nat config value, prioritizes service config file,
    if not found, it uses service mconfig value.
    """
    nat_enabled = get_service_config_value('mme', 'enable_nat', None)

    if nat_enabled is None:
        nat_enabled = service_mconfig.nat_enabled

    return nat_enabled


def _get_attached_enodeb_tacs(service_mconfig):
    # attachedEnodebTacs overrides 'tac', which is being deprecated, but for
    # now, both are supported
    tac = service_mconfig.tac
    attached_enodeb_tacs = service_mconfig.attached_enodeb_tacs
    if len(attached_enodeb_tacs) == 0:
        return [tac]
    return attached_enodeb_tacs


def _get_apn_correction_map_list(service_mconfig):
    if len(service_mconfig.apn_correction_map_list) != 0:
        return service_mconfig.apn_correction_map_list
    return get_service_config_value("mme", "apn_correction_map_list", None)


def _get_federated_mode_map(service_mconfig):
    if (
            service_mconfig.federated_mode_map
            and service_mconfig.federated_mode_map.enabled
            and len(service_mconfig.federated_mode_map.mapping) != 0
    ):
        return service_mconfig.federated_mode_map.mapping
    return {}


def _get_restricted_plmns(service_mconfig):
    if service_mconfig.restricted_plmns:
        return service_mconfig.restricted_plmns
    return {}


def _get_restricted_imeis(service_mconfig):
    if service_mconfig.restricted_imeis:
        return service_mconfig.restricted_imeis
    return {}


def _get_service_area_maps(service_mconfig):
    if not service_mconfig.service_area_maps:
        return {}
    service_area_map = []
    for sac, sam in service_mconfig.service_area_maps.items():
        service_area_map.append({'sac': sac, 'tac': sam.tac.copy()})
    return service_area_map


def _get_congestion_control_config(service_mconfig):
    """
    Retrieves congestion_control_enabled config value, it it does not exist
    it defaults to True.
    Args:
        service_mconfig:

    Returns: congestion control flag
    """
    congestion_control_enabled = get_service_config_value(
        'mme', 'congestion_control_enabled', None,
    )

    if congestion_control_enabled is not None:
        return congestion_control_enabled

    if service_mconfig.congestion_control_enabled is not None:
        return service_mconfig.congestion_control_enabled

    return True


def _get_context():
    """
    Create the context which has the interface IP and the OAI log level to use.
    """
    mme_service_config = load_service_mconfig('mme', MME())
    nat = _get_enable_nat(mme_service_config)
    if nat:
        iface_name = get_service_config_value(
            'spgw', 'sgw_s5s8_up_iface_name', '',
        )
    else:
        iface_name = get_service_config_value(
            'spgw', 'sgw_s5s8_up_iface_name_non_nat', '',
        )
    context = {
        "mme_s11_ip": _get_iface_ip("mme", "s11_iface_name"),
        "sgw_s11_ip": _get_iface_ip("spgw", "s11_iface_name"),
        'sgw_s5s8_up_iface_name': iface_name,
        "remote_sgw_ip": get_service_config_value("mme", "remote_sgw_ip", ""),
        "s1ap_ip": _get_iface_ip("mme", "s1ap_iface_name"),
        "oai_log_level": _get_oai_log_level(),
        "ipv4_dns": _get_primary_dns_ip(mme_service_config, "dns_iface_name"),
        "ipv4_sec_dns": _get_secondary_dns_ip(mme_service_config),
        "ipv4_p_cscf_address": _get_ipv4_pcscf_ip(mme_service_config),
        "ipv6_dns": _get_ipv6_dns_ip(mme_service_config),
        "ipv6_p_cscf_address": _get_ipv6_pcscf_ip(mme_service_config),
        "identity": _get_identity(),
        "relay_enabled": _get_relay_enabled(mme_service_config),
        "non_eps_service_control": _get_non_eps_service_control(
            mme_service_config,
        ),
        "csfb_mcc": _get_csfb_mcc(mme_service_config),
        "csfb_mnc": _get_csfb_mnc(mme_service_config),
        "lac": _get_lac(mme_service_config),
        "use_stateless": get_service_config_value("mme", "use_stateless", ""),
        "attached_enodeb_tacs": _get_attached_enodeb_tacs(mme_service_config),
        'enable_nat': nat,
        "federated_mode_map": _get_federated_mode_map(mme_service_config),
        "restricted_plmns": _get_restricted_plmns(mme_service_config),
        "restricted_imeis": _get_restricted_imeis(mme_service_config),
        "congestion_control_enabled": _get_congestion_control_config(
            mme_service_config,
        ),
        "service_area_map": _get_service_area_maps(mme_service_config),
        "sentry_config": mme_service_config.sentry_config,
    }

    context["s1u_ip"] = mme_service_config.ipv4_sgw_s1u_addr or _get_iface_ip(
        "spgw", "s1u_iface_name",
    )

    try:
        sgw_s5s8_up_ip = get_ip_from_if_cidr(iface_name, IpPreference.IPV4_ONLY)
    except ValueError:
        # ignore the error to avoid MME crash
        logging.warning("Could not read IP of interface: %s", iface_name)
        sgw_s5s8_up_ip = "127.0.0.1/8"
    context["sgw_s5s8_up_ip"] = sgw_s5s8_up_ip

    # set ovs params
    for key in (
            "ovs_bridge_name",
            "ovs_gtp_port_number",
            "ovs_mtr_port_number",
            "ovs_internal_sampling_port_number",
            "ovs_internal_sampling_fwd_tbl",
            "ovs_uplink_port_number",
            "ovs_uplink_mac",
            "pipelined_managed_tbl0",
    ):
        context[key] = get_service_config_value("spgw", key, "")
    context["enable_apn_correction"] = get_service_config_value(
        "mme", "enable_apn_correction", "",
    )
    context["apn_correction_map_list"] = _get_apn_correction_map_list(
        mme_service_config,
    )

    return context


def main():
    logging.basicConfig(
        level=logging.INFO,
        format="[%(asctime)s %(levelname)s %(name)s] %(message)s",
    )
    context = _get_context()
    generate_template_config(
        "spgw", "spgw", CONFIG_OVERRIDE_DIR,
        context.copy(),
    )
    generate_template_config("mme", "mme", CONFIG_OVERRIDE_DIR, context.copy())
    generate_template_config(
        "mme", "mme_fd", CONFIG_OVERRIDE_DIR,
        context.copy(),
    )
    cert_dir = get_service_config_value("mme", "cert_dir", "")
    generate_mme_certs(os.path.join(cert_dir, "freeDiameter"))


if __name__ == "__main__":
    main()

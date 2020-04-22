#!/usr/bin/env python3
"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

Pre-run script for services to generate a nghttpx config from a jinja template
and the config/mconfig for the service.
"""

import logging
import os
import socket
from create_oai_certs import generate_mme_certs
from generate_service_config import generate_template_config
from lte.protos.mconfig.mconfigs_pb2 import MME
from magma.common.misc_utils import get_ip_from_if, get_ip_from_if_cidr
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.configuration.service_configs import get_service_config_value

CONFIG_OVERRIDE_DIR = "/var/opt/magma/tmp"


def _get_iface_ip(service, iface_config):
    """
    Get the interface IP given its name.
    """
    iface_name = get_service_config_value(service, iface_config, "")
    return get_ip_from_if_cidr(iface_name)


def _get_dns_ip(iface_config):
    """
    Get dnsd interface IP without netmask.
    If caching is enabled, use the ip of interface that dnsd listens over.
    Otherwise, just use dns server in yml.
    """
    if load_service_mconfig("mme", MME()).enable_dns_caching:
        iface_name = get_service_config_value("dnsd", iface_config, "")
        return get_ip_from_if(iface_name)
    return get_service_config_value("spgw", "ipv4_dns", "")


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


def _get_relay_enabled():
    if load_service_mconfig("mme", MME()).relay_enabled:
        return "yes"
    return "no"


def _get_non_eps_service_control():
    non_eps_service_control = \
        load_service_mconfig("mme", MME()).non_eps_service_control
    if non_eps_service_control:
        if non_eps_service_control == 0:
            return "OFF"
        elif non_eps_service_control == 1:
            return "CSFB_SMS"
        elif non_eps_service_control == 2:
            return "SMS"
    return "OFF"


def _get_lac():
    lac = load_service_mconfig("mme", MME()).lac
    if lac:
        return lac
    return 0


def _get_csfb_mcc():
    csfb_mcc = load_service_mconfig("mme", MME()).csfb_mcc
    if csfb_mcc:
        return csfb_mcc
    return ""


def _get_csfb_mnc():
    csfb_mnc = load_service_mconfig("mme", MME()).csfb_mnc
    if csfb_mnc:
        return csfb_mnc
    return ""


def _get_identity():
    realm = get_service_config_value("mme", "realm", "")
    return "{}.{}".format(socket.gethostname(), realm)


def _get_attached_enodeb_tacs():
    mme_config = load_service_mconfig("mme", MME())
    # attachedEnodebTacs overrides 'tac', which is being deprecated, but for
    # now, both are supported
    tac = mme_config.tac
    attached_enodeb_tacs = mme_config.attached_enodeb_tacs
    if len(attached_enodeb_tacs) == 0:
        return [tac]
    return attached_enodeb_tacs


def _get_context():
    """
    Create the context which has the interface IP and the OAI log level to use.
    """
    context = {}
    context["s11_ip"] = _get_iface_ip("mme", "s11_iface_name")
    context["s11_sgw_ip"] = get_service_config_value("mme", key, "")
    context["s1ap_ip"] = _get_iface_ip("mme", "s1ap_iface_name")
    context["s1u_ip"] = _get_iface_ip("spgw", "s1u_iface_name")
    context["oai_log_level"] = _get_oai_log_level()
    context["ipv4_dns"] = _get_dns_ip("dns_iface_name")
    context["identity"] = _get_identity()
    context["relay_enabled"] = _get_relay_enabled()
    context["non_eps_service_control"] = _get_non_eps_service_control()
    context["csfb_mcc"] = _get_csfb_mcc()
    context["csfb_mnc"] = _get_csfb_mnc()
    context["lac"] = _get_lac()
    context["use_stateless"] = get_service_config_value("mme", "use_stateless", "")
    context["attached_enodeb_tacs"] = _get_attached_enodeb_tacs()
    # set ovs params
    for key in (
        "ovs_bridge_name",
        "ovs_gtp_port_number",
        "ovs_mtr_port_number",
        "ovs_uplink_port_number",
        "ovs_uplink_mac",
    ):
        context[key] = get_service_config_value("spgw", key, "")
    return context


def main():
    logging.basicConfig(
        level=logging.INFO, format="[%(asctime)s %(levelname)s %(name)s] %(message)s"
    )
    context = _get_context()
    generate_template_config("spgw", "spgw", CONFIG_OVERRIDE_DIR, context.copy())
    generate_template_config("mme", "mme", CONFIG_OVERRIDE_DIR, context.copy())
    generate_template_config("mme", "mme_fd", CONFIG_OVERRIDE_DIR, context.copy())
    cert_dir = get_service_config_value("mme", "cert_dir", "")
    generate_mme_certs(os.path.join(cert_dir, "freeDiameter"))


if __name__ == "__main__":
    main()

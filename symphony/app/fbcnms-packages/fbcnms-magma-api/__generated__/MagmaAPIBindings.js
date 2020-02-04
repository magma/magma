/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @generated
 */

export type aaa_server = {
    accounting_enabled ? : boolean,
    create_session_on_auth ? : boolean,
    idle_session_timeout_ms ? : number,
};
export type aggregation_logging_configs = {
    target_files_by_tag ? : {
        [string]: string,
    },
};
export type alert_bulk_upload_response = {
    errors: {
        [string]: string,
    },
    statuses: {
        [string]: string,
    },
};
export type alert_receiver_config = {
    email_configs ? : Array < email_receiver >
        ,
    name: string,
    slack_configs ? : Array < slack_receiver >
        ,
    webhook_configs ? : Array < webhook_receiver >
        ,
};
export type alert_routing_tree = {
    continue ?: boolean,
    group_by ? : Array < string >
        ,
    group_interval ? : string,
    group_wait ? : string,
    match ? : {
        label ? : string,
        value ? : string,
    },
    match_re ? : {
        label ? : string,
        value ? : string,
    },
    receiver: string,
    repeat_interval ? : string,
    routes ? : Array < alert_routing_tree >
        ,
};
export type alert_silence_status = {
    state: string,
};
export type alert_silencer = {
    comment: string,
    createdBy: string,
    endsAt: string,
    matchers: Array < matcher >
        ,
    startsAt: string,
};
export type allowed_gre_peer = {
    ip: string,
    key ? : number,
};
export type allowed_gre_peers = Array < allowed_gre_peer >
;
export type base_name = string;
export type base_name_record = {
    assigned_subscribers ? : Array < subscriber_id >
        ,
    name: base_name,
    rule_names: rule_names,
};
export type base_names = Array < base_name >
;
export type cambium_channel = {
    client_id ? : string,
    client_ip ? : string,
    client_mac ? : string,
    client_secret ? : string,
};
export type challenge_key = {
    key ? : string,
    key_type: "ECHO" | "SOFTWARE_ECDSA_SHA256",
};
export type channel_id = string;
export type config_info = {
    mconfig_created_at ? : number,
};
export type cwf_gateway = {
    carrier_wifi: gateway_cwf_configs,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type cwf_network = {
    carrier_wifi: network_carrier_wifi_configs,
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    federation: federated_network_configs,
    id: network_id,
    name: network_name,
    subscriber_config ? : network_subscriber_config,
};
export type cwf_subscriber_directory_record = {
    ipv4_addr ? : string,
    location_history: Array < string >
        ,
    mac_addr ? : string,
};
export type diameter_client_configs = {
    address ? : string,
    dest_host ? : string,
    dest_realm ? : string,
    disable_dest_host ? : boolean,
    host ? : string,
    local_address ? : string,
    overwrite_dest_host ? : boolean,
    product_name ? : string,
    protocol ? : "tcp" | "tcp4" | "tcp6" | "sctp" | "sctp4" | "sctp6",
    realm ? : string,
    retransmits ? : number,
    retry_count ? : number,
    watchdog_interval ? : number,
};
export type diameter_server_configs = {
    address ? : string,
    dest_host ? : string,
    dest_realm ? : string,
    local_address ? : string,
    protocol ? : "tcp" | "tcp4" | "tcp6" | "sctp" | "sctp4" | "sctp6",
};
export type disk_partition = {
    device ? : string,
    free ? : number,
    mount_point ? : string,
    total ? : number,
    used ? : number,
};
export type dns_config_record = {
    a_record ? : Array < string >
        ,
    aaaa_record ? : Array < string >
        ,
    cname_record ? : Array < string >
        ,
    domain: string,
};
export type eap_aka = {
    plmn_ids ? : Array < string >
        ,
    timeout ? : eap_aka_timeouts,
};
export type eap_aka_timeouts = {
    challenge_ms ? : number,
    error_notification_ms ? : number,
    session_authenticated_ms ? : number,
    session_ms ? : number,
};
export type elastic_hit = {
    _id: string,
    _index: string,
    _primary_term ? : string,
    _score ? : number,
    _seq_no ? : number,
    _sort ? : Array < number >
        ,
    _source: {
        [string]: string,
    },
    _type: string,
};
export type email_receiver = {
    auth_identity ? : string,
    auth_password ? : string,
    auth_secret ? : string,
    auth_username ? : string,
    from: string,
    headers ? : {
        [string]: string,
    },
    hello ? : string,
    html ? : string,
    send_resolved ? : boolean,
    smarthost: string,
    text ? : string,
    to: string,
};
export type enodeb = {
    attached_gateway_id ? : string,
    config: enodeb_configuration,
    name: string,
    serial: string,
};
export type enodeb_configuration = {
    bandwidth_mhz ? : 3 | 5 | 10 | 15 | 20,
    cell_id: number,
    device_class: "Baicells Nova-233 G2 OD FDD" | "Baicells Nova-243 OD TDD" | "Baicells Neutrino 224 ID FDD" | "Baicells ID TDD/FDD" | "NuRAN Cavium OC-LTE",
    earfcndl ? : number,
    pci ? : number,
    special_subframe_pattern ? : number,
    subframe_assignment ? : number,
    tac ? : number,
    transmit_enabled: boolean,
};
export type enodeb_serials = Array < string >
;
export type enodeb_state = {
    enodeb_configured: boolean,
    enodeb_connected: boolean,
    fsm_state: string,
    gps_connected: boolean,
    gps_latitude: string,
    gps_longitude: string,
    mme_connected: boolean,
    opstate_enabled: boolean,
    ptp_connected: boolean,
    reporting_gateway_id ? : string,
    rf_tx_desired: boolean,
    rf_tx_on: boolean,
    time_reported ? : number,
};
export type error = {
    message: string,
};
export type federated_network_configs = {
    feg_network_id: string,
};
export type federation_gateway = {
    description: gateway_description,
    device: gateway_device,
    federation: gateway_federation_configs,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type federation_gateway_health_status = {
    description: string,
    status: "HEALTHY" | "UNHEALTHY",
};
export type federation_network_cluster_status = {
    active_gateway: string,
};
export type feg_lte_network = {
    cellular: network_cellular_configs,
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    federation: federated_network_configs,
    id: network_id,
    name: network_name,
};
export type feg_network = {
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    federation: network_federation_configs,
    id: network_id,
    name: network_name,
    subscriber_config ? : network_subscriber_config,
};
export type feg_network_id = string;
export type flow_description = {
    action: "PERMIT" | "DENY",
    match: flow_match,
};
export type flow_match = {
    direction: "UPLINK" | "DOWNLINK",
    ip_proto: "IPPROTO_IP" | "IPPROTO_TCP" | "IPPROTO_UDP" | "IPPROTO_ICMP",
    ipv4_dst ? : string,
    ipv4_src ? : string,
    tcp_dst ? : number,
    tcp_src ? : number,
    udp_dst ? : number,
    udp_src ? : number,
};
export type flow_qos = {
    max_req_bw_dl: number,
    max_req_bw_ul: number,
};
export type frinx_channel = {
    authorization ? : string,
    device_type ? : string,
    device_version ? : string,
    frinx_port ? : number,
    host ? : string,
    password ? : string,
    port ? : number,
    transport_type ? : string,
    username ? : string,
};
export type gateway_cellular_configs = {
    epc: gateway_epc_configs,
    non_eps_service ? : gateway_non_eps_configs,
    ran: gateway_ran_configs,
};
export type gateway_cwf_configs = {
    allowed_gre_peers: allowed_gre_peers,
};
export type gateway_description = string;
export type gateway_device = {
    hardware_id: string,
    key: challenge_key,
};
export type gateway_epc_configs = {
    ip_block: string,
    nat_enabled: boolean,
};
export type gateway_federation_configs = {
    aaa_server: aaa_server,
    eap_aka: eap_aka,
    gx: gx,
    gy: gy,
    health: health,
    hss: hss,
    s6a: s6a,
    served_network_ids: served_network_ids,
    swx: swx,
};
export type gateway_id = string;
export type gateway_logging_configs = {
    aggregation ? : aggregation_logging_configs,
    log_level: "DEBUG" | "INFO" | "WARNING" | "ERROR" | "FATAL",
};
export type gateway_name = string;
export type gateway_non_eps_configs = {
    arfcn_2g ? : Array < number >
        ,
    csfb_mcc ? : string,
    csfb_mnc ? : string,
    csfb_rat ? : 0 | 1,
    lac ? : number,
    non_eps_service_control: 0 | 1 | 2,
};
export type gateway_ran_configs = {
    pci: number,
    transmit_enabled: boolean,
};
export type gateway_status = {
    cert_expiration_time ? : number,
    checkin_time ? : number,
    hardware_id ? : string,
    kernel_version ? : string,
    kernel_versions_installed ? : Array < string >
        ,
    machine_info ? : machine_info,
    meta ? : {
        [string]: string,
    },
    platform_info ? : platform_info,
    system_status ? : system_status,
    version ? : string,
    vpn_ip ? : string,
};
export type gateway_wifi_configs = {
    additional_props ? : {
        [string]: string,
    },
    client_channel ? : string,
    info ? : string,
    is_production ? : boolean,
    latitude ? : number,
    longitude ? : number,
    mesh_id ? : mesh_id,
    mesh_rssi_threshold ? : number,
    override_password ? : string,
    override_ssid ? : string,
    override_xwf_config ? : string,
    override_xwf_dhcp_dns1 ? : string,
    override_xwf_dhcp_dns2 ? : string,
    override_xwf_enabled ? : boolean,
    override_xwf_partner_name ? : string,
    override_xwf_radius_acct_port ? : number,
    override_xwf_radius_auth_port ? : number,
    override_xwf_radius_server ? : string,
    override_xwf_radius_shared_secret ? : string,
    override_xwf_uam_secret ? : string,
    use_override_ssid ? : boolean,
    use_override_xwf ? : boolean,
    wifi_disabled ? : boolean,
};
export type generic_command_params = {
    command: string,
    params ? : {
        [string]: {},
    },
};
export type generic_command_response = {
    response ? : {
        [string]: {},
    },
};
export type gettable_alert = {
    name: string,
};
export type gettable_alert_silencer = {
    comment: string,
    createdBy: string,
    endsAt: string,
    matchers: Array < matcher >
        ,
    startsAt: string,
    id: string,
    status: alert_silence_status,
    updatedAt: string,
};
export type gx = {
    server ? : diameter_client_configs,
};
export type gy = {
    init_method ? : 1 | 2,
    server ? : diameter_client_configs,
};
export type health = {
    cloud_disable_period_secs ? : number,
    cpu_utilization_threshold ? : number,
    health_services ? : Array < "S6A_PROXY" | "SESSION_PROXY" | "SWX_PROXY" >
        ,
    local_disable_period_secs ? : number,
    memory_available_threshold ? : number,
    minimum_request_threshold ? : number,
    request_failure_threshold ? : number,
    update_failure_threshold ? : number,
    update_interval_secs ? : number,
};
export type hss = {
    default_sub_profile ? : subscription_profile,
    lte_auth_amf ? : string,
    lte_auth_op ? : string,
    server ? : diameter_server_configs,
    stream_subscribers ? : boolean,
    sub_profiles ? : {
        [string]: subscription_profile,
    },
};
export type http_config = {
    basic_auth ? : http_config_basic_auth,
    bearer_token ? : string,
    proxy_url ? : string,
};
export type http_config_basic_auth = {
    password: string,
    username: string,
};
export type label_pair = {
    name: string,
    value: string,
};
export type lte_gateway = {
    cellular: gateway_cellular_configs,
    connected_enodeb_serials: enodeb_serials,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type lte_network = {
    cellular: network_cellular_configs,
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    id: network_id,
    name: network_name,
    subscriber_config ? : network_subscriber_config,
};
export type lte_subscription = {
    auth_algo: "MILENAGE",
    auth_key: string,
    auth_opc ? : string,
    state: "INACTIVE" | "ACTIVE",
    sub_profile: sub_profile,
};
export type machine_info = {
    cpu_info ? : {
        architecture ? : string,
        core_count ? : number,
        model_name ? : string,
        threads_per_core ? : number,
    },
    network_info ? : {
        network_interfaces ? : Array < network_interface >
            ,
        routing_table ? : Array < route >
            ,
    },
};
export type magmad_gateway = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type magmad_gateway_configs = {
    autoupgrade_enabled: boolean,
    autoupgrade_poll_interval: number,
    checkin_interval: number,
    checkin_timeout: number,
    dynamic_services ? : Array < string >
        ,
    feature_flags ? : {
        [string]: boolean,
    },
    logging ? : gateway_logging_configs,
};
export type managed_devices = Array < string >
;
export type matcher = {
    isRegex: boolean,
    name: string,
    value: string,
};
export type mesh_id = string;
export type mesh_name = string;
export type mesh_wifi_configs = {
    additional_props ? : {
        [string]: string,
    },
    mesh_channel_type ? : string,
    mesh_frequency ? : number,
    mesh_ssid ? : string,
    password ? : string,
    ssid ? : string,
    vl_ssid ? : string,
    xwf_enabled ? : boolean,
};
export type metric_datapoint = Array < string >
;
export type metric_datapoints = Array < metric_datapoint >
;
export type mutable_cwf_gateway = {
    carrier_wifi: gateway_cwf_configs,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    tier: tier_id,
};
export type mutable_federation_gateway = {
    description: gateway_description,
    device: gateway_device,
    federation: gateway_federation_configs,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    tier: tier_id,
};
export type mutable_lte_gateway = {
    cellular: gateway_cellular_configs,
    connected_enodeb_serials: enodeb_serials,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    tier: tier_id,
};
export type mutable_rating_group = {
    limit_type: "FINITE" | "INFINITE_UNMETERED" | "INFINITE_METERED",
};
export type mutable_symphony_agent = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    managed_devices: managed_devices,
    name: gateway_name,
    tier: tier_id,
};
export type mutable_symphony_device = {
    config: symphony_device_config,
    id: symphony_device_id,
    managing_agent ? : symphony_device_agent,
    name: symphony_device_name,
};
export type mutable_wifi_gateway = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    tier: tier_id,
    wifi: gateway_wifi_configs,
};
export type network = {
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    id: network_id,
    name: network_name,
    type ? : network_type,
};
export type network_carrier_wifi_configs = {
    aaa_server: aaa_server,
    default_rule_id: string,
    eap_aka: eap_aka,
    network_services: Array < "metering" | "dpi" | "policy_enforcement" >
        ,
};
export type network_cellular_configs = {
    epc: network_epc_configs,
    feg_network_id ? : feg_network_id,
    ran: network_ran_configs,
};
export type network_description = string;
export type network_dns_config = {
    enable_caching: boolean,
    local_ttl: number,
    records ? : network_dns_records,
};
export type network_dns_records = Array < dns_config_record >
;
export type network_epc_configs = {
    cloud_subscriberdb_enabled ? : boolean,
    default_rule_id ? : string,
    lte_auth_amf: string,
    lte_auth_op: string,
    mcc: string,
    mnc: string,
    mobility ? : {
        ip_allocation_mode: "NAT" | "STATIC" | "DHCP_PASSTHROUGH" | "DHCP_BROADCAST",
        nat ? : {
            ip_blocks ? : Array < string >
                ,
        },
        reserved_addresses ? : Array < string >
            ,
        static ? : {
            ip_blocks_by_tac ? : {
                [string]: Array < string >
                    ,
            },
        },
    },
    network_services ? : Array < "metering" | "dpi" | "policy_enforcement" >
        ,
    relay_enabled: boolean,
    sub_profiles ? : {
        [string]: {
            max_dl_bit_rate: number,
            max_ul_bit_rate: number,
        },
    },
    tac: number,
};
export type network_features = {
    features ? : {
        [string]: string,
    },
};
export type network_federation_configs = {
    aaa_server: aaa_server,
    eap_aka: eap_aka,
    gx: gx,
    gy: gy,
    health: health,
    hss: hss,
    s6a: s6a,
    served_network_ids: served_network_ids,
    swx: swx,
};
export type network_id = string;
export type network_interface = {
    ip_addresses ? : Array < string >
        ,
    ipv6_addresses ? : Array < string >
        ,
    mac_address ? : string,
    network_interface_id ? : string,
    status ? : "UP" | "DOWN" | "UNKNOWN",
};
export type network_name = string;
export type network_ran_configs = {
    bandwidth_mhz: 3 | 5 | 10 | 15 | 20,
    fdd_config ? : {
        earfcndl: number,
        earfcnul: number,
    },
    tdd_config ? : {
        earfcndl: number,
        special_subframe_pattern: number,
        subframe_assignment: number,
    },
};
export type network_subscriber_config = {
    network_wide_base_names ? : base_names,
    network_wide_rule_names ? : rule_names,
};
export type network_type = string;
export type network_wifi_configs = {
    additional_props ? : {
        [string]: string,
    },
    mgmt_vpn_enabled ? : boolean,
    mgmt_vpn_proto ? : string,
    mgmt_vpn_remote ? : string,
    openr_enabled ? : boolean,
    ping_host_list ? : Array < string >
        ,
    ping_num_packets ? : number,
    ping_timeout_secs ? : number,
    vl_auth_server_addr ? : string,
    vl_auth_server_port ? : number,
    vl_auth_server_shared_secret ? : string,
    xwf_config ? : string,
    xwf_dhcp_dns1 ? : string,
    xwf_dhcp_dns2 ? : string,
    xwf_partner_name ? : string,
    xwf_radius_acct_port ? : number,
    xwf_radius_auth_port ? : number,
    xwf_radius_server ? : string,
    xwf_radius_shared_secret ? : string,
    xwf_uam_secret ? : string,
};
export type other_channel = {
    channel_props ? : {
        [string]: string,
    },
};
export type package_type = {
    name ? : string,
    version ? : string,
};
export type ping_request = {
    hosts: Array < string >
        ,
    packets ? : number,
};
export type ping_response = {
    pings: Array < ping_result >
        ,
};
export type ping_result = {
    avg_response_ms ? : number,
    error ? : string,
    host_or_ip: string,
    num_packets: number,
    packets_received ? : number,
    packets_transmitted ? : number,
};
export type platform_info = {
    config_info ? : config_info,
    kernel_version ? : string,
    kernel_versions_installed ? : Array < string >
        ,
    packages ? : Array < package_type >
        ,
    vpn_ip ? : string,
};
export type policy_id = string;
export type policy_rule = {
    assigned_subscribers ? : Array < subscriber_id >
        ,
    flow_list: Array < flow_description >
        ,
    id: policy_id,
    monitoring_key ? : string,
    priority: number,
    qos ? : flow_qos,
    rating_group ? : number,
    redirect ? : redirect_information,
    tracking_type ? : "ONLY_OCS" | "ONLY_PCRF" | "OCS_AND_PCRF" | "NO_TRACKING",
};
export type policy_rule_config = {
    flow_list: Array < flow_description >
        ,
    monitoring_key ? : string,
    priority: number,
    qos ? : flow_qos,
    rating_group ? : number,
    redirect ? : redirect_information,
    tracking_type ? : "ONLY_OCS" | "ONLY_PCRF" | "OCS_AND_PCRF" | "NO_TRACKING",
};
export type prom_alert_config = {
    alert: string,
    annotations ? : prom_alert_labels,
    expr: string,
    for ? : string,
    labels ? : prom_alert_labels,
};
export type prom_alert_config_list = Array < prom_alert_config >
;
export type prom_alert_labels = {
    [string]: string,
};
export type prom_alert_status = {
    inhibitedBy: Array < string >
        ,
    silencedBy: Array < string >
        ,
    state: string,
};
export type prom_firing_alert = {
    annotations: prom_alert_labels,
    endsAt: string,
    fingerprint: string,
    generatorURL ? : string,
    labels: prom_alert_labels,
    receivers: gettable_alert,
    startsAt: string,
    status: prom_alert_status,
    updatedAt: string,
};
export type prometheus_labelset = {
    [string]: string,
};
export type promql_data = {
    result: promql_result,
    resultType: string,
};
export type promql_metric = {
    additionalProperties ? : string,
};
export type promql_metric_value = {
    metric: promql_metric,
    value ? : metric_datapoint,
    values ? : metric_datapoints,
};
export type promql_result = Array < promql_metric_value >
;
export type promql_return_object = {
    data: promql_data,
    status: string,
};
export type pushed_metric = {
    labels ? : Array < label_pair >
        ,
    metricName: string,
    timestamp ? : string,
    value: number,
};
export type rating_group = {
    id: rating_group_id,
    limit_type: "FINITE" | "INFINITE_UNMETERED" | "INFINITE_METERED",
};
export type rating_group_id = number;
export type redirect_information = {
    address_type: "IPv4" | "IPv6" | "URL" | "SIP_URI",
    server_address: string,
    support: "DISABLED" | "ENABLED",
};
export type release_channel = {
    id: channel_id,
    name ? : string,
    supported_versions: Array < string >
        ,
};
export type route = {
    destination_ip ? : string,
    gateway_ip ? : string,
    genmask ? : string,
    network_interface_id ? : string,
};
export type rule_id = string;
export type rule_names = Array < string >
;
export type s6a = {
    server ? : diameter_client_configs,
};
export type served_network_ids = Array < string >
;
export type slack_action = {
    confirm ? : slack_confirm_field,
    name ? : string,
    style ? : string,
    text: string,
    type: string,
    url: string,
    value ? : string,
};
export type slack_confirm_field = {
    dismiss_text: string,
    ok_text: string,
    text: string,
    title: string,
};
export type slack_field = {
    short ? : boolean,
    title: string,
    value: string,
};
export type slack_receiver = {
    actions ? : Array < slack_action >
        ,
    api_url: string,
    callback_id ? : string,
    channel ? : string,
    color ? : string,
    fallback ? : string,
    fields ? : Array < slack_field >
        ,
    footer ? : string,
    icon_emoji ? : string,
    icon_url ? : string,
    image_url ? : string,
    link_names ? : boolean,
    pretext ? : string,
    short_fields ? : boolean,
    text ? : string,
    thumb_url ? : string,
    title ? : string,
    username ? : string,
};
export type snmp_channel = {
    community ? : string,
    version ? : string,
};
export type sub_profile = string;
export type subscriber = {
    active_base_names ? : Array < base_name >
        ,
    active_policies ? : Array < policy_id >
        ,
    id: subscriber_id,
    lte: lte_subscription,
};
export type subscriber_id = string;
export type subscription_profile = {
    max_dl_bit_rate ? : number,
    max_ul_bit_rate ? : number,
};
export type swx = {
    cache_TTL_seconds ? : number,
    derive_unregister_realm ? : boolean,
    hlr_plmn_ids ? : Array < string >
        ,
    register_on_auth ? : boolean,
    server ? : diameter_client_configs,
    verify_authorization ? : boolean,
};
export type symphony_agent = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    managed_devices: managed_devices,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type symphony_device = {
    config: symphony_device_config,
    id: symphony_device_id,
    managing_agent: symphony_device_agent,
    name: symphony_device_name,
    state: symphony_device_state,
};
export type symphony_device_agent = string;
export type symphony_device_config = {
    channels ? : {
        cambium_channel ? : cambium_channel,
        frinx_channel ? : frinx_channel,
        other_channel ? : other_channel,
        snmp_channel ? : snmp_channel,
    },
    device_config ? : string,
    device_type ? : Array < string >
        ,
    host ? : string,
    platform ? : string,
};
export type symphony_device_id = string;
export type symphony_device_name = string;
export type symphony_device_state = {
    raw_state ? : string,
};
export type symphony_network = {
    description: network_description,
    features ? : network_features,
    id: network_id,
    name: network_name,
};
export type system_status = {
    cpu_idle ? : number,
    cpu_system ? : number,
    cpu_user ? : number,
    disk_partitions ? : Array < disk_partition >
        ,
    mem_available ? : number,
    mem_free ? : number,
    mem_total ? : number,
    mem_used ? : number,
    swap_free ? : number,
    swap_total ? : number,
    swap_used ? : number,
    time ? : number,
    uptime_secs ? : number,
};
export type tail_logs_request = {
    service ? : string,
};
export type tier = {
    gateways: tier_gateways,
    id: tier_id,
    images: tier_images,
    name ? : tier_name,
    version: tier_version,
};
export type tier_gateways = Array < gateway_id >
;
export type tier_id = string;
export type tier_image = {
    name: string,
    order: number,
};
export type tier_images = Array < tier_image >
;
export type tier_name = string;
export type tier_version = string;
export type webhook_receiver = {
    http_config ? : http_config,
    send_resolved ? : boolean,
    url: string,
};
export type wifi_gateway = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
    wifi: gateway_wifi_configs,
};
export type wifi_mesh = {
    config: mesh_wifi_configs,
    gateway_ids: Array < gateway_id >
        ,
    id: mesh_id,
    name: mesh_name,
};
export type wifi_network = {
    description: network_description,
    features ? : network_features,
    id: network_id,
    name: network_name,
    wifi: network_wifi_configs,
};

export default class MagmaAPIBindings {
    static request(
        path: string,
        method: 'POST' | 'GET' | 'PUT' | 'DELETE' | 'OPTIONS' | 'HEAD' | 'PATCH',
        query: {
            [string]: mixed
        },
        body ? : {
            [string]: any
        } | string | Array < any > ,
    ) {
        throw new Error("Must be implemented");
    }
    static async getChannels(): Promise < Array < channel_id >
        >
        {
            let path = '/channels';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postChannels(
        parameters: {
            'channel': release_channel,
        }
    ): Promise < "Success" > {
        let path = '/channels';
        let body;
        let query = {};
        if (parameters['channel'] === undefined) {
            throw new Error('Missing required  parameter: channel');
        }

        if (parameters['channel'] !== undefined) {
            body = parameters['channel'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteChannelsByChannelId(
        parameters: {
            'channelId': string,
        }
    ): Promise < "Success" > {
        let path = '/channels/{channel_id}';
        let body;
        let query = {};
        if (parameters['channelId'] === undefined) {
            throw new Error('Missing required  parameter: channelId');
        }

        path = path.replace('{channel_id}', `${parameters['channelId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getChannelsByChannelId(
            parameters: {
                'channelId': string,
            }
        ): Promise < release_channel >
        {
            let path = '/channels/{channel_id}';
            let body;
            let query = {};
            if (parameters['channelId'] === undefined) {
                throw new Error('Missing required  parameter: channelId');
            }

            path = path.replace('{channel_id}', `${parameters['channelId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putChannelsByChannelId(
        parameters: {
            'channelId': string,
            'releaseChannel': release_channel,
        }
    ): Promise < "Success" > {
        let path = '/channels/{channel_id}';
        let body;
        let query = {};
        if (parameters['channelId'] === undefined) {
            throw new Error('Missing required  parameter: channelId');
        }

        path = path.replace('{channel_id}', `${parameters['channelId']}`);

        if (parameters['releaseChannel'] === undefined) {
            throw new Error('Missing required  parameter: releaseChannel');
        }

        if (parameters['releaseChannel'] !== undefined) {
            body = parameters['releaseChannel'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwf(): Promise < Array < string >
        >
        {
            let path = '/cwf';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postCwf(
        parameters: {
            'cwfNetwork': cwf_network,
        }
    ): Promise < "Success" > {
        let path = '/cwf';
        let body;
        let query = {};
        if (parameters['cwfNetwork'] === undefined) {
            throw new Error('Missing required  parameter: cwfNetwork');
        }

        if (parameters['cwfNetwork'] !== undefined) {
            body = parameters['cwfNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteCwfByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getCwfByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < cwf_network >
        {
            let path = '/cwf/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkId(
        parameters: {
            'networkId': string,
            'cwfNetwork': cwf_network,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['cwfNetwork'] === undefined) {
            throw new Error('Missing required  parameter: cwfNetwork');
        }

        if (parameters['cwfNetwork'] !== undefined) {
            body = parameters['cwfNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteCwfByNetworkIdCarrierWifi(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/carrier_wifi';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getCwfByNetworkIdCarrierWifi(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_carrier_wifi_configs >
        {
            let path = '/cwf/{network_id}/carrier_wifi';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdCarrierWifi(
        parameters: {
            'networkId': string,
            'config': network_carrier_wifi_configs,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/carrier_wifi';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/cwf/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: cwf_gateway,
        } >
        {
            let path = '/cwf/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postCwfByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': mutable_cwf_gateway,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteCwfByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < cwf_gateway >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': mutable_cwf_gateway,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdCarrierWifi(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_cwf_configs >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/carrier_wifi';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdCarrierWifi(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_cwf_configs,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/carrier_wifi';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getCwfByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/cwf/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/cwf/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdSubscriberConfig(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_subscriber_config >
        {
            let path = '/cwf/{network_id}/subscriber_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdSubscriberConfig(
        parameters: {
            'networkId': string,
            'record': network_subscriber_config,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getCwfByNetworkIdSubscriberConfigBaseNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < base_names >
        {
            let path = '/cwf/{network_id}/subscriber_config/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdSubscriberConfigBaseNames(
        parameters: {
            'networkId': string,
            'record': base_names,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/base_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteCwfByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postCwfByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getCwfByNetworkIdSubscriberConfigRuleNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < rule_names >
        {
            let path = '/cwf/{network_id}/subscriber_config/rule_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putCwfByNetworkIdSubscriberConfigRuleNames(
        parameters: {
            'networkId': string,
            'record': rule_names,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/rule_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteCwfByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postCwfByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/cwf/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getCwfByNetworkIdSubscribersBySubscriberIdDirectoryRecord(
            parameters: {
                'networkId': string,
                'subscriberId': string,
            }
        ): Promise < cwf_subscriber_directory_record >
        {
            let path = '/cwf/{network_id}/subscribers/{subscriber_id}/directory_record';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['subscriberId'] === undefined) {
                throw new Error('Missing required  parameter: subscriberId');
            }

            path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getFeg(): Promise < Array < string >
        >
        {
            let path = '/feg';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postFeg(
        parameters: {
            'fegNetwork': feg_network,
        }
    ): Promise < "Success" > {
        let path = '/feg';
        let body;
        let query = {};
        if (parameters['fegNetwork'] === undefined) {
            throw new Error('Missing required  parameter: fegNetwork');
        }

        if (parameters['fegNetwork'] !== undefined) {
            body = parameters['fegNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteFegByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < feg_network >
        {
            let path = '/feg/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkId(
        parameters: {
            'networkId': string,
            'fegNetwork': feg_network,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['fegNetwork'] === undefined) {
            throw new Error('Missing required  parameter: fegNetwork');
        }

        if (parameters['fegNetwork'] !== undefined) {
            body = parameters['fegNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegByNetworkIdClusterStatus(
            parameters: {
                'networkId': string,
            }
        ): Promise < federation_network_cluster_status >
        {
            let path = '/feg/{network_id}/cluster_status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async deleteFegByNetworkIdFederation(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegByNetworkIdFederation(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_federation_configs >
        {
            let path = '/feg/{network_id}/federation';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkIdFederation(
        parameters: {
            'networkId': string,
            'config': network_federation_configs,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: federation_gateway,
        } >
        {
            let path = '/feg/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postFegByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': mutable_federation_gateway,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteFegByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < federation_gateway >
        {
            let path = '/feg/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': mutable_federation_gateway,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegByNetworkIdGatewaysByGatewayIdFederation(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways/{gateway_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegByNetworkIdGatewaysByGatewayIdFederation(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_federation_configs >
        {
            let path = '/feg/{network_id}/gateways/{gateway_id}/federation';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postFegByNetworkIdGatewaysByGatewayIdFederation(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_federation_configs,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways/{gateway_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putFegByNetworkIdGatewaysByGatewayIdFederation(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_federation_configs,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/gateways/{gateway_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegByNetworkIdGatewaysByGatewayIdHealthStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < federation_gateway_health_status >
        {
            let path = '/feg/{network_id}/gateways/{gateway_id}/health_status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getFegByNetworkIdSubscriberConfig(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_subscriber_config >
        {
            let path = '/feg/{network_id}/subscriber_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkIdSubscriberConfig(
        parameters: {
            'networkId': string,
            'record': network_subscriber_config,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegByNetworkIdSubscriberConfigBaseNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < base_names >
        {
            let path = '/feg/{network_id}/subscriber_config/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkIdSubscriberConfigBaseNames(
        parameters: {
            'networkId': string,
            'record': base_names,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/base_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postFegByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getFegByNetworkIdSubscriberConfigRuleNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < rule_names >
        {
            let path = '/feg/{network_id}/subscriber_config/rule_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegByNetworkIdSubscriberConfigRuleNames(
        parameters: {
            'networkId': string,
            'record': rule_names,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/rule_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postFegByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getFegLte(): Promise < Array < string >
        >
        {
            let path = '/feg_lte';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postFegLte(
        parameters: {
            'lteNetwork': feg_lte_network,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte';
        let body;
        let query = {};
        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteFegLteByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegLteByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < feg_lte_network >
        {
            let path = '/feg_lte/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegLteByNetworkId(
        parameters: {
            'networkId': string,
            'lteNetwork': feg_lte_network,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegLteByNetworkIdFederation(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getFegLteByNetworkIdFederation(
            parameters: {
                'networkId': string,
            }
        ): Promise < federated_network_configs >
        {
            let path = '/feg_lte/{network_id}/federation';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegLteByNetworkIdFederation(
        parameters: {
            'networkId': string,
            'config': federated_network_configs,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/federation';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegLteByNetworkIdSubscriberConfig(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_subscriber_config >
        {
            let path = '/feg_lte/{network_id}/subscriber_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegLteByNetworkIdSubscriberConfig(
        parameters: {
            'networkId': string,
            'record': network_subscriber_config,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFegLteByNetworkIdSubscriberConfigBaseNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < base_names >
        {
            let path = '/feg_lte/{network_id}/subscriber_config/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegLteByNetworkIdSubscriberConfigBaseNames(
        parameters: {
            'networkId': string,
            'record': base_names,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/base_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegLteByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postFegLteByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getFegLteByNetworkIdSubscriberConfigRuleNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < rule_names >
        {
            let path = '/feg_lte/{network_id}/subscriber_config/rule_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putFegLteByNetworkIdSubscriberConfigRuleNames(
        parameters: {
            'networkId': string,
            'record': rule_names,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/rule_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteFegLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postFegLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/feg_lte/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getFoo(): Promise < number >
        {
            let path = '/foo';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async getLte(): Promise < Array < string >
        >
        {
            let path = '/lte';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postLte(
        parameters: {
            'lteNetwork': lte_network,
        }
    ): Promise < "Success" > {
        let path = '/lte';
        let body;
        let query = {};
        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < lte_network >
        {
            let path = '/lte/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkId(
        parameters: {
            'networkId': string,
            'lteNetwork': lte_network,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellular(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_cellular_configs >
        {
            let path = '/lte/{network_id}/cellular';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellular(
        parameters: {
            'networkId': string,
            'config': network_cellular_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularEpc(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_epc_configs >
        {
            let path = '/lte/{network_id}/cellular/epc';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularEpc(
        parameters: {
            'networkId': string,
            'config': network_epc_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/epc';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularFegNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < string >
        {
            let path = '/lte/{network_id}/cellular/feg_network_id';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularFegNetworkId(
        parameters: {
            'networkId': string,
            'fegNetworkId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/feg_network_id';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['fegNetworkId'] === undefined) {
            throw new Error('Missing required  parameter: fegNetworkId');
        }

        if (parameters['fegNetworkId'] !== undefined) {
            body = parameters['fegNetworkId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularRan(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_ran_configs >
        {
            let path = '/lte/{network_id}/cellular/ran';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularRan(
        parameters: {
            'networkId': string,
            'config': network_ran_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/ran';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/lte/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDns(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_config >
        {
            let path = '/lte/{network_id}/dns';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDns(
        parameters: {
            'networkId': string,
            'config': network_dns_config,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDnsRecords(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < dns_config_record >
        >
        {
            let path = '/lte/{network_id}/dns/records';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDnsRecords(
        parameters: {
            'networkId': string,
            'records': Array < dns_config_record >
                ,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['records'] === undefined) {
            throw new Error('Missing required  parameter: records');
        }

        if (parameters['records'] !== undefined) {
            body = parameters['records'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdDnsRecordsByDomain(
            parameters: {
                'networkId': string,
                'domain': string,
            }
        ): Promise < dns_config_record >
        {
            let path = '/lte/{network_id}/dns/records/{domain}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['domain'] === undefined) {
                throw new Error('Missing required  parameter: domain');
            }

            path = path.replace('{domain}', `${parameters['domain']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdEnodebs(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: enodeb,
        } >
        {
            let path = '/lte/{network_id}/enodebs';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdEnodebs(
        parameters: {
            'networkId': string,
            'enodeb': enodeb,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodeb'] === undefined) {
            throw new Error('Missing required  parameter: enodeb');
        }

        if (parameters['enodeb'] !== undefined) {
            body = parameters['enodeb'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdEnodebsByEnodebSerial(
        parameters: {
            'networkId': string,
            'enodebSerial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodebSerial'] === undefined) {
            throw new Error('Missing required  parameter: enodebSerial');
        }

        path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdEnodebsByEnodebSerial(
            parameters: {
                'networkId': string,
                'enodebSerial': string,
            }
        ): Promise < enodeb >
        {
            let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['enodebSerial'] === undefined) {
                throw new Error('Missing required  parameter: enodebSerial');
            }

            path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdEnodebsByEnodebSerial(
        parameters: {
            'networkId': string,
            'enodebSerial': string,
            'enodeb': enodeb,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodebSerial'] === undefined) {
            throw new Error('Missing required  parameter: enodebSerial');
        }

        path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

        if (parameters['enodeb'] === undefined) {
            throw new Error('Missing required  parameter: enodeb');
        }

        if (parameters['enodeb'] !== undefined) {
            body = parameters['enodeb'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdEnodebsByEnodebSerialState(
            parameters: {
                'networkId': string,
                'enodebSerial': string,
            }
        ): Promise < enodeb_state >
        {
            let path = '/lte/{network_id}/enodebs/{enodeb_serial}/state';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['enodebSerial'] === undefined) {
                throw new Error('Missing required  parameter: enodebSerial');
            }

            path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getLteByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/lte/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'config': network_features,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: lte_gateway,
        } >
        {
            let path = '/lte/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': mutable_lte_gateway,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < lte_gateway >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': mutable_lte_gateway,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellular(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_cellular_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellular(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_cellular_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularEpc(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_epc_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/epc';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularEpc(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_epc_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/epc';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularNonEps(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_non_eps_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/non_eps';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularNonEps(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_non_eps_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/non_eps';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularRan(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_ran_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/ran';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularRan(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_ran_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/ran';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serial'] === undefined) {
            throw new Error('Missing required  parameter: serial');
        }

        if (parameters['serial'] !== undefined) {
            body = parameters['serial'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < enodeb_serials >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serial'] === undefined) {
            throw new Error('Missing required  parameter: serial');
        }

        if (parameters['serial'] !== undefined) {
            body = parameters['serial'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serials': enodeb_serials,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serials'] === undefined) {
            throw new Error('Missing required  parameter: serials');
        }

        if (parameters['serials'] !== undefined) {
            body = parameters['serials'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getLteByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/lte/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdSubscriberConfig(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_subscriber_config >
        {
            let path = '/lte/{network_id}/subscriber_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdSubscriberConfig(
        parameters: {
            'networkId': string,
            'record': network_subscriber_config,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdSubscriberConfigBaseNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < base_names >
        {
            let path = '/lte/{network_id}/subscriber_config/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdSubscriberConfigBaseNames(
        parameters: {
            'networkId': string,
            'record': base_names,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/base_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postLteByNetworkIdSubscriberConfigBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getLteByNetworkIdSubscriberConfigRuleNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < rule_names >
        {
            let path = '/lte/{network_id}/subscriber_config/rule_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdSubscriberConfigRuleNames(
        parameters: {
            'networkId': string,
            'record': rule_names,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/rule_names';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async postLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscriber_config/rule_names/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getLteByNetworkIdSubscribers(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: subscriber,
        } >
        {
            let path = '/lte/{network_id}/subscribers';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdSubscribers(
        parameters: {
            'networkId': string,
            'subscriber': subscriber,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriber'] === undefined) {
            throw new Error('Missing required  parameter: subscriber');
        }

        if (parameters['subscriber'] !== undefined) {
            body = parameters['subscriber'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdSubscribersBySubscriberId(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdSubscribersBySubscriberId(
            parameters: {
                'networkId': string,
                'subscriberId': string,
            }
        ): Promise < subscriber >
        {
            let path = '/lte/{network_id}/subscribers/{subscriber_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['subscriberId'] === undefined) {
                throw new Error('Missing required  parameter: subscriberId');
            }

            path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdSubscribersBySubscriberId(
        parameters: {
            'networkId': string,
            'subscriberId': string,
            'subscriber': subscriber,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        if (parameters['subscriber'] === undefined) {
            throw new Error('Missing required  parameter: subscriber');
        }

        if (parameters['subscriber'] !== undefined) {
            body = parameters['subscriber'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async postLteByNetworkIdSubscribersBySubscriberIdActivate(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}/activate';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async postLteByNetworkIdSubscribersBySubscriberIdDeactivate(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}/deactivate';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async putLteByNetworkIdSubscribersBySubscriberIdLteSubProfile(
        parameters: {
            'networkId': string,
            'subscriberId': string,
            'profileName': sub_profile,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}/lte/sub_profile';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        if (parameters['profileName'] === undefined) {
            throw new Error('Missing required  parameter: profileName');
        }

        if (parameters['profileName'] !== undefined) {
            body = parameters['profileName'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworks(): Promise < Array < string >
        >
        {
            let path = '/networks';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworks(
        parameters: {
            'network': network,
        }
    ): Promise < "Success" > {
        let path = '/networks';
        let body;
        let query = {};
        if (parameters['network'] === undefined) {
            throw new Error('Missing required  parameter: network');
        }

        if (parameters['network'] !== undefined) {
            body = parameters['network'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < network >
        {
            let path = '/networks/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkId(
        parameters: {
            'networkId': string,
            'network': network,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['network'] === undefined) {
            throw new Error('Missing required  parameter: network');
        }

        if (parameters['network'] !== undefined) {
            body = parameters['network'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdAlerts(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < prom_firing_alert >
        >
        {
            let path = '/networks/{network_id}/alerts';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async deleteNetworksByNetworkIdAlertsSilence(
            parameters: {
                'networkId': string,
                'silenceId': string,
            }
        ): Promise < string >
        {
            let path = '/networks/{network_id}/alerts/silence';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['silenceId'] === undefined) {
                throw new Error('Missing required  parameter: silenceId');
            }

            if (parameters['silenceId'] !== undefined) {
                query['silence_id'] = parameters['silenceId'];
            }

            return await this.request(path, 'DELETE', query, body);
        }
    static async getNetworksByNetworkIdAlertsSilence(
            parameters: {
                'networkId': string,
                'active' ? : boolean,
                'pending' ? : boolean,
                'expired' ? : boolean,
                'filter' ? : string,
            }
        ): Promise < Array < gettable_alert_silencer >
        >
        {
            let path = '/networks/{network_id}/alerts/silence';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['active'] !== undefined) {
                query['active'] = parameters['active'];
            }

            if (parameters['pending'] !== undefined) {
                query['pending'] = parameters['pending'];
            }

            if (parameters['expired'] !== undefined) {
                query['expired'] = parameters['expired'];
            }

            if (parameters['filter'] !== undefined) {
                query['filter'] = parameters['filter'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdAlertsSilence(
            parameters: {
                'networkId': string,
                'silencer': alert_silencer,
            }
        ): Promise < string >
        {
            let path = '/networks/{network_id}/alerts/silence';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['silencer'] === undefined) {
                throw new Error('Missing required  parameter: silencer');
            }

            if (parameters['silencer'] !== undefined) {
                body = parameters['silencer'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async getNetworksByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/networks/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdDns(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_config >
        {
            let path = '/networks/{network_id}/dns';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDns(
        parameters: {
            'networkId': string,
            'networkDns': network_dns_config,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['networkDns'] === undefined) {
            throw new Error('Missing required  parameter: networkDns');
        }

        if (parameters['networkDns'] !== undefined) {
            body = parameters['networkDns'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdDnsRecords(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_records >
        {
            let path = '/networks/{network_id}/dns/records';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDnsRecords(
        parameters: {
            'networkId': string,
            'records': network_dns_records,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['records'] === undefined) {
            throw new Error('Missing required  parameter: records');
        }

        if (parameters['records'] !== undefined) {
            body = parameters['records'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdDnsRecordsByDomain(
            parameters: {
                'networkId': string,
                'domain': string,
            }
        ): Promise < dns_config_record >
        {
            let path = '/networks/{network_id}/dns/records/{domain}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['domain'] === undefined) {
                throw new Error('Missing required  parameter: domain');
            }

            path = path.replace('{domain}', `${parameters['domain']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/networks/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'networkFeatures': network_features,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['networkFeatures'] === undefined) {
            throw new Error('Missing required  parameter: networkFeatures');
        }

        if (parameters['networkFeatures'] !== undefined) {
            body = parameters['networkFeatures'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: magmad_gateway,
        } >
        {
            let path = '/networks/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': magmad_gateway,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': magmad_gateway,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric(
            parameters: {
                'networkId': string,
                'gatewayId': string,
                'parameters': generic_command_params,
            }
        ): Promise < generic_command_response >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/command/generic';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            if (parameters['parameters'] === undefined) {
                throw new Error('Missing required  parameter: parameters');
            }

            if (parameters['parameters'] !== undefined) {
                body = parameters['parameters'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandPing(
            parameters: {
                'networkId': string,
                'gatewayId': string,
                'pingRequest': ping_request,
            }
        ): Promise < ping_response >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/command/ping';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            if (parameters['pingRequest'] === undefined) {
                throw new Error('Missing required  parameter: pingRequest');
            }

            if (parameters['pingRequest'] !== undefined) {
                body = parameters['pingRequest'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/command/reboot';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'services': Array < string >
                ,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/command/restart_services';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['services'] === undefined) {
            throw new Error('Missing required  parameter: services');
        }

        if (parameters['services'] !== undefined) {
            body = parameters['services'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdLogs(
            parameters: {
                'networkId': string,
                'simpleQuery' ? : string,
                'fields' ? : string,
                'filters' ? : string,
                'size' ? : string,
                'start' ? : string,
                'end' ? : string,
            }
        ): Promise < Array < elastic_hit >
        >
        {
            let path = '/networks/{network_id}/logs';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['simpleQuery'] !== undefined) {
                query['simple_query'] = parameters['simpleQuery'];
            }

            if (parameters['fields'] !== undefined) {
                query['fields'] = parameters['fields'];
            }

            if (parameters['filters'] !== undefined) {
                query['filters'] = parameters['filters'];
            }

            if (parameters['size'] !== undefined) {
                query['size'] = parameters['size'];
            }

            if (parameters['start'] !== undefined) {
                query['start'] = parameters['start'];
            }

            if (parameters['end'] !== undefined) {
                query['end'] = parameters['end'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdMetricsPush(
        parameters: {
            'networkId': string,
            'metrics': Array < pushed_metric >
                ,
        }
    ): Promise < "Submitted" > {
        let path = '/networks/{network_id}/metrics/push';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['metrics'] === undefined) {
            throw new Error('Missing required  parameter: metrics');
        }

        if (parameters['metrics'] !== undefined) {
            body = parameters['metrics'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/networks/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdPoliciesBaseNames(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < base_name >
        >
        {
            let path = '/networks/{network_id}/policies/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPoliciesBaseNames(
            parameters: {
                'networkId': string,
                'baseNameRecord': base_name_record,
            }
        ): Promise < base_name >
        {
            let path = '/networks/{network_id}/policies/base_names';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['baseNameRecord'] === undefined) {
                throw new Error('Missing required  parameter: baseNameRecord');
            }

            if (parameters['baseNameRecord'] !== undefined) {
                body = parameters['baseNameRecord'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async deleteNetworksByNetworkIdPoliciesBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/policies/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPoliciesBaseNamesByBaseName(
            parameters: {
                'networkId': string,
                'baseName': string,
            }
        ): Promise < base_name_record >
        {
            let path = '/networks/{network_id}/policies/base_names/{base_name}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['baseName'] === undefined) {
                throw new Error('Missing required  parameter: baseName');
            }

            path = path.replace('{base_name}', `${parameters['baseName']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdPoliciesBaseNamesByBaseName(
        parameters: {
            'networkId': string,
            'baseName': string,
            'baseNameRecord': base_name_record,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/policies/base_names/{base_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['baseName'] === undefined) {
            throw new Error('Missing required  parameter: baseName');
        }

        path = path.replace('{base_name}', `${parameters['baseName']}`);

        if (parameters['baseNameRecord'] === undefined) {
            throw new Error('Missing required  parameter: baseNameRecord');
        }

        if (parameters['baseNameRecord'] !== undefined) {
            body = parameters['baseNameRecord'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdPoliciesBaseNamesViewFull(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: base_name_record,
        } >
        {
            let path = '/networks/{network_id}/policies/base_names?view=full';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdPoliciesRules(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < rule_id >
        >
        {
            let path = '/networks/{network_id}/policies/rules';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPoliciesRules(
            parameters: {
                'networkId': string,
                'policyRule': policy_rule,
            }
        ): Promise < rule_id >
        {
            let path = '/networks/{network_id}/policies/rules';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['policyRule'] === undefined) {
                throw new Error('Missing required  parameter: policyRule');
            }

            if (parameters['policyRule'] !== undefined) {
                body = parameters['policyRule'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async deleteNetworksByNetworkIdPoliciesRulesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/policies/rules/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPoliciesRulesByRuleId(
            parameters: {
                'networkId': string,
                'ruleId': string,
            }
        ): Promise < policy_rule >
        {
            let path = '/networks/{network_id}/policies/rules/{rule_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['ruleId'] === undefined) {
                throw new Error('Missing required  parameter: ruleId');
            }

            path = path.replace('{rule_id}', `${parameters['ruleId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdPoliciesRulesByRuleId(
        parameters: {
            'networkId': string,
            'ruleId': string,
            'policyRule': policy_rule,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/policies/rules/{rule_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ruleId'] === undefined) {
            throw new Error('Missing required  parameter: ruleId');
        }

        path = path.replace('{rule_id}', `${parameters['ruleId']}`);

        if (parameters['policyRule'] === undefined) {
            throw new Error('Missing required  parameter: policyRule');
        }

        if (parameters['policyRule'] !== undefined) {
            body = parameters['policyRule'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdPoliciesRulesViewFull(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: policy_rule,
        } >
        {
            let path = '/networks/{network_id}/policies/rules?view=full';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async deleteNetworksByNetworkIdPrometheusAlertConfig(
        parameters: {
            'networkId': string,
            'alertName': string,
        }
    ): Promise < "Deleted" > {
        let path = '/networks/{network_id}/prometheus/alert_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertName'] === undefined) {
            throw new Error('Missing required  parameter: alertName');
        }

        if (parameters['alertName'] !== undefined) {
            query['alert_name'] = parameters['alertName'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertConfig(
            parameters: {
                'networkId': string,
                'alertName' ? : string,
            }
        ): Promise < prom_alert_config_list >
        {
            let path = '/networks/{network_id}/prometheus/alert_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['alertName'] !== undefined) {
                query['alert_name'] = parameters['alertName'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertConfig(
        parameters: {
            'networkId': string,
            'alertConfig': prom_alert_config,
        }
    ): Promise < "Created" > {
        let path = '/networks/{network_id}/prometheus/alert_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertConfig'] === undefined) {
            throw new Error('Missing required  parameter: alertConfig');
        }

        if (parameters['alertConfig'] !== undefined) {
            body = parameters['alertConfig'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertConfigByAlertName(
        parameters: {
            'networkId': string,
            'alertName': string,
            'alertConfig': prom_alert_config,
        }
    ): Promise < "Updated" > {
        let path = '/networks/{network_id}/prometheus/alert_config/{alert_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertName'] === undefined) {
            throw new Error('Missing required  parameter: alertName');
        }

        path = path.replace('{alert_name}', `${parameters['alertName']}`);

        if (parameters['alertConfig'] === undefined) {
            throw new Error('Missing required  parameter: alertConfig');
        }

        if (parameters['alertConfig'] !== undefined) {
            body = parameters['alertConfig'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertConfigBulk(
            parameters: {
                'networkId': string,
                'alertConfigs': prom_alert_config_list,
            }
        ): Promise < alert_bulk_upload_response >
        {
            let path = '/networks/{network_id}/prometheus/alert_config/bulk';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['alertConfigs'] === undefined) {
                throw new Error('Missing required  parameter: alertConfigs');
            }

            if (parameters['alertConfigs'] !== undefined) {
                body = parameters['alertConfigs'];
            }

            return await this.request(path, 'PUT', query, body);
        }
    static async deleteNetworksByNetworkIdPrometheusAlertReceiver(
        parameters: {
            'networkId': string,
            'receiver': string,
        }
    ): Promise < "Deleted" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiver'] === undefined) {
            throw new Error('Missing required  parameter: receiver');
        }

        if (parameters['receiver'] !== undefined) {
            query['receiver'] = parameters['receiver'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertReceiver(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < alert_receiver_config >
        >
        {
            let path = '/networks/{network_id}/prometheus/alert_receiver';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertReceiver(
        parameters: {
            'networkId': string,
            'receiverConfig': alert_receiver_config,
        }
    ): Promise < "Created" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiverConfig'] === undefined) {
            throw new Error('Missing required  parameter: receiverConfig');
        }

        if (parameters['receiverConfig'] !== undefined) {
            body = parameters['receiverConfig'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertReceiverByReceiver(
        parameters: {
            'networkId': string,
            'receiver': string,
            'receiverConfig': alert_receiver_config,
        }
    ): Promise < "Updated" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver/{receiver}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiver'] === undefined) {
            throw new Error('Missing required  parameter: receiver');
        }

        path = path.replace('{receiver}', `${parameters['receiver']}`);

        if (parameters['receiverConfig'] === undefined) {
            throw new Error('Missing required  parameter: receiverConfig');
        }

        if (parameters['receiverConfig'] !== undefined) {
            body = parameters['receiverConfig'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertReceiverRoute(
            parameters: {
                'networkId': string,
            }
        ): Promise < alert_routing_tree >
        {
            let path = '/networks/{network_id}/prometheus/alert_receiver/route';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertReceiverRoute(
        parameters: {
            'networkId': string,
            'route': alert_routing_tree,
        }
    ): Promise < "OK" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver/route';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['route'] === undefined) {
            throw new Error('Missing required  parameter: route');
        }

        if (parameters['route'] !== undefined) {
            body = parameters['route'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdPrometheusQuery(
            parameters: {
                'networkId': string,
                'query': string,
                'time' ? : string,
            }
        ): Promise < promql_return_object >
        {
            let path = '/networks/{network_id}/prometheus/query';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['query'] === undefined) {
                throw new Error('Missing required  parameter: query');
            }

            if (parameters['query'] !== undefined) {
                query['query'] = parameters['query'];
            }

            if (parameters['time'] !== undefined) {
                query['time'] = parameters['time'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdPrometheusQueryRange(
            parameters: {
                'networkId': string,
                'query': string,
                'start': string,
                'end' ? : string,
                'step' ? : string,
            }
        ): Promise < promql_return_object >
        {
            let path = '/networks/{network_id}/prometheus/query_range';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['query'] === undefined) {
                throw new Error('Missing required  parameter: query');
            }

            if (parameters['query'] !== undefined) {
                query['query'] = parameters['query'];
            }

            if (parameters['start'] === undefined) {
                throw new Error('Missing required  parameter: start');
            }

            if (parameters['start'] !== undefined) {
                query['start'] = parameters['start'];
            }

            if (parameters['end'] !== undefined) {
                query['end'] = parameters['end'];
            }

            if (parameters['step'] !== undefined) {
                query['step'] = parameters['step'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdPrometheusSeries(
            parameters: {
                'networkId': string,
                'match' ? : Array < string >
                    ,
                'start' ? : string,
                'end' ? : string,
            }
        ): Promise < Array < prometheus_labelset >
        >
        {
            let path = '/networks/{network_id}/prometheus/series';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['match'] !== undefined) {
                query['match'] = parameters['match'];
            }

            if (parameters['start'] !== undefined) {
                query['start'] = parameters['start'];
            }

            if (parameters['end'] !== undefined) {
                query['end'] = parameters['end'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdRatingGroups(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < rating_group >
        >
        {
            let path = '/networks/{network_id}/rating_groups';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdRatingGroups(
            parameters: {
                'networkId': string,
                'ratingGroup': rating_group,
            }
        ): Promise < rating_group_id >
        {
            let path = '/networks/{network_id}/rating_groups';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['ratingGroup'] === undefined) {
                throw new Error('Missing required  parameter: ratingGroup');
            }

            if (parameters['ratingGroup'] !== undefined) {
                body = parameters['ratingGroup'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async deleteNetworksByNetworkIdRatingGroupsByRatingGroupId(
        parameters: {
            'networkId': string,
            'ratingGroupId': number,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/rating_groups/{rating_group_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ratingGroupId'] === undefined) {
            throw new Error('Missing required  parameter: ratingGroupId');
        }

        path = path.replace('{rating_group_id}', `${parameters['ratingGroupId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdRatingGroupsByRatingGroupId(
            parameters: {
                'networkId': string,
                'ratingGroupId': number,
            }
        ): Promise < rating_group >
        {
            let path = '/networks/{network_id}/rating_groups/{rating_group_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['ratingGroupId'] === undefined) {
                throw new Error('Missing required  parameter: ratingGroupId');
            }

            path = path.replace('{rating_group_id}', `${parameters['ratingGroupId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdRatingGroupsByRatingGroupId(
        parameters: {
            'networkId': string,
            'ratingGroupId': number,
            'ratingGroup': mutable_rating_group,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/rating_groups/{rating_group_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['ratingGroupId'] === undefined) {
            throw new Error('Missing required  parameter: ratingGroupId');
        }

        path = path.replace('{rating_group_id}', `${parameters['ratingGroupId']}`);

        if (parameters['ratingGroup'] === undefined) {
            throw new Error('Missing required  parameter: ratingGroup');
        }

        if (parameters['ratingGroup'] !== undefined) {
            body = parameters['ratingGroup'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdTiers(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < tier_id >
        >
        {
            let path = '/networks/{network_id}/tiers';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiers(
        parameters: {
            'networkId': string,
            'tier': tier,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierId(
        parameters: {
            'networkId': string,
            'tierId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierId(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierId(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdGateways(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_gateways >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiersByTierIdGateways(
        parameters: {
            'networkId': string,
            'tierId': string,
            'gateway': gateway_id,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdTiersByTierIdGateways(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier_gateways,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'tierId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdImages(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_images >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/images';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiersByTierIdImages(
        parameters: {
            'networkId': string,
            'tierId': string,
            'image': tier_image,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['image'] === undefined) {
            throw new Error('Missing required  parameter: image');
        }

        if (parameters['image'] !== undefined) {
            body = parameters['image'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdTiersByTierIdImages(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier_images,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierIdImagesByImageName(
        parameters: {
            'networkId': string,
            'tierId': string,
            'imageName': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images/{image_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['imageName'] === undefined) {
            throw new Error('Missing required  parameter: imageName');
        }

        path = path.replace('{image_name}', `${parameters['imageName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdName(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_name >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierIdName(
        parameters: {
            'networkId': string,
            'tierId': string,
            'name': tier_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdVersion(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_version >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/version';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierIdVersion(
        parameters: {
            'networkId': string,
            'tierId': string,
            'version': tier_version,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/version';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['version'] === undefined) {
            throw new Error('Missing required  parameter: version');
        }

        if (parameters['version'] !== undefined) {
            body = parameters['version'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdType(
            parameters: {
                'networkId': string,
            }
        ): Promise < string >
        {
            let path = '/networks/{network_id}/type';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdType(
        parameters: {
            'networkId': string,
            'type': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/type';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['type'] === undefined) {
            throw new Error('Missing required  parameter: type');
        }

        if (parameters['type'] !== undefined) {
            body = parameters['type'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphony(): Promise < Array < string >
        >
        {
            let path = '/symphony';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postSymphony(
        parameters: {
            'symphonyNetwork': symphony_network,
        }
    ): Promise < "Success" > {
        let path = '/symphony';
        let body;
        let query = {};
        if (parameters['symphonyNetwork'] === undefined) {
            throw new Error('Missing required  parameter: symphonyNetwork');
        }

        if (parameters['symphonyNetwork'] !== undefined) {
            body = parameters['symphonyNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteSymphonyByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getSymphonyByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < symphony_network >
        {
            let path = '/symphony/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkId(
        parameters: {
            'networkId': string,
            'symphonyNetwork': symphony_network,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['symphonyNetwork'] === undefined) {
            throw new Error('Missing required  parameter: symphonyNetwork');
        }

        if (parameters['symphonyNetwork'] !== undefined) {
            body = parameters['symphonyNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgents(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: symphony_agent,
        } >
        {
            let path = '/symphony/{network_id}/agents';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postSymphonyByNetworkIdAgents(
        parameters: {
            'networkId': string,
            'symphonyAgent': mutable_symphony_agent,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['symphonyAgent'] === undefined) {
            throw new Error('Missing required  parameter: symphonyAgent');
        }

        if (parameters['symphonyAgent'] !== undefined) {
            body = parameters['symphonyAgent'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteSymphonyByNetworkIdAgentsByAgentId(
        parameters: {
            'networkId': string,
            'agentId': string,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentId(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < symphony_agent >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentId(
        parameters: {
            'networkId': string,
            'agentId': string,
            'agent': mutable_symphony_agent,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['agent'] === undefined) {
            throw new Error('Missing required  parameter: agent');
        }

        if (parameters['agent'] !== undefined) {
            body = parameters['agent'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdDescription(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdDescription(
        parameters: {
            'networkId': string,
            'agentId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdDevice(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdDevice(
        parameters: {
            'networkId': string,
            'agentId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdMagmad(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdMagmad(
        parameters: {
            'networkId': string,
            'agentId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdManagedDevices(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < managed_devices >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/managed_devices';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdManagedDevices(
        parameters: {
            'networkId': string,
            'agentId': string,
            'managedDevices': managed_devices,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/managed_devices';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['managedDevices'] === undefined) {
            throw new Error('Missing required  parameter: managedDevices');
        }

        if (parameters['managedDevices'] !== undefined) {
            body = parameters['managedDevices'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdName(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdName(
        parameters: {
            'networkId': string,
            'agentId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdAgentsByAgentIdTier(
            parameters: {
                'networkId': string,
                'agentId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/symphony/{network_id}/agents/{agent_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['agentId'] === undefined) {
                throw new Error('Missing required  parameter: agentId');
            }

            path = path.replace('{agent_id}', `${parameters['agentId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdAgentsByAgentIdTier(
        parameters: {
            'networkId': string,
            'agentId': string,
            'tier': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/agents/{agent_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['agentId'] === undefined) {
            throw new Error('Missing required  parameter: agentId');
        }

        path = path.replace('{agent_id}', `${parameters['agentId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/symphony/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdDevices(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: symphony_device,
        } >
        {
            let path = '/symphony/{network_id}/devices';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postSymphonyByNetworkIdDevices(
        parameters: {
            'networkId': string,
            'symphonyDevice': mutable_symphony_device,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/devices';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['symphonyDevice'] === undefined) {
            throw new Error('Missing required  parameter: symphonyDevice');
        }

        if (parameters['symphonyDevice'] !== undefined) {
            body = parameters['symphonyDevice'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteSymphonyByNetworkIdDevicesByDeviceId(
        parameters: {
            'networkId': string,
            'deviceId': string,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/devices/{device_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['deviceId'] === undefined) {
            throw new Error('Missing required  parameter: deviceId');
        }

        path = path.replace('{device_id}', `${parameters['deviceId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getSymphonyByNetworkIdDevicesByDeviceId(
            parameters: {
                'networkId': string,
                'deviceId': string,
            }
        ): Promise < symphony_device >
        {
            let path = '/symphony/{network_id}/devices/{device_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['deviceId'] === undefined) {
                throw new Error('Missing required  parameter: deviceId');
            }

            path = path.replace('{device_id}', `${parameters['deviceId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdDevicesByDeviceId(
        parameters: {
            'networkId': string,
            'deviceId': string,
            'symphonyDevice': mutable_symphony_device,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/devices/{device_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['deviceId'] === undefined) {
            throw new Error('Missing required  parameter: deviceId');
        }

        path = path.replace('{device_id}', `${parameters['deviceId']}`);

        if (parameters['symphonyDevice'] === undefined) {
            throw new Error('Missing required  parameter: symphonyDevice');
        }

        if (parameters['symphonyDevice'] !== undefined) {
            body = parameters['symphonyDevice'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdDevicesByDeviceIdConfig(
            parameters: {
                'networkId': string,
                'deviceId': string,
            }
        ): Promise < symphony_device_config >
        {
            let path = '/symphony/{network_id}/devices/{device_id}/config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['deviceId'] === undefined) {
                throw new Error('Missing required  parameter: deviceId');
            }

            path = path.replace('{device_id}', `${parameters['deviceId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdDevicesByDeviceIdConfig(
        parameters: {
            'networkId': string,
            'deviceId': string,
            'name': symphony_device_config,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/devices/{device_id}/config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['deviceId'] === undefined) {
            throw new Error('Missing required  parameter: deviceId');
        }

        path = path.replace('{device_id}', `${parameters['deviceId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdDevicesByDeviceIdName(
            parameters: {
                'networkId': string,
                'deviceId': string,
            }
        ): Promise < symphony_device_name >
        {
            let path = '/symphony/{network_id}/devices/{device_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['deviceId'] === undefined) {
                throw new Error('Missing required  parameter: deviceId');
            }

            path = path.replace('{device_id}', `${parameters['deviceId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdDevicesByDeviceIdName(
        parameters: {
            'networkId': string,
            'deviceId': string,
            'name': symphony_device_name,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/devices/{device_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['deviceId'] === undefined) {
            throw new Error('Missing required  parameter: deviceId');
        }

        path = path.replace('{device_id}', `${parameters['deviceId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdDevicesByDeviceIdState(
            parameters: {
                'networkId': string,
                'deviceId': string,
            }
        ): Promise < symphony_device_state >
        {
            let path = '/symphony/{network_id}/devices/{device_id}/state';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['deviceId'] === undefined) {
                throw new Error('Missing required  parameter: deviceId');
            }

            path = path.replace('{device_id}', `${parameters['deviceId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getSymphonyByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/symphony/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'config': network_features,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/symphony/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifi(): Promise < Array < string >
        >
        {
            let path = '/wifi';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postWifi(
        parameters: {
            'wifiNetwork': wifi_network,
        }
    ): Promise < "Success" > {
        let path = '/wifi';
        let body;
        let query = {};
        if (parameters['wifiNetwork'] === undefined) {
            throw new Error('Missing required  parameter: wifiNetwork');
        }

        if (parameters['wifiNetwork'] !== undefined) {
            body = parameters['wifiNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteWifiByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getWifiByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < wifi_network >
        {
            let path = '/wifi/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkId(
        parameters: {
            'networkId': string,
            'wifiNetwork': wifi_network,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['wifiNetwork'] === undefined) {
            throw new Error('Missing required  parameter: wifiNetwork');
        }

        if (parameters['wifiNetwork'] !== undefined) {
            body = parameters['wifiNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/wifi/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/wifi/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'config': network_features,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: wifi_gateway,
        } >
        {
            let path = '/wifi/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postWifiByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': mutable_wifi_gateway,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteWifiByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < wifi_gateway >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': mutable_wifi_gateway,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getWifiByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdGatewaysByGatewayIdWifi(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_wifi_configs >
        {
            let path = '/wifi/{network_id}/gateways/{gateway_id}/wifi';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdGatewaysByGatewayIdWifi(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_wifi_configs,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/gateways/{gateway_id}/wifi';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdMeshes(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < mesh_id >
        >
        {
            let path = '/wifi/{network_id}/meshes';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postWifiByNetworkIdMeshes(
            parameters: {
                'networkId': string,
                'wifiMesh': wifi_mesh,
            }
        ): Promise < mesh_id >
        {
            let path = '/wifi/{network_id}/meshes';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['wifiMesh'] === undefined) {
                throw new Error('Missing required  parameter: wifiMesh');
            }

            if (parameters['wifiMesh'] !== undefined) {
                body = parameters['wifiMesh'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async deleteWifiByNetworkIdMeshesByMeshId(
        parameters: {
            'networkId': string,
            'meshId': string,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/meshes/{mesh_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['meshId'] === undefined) {
            throw new Error('Missing required  parameter: meshId');
        }

        path = path.replace('{mesh_id}', `${parameters['meshId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getWifiByNetworkIdMeshesByMeshId(
            parameters: {
                'networkId': string,
                'meshId': string,
            }
        ): Promise < wifi_mesh >
        {
            let path = '/wifi/{network_id}/meshes/{mesh_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['meshId'] === undefined) {
                throw new Error('Missing required  parameter: meshId');
            }

            path = path.replace('{mesh_id}', `${parameters['meshId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdMeshesByMeshId(
            parameters: {
                'networkId': string,
                'meshId': string,
                'wifiMesh': wifi_mesh,
            }
        ): Promise < mesh_id >
        {
            let path = '/wifi/{network_id}/meshes/{mesh_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['meshId'] === undefined) {
                throw new Error('Missing required  parameter: meshId');
            }

            path = path.replace('{mesh_id}', `${parameters['meshId']}`);

            if (parameters['wifiMesh'] === undefined) {
                throw new Error('Missing required  parameter: wifiMesh');
            }

            if (parameters['wifiMesh'] !== undefined) {
                body = parameters['wifiMesh'];
            }

            return await this.request(path, 'PUT', query, body);
        }
    static async getWifiByNetworkIdMeshesByMeshIdConfig(
            parameters: {
                'networkId': string,
                'meshId': string,
            }
        ): Promise < mesh_wifi_configs >
        {
            let path = '/wifi/{network_id}/meshes/{mesh_id}/config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['meshId'] === undefined) {
                throw new Error('Missing required  parameter: meshId');
            }

            path = path.replace('{mesh_id}', `${parameters['meshId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdMeshesByMeshIdConfig(
        parameters: {
            'networkId': string,
            'meshId': string,
            'meshWifiConfigs': mesh_wifi_configs,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/meshes/{mesh_id}/config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['meshId'] === undefined) {
            throw new Error('Missing required  parameter: meshId');
        }

        path = path.replace('{mesh_id}', `${parameters['meshId']}`);

        if (parameters['meshWifiConfigs'] === undefined) {
            throw new Error('Missing required  parameter: meshWifiConfigs');
        }

        if (parameters['meshWifiConfigs'] !== undefined) {
            body = parameters['meshWifiConfigs'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdMeshesByMeshIdName(
            parameters: {
                'networkId': string,
                'meshId': string,
            }
        ): Promise < mesh_name >
        {
            let path = '/wifi/{network_id}/meshes/{mesh_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['meshId'] === undefined) {
                throw new Error('Missing required  parameter: meshId');
            }

            path = path.replace('{mesh_id}', `${parameters['meshId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdMeshesByMeshIdName(
        parameters: {
            'networkId': string,
            'meshId': string,
            'meshName': mesh_name,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/meshes/{mesh_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['meshId'] === undefined) {
            throw new Error('Missing required  parameter: meshId');
        }

        path = path.replace('{mesh_id}', `${parameters['meshId']}`);

        if (parameters['meshName'] === undefined) {
            throw new Error('Missing required  parameter: meshName');
        }

        if (parameters['meshName'] !== undefined) {
            body = parameters['meshName'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/wifi/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getWifiByNetworkIdWifi(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_wifi_configs >
        {
            let path = '/wifi/{network_id}/wifi';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putWifiByNetworkIdWifi(
        parameters: {
            'networkId': string,
            'config': network_wifi_configs,
        }
    ): Promise < "Success" > {
        let path = '/wifi/{network_id}/wifi';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
}

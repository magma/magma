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
import os

HIL = {
    "TEST_LIB_PATH": "TS",
}
TAS = {
    "tas_ip": os.getenv("TAS_IP", "192.168.0.142"),
    "username": os.getenv("TAS_USERNAME", "hil_user"),
    "password": os.getenv("TAS_PASSWORD"),
    "timeout": 1000,  # test run timeout value in sec
    "test_report_path": "/tmp/",
}

AWS = {
    "access_key": os.getenv("AWS_ACCESS_KEY"),
    "secret_key": os.getenv("AWS_SECRET_KEY"),
    "region": os.getenv("AWS_DEFAULT_REGION", "us-west-1"),
}
RDS = {
    "db_host": os.getenv("RDS_HOST"),
    "db_user": os.getenv("RDS_USER", "admin"),
    "db_pass": os.getenv("RDS_PASS"),
    "database": "MagmaAutomation",
    "table_name": "HilSanityResults",
}
MAGMA = {
    "magma_pkg_name": "magma",
    "magma_sctpd_pkg_name": "magma-sctpd",
    "username": os.getenv("MAGMA_USERNAME", "magma"),
    "password": os.getenv("MAGMA_PASSWORD"),
    "processes": ["magma@*"],
    "timezone": "US/Pacific",
    "sanity_res_badge": "hilsanityres.svg",
    "sanity_pass_badge": "hilsanitypass.svg",
    "core_path": "/var/core",
    "memory_delta_failed_pct": 75,
    "control_procedure_tolerance_pct": 5,
    # select same number of UE for gold/silver plan. max is 100.
    "gold_plan_UE": 100,
    "silver_plan_UE": 100,
    "gold_plan_qos": 1000000,  # bps
    "silver_plan_qos": 500000,  # bps
    "nms_apn_ambr": 100000000,  # 100mbps
    "pass_percentile": 0.95,
    "AGW": (
        "hil-1",
        "hil-2",
        "hil-3",
        "hil-4",
        "hil-5",
        "hil-6",
        "hil-7",
        "hil-8",
        "hil-9",
    ),
    "REL": ("ci", "1.6.0", "1.6.1", "test"),
    "EPC": ("magma"),
    "Test_suite": (
        "SANITY",
        "PERFORMANCE",
        "FEATURE",
        "AVAILABILITY",
        "STATIC",
        "IPV6",
    ),
}
SLACK = {
    "slack_webhook_path": os.getenv("SLACK_WEBHOOK_PATH"),
    "dashboard": "http://automation.yourdomain.xyz",
}

MAGMA_AGW = {
    "UE_STATE": {
        "mme_state": 0,
        "s1ap_state": 0,
        "spgw_state": 0,
        "table0_flows": 3,
        "table12_flows": 2,
        "table13_flows": 2,
    },
}

# merge this multiple spirent config under one dict

SPIRENT = {
    "1nodal_1nw_steps": lambda run_time: [
        {
            "delaySec": 0,
            "predecessorState": "",
            "predecessorTcIndex": -1,
            "predecessorTsIndex": -1,
            "tcActivity": "Init",
            "tcIndex": 0,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "",
            "predecessorTcIndex": -1,
            "predecessorTsIndex": -1,
            "tcActivity": "Init",
            "tcIndex": 1,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "",
            "predecessorTcIndex": -1,
            "predecessorTsIndex": -1,
            "tcActivity": "Start",
            "tcIndex": 1,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "RUNNING",
            "predecessorTcIndex": 1,
            "predecessorTsIndex": 0,
            "tcActivity": "Start",
            "tcIndex": 0,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "RUNNING",
            "predecessorTcIndex": 0,
            "predecessorTsIndex": 0,
            "tcActivity": "",
            "tcIndex": -1,
            "tsIndex": -1,
        },
        {
            "delaySec": run_time,
            "predecessorState": "RUNNING",
            "predecessorTcIndex": 0,
            "predecessorTsIndex": 0,
            "tcActivity": "Stop",
            "tcIndex": 0,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "STOPPED",
            "predecessorTcIndex": 0,
            "predecessorTsIndex": 0,
            "tcActivity": "",
            "tcIndex": -1,
            "tsIndex": -1,
        },
        {
            "delaySec": 0,
            "predecessorState": "",
            "predecessorTcIndex": -1,
            "predecessorTsIndex": -1,
            "tcActivity": "Stop",
            "tcIndex": 1,
            "tsIndex": 0,
        },
        {
            "delaySec": 0,
            "predecessorState": "STOPPED",
            "predecessorTcIndex": 1,
            "predecessorTsIndex": 0,
            "tcActivity": "Cleanup",
            "tcIndex": 1,
            "tsIndex": 0,
        },
        {
            "delaySec": 60,
            "predecessorState": "",
            "predecessorTcIndex": -1,
            "predecessorTsIndex": -1,
            "tcActivity": "Cleanup",
            "tcIndex": 0,
            "tsIndex": 0,
        },
    ],
}
SPIRENT_SEQUENCER_DELAY_PROFILE = {
    "short": [5, 10],
    "medium": [30, 80],
    "long": [100, 600],
}

SPIRENT_SEQUENCER_COMMAND = {
    "active_idle": '{% raw %}{  { Delay 5 }  { OnDemandCommand { ControlBearer { "op=Attach" "rate={% endraw %}{{rate}}" "start_sub=1" "end_sub={{total_subs}}" {% raw %}} } }  { LoopStart }  { Delay {% endraw %}{{inter_nodal_offset}} {% raw %}}  { StartDmf {% endraw %}{{rate}}  0 {% raw %}} { Delay 1 } { StartDmf {% endraw %}{{rate}}  1 {% raw %}}  { Wait }  { StopDmf 10000.0 0 }  { StopAllTraffic 10000.0  0 }  { Wait }  { LoopEnd {% endraw %}{{iterations}} {% raw %}} { Delay 5 } { OnDemandCommand { ControlBearer { "op=Detach" "rate={% endraw %}{{rate}}" "start_sub=1" "end_sub={{total_subs}}" {% raw %}} } } { Delay {% endraw %}{{detach_delay}} {% raw %}} }{% endraw %}',
}


SPIRENT_SLEEP_PROFILE = [90, 120, 180, 240]


HEADER_ENRICHMENT = {
    "capture_filter": "port 80 and tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420",
    "interface": "eth0",
    "output_file": "/tmp/header_enrich_no_cipher.pcapng",
}

IPFIX = {
    "capture_filter": "udp and port 65010",
    "interface": "eth0",
    "output_file": "/tmp/ipfix_hil.pcapng",
}

IPFIX_AGW_CONFIG = {
    "magmad": {"magma_services": "- connectiond"},
    "connectiond": {
        "log_level": "ERROR",
        "interface_name": "ipfix0",
        "zone": 897,
        "pkt_dst_mac": "33:aa:99:33:aa:00",
        "pkt_src_mac": "55:11:44:ee:00:00",
    },
    "pipelined": {
        "ipfix": {
            "enabled": "true",
            "probability": 65,
            "collector_set_id": 1,
            "collector_ip": "10.22.3.85",
            "collector_port": 60,
            "cache_timeout": 60,
            "obs_domain_id": 1,
            "obs_point_id": 1,
        },
        "static_services": ["conntrack", "ipfix"],
    },
}

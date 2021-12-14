/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <gtest/gtest.h>
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"

extern "C" {
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
}

namespace magma {
namespace lte {

const char* kEmptyConfig =
    R"libconfig(MME :
{
};)libconfig";

const char* kConfigMissingUpstreamSock =
    R"libconfig(MME :
{
    SCTP :
    {
        # Path to sctpd up and downstream unix domain sockets
        SCTP_DOWNSTREAM_SOCK = "unix:///tmp/sctpd_downstream_test.sock";
    };
};)libconfig";

const char* kHealthyConfig =
    R"libconfig(MME :
{
    REALM                                     = "magma.com"
    PID_DIRECTORY                             = "/var/run";
    # Define the limits of the system in terms of served eNB and served UE.
    # When the limits will be reached, overload procedure will take place.
    MAXENB                                    = 8;                              # power of 2
    MAXUE                                     = 16;                             # power of 2
    RELATIVE_CAPACITY                         = 11;

    EMERGENCY_ATTACH_SUPPORTED                     = "no";
    UNAUTHENTICATED_IMSI_SUPPORTED                 = "no";

    # EPS network feature support
    EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1      = "no";    # DO NOT CHANGE
    EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE = "no";    # DO NOT CHANGE
    EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC            = "no";    # DO NOT CHANGE
    EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST             = "no";    # DO NOT CHANGE

    # Report/Display MME statistics (expressed in seconds)
    STATS_TIMER_SEC                    = 60;

    USE_STATELESS = "True";
    USE_HA = "False";
    ENABLE_GTPU_PRIVATE_IP_CORRECTION = "False";
    ENABLE_CONVERGED_CORE = "False";

    # Congestion control configuration parameters
    CONGESTION_CONTROL_ENABLED = "True";
    # Congestion control thresholds (expressed in microseconds)
    S1AP_ZMQ_TH = 2000000;
    MME_APP_ZMQ_CONGEST_TH = 100000;
    MME_APP_ZMQ_AUTH_TH = 200000;
    MME_APP_ZMQ_IDENT_TH = 400000;
    MME_APP_ZMQ_SMC_TH = 1000000;

    INTERTASK_INTERFACE :
    {
        # max queue size per task
        ITTI_QUEUE_SIZE            = 2000000;
    };

    S6A :
    {
        S6A_CONF                   = "/var/opt/magma/tmp/mme_fd.conf"; # YOUR MME freeDiameter config file path
        HSS_HOSTNAME               = "hss"; # relevant for freeDiameter only
    };

    # ------- SCTP definitions
    SCTP :
    {
        # Path to sctpd up and downstream unix domain sockets
        SCTP_UPSTREAM_SOCK = "unix:///tmp/sctpd_upstream_test.sock";
        SCTP_DOWNSTREAM_SOCK = "unix:///tmp/sctpd_downstream_test.sock";
    };

    # ------- S1AP definitions
    S1AP :
    {
        # outcome drop timer value (seconds)
        S1AP_OUTCOME_TIMER = 10;
    };

    # ------- MME served GUMMEIs
    # MME code DEFAULT  size = 8 bits
    # MME GROUP ID size = 16 bits
    GUMMEI_LIST = (
         { MCC="001" ; MNC="01"; MME_GID="1" ; MME_CODE="1"; }
    );

    # ------- MME served TAIs
    # TA (mcc.mnc:tracking area code) DEFAULT = 208.34:1
    # max values = 999.999:65535
    # maximum of 16 TAIs, comma separated
    # !!! Actually use only one PLMN
    TAI_LIST = (
    
         { MCC="001" ; MNC="01" ; TAC="1"; }
    );

    TAC_LIST = (
         { MCC="001" ; MNC="01" ; TAC="1"; }
         
    );

    # List of restricted PLMNs
    # By default this list is empty
    # Max number of restricted plmn is 10
    RESTRICTED_PLMN_LIST = (
        # PlmnConfig values can be found at magma/lte/protos/mconfig/mconfigs.proto
        
    );

    # List of blocked IMEIs
    # By default this list is empty
    # Stored in a hash table on mme side
    # Length of IMEI=15 digits, length of IMEISV=16 digits
    BLOCKED_IMEI_LIST = (
        # Sample IMEI: TAC(8 digits) + SNR (6 digits)
        #{ IMEI_TAC="99000482"; SNR="351037"}
        # Sample IMEI without SNR: TAC(8 digits)
        #{ IMEI_TAC="99000482";}
        # ImeiConfig values can be found at magma/lte/protos/mconfig/mconfigs.proto
        
    );

    CSFB :
    {
        NON_EPS_SERVICE_CONTROL = "OFF";
        CSFB_MCC = "001";
        CSFB_MNC = "01";
        LAC = "1";
    };

    # NGAP definitions
    NGAP :
    {
      # outcome drop timer value (seconds)
      #  S1AP_OUTCOME_TIMER = 10;
      AMF_NAME = "MAGMAAMF1"
    };
    # AMF served GUAMFIs
    # AMF code DEFAULT  size = 8 bits
    # AMF GROUP ID size = 16 bits
    GUAMFI_LIST = (
     { MCC="001" ; MNC="01"; AMF_REGION_ID="1" ; AMF_SET_ID="1"; AMF_POINTER="3"}
    );

    NAS :
    {
        # 3GPP TS 33.401 section 7.2.4.3 Procedures for NAS algorithm selection
        # decreasing preference goes from left to right
        ORDERED_SUPPORTED_INTEGRITY_ALGORITHM_LIST = [ "EIA2" , "EIA1" , "EIA0" ];
        ORDERED_SUPPORTED_CIPHERING_ALGORITHM_LIST = [ "EEA0" , "EEA1" , "EEA2" ];

        # EMM TIMERS
        # T3402 start:
        # At attach failure and the attempt counter is equal to 5.
        # At tracking area updating failure and the attempt counter is equal to 5.
        # T3402 stop:
        # ATTACH REQUEST sent, TRACKING AREA REQUEST sent.
        # On expiry:
        # Initiation of the attach procedure, if still required or TAU procedure
        # attached for emergency bearer services.
        T3402                                 =  1                              # in minutes (default is 12 minutes)

        # T3412 start:
        # In EMM-REGISTERED, when EMM-CONNECTED mode is left.
        # T3412 stop:
        # When entering state EMM-DEREGISTERED or when entering EMM-CONNECTED mode.
        # On expiry:
        # Initiation of the periodic TAU procedure if the UE is not attached for
        # emergency bearer services. Implicit detach from network if the UE is
        # attached for emergency bearer services.
        T3412                                 =  54                             # in minutes (default is 54 minutes, network dependent)
        # T3422 start: DETACH REQUEST sent
        # T3422 stop: DETACH ACCEPT received
        # ON THE 1st, 2nd, 3rd, 4th EXPIRY: Retransmission of DETACH REQUEST
        T3422                                 =  6                              # in seconds (default is 6s)

        # T3450 start:
        # ATTACH ACCEPT sent, TRACKING AREA UPDATE ACCEPT sent with GUTI, TRACKING AREA UPDATE ACCEPT sent with TMSI,
        # GUTI REALLOCATION COMMAND sent
        # T3450 stop:
        # ATTACH COMPLETE received, TRACKING AREA UPDATE COMPLETE received, GUTI REALLOCATION COMPLETE received
        # ON THE 1st, 2nd, 3rd, 4th EXPIRY: Retransmission of the same message type
        T3450                                 =  6                              # in seconds (default is 6s)

        # T3460 start: AUTHENTICATION REQUEST sent, SECURITY MODE COMMAND sent
        # T3460 stop:
        # AUTHENTICATION RESPONSE received, AUTHENTICATION FAILURE received,
        # SECURITY MODE COMPLETE received, SECURITY MODE REJECT received
        # ON THE 1st, 2nd, 3rd, 4th EXPIRY: Retransmission of the same message type
        T3460                                 =  6                              # in seconds (default is 6s)

        # T3470 start: IDENTITY REQUEST sent
        # T3470 stop: IDENTITY RESPONSE received
        # ON THE 1st, 2nd, 3rd, 4th EXPIRY: Retransmission of IDENTITY REQUEST
        T3470                                 =  6                              # in seconds (default is 6s)

        # ESM TIMERS
        T3485                                 =  8                              # UNUSED in seconds (default is 8s)
        T3486                                 =  8                              # UNUSED in seconds (default is 8s)
        T3489                                 =  4                              # UNUSED in seconds (default is 4s)
        T3495                                 =  8                              # UNUSED in seconds (default is 8s)

        # APN CORRECTION FEATURE
        ENABLE_APN_CORRECTION                 = "False"
        APN_CORRECTION_MAP_LIST               = (
          {
            APN_CORRECTION_MAP_IMSI_PREFIX = "00101" ;
            APN_CORRECTION_MAP_APN_OVERRIDE = "magma.ipv4" ;
          }
         
        );
    };

    SGS :
    {
        # TS6_1 start: SGSAP LOCATION UPDATE REQUEST sent
        # TS6_1 stop: SGSAP LOCATION UPDATE ACCEPT received,SGSAP LOCATION UPDATE REJECT received
        TS6_1                                 =  10                             # in seconds (default is 10s)

        # TS8 start: SGSAP EPS DETACH INDICATION explicit detach sent for EPS services
        # TS8 stop: SGSAP EPS DETACH ACK  received
        TS8                                   =  4                              # in seconds (default is 4s)

        # TS9 start: SGSAP IMSI DETACH INDICATION explicit detach sent for non-EPS services
        # TS9 stop:  SGSAP IMSI DETACH ACK received
        # changed the Ts9 default value to 2s since the T3421 ue detach timer value is 5s
        # To avoid retransmission of UE detach message and small delay to wait for sgs detach retransmission
        TS9                                   =  2                              # in seconds (default is 4s)

        # TS10 start: SGSAP IMSI DETACH INDICATION implicit detach sent for non-EPS services
        # TS10 stop: SGSAP EPS DETACH ACK  received
        TS10                                   =  4                              # in seconds (default is 4s)

        # TS13 start: SGSAP EPS DETACH INDICATION implicit detach sent for EPS services
        # TS13 stop: SGSAP EPS DETACH ACK  received
        TS13                                   =  4                              # in seconds (default is 4s)


    };
    NETWORK_INTERFACES :
    {
        # MME binded interface for S1-C or S1-MME  communication (S1AP), can be ethernet interface, virtual ethernet interface,
        # we don't advise wireless interfaces
        MME_INTERFACE_NAME_FOR_S1_MME         = "eth1";
        MME_IPV4_ADDRESS_FOR_S1_MME           = "192.168.60.142/24";

        # MME binded interface for S11 communication (GTPV2-C)
        MME_INTERFACE_NAME_FOR_S11_MME        = "eth1";
        MME_IPV4_ADDRESS_FOR_S11_MME          = "192.168.60.142/24";
        MME_PORT_FOR_S11_MME                  = 2123;
    };

    LOGGING :
    {
        # OUTPUT choice in { "CONSOLE", "SYSLOG", `path to file`", "`IPv4@`:`TCP port num`"}
        # `path to file` must start with '.' or '/'
        # if TCP stream choice, then you can easily dump the traffic on the remote or local host: nc -l `TCP port num` > received.txt
        #OUTPUT            = "CONSOLE";
        #OUTPUT            = "SYSLOG";
        OUTPUT            = "/var/log/mme.log";
        #OUTPUT            = "127.0.0.1:5656";

        # THREAD_SAFE choice in { "yes", "no" } means use of thread safe intermediate buffer then a single thread pick each message log one
        # by one to flush it to the chosen output
        THREAD_SAFE       = "no";

        # COLOR choice in { "yes", "no" } means use of ANSI styling codes or no
        COLOR             = "no";

        # Log level choice in { "EMERGENCY", "ALERT", "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO", "DEBUG", "TRACE"}
        SCTP_LOG_LEVEL     = "INFO";
        GTPV1U_LOG_LEVEL   = "INFO";
        SPGW_APP_LOG_LEVEL = "INFO";
        UDP_LOG_LEVEL      = "INFO";
        S1AP_LOG_LEVEL     = "INFO";
        NAS_LOG_LEVEL      = "INFO";
        MME_APP_LOG_LEVEL  = "INFO";
        GTPV2C_LOG_LEVEL   = "INFO";
        S11_LOG_LEVEL      = "INFO";
        S6A_LOG_LEVEL      = "INFO";
        UTIL_LOG_LEVEL     = "INFO";
        SERVICE303_LOG_LEVEL     = "INFO";
        MSC_LOG_LEVEL      = "ERROR";
        ITTI_LOG_LEVEL     = "ERROR";
        MME_SCENARIO_PLAYER_LOG_LEVEL = "ERROR";

        # ASN1 VERBOSITY: none, info, annoying
        # for S1AP protocol
        # Won't be templatized because its value space is different
        ASN1_VERBOSITY    = "INFO";
    };
    TESTING :
    {
        # file should be copied here from source tree by following command: run_mme --install-mme-files ...
        SCENARIO_FILE = "/usr/local/share/oai/test/MME/no_regression.xml";
    };

    S-GW :
    {
        # S-GW binded interface for S11 communication (GTPV2-C), if none selected the ITTI message interface is used
        SGW_IPV4_ADDRESS_FOR_S11              = "192.168.60.153";
    };


    FEDERATED_MODE_MAP = (
        # ModeMapItem values can be found at magma/lte/protos/mconfig/mconfigs.proto
        
   );

   SRVC_AREA_CODE_2_TACS_MAP = (
     
   );

   SENTRY_CONFIG = {
     
       # Sentry.io configuration sent from the Orc8r
       SAMPLE_RATE      = 0.0;
       UPLOAD_MME_LOG   = "False";
       URL_NATIVE       = ""
   }
};

NGAP :
{
    # DNS address communicated to UEs
    DEFAULT_DNS_IPV4_ADDRESS     = "8.8.8.8";
    DEFAULT_DNS_SEC_IPV4_ADDRESS = "8.8.4.4";
};)libconfig";

constexpr std::array<int, 25> ncon_tac = {1,  2,  4,  6,  7,  8,  10,
                                          11, 12, 14, 15, 16, 19, 21,
                                          23, 26, 28, 31, 33, 37, 39};

// Test partial list with 1 TAI
TEST(MMEConfigTest, TestOneTai) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 1;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mcc[itr]     = 1;
  config_pP.served_tai.plmn_mnc[itr]     = 1;
  config_pP.served_tai.plmn_mnc_len[itr] = 2;
  config_pP.served_tai.tac[itr]          = 1;
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  // Check if consecutive tacs partial list is created
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);
  ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);

  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[itr].plmn);
  free(config_pP.partial_list[itr].tac);
  free(config_pP.partial_list);
}

// Test 1 partial list with Consecutive Tacs
TEST(MMEConfigTest, TestParTaiListWithConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 16;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Check if consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);

  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 2 partial lists with Consecutive Tacs
TEST(MMEConfigTest, TestTwoParTaiListsWithConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 20;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Check if consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    ASSERT_EQ(
        config_pP.partial_list[itr].list_type,
        TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  itr = 0;
  for (uint8_t idx = 0; idx < config_pP.num_par_lists; idx++) {
    for (uint8_t idx2 = 0; idx2 < config_pP.partial_list[idx].nb_elem; idx2++) {
      ASSERT_EQ(
          config_pP.partial_list[idx].tac[idx2],
          config_pP.served_tai.tac[itr++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 1 partial list with Non-consecutive Tacs
TEST(MMEConfigTest, TestParTaiListWithNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 16;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted non-consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = ncon_tac[itr];
  }
  // Check if non-consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);

  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 2 partial lists with Non-consecutive Tacs
TEST(MMEConfigTest, TestTwoParTaiListsWithNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 20;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted non-consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = ncon_tac[itr];
  }
  // Check if 2 non-consecutive tacs partial lists are created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    ASSERT_EQ(
        config_pP.partial_list[itr].list_type,
        TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  itr = 0;
  for (uint8_t idx = 0; idx < config_pP.num_par_lists; idx++) {
    for (uint8_t idx2 = 0; idx2 < config_pP.partial_list[idx].nb_elem; idx2++) {
      ASSERT_EQ(
          config_pP.partial_list[idx].tac[idx2],
          config_pP.served_tai.tac[itr++]);
    }
  }

  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 2 partial lists with 1-Consecutive tacs and 1-Non-consecutive Tacs
TEST(MMEConfigTest, TestTwoParTaiListsWithConsAndNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 24;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted consecutive TACs
  for (itr = 0; itr < 16; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Sorted non-consecutive TACs
  config_pP.served_tai.tac[itr++] = 19;
  config_pP.served_tai.tac[itr++] = 21;
  config_pP.served_tai.tac[itr++] = 23;
  config_pP.served_tai.tac[itr++] = 26;
  config_pP.served_tai.tac[itr++] = 28;
  config_pP.served_tai.tac[itr++] = 31;
  config_pP.served_tai.tac[itr++] = 33;
  config_pP.served_tai.tac[itr++] = 35;

  /* Check if 1 consecutive tacs partial list and 1 non-consecutive
   * tacs partial lists are created
   */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 8);
  ASSERT_EQ(
      config_pP.partial_list[0].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[1].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  uint8_t idx2 = 0;
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    for (uint8_t idx = 0; idx < config_pP.partial_list[itr].nb_elem; idx++) {
      ASSERT_EQ(
          config_pP.partial_list[itr].tac[idx],
          config_pP.served_tai.tac[idx2++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 1 partial list with many plmns
TEST(MMEConfigTest, TestParTaiListWithManyPlmns) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 6;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = itr + 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
    config_pP.served_tai.tac[itr]          = itr + 1;
  }
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  // Check if partial list with many plmns is created
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, itr + 1);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 3 partial lists, 1-consecutive tacs, 1-non consecutive tacs,1-many plmns
TEST(MMEConfigTest, TestMixedParTaiLists) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 35;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  // Fill the same PLMN for consecutive and non-consecutive tacs (16+16)
  for (itr = 0; itr < 32; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // First 16 sorted consecutive TACs
  for (itr = 0; itr < 16; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Next 16 sorted non-consecutive TACs
  for (uint8_t idx = 0; itr < 32; itr++, idx++) {
    config_pP.served_tai.tac[itr] = ncon_tac[idx];
  }
  // Next 3 many plmns with tacs
  for (uint8_t idx = 0; itr < 35; itr++, idx++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = idx + 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
    config_pP.served_tai.tac[itr]          = idx + 1;
  }

  // Check if 3 partial lists are created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 3);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[2].nb_elem, 3);

  ASSERT_EQ(
      config_pP.partial_list[0].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[1].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[2].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS);
  // Verify plmn for consecutive and non-consecutive tacs
  for (itr = 0; itr < config_pP.num_par_lists - 1; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }

  // Verify plmn for many plmns
  for (uint8_t idx = 0; idx < config_pP.partial_list[3].nb_elem; idx++) {
    EXPECT_FALSE(config_pP.partial_list[3].plmn == nullptr);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit2, idx + 1);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit3, 15);
  }
  // Verify TACs
  uint8_t idx2 = 0;
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    for (uint8_t idx = 0; idx < config_pP.partial_list[itr].nb_elem; idx++) {
      ASSERT_EQ(
          config_pP.partial_list[itr].tac[idx],
          config_pP.served_tai.tac[idx2++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

TEST(MMEConfigTest, TestParseHealthyConfig) {
  mme_config_t mme_config = {0};
  EXPECT_EQ(mme_config_parse_string(kHealthyConfig, &mme_config), 0);
  free_mme_config(&mme_config);
}

TEST(MMEConfigTest, TestParseHealthyConfigDisplay) {
  mme_config_t mme_config = {0};
  EXPECT_EQ(mme_config_parse_string(kHealthyConfig, &mme_config), 0);
  mme_config_display(&mme_config);
  free_mme_config(&mme_config);
}

TEST(MMEConfigTest, TestMissingSctpdConfig) {
  mme_config_t mme_config = {0};
  EXPECT_EQ(mme_config_parse_string(kEmptyConfig, &mme_config), 0);
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.upstream_sctp_sock)),
      "unix:///tmp/sctpd_upstream.sock");
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.downstream_sctp_sock)),
      "unix:///tmp/sctpd_downstream.sock");

  free_mme_config(&mme_config);
}

TEST(MMEConfigTest, TestMissingSctpdUpstreamSockConfig) {
  mme_config_t mme_config = {0};
  EXPECT_EQ(
      mme_config_parse_string(kConfigMissingUpstreamSock, &mme_config), 0);
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.upstream_sctp_sock)),
      "unix:///tmp/sctpd_upstream.sock");
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.downstream_sctp_sock)),
      "unix:///tmp/sctpd_downstream_test.sock");

  free_mme_config(&mme_config);
}

TEST(MMEConfigTest, TestHealthySctpdConfig) {
  mme_config_t mme_config = {0};
  EXPECT_EQ(mme_config_parse_string(kHealthyConfig, &mme_config), 0);
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.upstream_sctp_sock)),
      "unix:///tmp/sctpd_upstream_test.sock");
  EXPECT_EQ(
      std::string(bdata(mme_config.sctp_config.downstream_sctp_sock)),
      "unix:///tmp/sctpd_downstream_test.sock");

  free_mme_config(&mme_config);
}

TEST(MMEConfigTest, TestCopyAmfConfigFromMMEConfig) {
  mme_config_t mme_config = {0};
  amf_config_t amf_config = {0};

  EXPECT_EQ(mme_config_parse_string(kHealthyConfig, &mme_config), 0);

  copy_amf_config_from_mme_config(&amf_config, &mme_config);

  if (mme_config.log_config.output)
    EXPECT_EQ(
        0, strcmp(
               (char*) mme_config.log_config.output->data,
               (char*) amf_config.log_config.output->data));
  EXPECT_EQ(
      mme_config.log_config.is_output_thread_safe,
      amf_config.log_config.is_output_thread_safe);
  EXPECT_EQ(
      mme_config.log_config.mme_app_log_level,
      amf_config.log_config.amf_app_log_level);

  if (mme_config.realm)
    EXPECT_EQ(
        0,
        strcmp((char*) mme_config.realm->data, (char*) amf_config.realm->data));

  if (mme_config.full_network_name)
    EXPECT_EQ(
        0, strcmp(
               (char*) mme_config.full_network_name->data,
               (char*) amf_config.full_network_name->data));

  if (mme_config.short_network_name)
    EXPECT_EQ(
        0, strcmp(
               (char*) mme_config.short_network_name->data,
               (char*) amf_config.short_network_name->data));

  EXPECT_EQ(mme_config.daylight_saving_time, amf_config.daylight_saving_time);
  if (mme_config.pid_dir)
    EXPECT_EQ(
        0, strcmp(
               (char*) mme_config.pid_dir->data,
               (char*) amf_config.pid_dir->data));
  EXPECT_EQ(mme_config.max_enbs, amf_config.max_gnbs);
  EXPECT_EQ(mme_config.relative_capacity, amf_config.relative_capacity);

  EXPECT_EQ(mme_config.use_stateless, amf_config.use_stateless);
  EXPECT_EQ(
      mme_config.unauthenticated_imsi_supported,
      amf_config.unauthenticated_imsi_supported);

  EXPECT_EQ(mme_config.num_par_lists, amf_config.num_par_lists);
  for (uint8_t itr = 0;
       itr < mme_config.num_par_lists && mme_config.partial_list; ++itr) {
    EXPECT_EQ(
        mme_config.partial_list[itr].list_type,
        amf_config.partial_list[itr].list_type);
    EXPECT_EQ(
        mme_config.partial_list[itr].nb_elem,
        amf_config.partial_list[itr].nb_elem);

    for (uint8_t idx = 0; idx < mme_config.partial_list[itr].nb_elem; idx++) {
      if (mme_config.partial_list[itr].plmn &&
          mme_config.partial_list[itr].tac) {
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mcc_digit2,
            amf_config.partial_list[itr].plmn[idx].mcc_digit2);
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mcc_digit1,
            amf_config.partial_list[itr].plmn[idx].mcc_digit1);
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mnc_digit3,
            amf_config.partial_list[itr].plmn[idx].mnc_digit3);
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mcc_digit3,
            amf_config.partial_list[itr].plmn[idx].mcc_digit3);
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mnc_digit2,
            amf_config.partial_list[itr].plmn[idx].mnc_digit2);
        EXPECT_EQ(
            mme_config.partial_list[itr].plmn[idx].mnc_digit1,
            amf_config.partial_list[itr].plmn[idx].mnc_digit1);
        EXPECT_EQ(
            mme_config.partial_list[itr].tac[idx],
            amf_config.partial_list[itr].tac[idx]);
      }
    }
  }

  clear_amf_config(&amf_config);
  free_mme_config(&mme_config);
}
}  // namespace lte
}  // namespace magma

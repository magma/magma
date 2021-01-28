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
/*****************************************************************************

  Source      amf_common_defs.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_COMMON_DEFS_SEEN
#define AMF_COMMON_DEFS_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "common_types.h"
#ifdef __cplusplus
}
#endif
using namespace std;
#define NAS_MESSAGE_SECURITY_HEADER_SIZE 6
#define NAS_MESSAGE_SERVICE_REQUEST_SECURITY_HEADER_SIZE 4
typedef uint32_t context_identifier_t;
typedef uint32_t context_identifier_t;
typedef uint32_t in_address;
typedef uint32_t ard_t;
typedef uint32_t teid_t;
typedef uint32_t sctp_assoc_id_t;
typedef uint32_t gnb_ue_ngap_id_t;
typedef uint64_t gnb_ngap_id_key_t;
typedef uint32_t amf_ue_ngap_id_t;
typedef uint64_t bitrate_t;
typedef uint32_t rau_tau_timer_t;
typedef uint8_t pdn_type_t;
typedef uint8_t proc_pti_t;
typedef uint8_t amf_cause_t;

typedef uint64_t imsi64_t;

#define PDU_SESSION_PER_UE (11)
namespace magma5g {
// need to under stand *****************************
#define OFFSET_OF(TyPe, MeMBeR) ((size_t) & ((TyPe*) 0)->MeMBeR)
#define COUNT_OF(x)                                                            \
  ((sizeof(x) / sizeof(0 [x])) / ((size_t)(!(sizeof(x) % sizeof(0 [x])))))

#define PARENT_STRUCT(cOnTaiNeD, TyPe, MeMBeR)                                 \
  ({                                                                           \
    const typeof(((TyPe*) 0)->MeMBeR)* __MemBeR_ptr = (cOnTaiNeD);             \
    (TyPe*) ((char*) __MemBeR_ptr - OFFSET_OF(TyPe, MeMBeR));                  \
  })

#define OAI_MAX(a, b)                                                          \
  ({                                                                           \
    __typeof__(a) _a = (a);                                                    \
    __typeof__(b) _b = (b);                                                    \
    _a > _b ? _a : _b;                                                         \
  })

#define OAI_MIN(a, b)                                                          \
  ({                                                                           \
    __typeof__(a) _a = (a);                                                    \
    __typeof__(b) _b = (b);                                                    \
    _a < _b ? _a : _b;                                                         \
  })
//------------------------------------------------------------------------------
enum error_code_e {
  /* Fatal errors - received message should not be processed */
  TLV_MAC_MISMATCH                  = -14,
  TLV_BUFFER_NULL                   = -13,
  TLV_BUFFER_TOO_SHORT              = -12,
  TLV_PROTOCOL_NOT_SUPPORTED        = -11,
  TLV_WRONG_MESSAGE_TYPE            = -10,
  TLV_OCTET_STRING_TOO_LONG_FOR_IEI = -9,

  TLV_VALUE_DOESNT_MATCH          = -4,
  TLV_MANDATORY_FIELD_NOT_PRESENT = -3,
  TLV_UNEXPECTED_IEI              = -2,

  RETURNerror = -1,
  RETURNok    = 0,

  TLV_ERROR_OK = RETURNok,
  /* Defines error code limit below which received message should be discarded
   * because it cannot be further processed */
  TLV_FATAL_ERROR = TLV_VALUE_DOESNT_MATCH

};
enum all_apn_conf_ind_t {
  ALL_APN_CONFIGURATIONS_INCLUDED            = 0,
  MODIFIED_ADDED_APN_CONFIGURATIONS_INCLUDED = 1,
  ALL_APN_MAX,
};
class ip_address_t {
 public:
  pdn_type_value_t pdn_type;
  class address {
   public:
    in_address ipv4_address;
    class in6_addr ipv6_address;
  };
};

class ambr_t {
 public:
  bitrate_t br_ul;
  bitrate_t br_dl;
};

class apn_configuration_s {
 public:
  context_identifier_t context_identifier;

  /* Each APN configuration can have 0, 1, or 2 ip address:
   * - 0 means subscribed is dynamically allocated by P-GW depending on the
   * pdn_type
   * - 1 Only one type of IP address is returned by HSS
   * - 2 IPv4 and IPv6 address are returned by HSS and are statically
   * allocated
   */
  uint8_t nb_ip_address;
  ip_address_t ip_address[2];

#ifdef ACCESS_POINT_NAME_MAX_LENGTH
#define SERVICE_SELECTION_MAX_LENGTH ACCESS_POINT_NAME_MAX_LENGTH
#else
#define SERVICE_SELECTION_MAX_LENGTH 100
#endif
  pdn_type_t pdn_type;
  char service_selection[SERVICE_SELECTION_MAX_LENGTH];
  int service_selection_length;
  // eps_subscribed_qos_profile_t subscribed_qos;
  ambr_t ambr;
};

class apn_config_profile_t {
 public:
  context_identifier_t context_identifier;
  all_apn_conf_ind_t all_apn_conf_ind;
  /* Number of APNs provided */
  uint8_t nb_apns;
  /* List of APNs configuration 1 to n elements */
  apn_configuration_s apn_configuration[MAX_APN_PER_UE];
};
class in6_addr {
 public:
  union {
    uint8_t __u6_addr8[16];
    uint16_t __u6_addr16[8];
    uint32_t __u6_addr32[4];
  } __in6_u;
#define s6_addr __in6_u.__u6_addr8
#ifdef __USE_MISC
#define s6_addr16 __in6_u.__u6_addr16
#define s6_addr32 __in6_u.__u6_addr32
#endif
};

class pack_t {
 public:
  pdn_type_value_t pdn_type;
  in_address ipv4_address;
  class in6_addr ipv6_address;
  /* Note in rel.8 the ipv6 prefix length has a fixed value of /64 */
  uint8_t ipv6_prefix_length;
};
/* ESM procedure transaction states */
enum smf_pt_state_e {
  SMF_PROCEDURE_TRANSACTION_INACTIVE = 0,
  SMF_PROCEDURE_TRANSACTION_PENDING,
  SMF_PROCEDURE_TRANSACTION_MAX
};

/*
 * Structure of a PDN connection
 * -----------------------------
 * A PDN connection is the association between a UE represented by
 * one IPv4 address and/or one IPv6 prefix and a PDN represented by
 * an Access Point Name (APN).
 */
class esm_pdn_t {
 public:
  proc_pti_t pti;    /* Identity of the procedure transaction executed
                      * to activate the PDN connection entry     */
  bool is_emergency; /* Emergency bearer services indicator      */
  int ambr;          /* Aggregate Maximum Bit Rate of this APN   */

  int addr_realloc; /* Indicates whether the UE is allowed to subsequently
                     * request another PDN connectivity to the same APN
                     * using an address PDN type (IPv4 or IPv6) other
                     * than the one already activated       */
  int n_pdusession; /* Number of allocated pdu session;*/

  smf_pt_state_e pt_state;  // procedure transaction state
};

class pdn_context_t {
 public:
  context_identifier_t context_identifier;

  /* APN in Use:an ID at UPF through which a user can access the Subscribed APN
   *            This APN shall be composed of the APN Network
   *            Identifier and the default APN Operator Identifier,
   *            as specified in TS 23.003 [9],
   *            clause 9.1.2 (EURECOM: "mnc<MNC>.mcc<MCC>.gprs").
   *            Any received value in the APN OI Replacement field is not
   *            applied here.
   */
  bstring apn_in_use;

  /* APN Subscribed: The subscribed APN received from the HSS */
  bstring apn_subscribed;

  /* PDN Type: IPv4, IPv6 or IPv4v6 */
  pdn_type_t pdn_type;

  /* pack: IPv4 address and/or IPv6 prefix of UE set by
   *          N11 CREATE_SESSION_RESPONSE
   *          NOTE:
   *          The AMF might not have information on the allocated IPv4 address.
   *          Alternatively, following mobility involving a pre-release 8 SGSN,
   *          This IPv4 address might not be the one allocated to the UE.
   */
  pack_t pack;

  /* APN-OI Replacement: APN level APN-OI Replacement which has same role as
   *            UE level APN-OI Replacement but with higher priority than
   *            UE level APN-OI Replacement. This is and optional parameter.
   *            When available, it shall be used to construct the PDN GW
   *            FQDN instead of UE level APN-OI Replacement.
   */
  bstring apn_oi_replacement;

  /* PDN GW Address in Use(control plane): The IP address of the PDN GW
   *           currently
   *           used for sending control plane signalling.
   */
  ip_address_t p_gw_address_n4_cp;

  /* SMF to UPF TEID for  (control plane) */
  teid_t p_gw_teid_n4_cp;

  /* Pdu session subscribed QoS profile:
   *            The pdu session level QoS parameter values for that
   *            APN's default bearer's QCI and ARP (see clause 4.7.3).
   */
  // eps_subscribed_qos_profile_t default_bearer_eps_subscribed_qos_profile;

  /* Subscribed APN-AMBR: The Maximum Aggregated uplink and downlink MBR values
   *            to be shared across all Non-GBR bearers,
   *             which are established for this APN, according to the
   *            subscription of the user.
   */
  // ambr_t subscribed_apn_ambr;

  /* p_gw_apn_ambr: The Maximum Aggregated uplink and downlink MBR values to be
   *           shared across all Non-GBR bearers, which are established for this
   *           APN, as decided by the PDN GW.
   */
  // ambr_t p_gw_apn_ambr;

  /* default_ebi: Identifies the pdu session id Id of the default bearer
   * within the given PDN connection.
   */
  ebi_t default_ebi;

  /* bearer_contexts[]: contains bearer indexes in
   *           ue_m5gmm_context_s.bearer_contexts[], or -1
   */
  int bearer_contexts[PDU_SESSION_PER_UE];

  // set by N11 CREATE_SESSION_RESPONSE

  ip_address_t s_gw_address_s11_s4;
  teid_t s_gw_teid_s11_s4;

  esm_pdn_t esm_data;
  /* is_active == true indicates, pdu session is active */
  bool is_active;

  protocol_configuration_options_t* pco;
};

}  // namespace magma5g
#endif

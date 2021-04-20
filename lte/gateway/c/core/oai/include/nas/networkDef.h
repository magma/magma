/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************

Source      networkDef.h

Version     0.1

Date        2012/09/21

Product     NAS stack

Subsystem   include

Author      Frederic Maurel, Lionel GAUTHIER

Description Contains network's global definitions

*****************************************************************************/
#ifndef FILE_NETWORKDEF_SEEN
#define FILE_NETWORKDEF_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * ----------------------
 * Network selection mode
 * ----------------------
 */
#define NET_PLMN_AUTO 0
#define NET_PLMN_MANUAL 1

/*
 * ---------------------------
 * Network registration status
 * ---------------------------
 */
/* not registered, not currently searching an operator to register to */
#define NET_REG_STATE_OFF 0
/* registered, home network                       */
#define NET_REG_STATE_HN 1
/* not registered, currently trying to attach or searching an operator
 * to register to                         */
#define NET_REG_STATE_ON 2
/* registration denied                        */
#define NET_REG_STATE_DENIED 3
/* unknown (e.g. out of GERAN/UTRAN/E-UTRAN coverage)         */
#define NET_REG_STATE_UNKNOWN 4
/* registered, roaming                        */
#define NET_REG_STATE_ROAMING 5
/* registered for "SMS only", home network                */
#define NET_REG_STATE_SMS_HN 6
/* registered, for "SMS only", roaming                */
#define NET_REG_STATE_SMS_ROAMING 7
/* attached for emergency bearer services only (applicable to UTRAN)  */
#define NET_REG_STATE_EMERGENCY 8

/*
 * ------------------------------------
 * Network access technology indicators
 * ------------------------------------
 */
#define NET_ACCESS_UNAVAILABLE (-1) /* Not available        */
#define NET_ACCESS_GSM 0            /* GSM              */
#define NET_ACCESS_COMPACT 1        /* GSM Compact          */
#define NET_ACCESS_UTRAN 2          /* UTRAN            */
#define NET_ACCESS_EGPRS 3          /* GSM w/EGPRS          */
#define NET_ACCESS_HSDPA 4          /* UTRAN w/HSDPA        */
#define NET_ACCESS_HSUPA 5          /* UTRAN w/HSUPA        */
#define NET_ACCESS_HSDUPA 6         /* UTRAN w/HSDPA and HSUPA  */
#define NET_ACCESS_EUTRAN 7         /* E-UTRAN          */

/*
 * ---------------------------------------
 * Network operator representation formats
 * ---------------------------------------
 */
#define NET_FORMAT_LONG 0  /* long format alphanumeric */
#define NET_FORMAT_SHORT 1 /* short format alphanumeric    */
#define NET_FORMAT_NUM 2   /* numeric format       */

#define NET_FORMAT_MAX_SIZE NET_FORMAT_LONG_SIZE

/*
 * -----------------------------
 * Network operator availability
 * -----------------------------
 */
#define NET_OPER_UNKNOWN 0   /* unknown operator     */
#define NET_OPER_AVAILABLE 1 /* available operator       */
#define NET_OPER_CURRENT 2   /* currently selected operator  */
#define NET_OPER_FORBIDDEN 3 /* forbidden operator       */

/*
 * --------------------------------------
 * Network connection establishment cause
 * --------------------------------------
 */
#define NET_ESTABLISH_CAUSE_EMERGENCY 0x01
#define NET_ESTABLISH_CAUSE_HIGH_PRIO 0x02
#define NET_ESTABLISH_CAUSE_MT_ACCESS 0x03
#define NET_ESTABLISH_CAUSE_MO_SIGNAL 0x04
#define NET_ESTABLISH_CAUSE_MO_DATA 0x05
#define NET_ESTABLISH_CAUSE_V1020 0x06

/*
 * --------------------------------------
 * Network connection establishment type
 * --------------------------------------
 */
#define NET_ESTABLISH_TYPE_ORIGINATING_SIGNAL 0x10
#define NET_ESTABLISH_TYPE_EMERGENCY_CALLS 0x20
#define NET_ESTABLISH_TYPE_ORIGINATING_CALLS 0x30
#define NET_ESTABLISH_TYPE_TERMINATING_CALLS 0x40
#define NET_ESTABLISH_TYPE_MO_CS_FALLBACK 0x50

/*
 * -------------------
 * PDN connection type
 * -------------------
 */
#define NET_PDN_TYPE_IPV4 (0 + 1)
#define NET_PDN_TYPE_IPV6 (1 + 1)
#define NET_PDN_TYPE_IPV4V6 (2 + 1)

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * ---------------------
 * PDN connection status
 * ---------------------
 */
typedef enum {
  /* MT = The Mobile Terminal, NW = The Network               */
  NET_PDN_MT_DEFAULT_ACT = 1, /* MT has activated a PDN connection        */
  NET_PDN_NW_DEFAULT_DEACT,   /* NW has deactivated a PDN connection      */
  NET_PDN_MT_DEFAULT_DEACT,   /* MT has deactivated a PDN connection      */
  NET_PDN_NW_DEDICATED_ACT,   /* NW has activated an EPS bearer context   */
  NET_PDN_MT_DEDICATED_ACT,   /* MT has activated an EPS bearer context   */
  NET_PDN_NW_DEDICATED_DEACT, /* NW has deactivated an EPS bearer context */
  NET_PDN_MT_DEDICATED_DEACT, /* MT has deactivated an EPS bearer context */
} network_pdn_state_t;

/*
 * ---------------------------
 * Network operator identifier
 * ---------------------------
 */
typedef struct {
#define NET_FORMAT_LONG_SIZE 16 /* Long alphanumeric format     */
#define NET_FORMAT_SHORT_SIZE 8 /* Short alphanumeric format        */
#define NET_FORMAT_NUM_SIZE 6   /* Numeric format (PLMN identifier  */
  union {
    unsigned char alpha_long[NET_FORMAT_LONG_SIZE + 1];
    unsigned char alpha_short[NET_FORMAT_SHORT_SIZE + 1];
    unsigned char num[NET_FORMAT_NUM_SIZE + 1];
  } id;
} network_plmn_t;

/*
 * -------------------------------
 * EPS bearer level QoS parameters
 * -------------------------------
 */
typedef struct {
  int gbrUL; /* Guaranteed Bit Rate for uplink   */
  int gbrDL; /* Guaranteed Bit Rate for downlink */
  int mbrUL; /* Maximum Bit Rate for uplink      */
  int mbrDL; /* Maximum Bit Rate for downlink    */
  int qci;   /* QoS Class Identifier         */
} network_qos_t;

/*
 * -----------------------------
 * IPv4 packet filter parameters
 * -----------------------------
 */
typedef struct {
  unsigned char protocol; /* Protocol identifier      */
  unsigned char tos;      /* Type of service      */
#define NET_PACKET_FILTER_IPV4_ADDR_SIZE 4
  unsigned char addr[NET_PACKET_FILTER_IPV4_ADDR_SIZE];
  unsigned char mask[NET_PACKET_FILTER_IPV4_ADDR_SIZE];
} network_ipv4_data_t;

/*
 * -----------------------------
 * IPv6 packet filter parameters
 * -----------------------------
 */
typedef struct {
  unsigned char nh; /* Next header type     */
  unsigned char tf; /* Traffic class        */
#define NET_PACKET_FILTER_IPV6_ADDR_SIZE 16
  unsigned char addr[NET_PACKET_FILTER_IPV6_ADDR_SIZE];
  unsigned char mask[NET_PACKET_FILTER_IPV6_ADDR_SIZE];
  unsigned int ipsec; /* IPSec security parameter index */
  unsigned int fl;    /* Flow label             */
} network_ipv6_data_t;

/*
 * -------------
 * Packet Filter
 * -------------
 */
typedef struct {
  unsigned char id; /* Packet filter identifier */
#define NET_PACKET_FILTER_DOWNLINK 0x01
#define NET_PACKET_FILTER_UPLINK 0x02
#define NET_PACKET_FILTER_BIDIR 0x03
  unsigned char dir;        /* Packet filter direction  */
  unsigned char precedence; /* Evaluation precedence    */
  union {
    network_ipv4_data_t ipv4;
    network_ipv6_data_t ipv6;
  } data;
  unsigned short lport; /* Local (UE) port number   */
  unsigned short rport; /* Remote (network) port number */
} network_pkf_t;

/*
 * ---------------------
 * Traffic Flow Template
 * ---------------------
 */
typedef struct {
  int n_pkfs;
#define NET_PACKET_FILTER_MAX 16
  network_pkf_t* pkf[NET_PACKET_FILTER_MAX];
} network_tft_t;

/*
 * User notification callback, executed whenever a change of status with
 * respect of PDN connection or EPS bearer context is notified by the EPS
 * Session Management sublayer
 */
typedef int (*esm_indication_callback_t)(int, network_pdn_state_t);

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* FILE_NETWORKDEF_SEEN*/

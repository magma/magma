/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

#ifndef FILE_3GPP_24_301_SEEN
#define FILE_3GPP_24_301_SEEN

#include <stdint.h>
#include <stdbool.h>

//==============================================================================
// 9 General message format and information elements coding
//==============================================================================

//------------------------------------------------------------------------------
// 9.2 Protocol discriminator
//------------------------------------------------------------------------------

// 9.3.1 Security header type
#define SECURITY_HEADER_TYPE_NOT_PROTECTED 0b0000
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED 0b0001
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED 0b0010
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW 0b0011
#define SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW 0b0100
#define SECURITY_HEADER_TYPE_SERVICE_REQUEST 0b1100
#define SECURITY_HEADER_TYPE_RESERVED1 0b1101
#define SECURITY_HEADER_TYPE_RESERVED2 0b1110
#define SECURITY_HEADER_TYPE_RESERVED3 0b1111

// 9.3.2 EPS bearer identity
// see 24.007

//------------------------------------------------------------------------------
// 9.8 Message type
//------------------------------------------------------------------------------
// Table 9.8.1: Message types for EPS mobility management
/* Message identifiers for EPS Mobility Management     */
#define ATTACH_REQUEST 0b01000001                 /* 65 = 0x41 */
#define ATTACH_ACCEPT 0b01000010                  /* 66 = 0x42 */
#define ATTACH_COMPLETE 0b01000011                /* 67 = 0x43 */
#define ATTACH_REJECT 0b01000100                  /* 68 = 0x44 */
#define DETACH_REQUEST 0b01000101                 /* 69 = 0x45 */
#define DETACH_ACCEPT 0b01000110                  /* 70 = 0x46 */
#define TRACKING_AREA_UPDATE_REQUEST 0b01001000   /* 72 = 0x48 */
#define TRACKING_AREA_UPDATE_ACCEPT 0b01001001    /* 73 = 0x49 */
#define TRACKING_AREA_UPDATE_COMPLETE 0b01001010  /* 74 = 0x4a */
#define TRACKING_AREA_UPDATE_REJECT 0b01001011    /* 75 = 0x4b */
#define EXTENDED_SERVICE_REQUEST 0b01001100       /* 76 = 0x4c */
#define SERVICE_REJECT 0b01001110                 /* 78 = 0x4e */
#define GUTI_REALLOCATION_COMMAND 0b01010000      /* 80 = 0x50 */
#define GUTI_REALLOCATION_COMPLETE 0b01010001     /* 81 = 0x51 */
#define AUTHENTICATION_REQUEST 0b01010010         /* 82 = 0x52 */
#define AUTHENTICATION_RESPONSE 0b01010011        /* 83 = 0x53 */
#define AUTHENTICATION_REJECT 0b01010100          /* 84 = 0x54 */
#define AUTHENTICATION_FAILURE 0b01011100         /* 92 = 0x5c */
#define IDENTITY_REQUEST 0b01010101               /* 85 = 0x55 */
#define IDENTITY_RESPONSE 0b01010110              /* 86 = 0x56 */
#define SECURITY_MODE_COMMAND 0b01011101          /* 93 = 0x5d */
#define SECURITY_MODE_COMPLETE 0b01011110         /* 94 = 0x5e */
#define SECURITY_MODE_REJECT 0b01011111           /* 95 = 0x5f */
#define EMM_STATUS 0b01100000                     /* 96 = 0x60 */
#define EMM_INFORMATION 0b01100001                /* 97 = 0x61 */
#define DOWNLINK_NAS_TRANSPORT 0b01100010         /* 98 = 0x62 */
#define UPLINK_NAS_TRANSPORT 0b01100011           /* 99 = 0x63 */
#define CS_SERVICE_NOTIFICATION 0b01100100        /* 100 = 0x64 */
#define DOWNLINK_GENERIC_NAS_TRANSPORT 0b01101000 /* 104 = 0x68 */
#define UPLINK_GENERIC_NAS_TRANSPORT 0b01101001   /* 101 = 0x69 */

// Table 9.8.2: Message types for EPS session management
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST                            \
  0b11000001                                                  /* 193 = 0xc1    \
                                                               */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT 0b11000010 /* 194 = 0xc2 */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REJECT 0b11000011 /* 195 = 0xc3 */
#define ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST                          \
  0b11000101 /* 197 = 0xc5 */
#define ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_ACCEPT                           \
  0b11000110 /* 198 = 0xc6                                                     \
              */
#define ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REJECT                           \
  0b11000111                                             /* 199 = 0xc7         \
                                                          */
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST 0b11001001     /* 201 = 0xc9 */
#define MODIFY_EPS_BEARER_CONTEXT_ACCEPT 0b11001010      /* 202 = 0xca */
#define MODIFY_EPS_BEARER_CONTEXT_REJECT 0b11001011      /* 203 = 0xcb */
#define DEACTIVATE_EPS_BEARER_CONTEXT_REQUEST 0b11001101 /* 205 = 0xcd */
#define DEACTIVATE_EPS_BEARER_CONTEXT_ACCEPT 0b11001110  /* 206 = 0xce */
#define PDN_CONNECTIVITY_REQUEST 0b11010000              /* 208 = 0xd0 */
#define PDN_CONNECTIVITY_REJECT 0b11010001               /* 209 = 0xd1 */
#define PDN_DISCONNECT_REQUEST 0b11010010                /* 210 = 0xd2 */
#define PDN_DISCONNECT_REJECT 0b11010011                 /* 211 = 0xd3 */
#define BEARER_RESOURCE_ALLOCATION_REQUEST 0b11010100    /* 212 = 0xd4 */
#define BEARER_RESOURCE_ALLOCATION_REJECT 0b11010101     /* 213 = 0xd5 */
#define BEARER_RESOURCE_MODIFICATION_REQUEST 0b11010110  /* 214 = 0xd6 */
#define BEARER_RESOURCE_MODIFICATION_REJECT 0b11010111   /* 215 = 0xd7 */
#define ESM_INFORMATION_REQUEST 0b11011001               /* 217 = 0xd9 */
#define ESM_INFORMATION_RESPONSE 0b11011010              /* 218 = 0xda */
#define ESM_STATUS 0b11101000                            /* 232 = 0xe8 */

//------------------------------------------------------------------------------
// 9.9 OTHER INFORMATION ELEMENTS
//------------------------------------------------------------------------------

//..............................................................................
// 9.9.3 EPS Mobility Management (EMM) information elements
//..............................................................................

// 9.9.3.34 UE network capability
#define UE_NETWORK_CAPABILITY_MINIMUM_LENGTH 4
#define UE_NETWORK_CAPABILITY_MAXIMUM_LENGTH 15

typedef struct ue_network_capability_s {
  /* EPS encryption algorithms supported (octet 3) */
#define UE_NETWORK_CAPABILITY_EEA0 0b10000000
#define UE_NETWORK_CAPABILITY_EEA1 0b01000000
#define UE_NETWORK_CAPABILITY_EEA2 0b00100000
#define UE_NETWORK_CAPABILITY_EEA3 0b00010000
#define UE_NETWORK_CAPABILITY_EEA4 0b00001000
#define UE_NETWORK_CAPABILITY_EEA5 0b00000100
#define UE_NETWORK_CAPABILITY_EEA6 0b00000010
#define UE_NETWORK_CAPABILITY_EEA7 0b00000001
  uint8_t eea;
  /* EPS integrity algorithms supported (octet 4) */
#define UE_NETWORK_CAPABILITY_EIA0 0b10000000
#define UE_NETWORK_CAPABILITY_EIA1 0b01000000
#define UE_NETWORK_CAPABILITY_EIA2 0b00100000
#define UE_NETWORK_CAPABILITY_EIA3 0b00010000
#define UE_NETWORK_CAPABILITY_EIA4 0b00001000
#define UE_NETWORK_CAPABILITY_EIA5 0b00000100
#define UE_NETWORK_CAPABILITY_EIA6 0b00000010
#define UE_NETWORK_CAPABILITY_EIA7 0b00000001
  uint8_t eia;
  /* UMTS encryption algorithms supported (octet 5) */
#define UE_NETWORK_CAPABILITY_UEA0 0b10000000
#define UE_NETWORK_CAPABILITY_UEA1 0b01000000
#define UE_NETWORK_CAPABILITY_UEA2 0b00100000
#define UE_NETWORK_CAPABILITY_UEA3 0b00010000
#define UE_NETWORK_CAPABILITY_UEA4 0b00001000
#define UE_NETWORK_CAPABILITY_UEA5 0b00000100
#define UE_NETWORK_CAPABILITY_UEA6 0b00000010
#define UE_NETWORK_CAPABILITY_UEA7 0b00000001
  uint8_t uea;
  /* UCS2 support (octet 6, bit 8) */
#define UE_NETWORK_CAPABILITY_DEFAULT_ALPHABET 0
#define UE_NETWORK_CAPABILITY_UCS2_ALPHABET 1
  uint8_t ucs2 : 1;
  /* UMTS integrity algorithms supported (octet 6) */
#define UE_NETWORK_CAPABILITY_UIA1 0b01000000
#define UE_NETWORK_CAPABILITY_UIA2 0b00100000
#define UE_NETWORK_CAPABILITY_UIA3 0b00010000
#define UE_NETWORK_CAPABILITY_UIA4 0b00001000
#define UE_NETWORK_CAPABILITY_UIA5 0b00000100
#define UE_NETWORK_CAPABILITY_UIA6 0b00000010
#define UE_NETWORK_CAPABILITY_UIA7 0b00000001
  uint8_t uia : 7;

  /*Octet 7*/
  uint8_t prosedd : 1;
  uint8_t prose : 1;
  uint8_t h245ash : 1;
  /* eNodeB-based access class control for CSFB capability */
#define UE_NETWORK_CAPABILITY_CSFB 1
  uint8_t csfb : 1;
  /* LTE Positioning Protocol capability */
#define UE_NETWORK_CAPABILITY_LPP 1
  uint8_t lpp : 1;
  /* Location services notification mechanisms capability */
#define UE_NETWORK_CAPABILITY_LCS 1
  uint8_t lcs : 1;
  /* 1xSRVCC capability */
#define UE_NETWORK_CAPABILITY_SRVCC 1
  uint8_t srvcc : 1;
  /* NF notification procedure capability */
#define UE_NETWORK_CAPABILITY_NF 1
  uint8_t nf : 1;

  /*Octet 8*/
  uint8_t epco : 1;
  uint8_t hccpciot : 1;
  uint8_t erwfopdn : 1;
  /* S1-U data transfer enabled */
#define UE_NETWORK_CAPABILITY_S1UDATA 1
  uint8_t s1udata : 1;
  /* User plane CIoT EPS Optimization */
#define UE_NETWORK_CAPABILITY_UPCIOT 1
  uint8_t upciot : 1;
  /* Control plane CIoT EPS Optimization */
#define UE_NETWORK_CAPABILITY_CPCIOT 1
  uint8_t cpciot : 1;
  /* ProSe UE-to-network-relay */
#define UE_NETWORK_CAPABILITY_PROSERELAY 1
  uint8_t proserelay : 1;
  /* ProSe direct communication */
#define UE_NETWORK_CAPABILITY_PROSEDC 1
  uint8_t prosedc : 1;

  /*Octet 9*/
  uint8_t bearer : 1;
  uint8_t sgc : 1;
  uint8_t n1mod : 1;
  /* DCNR notification flag */
#define UE_NETWORK_CAPABILITY_DCNR 1
  uint8_t dcnr : 1;
  /* Control Plane data backoff support*/
#define UE_NETWORK_CAPABILITY_CPBACKOFF 1
  uint8_t cpbackoff : 1;
  /* Restriction o nuse of enhanced coversage support */
#define UE_NETWORK_CAPABILITY_RESTRICTEC 1
  uint8_t restrictec : 1;
  /* V2X communication ovre PC5 */
#define UE_NETWORK_CAPABILITY_V2XPC5 1
  uint8_t v2xpc5 : 1;
  /* Multiple DRB support */
#define UE_NETWORK_CAPABILITY_MULTIPLEDRB 1
  uint8_t multipledrb : 1;

  bool umts_present;
  uint32_t length;
} ue_network_capability_t;

// 9.9.3.36 UE security capability
#define UE_SECURITY_CAPABILITY_MINIMUM_LENGTH 4
#define UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH 7

typedef struct ue_security_capability_s {
/* EPS encryption algorithms supported (octet 3) */
#define UE_SECURITY_CAPABILITY_EEA0 0b10000000
#define UE_SECURITY_CAPABILITY_EEA1 0b01000000
#define UE_SECURITY_CAPABILITY_EEA2 0b00100000
#define UE_SECURITY_CAPABILITY_EEA3 0b00010000
#define UE_SECURITY_CAPABILITY_EEA4 0b00001000
#define UE_SECURITY_CAPABILITY_EEA5 0b00000100
#define UE_SECURITY_CAPABILITY_EEA6 0b00000010
#define UE_SECURITY_CAPABILITY_EEA7 0b00000001
  uint8_t eea;
  /* EPS integrity algorithms supported (octet 4) */
#define UE_SECURITY_CAPABILITY_EIA0 0b10000000
#define UE_SECURITY_CAPABILITY_EIA1 0b01000000
#define UE_SECURITY_CAPABILITY_EIA2 0b00100000
#define UE_SECURITY_CAPABILITY_EIA3 0b00010000
#define UE_SECURITY_CAPABILITY_EIA4 0b00001000
#define UE_SECURITY_CAPABILITY_EIA5 0b00000100
#define UE_SECURITY_CAPABILITY_EIA6 0b00000010
#define UE_SECURITY_CAPABILITY_EIA7 0b00000001
  uint8_t eia;
  bool umts_present;
  bool gprs_present;
  /* UMTS encryption algorithms supported (octet 5) */
#define UE_SECURITY_CAPABILITY_UEA0 0b10000000
#define UE_SECURITY_CAPABILITY_UEA1 0b01000000
#define UE_SECURITY_CAPABILITY_UEA2 0b00100000
#define UE_SECURITY_CAPABILITY_UEA3 0b00010000
#define UE_SECURITY_CAPABILITY_UEA4 0b00001000
#define UE_SECURITY_CAPABILITY_UEA5 0b00000100
#define UE_SECURITY_CAPABILITY_UEA6 0b00000010
#define UE_SECURITY_CAPABILITY_UEA7 0b00000001
  uint8_t uea;
  /* UMTS integrity algorithms supported (octet 6) */
#define UE_SECURITY_CAPABILITY_UIA1 0b01000000
#define UE_SECURITY_CAPABILITY_UIA2 0b00100000
#define UE_SECURITY_CAPABILITY_UIA3 0b00010000
#define UE_SECURITY_CAPABILITY_UIA4 0b00001000
#define UE_SECURITY_CAPABILITY_UIA5 0b00000100
#define UE_SECURITY_CAPABILITY_UIA6 0b00000010
#define UE_SECURITY_CAPABILITY_UIA7 0b00000001
  uint8_t uia : 7;
  /* GPRS encryption algorithms supported (octet 7) */
#define UE_SECURITY_CAPABILITY_GEA1 0b01000000
#define UE_SECURITY_CAPABILITY_GEA2 0b00100000
#define UE_SECURITY_CAPABILITY_GEA3 0b00010000
#define UE_SECURITY_CAPABILITY_GEA4 0b00001000
#define UE_SECURITY_CAPABILITY_GEA5 0b00000100
#define UE_SECURITY_CAPABILITY_GEA6 0b00000010
#define UE_SECURITY_CAPABILITY_GEA7 0b00000001
  uint8_t gea : 7;
} ue_security_capability_t;

// 9.9.3.53 UE additional security capability
typedef struct ue_additional_security_capability_s {
  /* NR encryption algorithms supported  */
  uint16_t _5g_ea;
  /* NR integrity algorithms supported */
  uint16_t _5g_ia;
} ue_additional_security_capability_t;

// 9.2.1.127 NR UE security capability
#define NR_UE_SECURITY_CAPABILITY_MINIMUM_LENGTH 4
#define NR_UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH 7

typedef struct nr_ue_security_capability_s {
/* NR encryption algorithms supported (octet 1) */
#define UE_SECURITY_CAPABILITY_NEA0 0b10000000
#define UE_SECURITY_CAPABILITY_NEA1 0b01000000
#define UE_SECURITY_CAPABILITY_NEA2 0b00100000
#define UE_SECURITY_CAPABILITY_NEA3 0b00010000
#define UE_SECURITY_CAPABILITY_NEA4 0b00001000
#define UE_SECURITY_CAPABILITY_NEA5 0b00000100
#define UE_SECURITY_CAPABILITY_NEA6 0b00000010
#define UE_SECURITY_CAPABILITY_NEA7 0b00000001
  uint8_t nea;
  /* NR integrity algorithms supported (octet 2) */
#define UE_SECURITY_CAPABILITY_NIA0 0b10000000
#define UE_SECURITY_CAPABILITY_NIA1 0b01000000
#define UE_SECURITY_CAPABILITY_NIA2 0b00100000
#define UE_SECURITY_CAPABILITY_NIA3 0b00010000
#define UE_SECURITY_CAPABILITY_NIA4 0b00001000
#define UE_SECURITY_CAPABILITY_NIA5 0b00000100
#define UE_SECURITY_CAPABILITY_NIA6 0b00000010
#define UE_SECURITY_CAPABILITY_NIA7 0b00000001
  uint8_t nia;
  bool nr_present;
} nr_ue_security_capability_t;

//------------------------------------------------------------------------------
// 10.2 Timers of EPS mobility management
//------------------------------------------------------------------------------

//..............................................................................
// Table 10.2.1: EPS mobility management timers – UE side
//..............................................................................

#define T3402_DEFAULT_VALUE 720
#define T3410_DEFAULT_VALUE 15
#define T3411_DEFAULT_VALUE 10
#define T3412_DEFAULT_VALUE 3240
#define T3416_DEFAULT_VALUE 30
#define T3417_DEFAULT_VALUE 5
#define T3417_EXT_DEFAULT_VALUE 10
#define T3420_DEFAULT_VALUE 15
#define T3421_DEFAULT_VALUE 15
#define T3423_DEFAULT_VALUE 0  // value provided by network
#define T3440_DEFAULT_VALUE 10
#define T3442_DEFAULT_VALUE 0  // value provided by network

//..............................................................................
// Table 10.2.2: EPS mobility management timers – network side
//..............................................................................
#define T3413_DEFAULT_VALUE 400 /* Network dependent    */
#define T3422_DEFAULT_VALUE 6
#define T3450_DEFAULT_VALUE 6
#define T3460_DEFAULT_VALUE 6
#define T3470_DEFAULT_VALUE 6

//------------------------------------------------------------------------------
// 10.3 Timers of EPS session management
//------------------------------------------------------------------------------

//..............................................................................
// Table 10.3.1: EPS session management timers – UE side
//..............................................................................
#define T3480_DEFAULT_VALUE 8
#define T3481_DEFAULT_VALUE 8
#define T3482_DEFAULT_VALUE 8
#define T3492_DEFAULT_VALUE 6

//..............................................................................
// Table 10.3.2: EPS session management timers – network side
//..............................................................................
#define T3485_DEFAULT_VALUE 8
#define T3486_DEFAULT_VALUE 8
#define T3489_DEFAULT_VALUE 4
#define T3495_DEFAULT_VALUE 8

//==============================================================================
// Annex A (informative): Cause values for EPS mobility management
//==============================================================================

//------------------------------------------------------------------------------
// A.1 Causes related to UE identification
//------------------------------------------------------------------------------
#define EMM_CAUSE_IMSI_UNKNOWN_IN_HSS 2
#define EMM_CAUSE_ILLEGAL_UE 3
#define EMM_CAUSE_ILLEGAL_ME 6
#define EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW 9
#define EMM_CAUSE_IMPLICITLY_DETACHED 10

//------------------------------------------------------------------------------
// A.2 Cause related to subscription options
//------------------------------------------------------------------------------
#define EMM_CAUSE_IMEI_NOT_ACCEPTED 5
#define EMM_CAUSE_EPS_NOT_ALLOWED 7
#define EMM_CAUSE_BOTH_NOT_ALLOWED 8
#define EMM_CAUSE_PLMN_NOT_ALLOWED 11
#define EMM_CAUSE_TA_NOT_ALLOWED 12
#define EMM_CAUSE_ROAMING_NOT_ALLOWED 13
#define EMM_CAUSE_EPS_NOT_ALLOWED_IN_PLMN 14
#define EMM_CAUSE_NO_SUITABLE_CELLS 15
#define EMM_CAUSE_CSG_NOT_AUTHORIZED 25
#define EMM_CAUSE_NOT_AUTHORIZED_IN_PLMN 35
#define EMM_CAUSE_NO_EPS_BEARER_CTX_ACTIVE 40

//------------------------------------------------------------------------------
// A.3 Causes related to PLMN specific network failures and
// congestion/authentication failures
//------------------------------------------------------------------------------
#define EMM_CAUSE_MSC_NOT_REACHABLE 16
#define EMM_CAUSE_NETWORK_FAILURE 17
#define EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE 18
#define EMM_CAUSE_ESM_FAILURE 19
#define EMM_CAUSE_MAC_FAILURE 20
#define EMM_CAUSE_SYNCH_FAILURE 21
#define EMM_CAUSE_CONGESTION 22
#define EMM_CAUSE_UE_SECURITY_MISMATCH 23
#define EMM_CAUSE_SECURITY_MODE_REJECTED 24
#define EMM_CAUSE_NON_EPS_AUTH_UNACCEPTABLE 26
#define EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE 39

//------------------------------------------------------------------------------
// A.4 Causes related to nature of request
//------------------------------------------------------------------------------
// NOTE: This subclause has no entries in this version of the specification

//------------------------------------------------------------------------------
// A.5 Causes related to invalid messages
//------------------------------------------------------------------------------
#define EMM_CAUSE_SEMANTICALLY_INCORRECT 95
#define EMM_CAUSE_INVALID_MANDATORY_INFO 96
#define AMF_CAUSE_INVALID_MANDATORY_INFO 96
#define EMM_CAUSE_MESSAGE_TYPE_NOT_IMPLEMENTED 97
#define EMM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE 98
#define EMM_CAUSE_IE_NOT_IMPLEMENTED 99
#define AMF_CAUSE_IE_NOT_IMPLEMENTED 99
#define EMM_CAUSE_CONDITIONAL_IE_ERROR 100
#define EMM_CAUSE_MESSAGE_NOT_COMPATIBLE 101
#define EMM_CAUSE_PROTOCOL_ERROR 111
#define AMF_CAUSE_PROTOCOL_ERROR 111

//==============================================================================
// Annex B (informative): Cause values for EPS session management
//==============================================================================

//------------------------------------------------------------------------------
// B.1 Causes related to nature of request
//------------------------------------------------------------------------------
#define ESM_CAUSE_OPERATOR_DETERMINED_BARRING 8
#define ESM_CAUSE_INSUFFICIENT_RESOURCES 26
#define ESM_CAUSE_UNKNOWN_ACCESS_POINT_NAME 27
#define ESM_CAUSE_UNKNOWN_PDN_TYPE 28
#define ESM_CAUSE_USER_AUTHENTICATION_FAILED 29
#define ESM_CAUSE_REQUEST_REJECTED_BY_GW 30
#define ESM_CAUSE_REQUEST_REJECTED_UNSPECIFIED 31
#define ESM_CAUSE_SERVICE_OPTION_NOT_SUPPORTED 32
#define ESM_CAUSE_REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED 33
#define ESM_CAUSE_SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER 34
#define ESM_CAUSE_PTI_ALREADY_IN_USE 35
#define ESM_CAUSE_REGULAR_DEACTIVATION 36
#define ESM_CAUSE_EPS_QOS_NOT_ACCEPTED 37
#define ESM_CAUSE_NETWORK_FAILURE 38
#define ESM_CAUSE_REACTIVATION_REQUESTED 39
#define ESM_CAUSE_SEMANTIC_ERROR_IN_THE_TFT_OPERATION 41
#define ESM_CAUSE_SYNTACTICAL_ERROR_IN_THE_TFT_OPERATION 42
#define ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY 43
#define ESM_CAUSE_SEMANTIC_ERRORS_IN_PACKET_FILTER 44
#define ESM_CAUSE_SYNTACTICAL_ERROR_IN_PACKET_FILTER 45
#define ESM_CAUSE_PTI_MISMATCH 47
#define ESM_CAUSE_LAST_PDN_DISCONNECTION_NOT_ALLOWED 49
#define ESM_CAUSE_PDN_TYPE_IPV4_ONLY_ALLOWED 50
#define ESM_CAUSE_PDN_TYPE_IPV6_ONLY_ALLOWED 51
#define ESM_CAUSE_SINGLE_ADDRESS_BEARERS_ONLY_ALLOWED 52
#define ESM_CAUSE_ESM_INFORMATION_NOT_RECEIVED 53
#define ESM_CAUSE_PDN_CONNECTION_DOES_NOT_EXIST 54
#define ESM_CAUSE_MULTIPLE_PDN_CONNECTIONS_NOT_ALLOWED 55
#define ESM_CAUSE_COLLISION_WITH_NETWORK_INITIATED_REQUEST 56
#define ESM_CAUSE_UNSUPPORTED_QCI_VALUE 59
#define ESM_CAUSE_BEARER_HANDLING_NOT_SUPPORTED 60
#define ESM_CAUSE_INVALID_PTI_VALUE 81
#define ESM_CAUSE_APN_RESTRICTION_VALUE_NOT_COMPATIBLE 112
#define ESM_CAUSE_REQUESTED_APN_NOT_SUPPORTED_IN_CURRENT_RAT 66

//------------------------------------------------------------------------------
// B.2 Protocol errors (e.g., unknown message) class
//------------------------------------------------------------------------------
#define ESM_CAUSE_SEMANTICALLY_INCORRECT 95
#define ESM_CAUSE_INVALID_MANDATORY_INFO 96
#define AMF_CAUSE_INVALID_MANDATORY_INFO 96
#define ESM_CAUSE_MESSAGE_TYPE_NOT_IMPLEMENTED 97
#define ESM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE 98
#define ESM_CAUSE_IE_NOT_IMPLEMENTED 99
#define ESM_CAUSE_CONDITIONAL_IE_ERROR 100
#define ESM_CAUSE_MESSAGE_NOT_COMPATIBLE 101
#define ESM_CAUSE_PROTOCOL_ERROR 111

#endif /* FILE_3GPP_24_301_SEEN */

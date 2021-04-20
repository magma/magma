/*----------------------------------------------------------------------------*
 *                                                                            *
 *                              n w - g t p v 2 c                             *
 *    G P R S   T u n n e l i n g    P r o t o c o l   v 2 c    S t a c k     *
 *                                                                            *
 *                                                                            *
 * Copyright (c) 2010-2011 Amit Chawre                                        *
 * All rights reserved.                                                       *
 *                                                                            *
 * Redistribution and use in source and binary forms, with or without         *
 * modification, are permitted provided that the following conditions         *
 * are met:                                                                   *
 *                                                                            *
 * 1. Redistributions of source code must retain the above copyright          *
 *    notice, this list of conditions and the following disclaimer.           *
 * 2. Redistributions in binary form must reproduce the above copyright       *
 *    notice, this list of conditions and the following disclaimer in the     *
 *    documentation and/or other materials provided with the distribution.    *
 * 3. The name of the author may not be used to endorse or promote products   *
 *    derived from this software without specific prior written permission.   *
 *                                                                            *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR       *
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES  *
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.    *
 * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,           *
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT   *
 * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,  *
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY      *
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT        *
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF   *
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.          *
 *----------------------------------------------------------------------------*/

/**
 * @file NwGtpv2cIe.h
 * @brief This header file contains Information Element definitions for GTPv2c
 * as per 3GPP TS 29274-930.
 */

#ifndef __NW_GTPV2C_IE_H__
#define __NW_GTPV2C_IE_H__

/*--------------------------------------------------------------------------*
 *   G T P V 2 C    I E  T Y P E     M A C R O     D E F I N I T I O N S    *
 *--------------------------------------------------------------------------*/

#define NW_GTPV2C_IE_RESERVED (0)
#define NW_GTPV2C_IE_IMSI (1)
#define NW_GTPV2C_IE_CAUSE (2)
#define NW_GTPV2C_IE_RECOVERY (3)
#define NW_GTPV2C_IE_APN (71)
#define NW_GTPV2C_IE_AMBR (72)
#define NW_GTPV2C_IE_EBI (73)
#define NW_GTPV2C_IE_IP_ADDRESS (74)
#define NW_GTPV2C_IE_MEI (75)
#define NW_GTPV2C_IE_MSISDN (76)
#define NW_GTPV2C_IE_INDICATION (77)
#define NW_GTPV2C_IE_PCO (78)
#define NW_GTPV2C_IE_PAA (79)
#define NW_GTPV2C_IE_BEARER_LEVEL_QOS (80)
#define NW_GTPV2C_IE_FLOW_QOS (81)
#define NW_GTPV2C_IE_RAT_TYPE (82)
#define NW_GTPV2C_IE_SERVING_NETWORK (83)
#define NW_GTPV2C_IE_BEARER_TFT (84)
#define NW_GTPV2C_IE_TAD (85)
#define NW_GTPV2C_IE_ULI (86)
#define NW_GTPV2C_IE_FTEID (87)
#define NW_GTPV2C_IE_TMSI (88)
#define NW_GTPV2C_IE_GLOBAL_CN_ID (89)
#define NW_GTPV2C_IE_S103PDF (90)
#define NW_GTPV2C_IE_S1UDF (91)
#define NW_GTPV2C_IE_DELAY_VALUE (92)
#define NW_GTPV2C_IE_BEARER_CONTEXT (93)
#define NW_GTPV2C_IE_CHARGING_ID (94)
#define NW_GTPV2C_IE_CHARGING_CHARACTERISTICS (95)
#define NW_GTPV2C_IE_TRACE_INFORMATION (96)
#define NW_GTPV2C_IE_BEARER_FLAGS (97)
#define NW_GTPV2C_IE_PDN_TYPE (99)
#define NW_GTPV2C_IE_PROCEDURE_TRANSACTION_ID (100)
#define NW_GTPV2C_IE_DRX_PARAMETER (101)
#define NW_GTPV2C_IE_UE_NETWORK_CAPABILITY (102)
#define NW_GTPV2C_IE_MM_CONTEXT (103)
#define NW_GTPV2C_IE_MM_EPS_CONTEXT (107)
#define NW_GTPV2C_IE_PDN_CONNECTION (109)
#define NW_GTPV2C_IE_PDU_NUMBERS (110)
#define NW_GTPV2C_IE_PTMSI (111)
#define NW_GTPV2C_IE_PTMSI_SIGNATURE (112)
#define NW_GTPV2C_IE_HOP_COUNTER (113)
#define NW_GTPV2C_IE_UE_TIME_ZONE (114)
#define NW_GTPV2C_IE_TRACE_REFERENCE (115)
#define NW_GTPV2C_IE_COMPLETE_REQUEST_MESSAGE (116)
#define NW_GTPV2C_IE_GUTI (117)
#define NW_GTPV2C_IE_F_CONTAINER (118)
#define NW_GTPV2C_IE_F_CAUSE (119)
#define NW_GTPV2C_IE_SELECTED_PLMN_ID (120)
#define NW_GTPV2C_IE_TARGET_IDENTIFICATION (121)
#define NW_GTPV2C_IE_PACKET_FLOW_ID (123)
#define NW_GTPV2C_IE_RAB_CONTEXT (124)
#define NW_GTPV2C_IE_SOURCE_RNC_PDCP_CONTEXT_INFO (125)
#define NW_GTPV2C_IE_UDP_SOURCE_PORT_NUMBER (126)
#define NW_GTPV2C_IE_APN_RESTRICTION (127)
#define NW_GTPV2C_IE_SELECTION_MODE (128)
#define NW_GTPV2C_IE_SOURCE_IDENTIFICATION (129)
#define NW_GTPV2C_IE_CHANGE_REPORTING_ACTION (131)
#define NW_GTPV2C_IE_FQ_CSID (132)
#define NW_GTPV2C_IE_CHANNEL_NEEDED (133)
#define NW_GTPV2C_IE_EMLPP_PRIORITY (134)
#define NW_GTPV2C_IE_NODE_TYPE (135)
#define NW_GTPV2C_IE_FQDN (136)
#define NW_GTPV2C_IE_TI (137)
#define NW_GTPV2C_IE_MBMS_SESSION_DURATION (138)
#define NW_GTPV2C_IE_MBMS_SERIVCE_AREA (139)
#define NW_GTPV2C_IE_MBMS_SESSION_IDENTIFIER (140)
#define NW_GTPV2C_IE_MBMS_FLOW_IDENTIFIER (141)
#define NW_GTPV2C_IE_MBMS_IP_MULTICAST_DISTRIBUTION (142)
#define NW_GTPV2C_IE_MBMS_IP_DISTRIBUTION_ACK (143)
#define NW_GTPV2C_IE_RFSP_INDEX (144)
#define NW_GTPV2C_IE_UCI (145)
#define NW_GTPV2C_IE_CSG_INFORMATION_REPORTING_ACTION (146)
#define NW_GTPV2C_IE_CSG_ID (147)
#define NW_GTPV2C_IE_CSG_MEMBERSHIP_INDICATION (148)
#define NW_GTPV2C_IE_SERVICE_INDICATOR (149)
#define NW_GTPV2C_IE_LDN (151)
#define NW_GTPV2C_IE_ADDITIONAL_MM_CTXT_FOR_SRVCC (159)
#define NW_GTPV2C_IE_ADDITIONAL_FLAGS_FOR_SRVCC (160)
#define NW_GTPV2C_IE_REMOTE_UE_CONTEXT (191)
#define NW_GTPV2C_IE_REMOTE_USER_ID (192)
#define NW_GTPV2C_IE_REMOTE_UE_IP_INFORMATION (193)
#define NW_GTPV2C_IE_PRIVATE_EXTENSION (255)
#define NW_GTPV2C_IE_TYPE_MAXIMUM (256)

/*--------------------------------------------------------------------------*
 *   G T P V 2 C      C A U S E      V A L U E     D E F I N I T I O N S    *
 *--------------------------------------------------------------------------*/

#define NW_GTPV2C_CAUSE_REQUEST_ACCEPTED (16)
#define NW_GTPV2C_CAUSE_INVALID_LENGTH (67)
#define NW_GTPV2C_CAUSE_MANDATORY_IE_INCORRECT (69)
#define NW_GTPV2C_CAUSE_MANDATORY_IE_MISSING (70)
#define NW_GTPV2C_CAUSE_SYSTEM_FAILURE (72)
#define NW_GTPV2C_CAUSE_REQUEST_REJECTED (94)
#define NW_GTPV2C_CAUSE_REMOTE_PEER_NOT_RESPONDING (100)
#define NW_GTPV2C_CAUSE_CONDITIONAL_IE_MISSING (103)

#define NW_GTPV2C_IE_INSTANCE_ZERO (0)
#define NW_GTPV2C_IE_INSTANCE_ONE (1)
#define NW_GTPV2C_IE_INSTANCE_TWO (2)
#define NW_GTPV2C_IE_INSTANCE_THREE (3)
#define NW_GTPV2C_IE_INSTANCE_FOUR (4)
#define NW_GTPV2C_IE_INSTANCE_MAXIMUM (NW_GTPV2C_IE_INSTANCE_FOUR)

#define NW_GTPV2C_IE_PRESENCE_MANDATORY (1)
#define NW_GTPV2C_IE_PRESENCE_CONDITIONAL (2)
#define NW_GTPV2C_IE_PRESENCE_CONDITIONAL_OPTIONAL (3)
#define NW_GTPV2C_IE_PRESENCE_OPTIONAL (4)

#endif /* __NW_GTPV2C_IE_H__ */

/*--------------------------------------------------------------------------*
 *                      E N D     O F    F I L E                            *
 *--------------------------------------------------------------------------*/

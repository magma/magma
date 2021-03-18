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

/*! \file gx_messages_types.h
 * \brief S11 definitions for interaction between MME and S11
 * 3GPP TS 29.274.
 * Messages are the same as for GTPv2-C but here we abstract the UDP layer
 * \author Sebastien ROUX <sebastien.roux@eurecom.fr>
 * \date 2013
 * \version 0.1
 */

#pragma once

#include "3gpp_24.007.h"
#include "3gpp_29.274.h"
#include "ip_forward_messages_types.h"

#define GX_NW_INITIATED_ACTIVATE_BEARER_REQ(mSGpTR)                            \
  (mSGpTR)->ittiMsg.gx_nw_init_actv_bearer_request
#define GX_NW_INITIATED_DEACTIVATE_BEARER_REQ(mSGpTR)                          \
  (mSGpTR)->ittiMsg.gx_nw_init_deactv_bearer_request

#define POLICY_RULE_NAME_MAXLEN                                                \
  100  // The policy name will be truncated to this

typedef enum {
  PCEF_STATUS_OK     = 0,
  PCEF_STATUS_FAILED = 1,
} PcefRpcStatus_t;

/**
 * PCEF Create Session response from PCEFClient, sent from GRPC task to SPGW
 * during processing of Create Session Request by SPGW task
 */
typedef struct itti_pcef_create_session_response_s {
  PcefRpcStatus_t rpc_status;
  teid_t teid;
  ebi_t eps_bearer_id;
  SGIStatus_t sgi_status;
} itti_pcef_create_session_response_t;

typedef struct itti_gx_nw_init_actv_bearer_request_s {
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;
  ebi_t lbi;
  traffic_flow_template_t ul_tft;
  traffic_flow_template_t dl_tft;
  bearer_qos_t eps_bearer_qos;
  char policy_rule_name[POLICY_RULE_NAME_MAXLEN + 1];
  uint8_t policy_rule_name_length;
} itti_gx_nw_init_actv_bearer_request_t;

typedef struct itti_gx_nw_init_deactv_bearer_request_s {
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;
  uint32_t no_of_bearers;
  ebi_t lbi;
  ebi_t ebi[BEARERS_PER_UE];
} itti_gx_nw_init_deactv_bearer_request_t;

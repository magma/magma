/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */
#ifndef FILE_S5_MESSAGES_TYPES_SEEN
#define FILE_S5_MESSAGES_TYPES_SEEN

#include "sgw_ie_defs.h"

#define S5_NW_INITIATED_ACTIVATE_BEARER_REQ(mSGpTR)                            \
  (mSGpTR)->ittiMsg.s5_nw_init_actv_bearer_request
#define S5_NW_INITIATED_ACTIVATE_BEARER_RESP(mSGpTR)                           \
  (mSGpTR)->ittiMsg.s5_nw_init_actv_bearer_response
#define S5_NW_INITIATED_DEACTIVATE_BEARER_REQ(mSGpTR)                          \
  (mSGpTR)->ittiMsg.s5_nw_init_deactv_bearer_request
#define S5_NW_INITIATED_DEACTIVATE_BEARER_RESP(mSGpTR)                         \
  (mSGpTR)->ittiMsg.s5_nw_init_deactv_bearer_response

typedef struct itti_s5_nw_init_actv_bearer_request_s {
  ebi_t lbi;///< linked Bearer ID
  teid_t mme_teid_S11;
  teid_t s_gw_teid_S11_S4;
  bearer_qos_t eps_bearer_qos;          ///< Bearer QoS
  traffic_flow_template_t ul_tft;       ///< UL TFT will be sent to UE
  traffic_flow_template_t dl_tft;       ///< DL TFT will be stored at SPGW
  protocol_configuration_options_t pco; ///< PCO protocol_configuration_options
} itti_s5_nw_init_actv_bearer_request_t;

typedef struct itti_s5_nw_init_actv_bearer_rsp_s {
  gtpv2c_cause_value_t cause;
  Imsi_t imsi;
  ebi_t lbi;
  ebi_t ebi;            ///<EPS Bearer ID
  teid_t S1_U_sgw_teid; ///< S1U sge TEID
  teid_t S1_U_enb_teid; ///< S1U enb TEID
} itti_s5_nw_init_actv_bearer_rsp_t;

typedef struct itti_s5_nw_init_deactv_bearer_request_s {
  uint8_t no_of_bearers;
  ebi_t ebi[BEARERS_PER_UE]; ///<EPS Bearer ID
  teid_t s11_mme_teid;
  bool delete_default_bearer; ///<True:Delete all bearers
                              ///<False:Delele ded bearer
} itti_s5_nw_init_deactv_bearer_request_t;

typedef struct itti_s5_nw_init_deactv_bearer_rsp_s {
  uint8_t no_of_bearers;
  ebi_t ebi[BEARERS_PER_UE]; ///<EPS Bearer ID
  bool default_bearer_deleted; ///<True:Delete all bearers
                              ///<False:Delele ded bearer
  gtpv2c_cause_t cause;
} itti_s5_nw_init_deactv_bearer_rsp_t;

#endif /* FILE_S5_MESSAGES_TYPES_SEEN*/

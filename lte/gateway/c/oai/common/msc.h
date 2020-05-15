/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#ifndef FILE_MSC_SEEN
#define FILE_MSC_SEEN

#ifdef __cplusplus
extern "C" {
#endif

typedef enum {
  MIN_MSC_ENV = 0,
  MSC_E_UTRAN = MIN_MSC_ENV,
  MSC_E_UTRAN_LIPA,
  MSC_MME_GW,
  MSC_MME,
  MSC_SP_GW,
  MAX_MSC_ENV
} msc_env_t;

typedef enum {
  MIN_MSC_PROTOS = 0,
  MSC_NAS_UE = MIN_MSC_PROTOS,
  MSC_S1AP_ENB,
  MSC_GTPU_ENB,
  MSC_GTPU_SGW,
  MSC_GTPC_SGW,
  MSC_GTPC_MME,
  MSC_S1AP_MME,
  MSC_MMEAPP_MME,
  MSC_NAS_MME,
  MSC_NAS_EMM_MME,
  MSC_NAS_ESM_MME,
  MSC_SP_GWAPP_MME,
  MSC_S11_MME,
  MSC_S6A_MME,
  MSC_SGW,
  MSC_HSS,
  MAX_MSC_PROTOS,
} msc_proto_t;

typedef struct msc_private_s {
  uint8_t                                *message_bin;
  uint32_t                                message_bin_size;
} msc_private_t;


// Access stratum
#define MSC_AS_TIME_FMT "%05u:%02u"

#define MSC_AS_TIME_ARGS(CTXT_Pp) \
    (CTXT_Pp)->frame, \
    (CTXT_Pp)->subframe
#if MESSAGE_CHART_GENERATOR
int msc_init(const msc_env_t envP, const int max_threadsP);
void msc_start_use(void);
void msc_flush_messages(void);
void msc_end(void);
void msc_log_declare_proto(const msc_proto_t  protoP);
void msc_log_event(const msc_proto_t  protoP,char *format, ...);
void msc_log_message(
	const char * const message_operationP,
    const msc_proto_t  receiverP,
    const msc_proto_t  senderP,
    uint8_t* bytesP,
    const unsigned int num_bytes,
    char *format, ...);
struct shared_log_queue_item_s;

void msc_flush_message (struct shared_log_queue_item_s *item_p);

#define MSC_INIT(arg1,arg2)                                      msc_init(arg1,arg2)
#define MSC_END                                                  msc_end
#define MSC_LOG_EVENT(mScPaRaMs, fORMAT, aRGS...)                msc_log_event(mScPaRaMs, fORMAT, ##aRGS)
#define MSC_LOG_RX_MESSAGE(rECEIVER, sENDER, bYTES, nUMbYTES, fORMAT, aRGS...)           msc_log_message("<-",rECEIVER, sENDER, bYTES, nUMbYTES, fORMAT, ##aRGS)
#define MSC_LOG_RX_DISCARDED_MESSAGE(rECEIVER, sENDER, bYTES, nUMbYTES, fORMAT, aRGS...) msc_log_message("x-",rECEIVER, sENDER, bYTES, nUMbYTES, fORMAT, ##aRGS)
#define MSC_LOG_TX_MESSAGE(sENDER, rECEIVER, bYTES, nUMbYTES, fORMAT, aRGS...)           msc_log_message("->",sENDER, rECEIVER, bYTES, nUMbYTES, fORMAT, ##aRGS)
#define MSC_LOG_TX_MESSAGE_FAILED(sENDER, rECEIVER, bYTES, nUMbYTES, fORMAT, aRGS...)    msc_log_message("-x",sENDER, rECEIVER, bYTES, nUMbYTES, fORMAT, ##aRGS)
#include "shared_ts_log.h"
#else
#define MSC_INIT(arg1,arg2)
#define MSC_END(mScPaRaMs)
#define MSC_LOG_EVENT(mScPaRaMs, fORMAT, aRGS...)
#define MSC_LOG_RX_MESSAGE(mScPaRaMs, fORMAT, aRGS...)
#define MSC_LOG_RX_DISCARDED_MESSAGE(mScPaRaMs, fORMAT, aRGS...)
#define MSC_LOG_TX_MESSAGE(mScPaRaMs, fORMAT, aRGS...)
#define MSC_LOG_TX_MESSAGE_FAILED(mScPaRaMs, fORMAT, aRGS...)
#endif

#ifdef __cplusplus
}
#endif

#endif

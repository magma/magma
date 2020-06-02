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

/*! \file msc.c
   \brief Message chart generator logging utility (generate files to be processed by a script to produce a mscgen input file for generating a sequence diagram document)
   \author  Lionel GAUTHIER
   \date 2015
   \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>
#include <limits.h>
#include <inttypes.h>
#include <sys/time.h>

#include "bstrlib.h"

#include "hashtable.h"
#include "obj_hashtable.h"
#include "log.h"
#include "msc.h"
#include "assertions.h"
#include "conversions.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "timer.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "shared_ts_log.h"
#include "log.h"

void msc_log_declare_proto (const msc_proto_t protoP);

//-------------------------------
#define MSC_MAX_QUEUE_ELEMENTS    1024
#define MSC_MAX_PROTO_NAME_LENGTH 16
#define MSC_MAX_MESSAGE_LENGTH    512

#ifdef __cplusplus
extern "C" {
#endif

//-------------------------------

FILE                                   *g_msc_fd = NULL;
char                                    g_msc_proto2str[MAX_MSC_PROTOS][MSC_MAX_PROTO_NAME_LENGTH];
int                                     g_msc_start_time_second = 0;



typedef unsigned long                   msc_message_number_t;


msc_message_number_t                    g_message_number = 0;


//------------------------------------------------------------------------------
int
msc_init (
  const msc_env_t envP,
  const int max_threadsP)
{
  int                                     i = 0;
  int                                     rv = 0;
  char                                    msc_filename[NAME_MAX+1];


  OAI_FPRINTF_INFO ("Initializing MSC logs\n");
  g_msc_start_time_second = shared_log_get_start_time_sec();
  rv = snprintf (msc_filename, NAME_MAX, "/tmp/openair.msc.%u.log", envP);   // TODO NAME

  if ((0 >= rv) || (256 < rv)) {
    OAI_FPRINTF_ERR ("Error in MSC log file name");
  }

  g_msc_fd = fopen (msc_filename, "w");
  AssertFatal (g_msc_fd != NULL, "Could not open MSC log file %s : %s", msc_filename, strerror (errno));


  for (i = MIN_MSC_PROTOS; i < MAX_MSC_PROTOS; i++) {
    switch (i) {
    case MSC_NAS_UE:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "NAS_UE");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      break;


    case MSC_S1AP_ENB:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "S1AP_ENB");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_GTPU_ENB:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "GTPU_ENB");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }
      break;

    case MSC_GTPU_SGW:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "GTPU_SGW");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_SP_GW)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_S1AP_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "S1AP_MME");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_MMEAPP_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "MME_APP");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_NAS_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "NAS_MME");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME) || (envP == MSC_E_UTRAN)) {
        msc_log_declare_proto (i);
      }
      break;

    case MSC_NAS_EMM_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "NAS_EMM");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_NAS_ESM_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "NAS_ESM");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_SP_GWAPP_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "SP_GW_MME");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if (envP == MSC_MME_GW) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_S11_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "S11_MME");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if (envP == MSC_MME) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_S6A_MME:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "S6A");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_SGW:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "SGW");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if (envP == MSC_MME) {
        msc_log_declare_proto (i);
      }

      break;

    case MSC_HSS:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "HSS");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }

      if ((envP == MSC_MME_GW) || (envP == MSC_MME)) {
        msc_log_declare_proto (i);
      }

      break;

    default:
      rv = snprintf (&g_msc_proto2str[i][0], MSC_MAX_PROTO_NAME_LENGTH, "UNKNOWN");

      if (rv >= MSC_MAX_PROTO_NAME_LENGTH) {
        g_msc_proto2str[i][MSC_MAX_PROTO_NAME_LENGTH - 1] = 0;
      }
    }
  }

  OAI_FPRINTF_INFO ("Initializing MSC logs Done\n");
  return 0;
}


//------------------------------------------------------------------------------
void msc_flush_messages (void)
{
  shared_log_flush_messages();

  fflush (g_msc_fd);
}


//------------------------------------------------------------------------------
void msc_end (void)
{
  int                                     rv = 0;

  if (NULL != g_msc_fd) {
    msc_flush_messages ();
    rv = fflush (g_msc_fd);

    if (rv != 0) {
      OAI_FPRINTF_ERR ("Error while flushing stream of MSC log file: %s", strerror (errno));
    }

    rv = fclose (g_msc_fd);

    if (rv != 0) {
      OAI_FPRINTF_ERR ("Error while closing MSC log file: %s", strerror (errno));
    }
  }
}

//------------------------------------------------------------------------------
void msc_log_declare_proto (const msc_proto_t protoP)
{
  int                                     rv = 0;
  shared_log_queue_item_t                *new_item_p = NULL;

  if ((MIN_MSC_PROTOS <= protoP) && (MAX_MSC_PROTOS > protoP)) {
    // may be build a memory pool for that also ?
    new_item_p = get_new_log_queue_item (SH_TS_LOG_MSC);

    if (NULL != new_item_p) {
      rv = bassignformat (new_item_p->bstr, "%" PRIu64 " [PROTO] %d %s\n",
          __sync_fetch_and_add (&g_message_number, 1), protoP, &g_msc_proto2str[protoP][0]);

      if (BSTR_ERR ==  rv) {
        OAI_FPRINTF_ERR ("Error while declaring new protocol in MSC: %d", protoP);

        return;
      }

      new_item_p->u_app_log.msc.message_bin = NULL;
      new_item_p->u_app_log.msc.message_bin_size = 0;
      shared_log_item (new_item_p);
    }
  }
}

//------------------------------------------------------------------------------
void msc_log_event (const msc_proto_t protoP, char *format, ...)
{
  va_list                                 args;
  int                                     rv = 0;
  shared_log_queue_item_t                *new_item_p = NULL;

  if ((MIN_MSC_PROTOS > protoP) || (MAX_MSC_PROTOS <= protoP)) {
    return;
  }

  new_item_p = get_new_log_queue_item (SH_TS_LOG_MSC);

  if (NULL != new_item_p) {
    struct timeval elapsed_time;
    shared_log_get_elapsed_time_since_start(&elapsed_time);

    rv = bassignformat (new_item_p->bstr, "%" PRIu64 " [EVENT] %d %04ld:%06ld",
          __sync_fetch_and_add (&g_message_number, 1), protoP, elapsed_time.tv_sec, elapsed_time.tv_usec);

    if (BSTR_ERR ==  rv) {
      OAI_FPRINTF_ERR ("Error while logging MSC event : %s", &g_msc_proto2str[protoP][0]);
      return;
    }

    va_start (args, format);
    rv = bvcformata (new_item_p->bstr, MSC_MAX_MESSAGE_LENGTH - rv, format, args);
    va_end (args);

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR("Error while logging MSC event : %s", &g_msc_proto2str[protoP][0]);
      return;
    }

    bcatcstr(new_item_p->bstr, "\n");

    new_item_p->u_app_log.msc.message_bin = NULL;
    new_item_p->u_app_log.msc.message_bin_size = 0;

    shared_log_item(new_item_p);
  }
}

//------------------------------------------------------------------------------
void
msc_log_message (
  const char *const message_operationP,
  const msc_proto_t proto1P,
  const msc_proto_t proto2P,
  uint8_t * bytesP,
  const unsigned int num_bytes,
  char *format,
  ...)
{
  va_list                                 args;
  uint64_t                                mac = 0;      // TO DO mac on bytesP param
  int                                     rv = 0;
  shared_log_queue_item_t                *new_item_p = NULL;

  if ((MIN_MSC_PROTOS > proto1P) || (MAX_MSC_PROTOS <= proto1P) || (MIN_MSC_PROTOS > proto2P) || (MAX_MSC_PROTOS <= proto2P)) {
    return;
  }

  new_item_p = get_new_log_queue_item (SH_TS_LOG_MSC);

  if (NULL != new_item_p) {
    struct timeval elapsed_time;
    shared_log_get_elapsed_time_since_start(&elapsed_time);
    rv = bassignformat (new_item_p->bstr, "%" PRIu64 " [MESSAGE] %d %s %d %" PRIu64 " %04ld:%06ld",
          __sync_fetch_and_add (&g_message_number, 1), proto1P, message_operationP, proto2P, mac, elapsed_time.tv_sec, elapsed_time.tv_usec);

    if (BSTR_ERR ==  rv) {
      OAI_FPRINTF_ERR ("Error while logging MSC message : %s/%s", &g_msc_proto2str[proto1P][0], &g_msc_proto2str[proto2P][0]);
      return;
    }

    va_start (args, format);
    rv = bvcformata (new_item_p->bstr, MSC_MAX_MESSAGE_LENGTH - rv, format, args);
    va_end (args);

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR("Error while logging MSC message : %s/%s", &g_msc_proto2str[proto1P][0], &g_msc_proto2str[proto2P][0]);
      return;
    }

    bcatcstr(new_item_p->bstr, "\n");

    new_item_p->u_app_log.msc.message_bin = bytesP;
    new_item_p->u_app_log.msc.message_bin_size = num_bytes;

    shared_log_item(new_item_p);
  }
}
//------------------------------------------------------------------------------
void msc_flush_message (struct shared_log_queue_item_s *item_p)
{
  int                                     rv_put = 0;

  if (blength(item_p->bstr) > 0) {
    if (g_msc_fd) {
      rv_put = fputs ((const char *)item_p->bstr->data, g_msc_fd);

      if (rv_put < 0) {
        // error occured
        OAI_FPRINTF_ERR("Error while writing msc %d\n", rv_put);
      }
      fflush (g_msc_fd);
    }
  }
}

#ifdef __cplusplus
}
#endif

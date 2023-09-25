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

#include "lte/gateway/c/core/oai/oai_mme/oai_mme.h"

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <stdbool.h>
#include <string.h>
#include <systemd/sd-daemon.h>

#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_events.hpp"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/shared_ts_log.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/mme_init.hpp"
#include "orc8r/gateway/c/common/sentry/SentryWrapper.hpp"

#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_init.h"
#include "lte/gateway/c/core/oai/tasks/sctp/sctp_primitives_server.hpp"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_extern.hpp"
/* FreeDiameter headers for support of S6A interface */
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sgs/sgs_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sms_orc8r/sms_orc8r_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/ha/ha_defs.hpp"
#include "lte/gateway/c/core/oai/common/pid_file.h"
#include "lte/gateway/c/core/oai/lib/message_utils/service303_message_utils.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#if EMBEDDED_SGW
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_embedded_spgw.hpp"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_defs.hpp"
#endif
#include "lte/gateway/c/core/oai/include/udp_primitives_server.h"
#include "lte/gateway/c/core/oai/include/s11_mme.hpp"
#include "lte/gateway/c/core/oai/include/service303.hpp"
#include "lte/gateway/c/core/oai/common/shared_ts_log.h"
#include "lte/gateway/c/core/oai/include/grpc_service.hpp"

static void send_timer_recovery_message(void);

task_zmq_ctx_t main_zmq_ctx;

static int main_init(void) {
  // Initialize main thread ZMQ context
  // We dont use the PULL socket nor the ZMQ loop
  // Don't include optional services such as CSFB, SMS, HA
  // into target task list (i.e., they will not receive any
  // broadcast messages or timer messages)
  init_task_context(
      TASK_MAIN,
      (task_id_t[]){TASK_MME_APP, TASK_SERVICE303, TASK_SERVICE303_SERVER,
                    TASK_S6A, TASK_S1AP, TASK_SCTP, TASK_SPGW_APP, TASK_SGW_S8,
                    TASK_GRPC_SERVICE, TASK_LOG, TASK_SHARED_TS_LOG,
                    TASK_ASYNC_GRPC_SERVICE},
      12, NULL, &main_zmq_ctx);

  return RETURNok;
}

static void main_exit(void) { destroy_task_context(&main_zmq_ctx); }

int main(int argc, char* argv[]) {
  srand(time(NULL));
  char* pid_file_name;

  CHECK_INIT_RETURN(OAILOG_INIT(MME_CONFIG_STRING_MME_CONFIG,
                                OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(shared_log_init(MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info,
                              messages_info, NULL, NULL));

  /*
   * Parse the command line for options and set the mme_config accordingly.
   */
#if EMBEDDED_SGW
  CHECK_INIT_RETURN(mme_config_embedded_spgw_parse_opt_line(
      argc, argv, &mme_config, &amf_config, &spgw_config));
#else
  CHECK_INIT_RETURN(mme_config_parse_opt_line(argc, argv, &mme_config));
#endif
  // Initialize Sentry error collection
  // We have to initialize here for now since itti_init asserts on there being
  // only 1 thread
  sentry_config_t sentry_config = construct_sentry_config_from_mconfig();
  initialize_sentry(SENTRY_TAG_MME, &sentry_config);

  // Could not be launched before ITTI initialization
  shared_log_itti_connect();
  OAILOG_ITTI_CONNECT();
  CHECK_INIT_RETURN(main_init());

  pid_file_name = get_pid_file_name(mme_config.pid_dir);

  if (!pid_file_lock(pid_file_name)) {
    exit(-EDEADLK);
  }
  free_wrapper((void**)&pid_file_name);

  /*
   * Calling each layer init function
   */
  // Intialize loggers and configured log levels.
  OAILOG_LOG_CONFIGURE(&mme_config.log_config);
  CHECK_INIT_RETURN(service303_init(&(mme_config.service303_config)));

  event_client_init();

  CHECK_INIT_RETURN(mme_app_init(&mme_config));
  if (mme_config.enable5g_features) {
    CHECK_INIT_RETURN(amf_app_init(&amf_config));
  }
  CHECK_INIT_RETURN(sctp_init(&mme_config));
#if EMBEDDED_SGW
  CHECK_INIT_RETURN(spgw_app_init(&spgw_config, mme_config.use_stateless));
  CHECK_INIT_RETURN(sgw_s8_init(&spgw_config.sgw_config));
#else
  CHECK_INIT_RETURN(udp_init());
  CHECK_INIT_RETURN(s11_mme_init(&mme_config));
#endif
  CHECK_INIT_RETURN(s1ap_mme_init(&mme_config));

  if (mme_config.enable5g_features) {
    CHECK_INIT_RETURN(ngap_amf_init(&amf_config));
  }
  CHECK_INIT_RETURN(s6a_init(&mme_config));

  // Create SGS Task only if non_eps_service_control is not set to OFF
  char* non_eps_service_control = bdata(mme_config.non_eps_service_control);
  if (!(strcmp(non_eps_service_control, "SMS")) ||
      !(strcmp(non_eps_service_control, "CSFB_SMS"))) {
    CHECK_INIT_RETURN(sgs_init(&mme_config));
    OAILOG_DEBUG(LOG_MME_APP, "SGS Task initialized\n");
  } else if (!(strcmp(non_eps_service_control, "SMS_ORC8R"))) {
    CHECK_INIT_RETURN(sms_orc8r_init(&mme_config));
    OAILOG_DEBUG(LOG_MME_APP, "SMS_ORC8R Task initialized\n");
  }
  CHECK_INIT_RETURN(grpc_service_init(GRPCSERVICES_SERVER_ADDRESS));
  if (mme_config.use_ha) {
    CHECK_INIT_RETURN(ha_init(&mme_config));
  }
  CHECK_INIT_RETURN(grpc_async_service_init());
  OAILOG_DEBUG(LOG_MME_APP, "MME app initialization complete\n");

  sd_notify(0, "READY=1");

#if EMBEDDED_SGW
  /*
   * Display the configuration
   */
  mme_config_display(&mme_config);
  spgw_config_display(&spgw_config);
#endif
  if (mme_config.use_stateless) {
    send_timer_recovery_message();
  }

  /*
   * Handle signals here
   */
  itti_wait_tasks_end(&main_zmq_ctx);
#if EMBEDDED_SGW
  free_spgw_config(&spgw_config);
#endif
  shutdown_sentry();
  main_exit();
  pid_file_unlock();

  return 0;
}

static void send_timer_recovery_message(void) {
  MessageDef* recovery_message_p;

  recovery_message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_UNKNOWN, RECOVERY_MESSAGE);
  send_broadcast_msg(&main_zmq_ctx, recovery_message_p);
  return;
}

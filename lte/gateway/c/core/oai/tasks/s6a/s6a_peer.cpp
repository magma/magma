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

/*! \file s6a_peer.cpp
   \brief Add a new entity to the list of peers to connect
   \author Sebastien ROUX <sebastien.roux@eurecom.fr>
   \date 2013
   \version 0.1
*/

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"

#if !S6A_OVER_GRPC
#define NB_MAX_TRIES (8)

extern __pid_t g_pid;

void s6a_peer_connected_cb(struct peer_info* info, void* arg) {
  if (info == NULL) {
    OAILOG_ERROR(LOG_S6A, "Failed to connect to HSS entity\n");
  } else {
    OAILOG_DEBUG(LOG_S6A, "Peer %*s is now connected...\n",
                 (int)info->pi_diamidlen, info->pi_diamid);

    send_activate_messages();
  }
}
#endif

status_code_e s6a_fd_new_peer(void) {
// We need to expose the function definition here because the declaration is
// outside of the guard (GH11646)
#if !S6A_OVER_GRPC
  int ret = 0;

  if (mme_config_read_lock(&mme_config)) {
    OAILOG_ERROR(LOG_S6A, "Failed to lock configuration for reading\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_S6A, "Diameter identity of MME: %s with length: %zd\n",
               fd_g_config->cnf_diamid, fd_g_config->cnf_diamid_len);
  bstring hss_name = bstrcpy(mme_config.s6a_config.hss_host_name);

  if (mme_config_unlock(&mme_config)) {
    OAILOG_ERROR(LOG_S6A, "Failed to unlock configuration\n");
    return RETURNerror;
  }
  DiamId_t diamid = bdata(hss_name);
  size_t diamidlen = blength(hss_name);
  struct peer_hdr* peer = NULL;
  int nb_tries = 0;
  int timeout = fd_g_config->cnf_timer_tc;
  for (nb_tries = 0; nb_tries < NB_MAX_TRIES; nb_tries++) {
    OAILOG_DEBUG(LOG_S6A, "S6a peer connection attempt %d / %d\n", 1 + nb_tries,
                 NB_MAX_TRIES);
    ret = fd_peer_getbyid(diamid, diamidlen, 0, &peer);
    if (peer && peer->info.config.pic_tctimer != 0) {
      timeout = peer->info.config.pic_tctimer;
    }
    if (!ret) {
      if (peer) {
        ret = fd_peer_get_state(peer);
        if (STATE_OPEN == ret) {
          OAILOG_DEBUG(LOG_S6A, "Peer %*s is now connected...\n",
                       (int)diamidlen, diamid);

          send_activate_messages();

          {
            FILE* fp = NULL;
            bstring filename = bformat("/tmp/mme_%d.status", g_pid);
            fp = fopen(bdata(filename), "w+");
            bdestroy(filename);
            fflush(fp);
            fclose(fp);
          }
          bdestroy_wrapper(&hss_name);
          return RETURNok;
        } else {
          OAILOG_DEBUG(LOG_S6A, "S6a peer state is %d\n", ret);
        }
      }
    } else {
      OAILOG_DEBUG(LOG_S6A, "Could not get S6a peer\n");
    }
    sleep(timeout);
  }
  bdestroy(hss_name);
  free_wrapper((void**)&fd_g_config->cnf_diamid);
  fd_g_config->cnf_diamid_len = 0;
#endif
  return RETURNerror;
}

/*
 * Inform S1AP and MME that connection to HSS is established
 */
void send_activate_messages(void) {
  MessageDef* message_p;
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, ACTIVATE_MESSAGE);
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, ACTIVATE_MESSAGE);
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_S1AP, message_p);
}

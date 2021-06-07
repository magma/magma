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

/*! \file s6a_peers.c
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

#include "bstrlib.h"
#include "log.h"
#include "intertask_interface.h"
#include "common_defs.h"
#include "s6a_defs.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "mme_config.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

#if !S6A_OVER_GRPC

#define NB_MAX_TRIES (8)

extern __pid_t g_pid;

void s6a_peer_connected_cb(struct peer_info* info, void* arg) {
  if (info == NULL) {
    OAILOG_ERROR(LOG_S6A, "Failed to connect to HSS entity\n");
  } else {
    OAILOG_DEBUG(
        LOG_S6A, "Peer %*s is now connected...\n", (int) info->pi_diamidlen,
        info->pi_diamid);

    send_activate_messages();
  }

  /*
   * For test
   */
#if 0
  s6a_auth_info_req_t                     s6a_air;

  memset (&s6a_air, 0, sizeof (s6a_auth_info_req_t));
  snprintf (s6a_air.imsi, sizeof(s6a_air.imsi), "%14llu", 20834123456789ULL);
  s6a_air.nb_of_vectors = 1;
  s6a_air.visited_plmn.mcc_digit2 = 0,
    s6a_air.visited_plmn.mcc_digit1 = 8, s6a_air.visited_plmn.mcc_digit3 = 2, s6a_air.visited_plmn.mnc_digit1 = 0, s6a_air.visited_plmn.mnc_digit2 = 3, s6a_air.visited_plmn.mnc_digit3 = 4, s6a_generate_authentication_info_req (&s6a_air);
  // #else
  //     s6a_update_location_req_t s6a_ulr;
  //
  //     memset(&s6a_ulr, 0, sizeof(s6a_update_location_req_t));
  //
  //     snprintf(s6a_ulr.imsi, sizeof(s6a_air.imsi), "%14llu", 20834123456789ULL);
  //     s6a_ulr.initial_attach = INITIAL_ATTACH;
  //     s6a_ulr.rat_type = RAT_EUTRAN;
  //     s6a_generate_update_location(&s6a_ulr);
#endif
}

int s6a_fd_new_peer(void) {
  int ret = 0;
#if FD_CONF_FILE_NO_CONNECT_PEERS_CONFIGURED
  struct peer_info info = {0};
#endif

  if (mme_config_read_lock(&mme_config)) {
    OAILOG_ERROR(LOG_S6A, "Failed to lock configuration for reading\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(
      LOG_S6A, "Diameter identity of MME: %s with length: %zd\n",
      fd_g_config->cnf_diamid, fd_g_config->cnf_diamid_len);
  bstring hss_name = bstrcpy(mme_config.s6a_config.hss_host_name);

  if (mme_config_unlock(&mme_config)) {
    OAILOG_ERROR(LOG_S6A, "Failed to unlock configuration\n");
    return RETURNerror;
  }
#if FD_CONF_FILE_NO_CONNECT_PEERS_CONFIGURED
  info.pi_diamid    = bdata(hss_name);
  info.pi_diamidlen = blength(hss_name);
  OAILOG_DEBUG(
      LOG_S6A, "Diameter identity of HSS: %s with length: %zd\n",
      info.pi_diamid, info.pi_diamidlen);
  info.config.pic_flags.sec     = PI_SEC_NONE;
  info.config.pic_flags.pro3    = PI_P3_DEFAULT;
  info.config.pic_flags.pro4    = PI_P4_TCP;
  info.config.pic_flags.alg     = PI_ALGPREF_TCP;
  info.config.pic_flags.exp     = PI_EXP_INACTIVE;
  info.config.pic_flags.persist = PI_PRST_NONE;
  info.config.pic_port          = 3868;
  info.config.pic_lft           = 3600;
  info.config.pic_tctimer       = 7;   // retry time-out connection
  info.config.pic_twtimer       = 60;  // watchdog
  CHECK_FCT(fd_peer_add(&info, "", s6a_peer_connected_cb, NULL));

  return ret;
#else
  DiamId_t diamid       = bdata(hss_name);
  size_t diamidlen      = blength(hss_name);
  struct peer_hdr* peer = NULL;
  int nb_tries          = 0;
  int timeout           = fd_g_config->cnf_timer_tc;
  for (nb_tries = 0; nb_tries < NB_MAX_TRIES; nb_tries++) {
    OAILOG_DEBUG(
        LOG_S6A, "S6a peer connection attempt %d / %d\n", 1 + nb_tries,
        NB_MAX_TRIES);
    ret = fd_peer_getbyid(diamid, diamidlen, 0, &peer);
    if (peer && peer->info.config.pic_tctimer != 0) {
      timeout = peer->info.config.pic_tctimer;
    }
    if (!ret) {
      if (peer) {
        ret = fd_peer_get_state(peer);
        if (STATE_OPEN == ret) {
          OAILOG_DEBUG(
              LOG_S6A, "Peer %*s is now connected...\n", (int) diamidlen,
              diamid);

          send_activate_messages();

          {
            FILE* fp         = NULL;
            bstring filename = bformat("/tmp/mme_%d.status", g_pid);
            fp               = fopen(bdata(filename), "w+");
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
  free_wrapper((void**) &fd_g_config->cnf_diamid);
  fd_g_config->cnf_diamid_len = 0;
  return RETURNerror;
#endif
}

#endif
/*
 * Inform S1AP and MME that connection to HSS is established
 */
void send_activate_messages(void) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(TASK_S6A, ACTIVATE_MESSAGE);
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);

  message_p = itti_alloc_new_message(TASK_S6A, ACTIVATE_MESSAGE);
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_S1AP, message_p);
}

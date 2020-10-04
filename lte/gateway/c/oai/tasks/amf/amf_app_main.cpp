
/****************************************************************************
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ****************************************************************************/
/*****************************************************************************

  Source      amf_app_main.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_config.h"
#include "../common/log.h"
#include "amf_app_state.h"
#include "nas_network.h"
#include "amf_app_defs.h"
#include "itti/intertask_interface_types.h"
#include "itti/intertask_interface.h"
#include "ngap_messages_types.h"

task_zmq_ctx_t amf_app_task_zmq_ctx;

using namespace std;

namespace magma5g
{
    class amf_app_main
    {
        public:
        amf_app_main();
        ~amf_app_main();
        int amf_app_init(const amf_config_t *amf_config_p);
        static void *amf_app_thread(void *args);
        int handle_message(zloop_t* loop, zsock_t* reader, void* arg);
        void amf_app_exit(void);

    };
    int amf_app_main::amf_app_init(const amf_config_t *amf_config_p)
    {
        //amf_nas_state_init(amf_config_p); //nees to crete amf_app_state
       // amf_app_edns_init(amf_config_p);
       //nas_network_initialize(amf_config_p); // needs to create initialization part
       itti_create_task(TASK_AMF_APP, &amf_app_thread, NULL); 
       //can we use MME ITTI code or wright new code for this.

    }
    static void *amf_app_thread(void *args ){

        itti_mark_task_ready(TASK_AMF_APP);
        init_task_context(TASK_AMF_APP, (task_id_t[]){TASK_NGAP, TASK_SERVICE303},2, handle_message, &amf_app_task_zmq_ctx);

      // Service started, but not healthy yet
        send_app_health_to_service303(&amf_app_task_zmq_ctx, TASK_AMF_APP, false);

        zloop_start(amf_app_task_zmq_ctx.event_loop);
        amf_app_exit();
        return NULL;
    }
    int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
        zframe_t* msg_frame = zframe_recv(reader);
        assert(msg_frame);
        MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);

        imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
        amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
               
        switch (ITTI_MSG_ID(received_message_p)) {
            case NGAP_INITIAL_UE_MESSAGE: {
                 imsi64 = amf_app_defs::amf_app_handle_initial_ue_message(amf_app_desc_p,&NGAP_INITIAL_UE_MESSAGE(received_message_p));
            }break;  
             case TERMINATE_MESSAGE: {
                itti_free_msg_content(received_message_p);
                zframe_destroy(&msg_frame);
                amf_app_exit();
            } break;                    
        // more cases.....
        }




        



    }
    amf_app_main::static void amf_app_exit(void) {
        destroy_task_context(&amf_app_task_zmq_ctx);
        //put_amf_nas_state();
        //amf_app_edns_exit();
        //clear_amf_nas_state();
        // Clean-up NAS module
        //nas_network_cleanup();
        //amf_config_exit();

        OAI_FPRINTF_INFO("TASK_MME_APP terminated\n");
        pthread_exit(NULL);
    }
}
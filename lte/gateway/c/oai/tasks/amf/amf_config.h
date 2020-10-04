/**
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
 */
/*****************************************************************************

  Source      amf_message.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "../bstr/bstrlib.h"
using namespace std;
#pragma once
#define MIN_GUMMEI 1
#define MAX_GUMMEI 5
typedef uint64_t imsi64_t;
typedef uint32_t amf_ue_ngap_id_t;
#include "amf_app_desc.h"
namespace magma5g
{
    class guamfi_config_t 
    {
        public:
        int nb;
        guamfi_t guamfi[MAX_GUMMEI]; //
    } ;
    typedef enum { RUN_MODE_TEST = 0, RUN_MODE_OTHER } run_mode_t;
    class amf_config_t 
    {
        public:
        /* Reader/writer lock for this configuration */
        bstring config_file;
        bstring pid_dir;
        bstring realm;
        bstring full_network_name;
        bstring short_network_name;
        uint8_t daylight_saving_time;

        run_mode_t run_mode;

        uint32_t max_enbs;
        uint32_t max_ues;

        uint8_t relative_capacity;

        uint32_t amf_statistic_timer;

        bstring ip_capability;
        bstring non_eps_service_control;

        uint8_t unauthenticated_imsi_supported;

        //eps_network_feature_config_t eps_network_feature_support;

        class guamfi_config_t guamfi;

        served_tai_t served_tai;

        service303_data_t service303_config;
        sctp_config_t sctp_config;
        ngap_config_t ngap_config;
        n8_config_t  n8_config;
        itti_config_t itti_config;
        nas_config_t nas_config;
        sgs_config_t sgs_config;
        log_config_t log_config;
        e_dns_config_t e_dns_emulation;

        ipv4_t ipv4;

        lai_t lai;

        bool use_stateless;
    };

    class amf_config_t amf_config = {.rw_lock = PTHREAD_RWLOCK_INITIALIZER, 0};
    #define amf_config_read_lock(aMFcONFIG)                                        \
    pthread_rwlock_rdlock(&(aMFcONFIG)->rw_lock)
    #define amf_config_write_lock(aMFcONFIG)                                       \
    pthread_rwlock_wrlock(&(aMFcONFIG)->rw_lock)
    #define amf_config_unlock(aMFcONFIG)                                           \
    pthread_rwlock_unlock(&(aMFcONFIG)->rw_lock)

}

        
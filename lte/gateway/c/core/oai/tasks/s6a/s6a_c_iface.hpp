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

#pragma once

#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"

bool s6a_viface_open(const s6a_config_t* config);
void s6a_viface_close(void);
bool s6a_viface_update_location_req(s6a_update_location_req_t* ulr_p);
bool s6a_viface_authentication_info_req(s6a_auth_info_req_t* air_p);
bool s6a_viface_send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP);
bool s6a_viface_purge_ue(const char* imsi);

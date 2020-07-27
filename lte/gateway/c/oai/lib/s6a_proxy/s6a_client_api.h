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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <gmp.h>

#ifdef __cplusplus
extern "C" {
#endif

#include "intertask_interface.h"
#include "s6a_messages_types.h"

/**
 * s6a_purge is an asynchronous call that forwards S6a PU to Federation Gateway
 * if S6a Relay is enabled by mconfig
 */
bool s6a_purge_ue(const char* imsi);

/**
 * s6a_authentication_info_req is an asynchronous call that forwards S6a AIR to
 * Federation Gateway, if S6a Relay is enabled by mconfig
 */
bool s6a_authentication_info_req(const s6a_auth_info_req_t* air_p);

/**
 * s6a_update_location_req is an asynchronous call that forwards S6a ULR to
 * Federation Gateway, if S6a Relay is enabled by mconfig
 */
bool s6a_update_location_req(const s6a_update_location_req_t* const ulr_p);

#ifdef __cplusplus
}
#endif

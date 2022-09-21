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
#ifndef FILE_GTPV1U_SGW_DEFS_SEEN
#define FILE_GTPV1U_SGW_DEFS_SEEN

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/include/gtpv1u_types.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#ifdef __cplusplus
}
#endif

int gtpv1u_init(gtpv1u_data_t* gtpv1u_data, spgw_config_t* spgw_config,
                bool persist_state);

void gtpv1u_exit(void);

#endif /* FILE_GTPV1U_SGW_DEFS_SEEN */

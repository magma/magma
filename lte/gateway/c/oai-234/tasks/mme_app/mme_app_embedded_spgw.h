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
#ifndef FILE_MME_APP_SPGW_SEEN
#define FILE_MME_APP_SPGW_SEEN
#include "mme_config.h"
#include "sgw_config.h"
#include "sgw_defs.h"

int mme_config_embedded_spgw_parse_opt_line(
    int argc, char* argv[], mme_config_t*, spgw_config_t*);

#endif /* ifndef FILE_MME_APP_SPGW_SEEN */

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

/*! \file spgw_config.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_SPGW_CONFIG_SEEN
#define FILE_SPGW_CONFIG_SEEN

#include "sgw_config.h"
#include "pgw_config.h"
#include "service303.h"
#include "bstrlib.h"

#define MAGMA_CONFIG_STRING "MAGMA"
#define SPGW_CONFIG_STRING_SERVICE303_CONFIG "SERVICE303"
#define SPGW_CONFIG_STRING_SERVICE303_CONF_SERVER_ADDRESS "SERVER_ADDRESS"

typedef struct spgw_config_s {
  sgw_config_t sgw_config;
  pgw_config_t pgw_config;
  service303_data_t service303_config;
  bstring config_file;
} spgw_config_t;

#ifndef SGW
extern spgw_config_t spgw_config;
#endif

void spgw_config_init(spgw_config_t*);

int spgw_config_parse_file(spgw_config_t*);

void spgw_config_display(spgw_config_t*);

void free_spgw_config(spgw_config_t* spgw_config_p);

#endif /* FILE_SPGW_CONFIG_SEEN */

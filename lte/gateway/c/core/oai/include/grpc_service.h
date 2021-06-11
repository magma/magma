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

#include "bstrlib.h"
#include "mme_default_values.h"

typedef struct grpc_service_data_s {
  bstring server_address;
} grpc_service_data_t;

/*
  Init GRPC Service for MME
*/
int grpc_service_init(void);

#ifdef __cplusplus
extern "C" {
#endif
/**
 * Start the GRPC Server and blocks
 *
 * @param server_address: the address and port to bind to ex "0.0.0.0:50051"
 */
void start_grpc_service(bstring server_address);

/**
 * Stop the GRPC server and clean up
 *
 */
void stop_grpc_service(void);

#ifdef __cplusplus
}
#endif

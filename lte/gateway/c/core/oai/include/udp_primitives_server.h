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

/*! \file udp_primitives_server.c
 * \brief
 * \author Sebastien ROUX
 * \company Eurecom
 * \email:
 */
#include "mme_config.h"

#ifndef UDP_PRIMITIVES_SERVER_H_
#define UDP_PRIMITIVES_SERVER_H_

/** \brief UDP task init function.
 @returns -1 on error, 0 otherwise.
 **/
int udp_init(void);
void udp_exit(void);

#endif /* UDP_PRIMITIVES_SERVER_H_ */

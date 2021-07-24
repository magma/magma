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

/*! \file sms.h
 * \brief
 * \author
 * \company
 * \email:
 */

#ifndef FILE_SMS_ORC8R_SEEN
#define FILE_SMS_ORC8R_SEEN
#include <stdint.h>
#include <netinet/in.h>
#include "bstrlib.h"
#include "hashtable.h"
#include "queue.h"
#include "nas/commonDef.h"
#include "common_types.h"

int sms_orc8r_init(const mme_config_t* mme_config);

#endif

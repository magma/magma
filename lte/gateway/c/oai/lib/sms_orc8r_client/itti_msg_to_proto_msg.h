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

#include <gmp.h>

#include "lte/protos/sms_orc8r.grpc.pb.h"
#include "lte/protos/sms_orc8r.pb.h"
#include "sgs_messages_types.h"

extern "C" {
#include "intertask_interface.h"
}

namespace magma {
using namespace lte;

SMOUplinkUnitdata convert_itti_sgsap_uplink_unitdata_to_proto_msg(
    const itti_sgsap_uplink_unitdata_t* msg);

}  // namespace magma

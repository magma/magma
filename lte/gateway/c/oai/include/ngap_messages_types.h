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
/****************************************************************************
  Source      ngap_messages_types.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_38.331.h"
#include "3gpp_38.401.h"
#include "3gpp_38.413.h"
#include "Ngap_Cause.h"
#include "TrackingAreaIdentity.h"

#define NGAP_PDUSESSION_RESOURCE_SETUP_REQ(mSGpTR)                             \
  (mSGpTR)->ittiMsg.ngap_pdusession_resource_setup_req


typedef struct itti_ngap_pdusession_resource_setup_req_s {
  gnb_ue_ngap_id_t gnb_ue_ngap_id;
  amf_ue_ngap_id_t amf_ue_ngap_id;
  bstring nas_pdu;  // optional
  ngap_ue_aggregate_maximum_bit_rate_t ue_aggregate_maximum_bit_rate;  // optional
  Ngap_PDUSession_Resource_Setup_Request_List_t pduSessionResource_setup_list;
} itti_ngap_pdusession_resource_setup_req_t;



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
  Source      ngap_amf_nas_procedures.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include "common_defs.h"
#include "3gpp_38.401.h"
#include "bstrlib.h"
#include "common_types.h"
#include "amf_app_messages_types.h"
#include "ngap_messages_types.h"
#include "ngap_state.h"
struct ngap_message_s;

int ngap_generate_ngap_pdusession_resource_setup_req(
    ngap_state_t* state, itti_ngap_pdusession_resource_setup_req_t* const
                             pdusession_resource_setup_req);


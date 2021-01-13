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
  Source      ngap_amf_handlers.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages Handlers

*****************************************************************************/
#pragma once

#include <stdbool.h>

#include "Ngap_Cause.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "ngap_amf.h"
#include "ngap_messages_types.h"
#include "sctp_messages_types.h"

int ngap_amf_handle_ue_context_release_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

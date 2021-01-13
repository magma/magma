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
  Source      amf_app_ue_context.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#pragma once

#include <inttypes.h> /* For sscanf formats */
#include <stdint.h>
#include <time.h> /* to provide time_t */

#include "amf_app_sgs_fsm.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "common_types.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "ngap_messages_types.h"
#include "obj_hashtable.h"
#include "queue.h"
#include "security_types.h"
#include "sgw_ie_defs.h"
#include "tree.h"

void amf_app_handle_ngap_ue_context_release_req(
    const itti_ngap_ue_context_release_req_t* ngap_ue_context_release_req);

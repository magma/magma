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

#pragma once
#include <sstream>
#include <thread>
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"

/* Api for deleting core network procedures */
void nas5g_delete_cn_procedure(
    struct amf_context_s* amf_context, struct nas5g_cn_proc_s* cn_proc);

/* Api for getting the Core Network procedures */
struct nas5g_cn_proc_s* get_nas5g_cn_procedure(
    const struct amf_context_s* amf_context, cn5g_proc_type_t proc_type);

/* Api for deleting common procedures */
void amf_delete_common_procedure(
    struct amf_context_s* amf_context, struct nas_amf_common_proc_s** proc);

/* Api for deleting child procedures */
void amf_delete_child_procedures(
    struct amf_context_s* amf_context,
    struct nas5g_base_proc_s* const parent_proc);

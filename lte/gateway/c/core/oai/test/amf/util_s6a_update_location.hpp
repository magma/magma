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

#include <iostream>

#include "lte/gateway/c/core/oai/test/amf/util_nas5g_pkt.hpp"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"

namespace magma5g {
// api to mock handling s6a_update_location_ans

s6a_update_location_ans_t util_amf_send_s6a_ula(const std::string& imsi);

}  // namespace magma5g

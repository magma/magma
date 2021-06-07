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

/* File : gtp_tunnel_openflow.h
 */

#pragma once

#include <arpa/inet.h>
#include <net/if.h>

// openflow init framework
int openflow_init(
    struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
    bool persist_state);

// openflow uninit flows
int openflow_uninit(void);

// openflow reset
int openflow_reset(void);

// Send end marker pdu
int openflow_send_end_marker(struct in_addr enb, uint32_t tei);

// get dev by name
const char* openflow_get_dev_name(void);

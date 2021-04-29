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

#define UPSTREAM_SOCK "unix:///tmp/sctpd_upstream.sock"
#define DOWNSTREAM_SOCK "unix:///tmp/sctpd_downstream.sock"

#define SCTP_OUT_STREAMS (16)
#define SCTP_IN_STREAMS (16)
#define SCTP_MAX_ATTEMPTS (2)
#define SCTP_TIMEOUT (5)
#define SCTP_RECV_BUFFER_SIZE (4096)

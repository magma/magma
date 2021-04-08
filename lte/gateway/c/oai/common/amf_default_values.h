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
  Source      amf_default_values.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines Access Management default values

*****************************************************************************/
#pragma once

/*******************************************************************************
 * NGAP Constants
 ******************************************************************************/
#define NGAP_OUTCOME_TIMER_DEFAULT (5)
#define AMF_STATISTIC_TIMER_S (60)
#define NGAP_PORT_NUMBER (38412)
///< IANA assigned port number for NGAP payloads on SCTP endpoint
#define NGAP_SCTP_PPID (60)  ///< NGAP SCTP Payload Protocol Identifier (PPID)

/*******************************************************************************
 * SCTP Constants
 ******************************************************************************/

#define SCTP_RECV_BUFFER_SIZE (1 << 16)
#define SCTP_OUT_STREAMS (32)
#define SCTP_IN_STREAMS (32)
#define SCTP_MAX_ATTEMPTS (5)

/*******************************************************************************
 * AMF global definitions
 ******************************************************************************/

#define AMFC (1)
#define AMFGID (1)
#define AMFPOINTER (1)
#define PLMN_MCC (208)
#define PLMN_MNC (34)
#define PLMN_MNC_LEN (2)
#define PLMN_TAC (1)

#define RELATIVE_CAPACITY (15)

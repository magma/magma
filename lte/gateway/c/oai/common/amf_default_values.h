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
 * Timer Constants
 ******************************************************************************/
#define AMF_STATISTIC_TIMER_S (60)

/*******************************************************************************
 * NGAP Constants
 ******************************************************************************/

#define NGAP_PORT_NUMBER                   (38412)  
///< IANA assigned port number for NGAP payloads on SCTP endpoint
#define NGAP_SCTP_PPID (60)  ///< NGAP SCTP Payload Protocol Identifier (PPID)

#define NGAP_OUTCOME_TIMER_DEFAULT (5)  ///< NGAP Outcome drop timer (s)


/*******************************************************************************
 * AMF global definitions
 ******************************************************************************/

#define AMFC (0)
#define AMFGID (0)
#define PLMN_MCC (208)
#define PLMN_MNC (34)
#define PLMN_MNC_LEN (2)
#define PLMN_TAC (1)

#define RELATIVE_CAPACITY (15)

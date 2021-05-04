/*
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

// could be extracted with asn1 tool

typedef enum m5g_EstablishmentCause {
  M5G_EMERGENCY = 1,
  M5G_HIGH_PRIORITY_ACCESS,
  M5G_MT_ACCESS,
  M5G_MO_SIGNALLING,
  M5G_MO_DATA,
  M5G_MO_VOICE_CALL,
  M5G_MO_VIDEOCALL,
  M5G_MO_SMS,
  M5G_MPS_PRIORITYACCESS,
  M5G_MCS_PRIORITYACCESS,
  M5G_SPARE6,
  M5G_SPARE5,
  M5G_SPARE4,
  M5G_SPARE3,
  M5G_SPARE2,
  M5G_SPARE1,
} m5g_rrc_establishment_cause_t;

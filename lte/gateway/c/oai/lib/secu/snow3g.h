/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file snow3g.h
 * \brief
 * \author Open source Adapted from Specification of the 3GPP Confidentiality
 * and Integrity Algorithms UEA2 & UIA2. Document 2: SNOW 3G Specification
 * \integrators  Kharbach Othmane, GAUTHIER Lionel.
 * \date 2014
 * \version
 * \note
 * \bug
 * \warning
 */
#ifndef FILE_SNOW3G_SEEN
#define FILE_SNOW3G_SEEN

#include <stdint.h>

typedef struct snow_3g_context_s {
  uint32_t LFSR_S0;
  uint32_t LFSR_S1;
  uint32_t LFSR_S2;
  uint32_t LFSR_S3;
  uint32_t LFSR_S4;
  uint32_t LFSR_S5;
  uint32_t LFSR_S6;
  uint32_t LFSR_S7;
  uint32_t LFSR_S8;
  uint32_t LFSR_S9;
  uint32_t LFSR_S10;
  uint32_t LFSR_S11;
  uint32_t LFSR_S12;
  uint32_t LFSR_S13;
  uint32_t LFSR_S14;
  uint32_t LFSR_S15;

  /* FSM : The Finite State Machine has three 32-bit registers R1, R2 and R3.
   */
  uint32_t FSM_R1;
  uint32_t FSM_R2;
  uint32_t FSM_R3;
} snow_3g_context_t;

/* Initialization.
 * Input k[4]: Four 32-bit words making up 128-bit key.
 * Input IV[4]: Four 32-bit words making 128-bit initialization variable.
 * Output: All the LFSRs and FSM are initialized for key generation.
 */
void snow3g_initialize(
    uint32_t k[4], uint32_t IV[4], snow_3g_context_t* snow_3g_context_pP);

/* Generation of Keystream.
 * input n: number of 32-bit words of keystream.
 * input z: space for the generated keystream, assumes
 * memory is allocated already.
 * output: generated keystream which is filled in z
 */

void snow3g_generate_key_stream(
    uint32_t n, uint32_t* z, snow_3g_context_t* snow_3g_context_pP);

#endif

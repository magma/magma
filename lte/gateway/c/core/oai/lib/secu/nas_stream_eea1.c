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

#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "assertions.h"
#include "conversions.h"
#include "secu_defs.h"
#include "snow3g.h"
#include "dynamic_memory_check.h"

int nas_stream_encrypt_eea1(
    nas_stream_cipher_t* const stream_cipher, uint8_t* const out) {
  snow_3g_context_t snow_3g_context;
  int n;
  int i             = 0;
  uint32_t zero_bit = 0;
  // uint32_t                                byte_length;
  uint32_t* KS;
  uint32_t K[4], IV[4];

  DevAssert(stream_cipher != NULL);
  DevAssert(stream_cipher->key != NULL);
  DevAssert(stream_cipher->key_length == 16);
  DevAssert(out != NULL);
  n        = (stream_cipher->blength + 31) / 32;
  zero_bit = stream_cipher->blength & 0x7;
  // byte_length = stream_cipher->blength >> 3;
  memset(&snow_3g_context, 0, sizeof(snow_3g_context));
  /*
   * Initialisation
   */
  /*
   * Load the confidentiality key for SNOW 3G initialization as in section
   * 3.4.
   */
  memcpy(
      K + 3, stream_cipher->key + 0,
      4);                                   /*K[3] = key[0]; we assume
                                             * K[3]=key[0]||key[1]||...||key[31] , with key[0] the
                                             * * * * most important bit of key */
  memcpy(K + 2, stream_cipher->key + 4, 4); /*K[2] = key[1]; */
  memcpy(K + 1, stream_cipher->key + 8, 4); /*K[1] = key[2]; */
  memcpy(
      K + 0, stream_cipher->key + 12,
      4); /*K[0] = key[3]; we assume
           * K[0]=key[96]||key[97]||...||key[127] , with key[127] the
           * * * * least important bit of key */
  K[3] = hton_int32(K[3]);
  K[2] = hton_int32(K[2]);
  K[1] = hton_int32(K[1]);
  K[0] = hton_int32(K[0]);
  /*
   * Prepare the initialization vector (IV) for SNOW 3G initialization as in
   * section 3.4.
   */
  IV[3] = stream_cipher->count;
  IV[2] = ((((uint32_t) stream_cipher->bearer) << 3) |
           ((((uint32_t) stream_cipher->direction) & 0x1) << 2))
          << 24;
  IV[1] = IV[3];
  IV[0] = IV[2];
  /*
   * Run SNOW 3G algorithm to generate sequence of key stream bits KS
   */
  snow3g_initialize(K, IV, &snow_3g_context);
  KS = (uint32_t*) calloc(1, 4 * n);
  snow3g_generate_key_stream(n, (uint32_t*) KS, &snow_3g_context);

  if (zero_bit > 0) {
    KS[n - 1] = KS[n - 1] & (uint32_t)(0xFFFFFFFF << (8 - zero_bit));
  }

  for (i = 0; i < n; i++) {
    KS[i] = hton_int32(KS[i]);
  }

  /*
   * Exclusive-OR the input data with keystream to generate the output bit
   * stream
   */
  for (i = 0; i < n * 4; i++) {
    stream_cipher->message[i] ^= *(((uint8_t*) KS) + i);
  }

  int ceil_index = 0;

  if (zero_bit > 0) {
    ceil_index = (stream_cipher->blength + 7) >> 3;
    stream_cipher->message[ceil_index - 1] =
        stream_cipher->message[ceil_index - 1] &
        (uint8_t)(0xFF << (8 - zero_bit));
  }

  free_wrapper((void**) &KS);
  memcpy(out, stream_cipher->message, n * 4);

  if (zero_bit > 0) {
    out[ceil_index - 1] = stream_cipher->message[ceil_index - 1];
  }

  return 0;
}

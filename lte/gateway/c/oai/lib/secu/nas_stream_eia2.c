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
#include <openssl/cmac.h>
#include <openssl/evp.h>
#include <openssl/ossl_typ.h>

#include "secu_defs.h"
#include "assertions.h"
#include "conversions.h"
#include "dynamic_memory_check.h"
#include "log.h"

/*!
   @brief Create integrity cmac t for a given message.
   @param[in] stream_cipher Structure containing various variables to setup
   encoding
   @param[out] out For EIA2 the output string is 32 bits long
*/
int nas_stream_encrypt_eia2(
    nas_stream_cipher_t* const stream_cipher, uint8_t const out[4]) {
  uint8_t* m               = NULL;
  uint32_t local_count     = 0;
  size_t size              = 4;
  uint8_t data[16]         = {0};
  CMAC_CTX* cmac_ctx       = NULL;
  const EVP_CIPHER* cipher = EVP_aes_128_cbc();
  uint32_t zero_bit        = 0;
  uint32_t m_length;

  DevAssert(stream_cipher != NULL);
  DevAssert(stream_cipher->key != NULL);
  DevAssert(stream_cipher->key_length > 0);
  DevAssert(out != NULL);
  zero_bit = stream_cipher->blength & 0x7;
  m_length = stream_cipher->blength >> 3;

  if (zero_bit > 0) m_length += 1;

  local_count = hton_int32(stream_cipher->count);
  m           = calloc(1, m_length + 8);
  memcpy(&m[0], &local_count, 4);
  m[4] = ((stream_cipher->bearer & 0x1F) << 3) |
         ((stream_cipher->direction & 0x01) << 2);
  memcpy(&m[8], stream_cipher->message, m_length);

  OAILOG_TRACE(
      LOG_NAS, "Byte length: %u, Zero bits: %u:\n", m_length + 8, zero_bit);
  OAILOG_STREAM_HEX(OAILOG_LEVEL_TRACE, LOG_NAS, "m:", m, m_length + 8);
  OAILOG_STREAM_HEX(
      OAILOG_LEVEL_TRACE, LOG_NAS, "Key:", stream_cipher->key,
      stream_cipher->key_length);
  OAILOG_STREAM_HEX(
      OAILOG_LEVEL_TRACE, LOG_NAS, "Message:", stream_cipher->message,
      m_length);

  cmac_ctx = CMAC_CTX_new();
  CMAC_Init(
      cmac_ctx, stream_cipher->key, stream_cipher->key_length, cipher, NULL);
  CMAC_Update(cmac_ctx, m, m_length + 8);
  CMAC_Final(cmac_ctx, data, &size);
  CMAC_CTX_free(cmac_ctx);
  OAILOG_STREAM_HEX(OAILOG_LEVEL_TRACE, LOG_NAS, "Out:", data, size);
  memcpy((void*) out, data, 4);
  free_wrapper((void**) &m);
  return 0;
}

/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file digest.c
   \brief
   \author  Lionel GAUTHIER
   \date 2017
   \email: lionel.gauthier@eurecom.fr
*/

#include <stdlib.h>
#include <openssl/crypto.h>
#include <openssl/evp.h>

#include "log.h"
#include "common_defs.h"
#include "digest.h"

//------------------------------------------------------------------------------
// evp_x can be EVP_sha256, ...
int digest_buffer(
    const EVP_MD* (*evp_x)(void), const unsigned char* buffer,
    size_t buffer_len, unsigned char** digest, unsigned int* digest_len) {
  EVP_MD_CTX* mdctx = NULL;

  if ((mdctx = EVP_MD_CTX_create())) {
    if (1 == EVP_DigestInit_ex(mdctx, (*evp_x)(), NULL)) {
      if (1 == EVP_DigestUpdate(mdctx, buffer, buffer_len)) {
        if ((*digest =
                 (unsigned char*) OPENSSL_malloc(EVP_MD_size((*evp_x)())))) {
          if (1 == EVP_DigestFinal_ex(mdctx, *digest, digest_len)) {
            EVP_MD_CTX_destroy(mdctx);
            return RETURNok;
          } else {
            OPENSSL_free(*digest);
            OAILOG_ERROR(LOG_UTIL, "Digest EVP_DigestFinal_ex()\n");
          }
        } else {
          OAILOG_ERROR(LOG_UTIL, "Digest OPENSSL_malloc()\n");
        }
      } else {
        OAILOG_ERROR(LOG_UTIL, "Digest EVP_DigestUpdate()\n");
      }
    } else {
      OAILOG_ERROR(LOG_UTIL, "Digest EVP_DigestInit_ex()\n");
    }
    EVP_MD_CTX_destroy(mdctx);
  } else {
    OAILOG_ERROR(LOG_UTIL, "Digest EVP_MD_CTX_create()\n");
  }
  *digest     = NULL;
  *digest_len = 0;
  return RETURNerror;
}

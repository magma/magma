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

/*! \file TLVDecoder.c
  \brief
  \author Philippe MOREL, Sebastien ROUX, Lionel GAUTHIER
  \company Eurecom
*/

#include <stdio.h>
#include <stdint.h>

#include "bstrlib.h"
#include "TLVDecoder.h"

int errorCodeDecoder = 0;

//------------------------------------------------------------------------------
int decode_bstring(
    bstring* bstr, const uint16_t pdulen, const uint8_t* const buffer,
    const uint32_t buflen) {
  if (buflen < pdulen) {
    return TLV_BUFFER_TOO_SHORT;
  }

  if ((bstr) && (buffer)) {
    *bstr = blk2bstr(buffer, pdulen);
    return pdulen;
  } else {
    return TLV_BUFFER_TOO_SHORT;
  }
}

//------------------------------------------------------------------------------
bstring dump_bstring_xml(const bstring bstr) {
  if (bstr) {
    int i;

    bstring b = bformat("<Length>%u</Length>\n\t<values>", bstr->slen);
    for (i = 0; i < bstr->slen; i++) {
      bformata(b, "0x%x ", bstr->data[i]);
    }
    bcatcstr(b, "</values>\n");
    return b;
  } else {
    bstring b = bfromcstr("<Length>0</Length>\n");
    return b;
  }
}

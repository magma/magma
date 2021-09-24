/*
Copyright 2020 The Magma Authors.
This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
#include <sstream>
#include <cstdint>
#include <string.h>
#include "M5GNSSAI.h"
#include "M5GCommonDefs.h"
namespace magma5g {
NSSAIMsg::NSSAIMsg(){};

NSSAIMsg::~NSSAIMsg(){};

int NSSAIMsg::EncodeNSSAIMsg(
    NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  MLOG(MDEBUG) << "EncodeNSSAIMsg ";
  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) NSSAI->iei);
	ENCODE_U8(buffer, iei, encoded);
	MLOG(MDEBUG) << "iei: " << std::hex << (int)iei ;
  }

#if 0
  uint8_t l  = 0;
  l++; //SST

  if(NSSAI->sd && NSSAI->hplmn_sst && NSSAI->hplmn_sd) {
  	l += 7;
  }
  else if (NSSAI->sd && NSSAI->hplmn_sst) {
  	l += 5;
  }
  else if ( NSSAI->sd ) {
  	l += 3;
  }
  else if ( NSSAI->hplmn_sst ) {
  	l += 1;	
  }
  NSSAI->len = l;

 #endif 
  ENCODE_U8(buffer, NSSAI->len, encoded);
  MLOG(MDEBUG) << "len: " <<  (int)NSSAI->len ;
  
  switch (NSSAI->len)
  {
  case 0b00000001: // SST
	  ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
	   MLOG(MDEBUG) << "sst: " <<  (int)NSSAI->sst ;
	  break;
  case 0b00000010: // SST and mapped HPLMN SST
	  ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);
	  break;
  case 0b00000100: // SST and SD
	  ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);
	  break;
  case 0b00000101: // SST, SD and mapped HPLMN SST
	  ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);
	  break;
  case 0b00001000: // SST, SD, mapped HPLMN SST and mapped HPLMN SD
	  ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);
	  
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[0], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[1], encoded);
	  ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[2], encoded);
	  break;
  default: // All other values are reserved
	  break;
  	}

  return (encoded);
};

int NSSAIMsg::DecodeNSSAIMsg(
    NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer, uint32_t len) {
	int decoded   = 0;

	MLOG(MDEBUG) << "DecodeNSSAIMsg ";
	if (iei > 0) {
	  CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
	  NSSAI->iei = *(buffer + decoded);
	  MLOG(MDEBUG) << "iei: " << std::hex << (int)iei ;
	  decoded++;	  
	}
	DECODE_U8(buffer + decoded, NSSAI->len, decoded);
	CHECK_LENGTH_DECODER(len - decoded, NSSAI->len);
    MLOG(MDEBUG) << "len: " << (int)NSSAI->len;
	
	switch (NSSAI->len)
    {
    case 0b00000001: // SST
    	DECODE_U8(buffer + decoded, NSSAI->sst, decoded);
        break;
    case 0b00000010: // SST and mapped HPLMN SST
        DECODE_U8(buffer + decoded, NSSAI->sst, decoded);
		DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);
        break;
    case 0b00000100: // SST and SD
        DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

		DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);
        break;
    case 0b00000101: // SST, SD and mapped HPLMN SST
        DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

		DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);
		
		DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);
        break;
    case 0b00001000: // SST, SD, mapped HPLMN SST and mapped HPLMN SD
        DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

		DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);
		
		DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);

		DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[0], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[1], decoded);
	  	DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[2], decoded);
        break;
    default: // All other values are reserved
        break;
    }

	 MLOG(MDEBUG) << "sst: " << (int)NSSAI->sst;
	
	return decoded;

};
}  // namespace magma5g

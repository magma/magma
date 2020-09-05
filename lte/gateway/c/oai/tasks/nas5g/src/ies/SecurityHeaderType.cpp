#include <iostream>
#include <sstream>
#include <cstdint>
#include "SecurityHeaderType.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  SecurityHeaderTypeMsg::SecurityHeaderTypeMsg()
  {
  };

  SecurityHeaderTypeMsg::~SecurityHeaderTypeMsg()
  {
  };

  // Decode SecurityHeaderType IE
  int SecurityHeaderTypeMsg::DecodeSecurityHeaderTypeMsg(SecurityHeaderTypeMsg *securityheadertype, uint8_t iei, uint8_t *buffer, uint32_t len) 
  {
    int decoded = 0;

    MLOG(MDEBUG) << "   DecodeSecurityHeaderTypeMsg : \n";
    securityheadertype->securityhdr = *(buffer) & 0xf;
    decoded++;
    MLOG(MDEBUG) << "      Security hdr type = 0x" << hex << int(securityheadertype->securityhdr)<< "\n";
    return (decoded);
  };

  // Encode SecurityHeaderType IE
  int SecurityHeaderTypeMsg::EncodeSecurityHeaderTypeMsg(SecurityHeaderTypeMsg *securityheadertype, uint8_t iei, uint8_t * buffer, uint32_t len)
  {
    int encoded = 0;

    MLOG(MDEBUG) << "EncodeSecurityHeaderTypeMsg:";
    *(buffer) = securityheadertype->securityhdr & 0xf;
    MLOG(MDEBUG) << "Security hdr type 0x" << hex << int(*(buffer))<< endl;
    encoded++;
    return (encoded);
  };
}


#include <iostream>
#include <sstream>
#include <cstdint>
#include "SpareHalfOctet.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  SpareHalfOctetMsg::SpareHalfOctetMsg()
  {
  };

  SpareHalfOctetMsg::~SpareHalfOctetMsg()
  {
  };

  // Decode SpareHalfOctet IE
  int SpareHalfOctetMsg::DecodeSpareHalfOctetMsg(SpareHalfOctetMsg *sparehalfoctet, uint8_t iei, uint8_t *buffer, uint32_t len) 
  {
    int decoded = 0;

    MLOG(MDEBUG) << "   DecodeSpareHalfOctetMsg : ";
    sparehalfoctet->spare = (*buffer & 0xf0) >> 4;
    decoded++;
    MLOG(MDEBUG) << "Spare = 0x" << hex << int(sparehalfoctet->spare)<< "\n";
    return (decoded);
  };

  // Encode SpareHalfOctet IE
  int SpareHalfOctetMsg::EncodeSpareHalfOctetMsg(SpareHalfOctetMsg *sparehalfoctet, uint8_t iei, uint8_t * buffer, uint32_t len)
  {
    int encoded = 0;

    MLOG(MDEBUG) << " EncodeSpareHalfOctetMsg : ";
    *(buffer) = 0x00 | (sparehalfoctet->spare & 0xf) << 4;
    MLOG(MDEBUG) << "   Spare = 0x" << hex << int(*(buffer))<< endl;
    encoded++;
    return (encoded);
  };
}


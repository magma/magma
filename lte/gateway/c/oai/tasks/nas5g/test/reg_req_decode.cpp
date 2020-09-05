#include <iostream>
#include <RegistrationRequest.h>

using namespace std;
using namespace magma5g;

namespace magma5g
{
   int decode(void)
   {
      int ret = 0;
      uint8_t buffer[] = {0x7E, 0x00, 0x41, 0x71, 0x00, 0x0D, 0x01, 0x13, 0x00, 0x14,0xF0, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xF1, 0x2E, 0x02, 0x80, 0x40};
      int len = 23;
      RegistrationRequestMsg Req;
      ret = Req.DecodeRegistrationRequestMsg( &Req, buffer, len);
      return 0;
   }
}  

int main(void)
{
   int ret;
   ret = magma5g::decode();
   return 0;
   }

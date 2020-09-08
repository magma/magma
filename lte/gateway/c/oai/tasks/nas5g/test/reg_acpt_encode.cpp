#include <iostream>
#include <RegistrationAccept.h>
#include <AmfMessage.h>
#include <CommonDefs.h>

using namespace std;
using namespace magma5g;

namespace magma5g
{
   int encode(void)
   {
      int enc_r = 0;
      uint8_t buffer[5] = {};
      int len = 5;
      RegistrationAcceptMsg msg;

      // Message
      msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 126;
      msg.securityheadertype.securityhdr = 0;
      msg.messagetype.msgtype = 0x42;
      msg.m5gsregistrationresult.spare = 0;
      msg.m5gsregistrationresult.smsallowed = 0;
      msg.m5gsregistrationresult.registrationresultval = 1;
     
      MLOG(MDEBUG) << "\n---Encoding Registration Accept Message---";
      enc_r = msg.EncodeRegistrationAcceptMsg(&msg, buffer, len);

      MLOG(MDEBUG) << "\n\n Encoded Message : ";
      MLOG(MDEBUG) << hex << int(buffer[0]) << "0" << hex << int(buffer[1]) << hex << int(buffer[2]) << "0" << hex << int(buffer[3]) << "0" << hex<< int(buffer[4]);
      MLOG(MDEBUG) << "\n\n";

      MLOG(MDEBUG) << "---Decoding Encoded Registration Accept Message---\n\n";
      int dec_r = 0;
      dec_r = msg.DecodeRegistrationAcceptMsg(&msg, buffer, len);
      MLOG(MDEBUG) << "\n\n";

      return 0;
   }
}  

int main(void)
{
   int ret;
   ret = magma5g::encode();
   return 0;
}


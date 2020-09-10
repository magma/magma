#include <iostream>
#include "DeRegistrationAcceptUEInit.h"
#include "AmfMessage.h"
#include "CommonDefs.h"
using namespace std;
using namespace magma5g;

namespace magma5g
{
   int encode(void)
   {
      int ret = 0;
      uint8_t buffer[3] = {};
      int len = 5;
      DeRegistrationAcceptUEInitMsg msg;

      msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 126;
      msg.securityheadertype.securityhdr = 0;
      msg.messagetype.msgtype = 0x46;

      MLOG(MDEBUG) << "\n\n---Encoding message--- \n\n"; 
      ret = msg.EncodeDeRegistrationAcceptUEInitMsg( &msg, buffer, len);
      
      MLOG(MDEBUG) << "\n\n ENCODED MESSAGE : " ;
      MLOG(MDEBUG) << hex << int(buffer[0])<<"0"<< hex << int(buffer[1])<< hex << int(buffer[2]) << "\n\n";

      MLOG(MDEBUG) << "---Decoding encoded message--- \n\n" ;
      int ret2 =0;
      ret2 = msg.DecodeDeRegistrationAcceptUEInitMsg(&msg, buffer, len);

      MLOG(BEBUG) << "\n ---DECODED MESSAGE ---\n\n";
      MLOG(MDEBUG) << " Extended Protocol Discriminator :" << dec << int(msg.extendedprotocoldiscriminator.extendedprotodiscriminator) << endl;
      MLOG(MDEBUG) << " Spare half octet : 0" << endl;
      MLOG(MDEBUG) << " Security Header Type : " << dec << int(msg.securityheadertype.securityhdr) << endl;
      MLOG(MDEBUG) << " Message Type : 0x" << hex << int(msg.messagetype.msgtype) << endl;
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

#include <iostream>
#include "DeRegistrationRequestUEInit.h"
#include "CommonDefs.h"
using namespace std;
using namespace magma5g;

namespace magma5g
{
   int decode(void)
   {
      int ret = 0;
      uint8_t buffer[] = {0x7E, 0x00, 0x45, 0x01, 0x00, 0x0B, 0xF2, 0x13, 0x00, 0x14,0x44, 0x33, 0x12, 0x00, 0x00, 0x00, 0x01};
      int len = 17;
      DeRegistrationRequestUEInitMsg De_Req;
      MLOG(MDEBUG) << "\n\n---Decoding De-registration request (UE originating) Message---\n\n";
      ret = De_Req.DecodeDeRegistrationRequestUEInitMsg( &De_Req, buffer, len);


      MLOG(BEBUG) << "\n ---DECODED MESSAGE ---\n\n";
      MLOG(MDEBUG) << " Extended Protocol Discriminator :" << dec << int(De_Req.extendedprotocoldiscriminator.extendedprotodiscriminator) << endl;
      MLOG(MDEBUG) << " Spare half octet : 0" << endl;
      MLOG(MDEBUG) << " Security Header Type : " << dec << int(De_Req.securityheadertype.securityhdr) << endl;
      MLOG(MDEBUG) << " Message Type : 0x" << hex << int(De_Req.messagetype.msgtype) << endl;
      MLOG(MDEBUG) << " M5GS De-Registration Type :\n";
      MLOG(DEBUG)  << "   Switch off = "<< dec << int(De_Req.m5gsderegistrationtype.switchoff)<<"\n";
      MLOG(MDEBUG) << "   Re-registration required = "<< dec << int(De_Req.m5gsderegistrationtype.reregistrationrequired) << "\n";
      MLOG(MDEBUG) << "   Access Type = "<< dec << int(De_Req.m5gsderegistrationtype.accesstype) <<"\n";
      MLOG(MDEBUG) << " NAS key set identifier : \n";
      MLOG(MDEBUG) << "   Type of security context flag = " << dec << int(De_Req.naskeysetidentifier.tsc) << "\n";
      MLOG(MDEBUG) << "   NAS key set identifier = " << dec << int(De_Req.naskeysetidentifier.naskeysetidentifier) << "\n";
      MLOG(MDEBUG) << " M5GS mobile identity : \n";
      MLOG(MDEBUG) << "   Odd/even Indication = "<< dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.oddeven) << "\n";
      MLOG(MDEBUG) << "   Type of identity = " << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.typeofidentity) << "\n";
      MLOG(MDEBUG) << "   Mobile Country Code (MCC) = "<< dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mcc_digit1);
      MLOG(MDEBUG) << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mcc_digit2); 
      MLOG(MDEBUG) << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mcc_digit3) << "\n";
      MLOG(MDEBUG) << "   Mobile NetWork Code (MNC) = "<< dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mnc_digit1);
      MLOG(MDEBUG) << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mnc_digit2); 
      MLOG(MDEBUG) << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.mnc_digit3) << "\n";
      MLOG(MDEBUG) << " AMF Region ID = " << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.amfregionid)<<endl;
      MLOG(MDEBUG) << " AMF Set ID = " << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.amfsetid)<<endl;
      MLOG(MDEBUG) << " AMF Pointer = " << dec << int(De_Req.m5gsmobileidentity.mobileidentity.guti.amfpointer)<<endl;
      MLOG(MDEBUG) << " 5G-TMSI = 0x0" << hex << int(De_Req.m5gsmobileidentity.mobileidentity.guti.tmsi1);
      MLOG(MDEBUG) << "0" << hex << int(De_Req.m5gsmobileidentity.mobileidentity.guti.tmsi2);
      MLOG(MDEBUG) << "0" << hex << int(De_Req.m5gsmobileidentity.mobileidentity.guti.tmsi3);
      MLOG(MDEBUG) << "0" << hex << int(De_Req.m5gsmobileidentity.mobileidentity.guti.tmsi4)<<"\n\n";

      return 0;
   }
}  

int main(void)
{
   int ret;
   ret = magma5g::decode();
   //cout<<"Decoded messgae :"<<ret;
   return 0;
   }

#include <string>
#include "NasMessage.h"

using namespace nas;
int main() {
  std::string initialUEMsgHexStr = "7e004179000d0109f1070000000000000000102e04f0f0f0f0";
  //std::string initialUEMsgHexStr = "2e0101c1ffff91a12801007b000780000a00000d00";
  NasMessage* pNasMsg = nullptr;
  NasCause cause = NasCause::NAS_CAUSE_SUCCESS;
  pNasMsg = NasMessageFactory::DecodeNasMessage(initialUEMsgHexStr, cause);
  NasBuffer nasBuffer;
  pNasMsg->Encode(nasBuffer);
  std::cout << nasBuffer.ToHexString() << std::endl;
  std::cout << pNasMsg->ToHexString() << std::endl;
  return 0;
}

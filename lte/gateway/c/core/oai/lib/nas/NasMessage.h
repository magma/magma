#pragma once
//https://www.etsi.org/deliver/etsi_ts/124500_124599/124501/16.05.01_60/ts_124501v160501p.pdf
#include <unordered_map>
#include <unordered_set>


#include "NasInformationElement.h"
#include "NasUtils.h"

namespace nas {

//**************************Nas Messages *********************/

#define NAS_SECURITY_HEADER_SIZE 6 // SecurityHdrtype+MAC+Sqn

class NasSecurity  {

private:
    SecurityHeaderTypeIE        m_sht;      
	  SpareHalfOctetIE            m_sho;
    MessageAuthenticationCodeIE m_mac;
    SequenceNumberIE            m_sqn;
public:
    NasSecurity() {}

    void SetSecurityHeaderType(SecurityHeaderType sht) { m_sht.SetSecurityHeaderType(sht); }
    SecurityHeaderType GetSecurityHeaderType() const{ return m_sht.GetSecurityHeaderType(); }
   
    void SetMessageAuthenticationCode(std::vector<uint8_t> mac) {
      m_mac.SetMessageAuthenticationCode(mac);      
    }
    std::vector<uint8_t> GetMessageAuthenticationCode() const{
      return m_mac.GetMessageAuthenticationCode();
    }

    void SetSequenceNumber(uint8_t sqn) {
      m_sqn.SetSequenceNumber(sqn);
    }
    uint8_t GetSequenceNumber() const{
      return m_sqn.GetSequenceNumber();
    }

    void Decode(const NasBuffer& nasBuffer) {
      m_sho.Decode(nasBuffer);  
      m_sht.Decode(nasBuffer);
      m_mac.Decode(nasBuffer);
      m_sqn.Decode(nasBuffer);
    }
    void Encode(NasBuffer& nasBuffer) const{
       m_sho.Encode(nasBuffer);
       m_sht.Encode(nasBuffer);
       m_mac.Encode(nasBuffer);
       m_sqn.Encode(nasBuffer);
    }
};

class NasMessage {
 private:
  NasSecurity*                    m_nasSecurity = nullptr;
 
 protected: 
  ExtendedProtocolDiscriminatorIE m_epd ;
  MessageTypeIE                   m_msgType ;
 
  std::unordered_map<InformationElementType, InformationElement*> m_optionalIEIs;

  bool DecodeOptionalIEI(const NasBuffer& nasBuffer) {

    InformationElementType ieiType = 
        static_cast<InformationElementType>(nasBuffer.GetCurrentOctet());

    if(IsOptionalIESupport(ieiType)){
        ieiType = static_cast<InformationElementType>(nasBuffer.DecodeU8());
        InformationElement* pIEI = 
                InformationElementFactory::AllocInformationElement(ieiType);
        if(pIEI) {
            m_optionalIEIs.emplace(ieiType, pIEI);
            NasCause cause =  pIEI->Decode(nasBuffer);
            return true;  
        }
        else {
            std::cout << NasUtils::Enum2String(ieiType) 
                      << " Optional IEI  not implemented" << std::endl;
            return false;
        }     
    }
    else {
            std::cout << NasUtils::Enum2String(ieiType) 
                      << " Optional IEI does not supported" << std::endl;
          return false;
    }
  }
  void DecodeNasSecurity(const NasBuffer& nasBuffer) {
      m_nasSecurity = new NasSecurity();
      m_nasSecurity->Decode(nasBuffer);
  }
  void EncodeNasSecurity(NasBuffer& nasBuffer) const{
      if(m_nasSecurity) {
        m_nasSecurity->Encode(nasBuffer);
      }
  }

  void DecodeOptionalIEIs(const NasBuffer& nasBuffer) {
      while(!nasBuffer.EndOfBuffer()) {
          if(!DecodeOptionalIEI(nasBuffer)) {
            break;  
          }
      }
  }
  void EncodeOptionalIEIs(NasBuffer& nasBuffer) const {
    for(auto& iei: m_optionalIEIs) {
        InformationElement* pElement = iei.second;
        if(pElement) { pElement->Encode(nasBuffer); }
    }
  }

  virtual void DecodeNasHeader(const NasBuffer& nasBuffer) = 0;
  virtual void EncodeNasHeader(NasBuffer& nasBuffer) const = 0;

 public:
  NasMessage() {
    
  }
  virtual ~NasMessage() {
    if(m_nasSecurity) {
      delete m_nasSecurity;
      m_nasSecurity = nullptr;
    }
    for(auto it: m_optionalIEIs) {
      InformationElementFactory::DeallocInformatonElement(it.second); 
    }
    m_optionalIEIs.clear();
  }
  void SetExtendedProtocolDiscriminator(ExtendedProtocolDiscriminator epd) {
    m_epd.SetExtendedProtocolDiscriminator(epd);
  }
  ExtendedProtocolDiscriminator GetExtendedProtocolDiscriminator() const {
    return m_epd.GetExtendedProtocolDiscriminator();
  }

  void SetNasSecurity(NasSecurity* nasSecurity) {
      m_nasSecurity = nasSecurity;
  }

  NasSecurity* GetNasSecurity() const{
      return m_nasSecurity;
  }

  void SetMessageType(MessageType type) { m_msgType.SetMessageType(type); }
  MessageType GetMessageType() const { return m_msgType.GetMessageType(); }

  std::string ToHexString() {
      NasBuffer nasBuffer;
      this->Encode(nasBuffer);
      return nasBuffer.ToHexString();
  }

  virtual void Decode(const NasBuffer& nasBuffer) = 0;
  virtual void Encode(NasBuffer& nasBuffer) const = 0;
  virtual bool IsOptionalIESupport(InformationElementType type) const = 0;
};

class MobilityManagementMessage : public NasMessage {
 protected:
  SecurityHeaderTypeIE m_sht;
  SpareHalfOctetIE     m_sho;

  void DecodeNasHeader(const NasBuffer& nasBuffer) override{
    m_epd.Decode(nasBuffer);
    if(SecurityHeaderType::NOT_PROTECTED != nasBuffer.GetSecurityHeaderType()) {
      DecodeNasSecurity(nasBuffer);
      m_epd.Decode(nasBuffer);    
    }
    m_sho.Decode(nasBuffer);
    m_sht.Decode(nasBuffer);
    m_msgType.Decode(nasBuffer);
  }

  void EncodeNasHeader(NasBuffer& nasBuffer) const override{
      m_epd.Encode(nasBuffer);
      if(GetNasSecurity()) {
          EncodeNasSecurity(nasBuffer);
          m_epd.Encode(nasBuffer);
      } 
      m_sho.Encode(nasBuffer);
      m_sht.Encode(nasBuffer);
      m_msgType.Encode(nasBuffer);
  }

 public:
  MobilityManagementMessage() {}
  
  uint8_t GetSpareHalfOctet() const{
      return m_sho.GetSpareHalfOctet();
  }
  void SetSecurityHeaderType(SecurityHeaderType sht) { m_sht.SetSecurityHeaderType(sht); }
  SecurityHeaderType GetSecurityHeaderType() const { return m_sht.GetSecurityHeaderType(); }
};

class RegistrationRequest : public MobilityManagementMessage {
 private:
  FiveGSRegistrationTypeIE    m_fgsRegType;
  NasKeySetIdentifierIE       m_ngKSI;
  MobileIdentityIE            m_fgsMobileIdentity;

public:
  RegistrationRequest(){}

  bool IsOptionalIESupport(InformationElementType type) const {
    std::unordered_set<InformationElementType> supportedOptionalIEIs;
    const auto& it =  supportedOptionalIEIs.find(type) ;
	  if (it != supportedOptionalIEIs.end()) {
      return true;
    }
    std::cout << "unsupported Optional IEI present\n";
    return false;
  }

  void SetRegistrationType(RegistrationType regType) {
    m_fgsRegType.SetRegistrationType(regType);
  }
  RegistrationType GetRegistrationType() const{
    return m_fgsRegType.GetRegistrationType();
  }

  void  Decode(const NasBuffer& nasBuffer) override {   
    
    DecodeNasHeader(nasBuffer);
    m_ngKSI.Decode(nasBuffer);
    m_fgsRegType.Decode(nasBuffer); 
    m_fgsMobileIdentity.Decode(nasBuffer);
    DecodeOptionalIEIs(nasBuffer);
  }
  void Encode(NasBuffer& nasBuffer) const override {

    EncodeNasHeader(nasBuffer);
    m_ngKSI.Encode(nasBuffer);
    m_fgsRegType.Encode(nasBuffer); 
    m_fgsMobileIdentity.Encode(nasBuffer);
    EncodeOptionalIEIs(nasBuffer);
  }
};

class SessionManagementMessage : public NasMessage {
 private:
  PduSessionIdentityIE           m_psi;
  ProcedureTransactionIdentityIE m_pti;

  public:
  SessionManagementMessage() {
  }
  void SetPduSessionIdentity(uint8_t psi) {
    m_psi.SetPduSessionIdentity(psi);
  }
  uint8_t GetPduSessionIdentity() const {
    return m_psi.GetPduSessionIdentity();
  }
  void SetProcedureTransactionIdentity(uint8_t pti) {
    m_pti.SetProcedureTransactionIdentity(pti);
  }
  uint8_t GetProcedureTransactionIdentity() const {
    return m_pti.GetProcedureTransactionIdentity();
  }

  void DecodeNasHeader(const NasBuffer& nasBuffer) override{
    m_epd.Decode(nasBuffer);
    m_psi.Decode(nasBuffer);
    m_pti.Decode(nasBuffer);
    m_msgType.Decode(nasBuffer);
  }

  void EncodeNasHeader(NasBuffer& nasBuffer) const override{
      m_epd.Encode(nasBuffer);
      m_psi.Encode(nasBuffer);
      m_pti.Encode(nasBuffer);
      m_msgType.Encode(nasBuffer);
  }

};

class PDUSessionEstablishmentRequest : public SessionManagementMessage {
 private:
  IntegrityProtectionMaximumDataRateIE m_ipmdr;

  bool IsOptionalIESupport(InformationElementType type) const {
    std::unordered_set<InformationElementType> supportedOptionalIEIs;
    const auto& it =  supportedOptionalIEIs.find(type);
	  if (it != supportedOptionalIEIs.end()) {
      return true;
    }
    return false;
  }

  void SetMaxDRPUIForDownlink(MaxDRPUIForDownlink mipd) {
    m_ipmdr.SetMaxDRPUIForDownlink(mipd);
  }
  MaxDRPUIForDownlink GetMaxDRPUIForDownlink() const{
    return m_ipmdr.GetMaxDRPUIForDownlink();
  }

  void SetMaxDRPUIForUplink(MaxDRPUIForUplink mipu) {
    m_ipmdr.SetMaxDRPUIForUplink(mipu);
  }
  MaxDRPUIForUplink GetMaxDRPUIForUplink() const{
    return m_ipmdr.GetMaxDRPUIForUplink();
  }


  void Decode(const NasBuffer& nasBuffer) override {
    DecodeNasHeader(nasBuffer);
    m_ipmdr.Decode(nasBuffer);
    DecodeOptionalIEIs(nasBuffer);
  }
  void Encode(NasBuffer& nasBuffer) const override {
    EncodeNasHeader(nasBuffer);
    m_ipmdr.Encode(nasBuffer);
    EncodeOptionalIEIs(nasBuffer);
  }
};


class NasMessageFactory {
public:
  static void DeallocNasMessage(NasMessage** pNasMessage) {
    if(pNasMessage && *pNasMessage) {
      delete (*pNasMessage);
      *pNasMessage = nullptr;
    }
  }
  static NasMessage* AllocNasMessage( ExtendedProtocolDiscriminator epd,
                              MessageType type) {
    NasMessage* pNasMsg = nullptr;

    switch (epd) {
      case ExtendedProtocolDiscriminator::MOBILITY_MANAGEMENT_MESSAGES: {
        pNasMsg = NasMessageFactory::AllocMobilityManagementMessages(type);
        break;
      }
      case ExtendedProtocolDiscriminator::SESSION_MANAGEMENT_MESSAGES: {
        pNasMsg = NasMessageFactory::AllocSessionManagementMessages(type);
        break;
      }
      default: {
        std::cout << "ExtendedProtocolDiscriminator is not Supported" << std::endl;
        break;
      }
    }
    if(pNasMsg) {
      pNasMsg->SetExtendedProtocolDiscriminator(epd);
    }
    return pNasMsg;
  }

  static NasMessage* AllocMobilityManagementMessages(MessageType type) {
    NasMessage* pNasMsg = nullptr;
    switch (type) {
      case MessageType::REGISTRATION_REQUEST...MessageType::DL_NAS_TRANSPORT: {
        pNasMsg = AllocNasMessage(type);
        break;
      }
      default: {
        std::cout << "MobilityManagementMessage is not Supported" << std::endl;
        break;
      }
    }
    return pNasMsg;
  }
  static NasMessage* AllocSessionManagementMessages(MessageType type) {
    NasMessage* pNasMsg = nullptr;
    switch (type) {
      case MessageType::PDU_SESSION_ESTABLISHMENT_REQUEST...MessageType::FIVEG_SM_STATUS: {
        pNasMsg = AllocNasMessage(type);
        break;
      }
      default: {
        std::cout << "SessionManagementMessage is not Supported" << std::endl;
        break;
      }
    }
    return pNasMsg;
  }
private:
  static NasMessage* AllocNasMessage(MessageType type) {
    NasMessage* pNasMsg = nullptr;
    switch (type) {
      case  MessageType::REGISTRATION_REQUEST:
      {
        pNasMsg = new RegistrationRequest();
        break;
      }
      case  MessageType::REGISTRATION_ACCEPT:
      {
      break;
      }
      case  MessageType::REGISTRATION_COMPLETE:
      {
      break;
      }
      case  MessageType::REGISTRATION_REJECT:
      {
      break;
      }
      case  MessageType::DEREGISTRATION_REQUEST_UE_ORIGINATING:
      {
      break;
      }
      case  MessageType::DEREGISTRATION_ACCEPT_UE_ORIGINATING:
      {
      break;
      }
      case  MessageType::DEREGISTRATION_REQUEST_UE_TERMINATED:
      {
      break;
      }
      case  MessageType::DEREGISTRATION_ACCEPT_UE_TERMINATED:
      {
      break;
      }
      case  MessageType::SERVICE_REQUEST:
      {
      break;
      }
      case  MessageType::SERVICE_REJECT:
      {
      break;
      }
      case  MessageType::SERVICE_ACCEPT:
      {
      break;
      }
      case  MessageType::CONFIGURATION_UPDATE_COMMAND:
      {
      break;
      }
      case  MessageType::CONFIGURATION_UPDATE_COMPLETE:
      {
      break;
      }
      case  MessageType::AUTHENTICATION_REQUEST:
      {
      break;
      }
      case  MessageType::AUTHENTICATION_RESPONSE:
      {
      break;
      }
      case  MessageType::AUTHENTICATION_REJECT:
      {
      break;
      }
      case  MessageType::AUTHENTICATION_FAILURE:
      {
      break;
      }
      case  MessageType::AUTHENTICATION_RESULT:
      {
      break;
      }
      case  MessageType::IDENTITY_REQUEST:
      {
      break;
      }
      case  MessageType::IDENTITY_RESPONSE:
      {
      break;
      }
      case  MessageType::SECURITY_MODE_COMMAND:
      {
      break;
      }
      case  MessageType::SECURITY_MODE_COMPLETE:
      {
      break;
      }
      case  MessageType::SECURITY_MODE_REJECT:
      {
      break;
      }
      case  MessageType::FIVEG_MM_STATUS:
      {
      break;
      }
      case  MessageType::NOTIFICATION:
      {
      break;
      }
      case  MessageType::NOTIFICATION_RESPONSE:
      {
      break;
      }
      case  MessageType::UL_NAS_TRANSPORT:
      {
      break;
      }
      case  MessageType::DL_NAS_TRANSPORT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
      {
        pNasMsg = new PDUSessionEstablishmentRequest();
        break;
      }
      case  MessageType::PDU_SESSION_ESTABLISHMENT_ACCEPT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_ESTABLISHMENT_REJECT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_AUTHENTICATION_COMMAND:
      {
      break;
      }
      case  MessageType::PDU_SESSION_AUTHENTICATION_COMPLETE:
      {
      break;
      }
      case  MessageType::PDU_SESSION_AUTHENTICATION_RESULT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_MODIFICATION_REQUEST:
      {
      break;
      }
      case  MessageType::PDU_SESSION_MODIFICATION_REJECT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_MODIFICATION_COMMAND:
      {
      break;
      }
      case  MessageType::PDU_SESSION_MODIFICATION_COMPLETE:
      {
      break;
      }
      case  MessageType::PDU_SESSION_MODIFICATION_COMMAND_REJECT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_RELEASE_REQUEST:
      {
      break;
      }
      case  MessageType::PDU_SESSION_RELEASE_REJECT:
      {
      break;
      }
      case  MessageType::PDU_SESSION_RELEASE_COMMAND:
      {
      break;
      }
      case  MessageType::PDU_SESSION_RELEASE_COMPLETE:
      {
      break;
      }
      case  MessageType::FIVEG_SM_STATUS:
      {
      break;
      }
      default: 
      {
        std::cout << "unsupported nasmessage received\n";
        break;
      }
    }
    if(pNasMsg) {
      std::cout << NasUtils::Enum2String(type) << " msg created.\n";
      pNasMsg->SetMessageType(type);
    }
    return pNasMsg;
  }
public:
static  NasMessage* DecodeNasMessage(const std::string& nasMsgHex, NasCause& cause) {
   NasMessage* pNasMsg = nullptr; 

    std::vector<uint8_t> nasHexBuffer = NasUtils::HexStringToVector(nasMsgHex);
    NasBuffer nasBuffer(nasHexBuffer);
    std::cout << "Hex Buffer: \n" << nasBuffer.ToHexString() << std::endl;
    
    ExtendedProtocolDiscriminator epd = nasBuffer.GetExtendedProtocolDiscriminator();
    MessageType msgType = nasBuffer.GetMessageType(epd);

    std::cout << "ExtendedProtocolDiscriminator :" << NasUtils::Enum2String(epd) 
              << " MessageType :" << NasUtils::Enum2String(msgType) << std::endl;

    pNasMsg = NasMessageFactory::AllocNasMessage(epd, msgType);

    if(!pNasMsg) {
      std::cout << "msg allocation failed" << std::endl;
      cause =  NasCause::NAS_CAUSE_FAILURE;
      return pNasMsg;
    }

    pNasMsg->Decode(nasBuffer);

    return pNasMsg; 
} 

};

}




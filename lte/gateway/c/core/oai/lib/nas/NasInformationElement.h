#pragma once

#include<iostream>

#include "NasBuffer.h"

namespace nas {
    
//**************************Information Elements *********************/

// Interface all NAS Messages
class InformationElement {

 public:
  virtual InformationElementType getInformationElementType() = 0;
  virtual NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false)  = 0;
  virtual void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const  = 0;
};

class ExtendedProtocolDiscriminatorIE: public InformationElement {
 private:
    ExtendedProtocolDiscriminator m_extendedProtocolDiscriminator;
 public:
  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetExtendedProtocolDiscriminator(ExtendedProtocolDiscriminator epd) {
    m_extendedProtocolDiscriminator = epd;
  }
  ExtendedProtocolDiscriminator GetExtendedProtocolDiscriminator() const {
    return m_extendedProtocolDiscriminator;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_extendedProtocolDiscriminator = 
    		static_cast<ExtendedProtocolDiscriminator>(nasBuffer.DecodeU8());
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8(static_cast<uint8_t>(m_extendedProtocolDiscriminator));
  }
};
class SpareHalfOctetIE: public InformationElement {
 private:
    uint8_t m_spareHalfOctet = 0x0;
 public:
  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetSpareHalfOctet(uint8_t sho) {
      m_spareHalfOctet = sho;
  }
  uint8_t GetSpareHalfOctet() const {
      return m_spareHalfOctet;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_spareHalfOctet =  nasBuffer.DecodeU8UpperNibble();
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8UpperNibble(m_spareHalfOctet);
  }
};

class  SecurityHeaderTypeIE: public InformationElement {
  SecurityHeaderType    m_securityHdrType = SecurityHeaderType::NOT_PROTECTED;

 public:
  InformationElementType getInformationElementType() override {
      return InformationElementType::IEI_UNKNOWN;
  }
  
  void SetSecurityHeaderType(SecurityHeaderType sht) { m_securityHdrType = sht; }
  SecurityHeaderType GetSecurityHeaderType() const{ return m_securityHdrType; }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override {
    m_securityHdrType = static_cast<SecurityHeaderType>(nasBuffer.DecodeU8LowerNibble());
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8LowerNibble(static_cast<uint8_t>(m_securityHdrType));
  }
};

class MessageTypeIE: public InformationElement {
 private:
    MessageType         m_msgType ;
 public:
 
  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetMessageType(MessageType type) { m_msgType = type; }
  MessageType GetMessageType() const { return m_msgType; }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_msgType =  static_cast<MessageType>(nasBuffer.DecodeU8());
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8(static_cast<uint8_t>(m_msgType));
  }
};

#define NAS_SECURITY_MAC_LENGTH 4
class MessageAuthenticationCodeIE: public InformationElement {
 private:
    std::vector<uint8_t>      m_messageAuthenticationCode;
 public:
  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetMessageAuthenticationCode(std::vector<uint8_t> mac) {
    m_messageAuthenticationCode = mac;      
  }
  std::vector<uint8_t> GetMessageAuthenticationCode() const{
    return m_messageAuthenticationCode;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_messageAuthenticationCode = nasBuffer.DecodeU8Vector(NAS_SECURITY_MAC_LENGTH);
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8Vector(m_messageAuthenticationCode);
  }
};

class SequenceNumberIE: public InformationElement {
 private:
    uint8_t m_SequenceNumber = 0x0;
 public:
  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetSequenceNumber(uint8_t sqn) {
    m_SequenceNumber = sqn;
  }
  uint8_t GetSequenceNumber() const{
    return m_SequenceNumber;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_SequenceNumber = nasBuffer.DecodeU8();
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8(m_SequenceNumber);
  }
};

class NasKeySetIdentifierIE : public InformationElement {
private:
	bool         m_securityContextTypeFlag = false;
  uint8_t 		 m_ngKSI = 0;

public:
  NasKeySetIdentifierIE(){}

	InformationElementType getInformationElementType() override {
		return InformationElementType::IEI_UNKNOWN;
	}
    void SetSecurityContextTypeFlag(bool tsc) {
		m_securityContextTypeFlag = tsc;
	}
	bool GetSecurityContextTypeFlag() const {
		return m_securityContextTypeFlag;
	}

	void SetNasKeySetIdentifier(uint8_t ngKSI) {
		m_ngKSI = ngKSI;
	}
	uint8_t GetNasKeySetIdentifier() const{
		return m_ngKSI;
	}
  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
      uint8_t upperNibble = nasBuffer.DecodeU8UpperNibble();
      SetSecurityContextTypeFlag((upperNibble >> 3) & 1 );
      SetNasKeySetIdentifier(upperNibble & 0x7);
      return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override{
      uint8_t ngKsi =  (m_securityContextTypeFlag << 3) | m_ngKSI ; 
      nasBuffer.EncodeU8UpperNibble(ngKsi);
  }
};

class MobileIdentityIE : public InformationElement {
private:
	MobileIdentityType    m_IdentityType = MobileIdentityType::NO_IDENTITY;
	std::vector<uint8_t>  m_IdentityBuf;
public:
  MobileIdentityIE(){}

	InformationElementType getInformationElementType() override {
		return InformationElementType::IEI_UNKNOWN;
	}
	void SetMobileIdenityType(MobileIdentityType type) {
		m_IdentityType = type;
	}
	MobileIdentityType GetMobileIdenityType() {
		return m_IdentityType;
	}

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    uint16_t l = nasBuffer.DecodeU16();
    m_IdentityBuf = nasBuffer.DecodeU8Vector(l);
    SetMobileIdenityType( static_cast<MobileIdentityType> (m_IdentityBuf[0] & 0x07) );
    return NasCause::NAS_CAUSE_SUCCESS;
  } 
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override{
      uint16_t l = m_IdentityBuf.size();
      nasBuffer.EncodeU16(l);
      nasBuffer.EncodeU8Vector(m_IdentityBuf);
  }
};

class FiveGSRegistrationTypeIE  : public InformationElement {
  private:
    bool             m_followOnRequestBit = false;
    RegistrationType m_registrationType   = RegistrationType::INITIAL_REGISTRATION;

  public:
  FiveGSRegistrationTypeIE(){}

  InformationElementType getInformationElementType() override {
		return InformationElementType::IEI_UNKNOWN;
	}
  
  void SetFollowOnRequestBit(bool forb) {
     m_followOnRequestBit = forb;
  }
  bool GetFollowOnRequestBit() const{
      return m_followOnRequestBit;
  }

  void SetRegistrationType(RegistrationType regType) {
      m_registrationType = regType;
  }
  RegistrationType GetRegistrationType() const {
      return m_registrationType;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override {
      uint8_t lowerNibble = nasBuffer.DecodeU8LowerNibble();
      SetFollowOnRequestBit((lowerNibble >> 3) & 1 );
      SetRegistrationType(static_cast<RegistrationType>(lowerNibble & 0x7));
      return NasCause::NAS_CAUSE_SUCCESS;
  }

  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
      uint8_t regType =  (m_followOnRequestBit << 3) | 
                        (static_cast<uint8_t>(m_registrationType)); 
      nasBuffer.EncodeU8LowerNibble(regType);
  }
};
class PduSessionIdentityIE: public InformationElement {
 private:
    uint8_t  m_pduSessionIdentity;
 public:

  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetPduSessionIdentity(uint8_t psi) {
    m_pduSessionIdentity = psi;
  }
  uint8_t GetPduSessionIdentity() const {
    return m_pduSessionIdentity;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_pduSessionIdentity =  nasBuffer.DecodeU8();
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8(m_pduSessionIdentity);
  }
};


class ProcedureTransactionIdentityIE: public InformationElement {
 private:
    uint8_t  m_procedureTransactionIdentity;
 public:

  InformationElementType getInformationElementType() {
    return InformationElementType::IEI_UNKNOWN;
  }

  void SetProcedureTransactionIdentity(uint8_t pti) {
    m_procedureTransactionIdentity = pti;
  }
  uint8_t GetProcedureTransactionIdentity() const {
    return m_procedureTransactionIdentity;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override{
    m_procedureTransactionIdentity =  nasBuffer.DecodeU8();
    return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
    nasBuffer.EncodeU8(m_procedureTransactionIdentity);
  }
};

class IntegrityProtectionMaximumDataRateIE : public InformationElement {
  private:
  MaxDRPUIForDownlink m_mipd = MaxDRPUIForDownlink::NULL_NOTE;
  MaxDRPUIForUplink   m_mipu = MaxDRPUIForUplink::NULL_NOTE;

  public:
  IntegrityProtectionMaximumDataRateIE() {}

	InformationElementType getInformationElementType() override {
		return InformationElementType::IEI_UNKNOWN;
	}
  void SetMaxDRPUIForDownlink(MaxDRPUIForDownlink mipd) {
    m_mipd = mipd;
  }
  MaxDRPUIForDownlink GetMaxDRPUIForDownlink() const{
    return m_mipd;
  }

  void SetMaxDRPUIForUplink(MaxDRPUIForUplink mipu) {
    m_mipu = mipu;
  }
  MaxDRPUIForUplink GetMaxDRPUIForUplink() const{
    return m_mipu;
  }

  NasCause Decode(const NasBuffer& nasBuffer, bool isOptionalIE = false) override {
      m_mipd = static_cast<MaxDRPUIForDownlink>(nasBuffer.DecodeU8());
      m_mipu = static_cast<MaxDRPUIForUplink>(nasBuffer.DecodeU8());
      return NasCause::NAS_CAUSE_SUCCESS;
  }
  void Encode(NasBuffer& nasBuffer, bool isOptionalIE = false) const override {
      nasBuffer.EncodeU8(static_cast<uint8_t>(m_mipd));
      nasBuffer.EncodeU8(static_cast<uint8_t>(m_mipu));
  } 
};

class InformationElementFactory {
public:
  static void DeallocInformatonElement(InformationElement* pIEI) {
      if(pIEI ) {
        delete pIEI;
        pIEI = nullptr;
      }
  }
  static InformationElement* AllocInformationElement(InformationElementType type) {
    InformationElement* pIEI = nullptr;
    switch (type) {
      case InformationElementType::IEI_UE_SECURITY_CAPABILITIES :
      {
        //pIEI = new UESecurityCapability();
        break;
      }
      default: {
        break;
      }
    }
    return pIEI;
  }
};

//**************************Information Elements *********************/
}
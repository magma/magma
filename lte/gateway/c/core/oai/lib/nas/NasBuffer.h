#pragma once

#include <vector>
#include <iostream>
#include <sstream>
#include <iomanip>

#include "NasEnum.h"
namespace nas {

class OctetBuffer {
 private:
  std::vector<uint8_t> m_buffer;
  mutable uint32_t m_currentIndex = 0;
 protected:
  uint8_t GetOctet(size_t index) const{
    if(index < m_buffer.size())
      return m_buffer[index];
    return 0;
  }
 public:
  /*
   * TODO
          if constexpr (std::endian::native == std::endian::big)
          {
              // Big endian system
          }
          else
          {
                  //constexpr (std::endian::native == std::endian::little)
              // Little endian system

          }
  */

  OctetBuffer() {}
  OctetBuffer(uint8_t* pBuffer, size_t size) {
    if (pBuffer && size) {
      for (uint32_t i = 0; i < size; ++i) {
        m_buffer.emplace_back(pBuffer[i]);
      }
    }
  }
  OctetBuffer(const std::vector<uint8_t>& nasHexBuffer) {
    m_buffer = nasHexBuffer;
  }

  uint8_t GetCurrentOctet() const{
      return GetOctet(m_currentIndex);
  }

  bool EndOfBuffer() const {
    return m_currentIndex  > m_buffer.size();
  }

  void EncodeU8(const uint8_t& v) {
    m_buffer.emplace_back(v);
  }

  uint8_t DecodeU8() const {
    return m_buffer[m_currentIndex++];
  }

  void EncodeU16(const uint16_t& v) {
    EncodeU8(static_cast<uint8_t>(v >> 8));
    EncodeU8(static_cast<uint8_t>(v));
  }
  uint16_t DecodeU16() const {
    uint16_t v = 0x0;

    v = ((static_cast<uint16_t>(m_buffer[m_currentIndex++]) << 8) |
         static_cast<uint16_t>(m_buffer[m_currentIndex++]));

    return v;
  }

  void EncodeU32(const uint32_t& v) {
    EncodeU8(static_cast<uint8_t>(v >> 24));
    EncodeU8(static_cast<uint8_t>(v >> 16));
    EncodeU8(static_cast<uint8_t>(v >> 8));
    EncodeU8(static_cast<uint8_t>(v));
  }
  uint32_t DecodeU32() {
    uint32_t v = 0x0;
    v = ((static_cast<uint32_t>(m_buffer[m_currentIndex++]) << 24) |
         (static_cast<uint32_t>(m_buffer[m_currentIndex++]) << 16) |
         (static_cast<uint32_t>(m_buffer[m_currentIndex++]) << 8) |
         static_cast<uint32_t>(m_buffer[m_currentIndex++]));

    return v;
  }

  void EncodeU64(const uint64_t& v) {
    EncodeU8(static_cast<uint8_t>(v >> 56U));
    EncodeU8(static_cast<uint8_t>(v >> 48U));
    EncodeU8(static_cast<uint8_t>(v >> 40U));
    EncodeU8(static_cast<uint8_t>(v >> 32U));
    EncodeU8(static_cast<uint8_t>(v >> 24U));
    EncodeU8(static_cast<uint8_t>(v >> 16U));
    EncodeU8(static_cast<uint8_t>(v >> 8U));
    EncodeU8(static_cast<uint8_t>(v));

  }
  uint64_t DecodeU64() {

    uint64_t v =  0x0;
    v = ((static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 56U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 48U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 40U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 32U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 24U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 16U) |
         (static_cast<uint64_t>(m_buffer[m_currentIndex++]) << 8U) |
         static_cast<uint64_t>(m_buffer[m_currentIndex++]));

    return v;
  }

  void EncodeUtf8(const std::string& v) {
    m_buffer.insert(m_buffer.end(), v.begin(), v.end());
  }
  std::string DecodeUtf8(size_t l) const {

    std::string v;
    auto begin = m_buffer.begin() + m_currentIndex;
    auto end = begin + l;
    v.assign(begin, end);

    m_currentIndex += l;
    return v;
  }

  std::vector<uint8_t> DecodeU8Vector(size_t l) const{
      std::vector<uint8_t> v ;
      auto begin = m_buffer.begin() + m_currentIndex;
      auto end   = begin + l;
      v.assign(begin, end);

      m_currentIndex += l;
      return v;
  }
  
  void EncodeU8Vector(const std::vector<uint8_t>& v) {
    m_buffer.insert(m_buffer.end(), v.begin(), v.end());
  }

  /*
    <------Big Endian Byte-------->
    MSBit.....................LSBit
    UpplerNibble     LowerNibble
  */
  uint8_t DecodeU8UpperNibble() const{
      return  (m_buffer[m_currentIndex] >> 4) & 0x0F; 
  }
  void EncodeU8UpperNibble(uint8_t v) {
     m_buffer.emplace_back( (v & 0x0F) << 4 );
  }

  uint8_t DecodeU8LowerNibble() const {
    return  (m_buffer[m_currentIndex++]) & 0x0F;
  }
  void EncodeU8LowerNibble(uint8_t v) {
    uint8_t& b = m_buffer.back();
    b = ((b & 0xF0) | (v & 0x0F));
  }

  std::string ToHexString() {
    std::stringstream ss;

    if (m_buffer.size() <= 0) return ss.str();

    ss << std::hex << std::setfill('0');
    uint32_t i = 0;
    for (auto& ch : m_buffer) {
      ss << std::hex << std::setw(2) << static_cast<int>(ch) << " ";
      if ((i + 1) % 8 == 0) ss << " ";
      if ((i + 1) % 16 == 0) ss << "\n";
      ++i;
    }
    return ss.str();
  }

  void clear() {
    m_buffer.clear();
    m_currentIndex = 0;
  }

};

/*
The different formats (V, LV, T, TV, TLV, LV-E, TLV-E) and
the five categories of information elements (type 1, 2, 3, 4 and 6)

Totally four categories of standard information elements are defined:
- information elements of format V or TV with value part consisting of 1/2 octet
(type 1);
- information elements of format T with value part consisting of 0 octets (type
2);
- information elements of format V or TV with value part that has fixed length
of at least one octet (type 3);
- information elements of format LV or TLV with value part consisting of zero,
one or more octets (type 4)
- information elements of format LV-E or TLV-E with value part consisting of
zero, one or more octets and a maximum of 65535 octets (type 6). This category
is used in EPS only
*/

class NasBuffer : public OctetBuffer {
  public:
    NasBuffer(){}
    NasBuffer(const std::vector<uint8_t>& nasHexBuffer):OctetBuffer(nasHexBuffer) {
    }

    ExtendedProtocolDiscriminator GetExtendedProtocolDiscriminator() {
      ExtendedProtocolDiscriminator epd = 
          static_cast<ExtendedProtocolDiscriminator>(GetOctet(0));
      return epd;
    }
    SecurityHeaderType GetSecurityHeaderType() const{ 
      SecurityHeaderType sht = 
          static_cast<SecurityHeaderType>(GetOctet(1) & 0x0F);
      return sht; 
    }
    MessageType GetMessageType(ExtendedProtocolDiscriminator epd) { 
      MessageType msgType;
      uint8_t index = 0;
     if(ExtendedProtocolDiscriminator::MOBILITY_MANAGEMENT_MESSAGES == epd) { 
        if(SecurityHeaderType::NOT_PROTECTED == GetSecurityHeaderType()) { 
            index = 2;
        }
        else {
            index = 9;
        }
      }
      else {
            index = 3 ;
      }
      msgType = static_cast<MessageType>(GetOctet(index));
      return msgType;
    } 
};
}
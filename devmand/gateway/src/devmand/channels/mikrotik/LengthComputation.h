// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <iostream>
#include <string>

namespace devmand {
namespace channels {
namespace mikrotik {

constexpr bool isLittleEndian() {
  const int value{0x01};
  const void* address = static_cast<const void*>(&value);
  const unsigned char* leastSignificantAddress =
      static_cast<const unsigned char*>(address);
  return *leastSignificantAddress == 0x01;
}

static constexpr unsigned char b2 = 0x80;
static constexpr unsigned char b3 = 0xC0;
static constexpr unsigned char b4 = 0xE0;
static constexpr unsigned char b5 = 0xF0;

// This length computation is defined here:
//   https://wiki.mikrotik.com/wiki/Manual:API
//
// TODO I'm sure this could be done more elegantly so let's visit that in the
// future. htonl etc would prob. be much cleaner. The compiler should get ride
// of the is endian checks at least.
static inline std::string computeLength(uint32_t messageLength) {
  char buffer[5];
  char* len = reinterpret_cast<char*>(&messageLength);
  unsigned int writeLength{0};
  if (messageLength < 0x80) {
    buffer[0] = len[0];
    writeLength = 1;
  } else if (messageLength < 0x4000) {
    if (isLittleEndian()) {
      buffer[0] = static_cast<char>(len[1] | b2);
      buffer[1] = static_cast<char>(len[0]);
    } else {
      buffer[0] = static_cast<char>(len[2] | b2);
      buffer[1] = static_cast<char>(len[3]);
    }
    writeLength = 2;
  } else if (messageLength < 0x200000) {
    if (isLittleEndian()) {
      buffer[0] = static_cast<char>(len[2] | b3);
      buffer[1] = static_cast<char>(len[1]);
      buffer[2] = static_cast<char>(len[0]);
    } else {
      buffer[0] = static_cast<char>(len[1] | b3);
      buffer[1] = static_cast<char>(len[2]);
      buffer[2] = static_cast<char>(len[3]);
    }
    writeLength = 3;
  } else if (messageLength < 0x10000000) {
    if (isLittleEndian()) {
      buffer[0] = static_cast<char>(len[3] | b4);
      buffer[1] = static_cast<char>(len[2]);
      buffer[2] = static_cast<char>(len[1]);
      buffer[3] = static_cast<char>(len[0]);
    } else {
      buffer[0] = static_cast<char>(len[0] | b4);
      buffer[1] = static_cast<char>(len[1]);
      buffer[2] = static_cast<char>(len[2]);
      buffer[3] = static_cast<char>(len[3]);
    }
    writeLength = 4;
  } else {
    if (isLittleEndian()) {
      buffer[0] = static_cast<char>(b5);
      buffer[1] = static_cast<char>(len[3]);
      buffer[2] = static_cast<char>(len[2]);
      buffer[3] = static_cast<char>(len[1]);
      buffer[4] = static_cast<char>(len[0]);
    } else {
      buffer[0] = static_cast<char>(b5);
      buffer[1] = static_cast<char>(len[0]);
      buffer[2] = static_cast<char>(len[1]);
      buffer[3] = static_cast<char>(len[2]);
      buffer[4] = static_cast<char>(len[3]);
    }
    writeLength = 5;
  }
  return std::string(buffer, writeLength);
}

struct ReadLength {
  size_t lengthSize;
  size_t contentLength;
};

static inline ReadLength readLength(const char* buffer, size_t bufferLength) {
  ReadLength length{0, 0};
  char* lengthBytes = reinterpret_cast<char*>(&length.contentLength);

  if (bufferLength >= 1) {
    if ((buffer[0] & b5) == b5) {
      if (bufferLength >= 5) {
        if (isLittleEndian()) {
          lengthBytes[0] = buffer[4];
          lengthBytes[1] = buffer[3];
          lengthBytes[2] = buffer[2];
          lengthBytes[3] = buffer[1];
        } else {
          lengthBytes[0] = buffer[1];
          lengthBytes[1] = buffer[2];
          lengthBytes[2] = buffer[3];
          lengthBytes[3] = buffer[4];
        }
        length.lengthSize = 5;
      }
    } else if ((buffer[0] & b4) == b4) {
      if (bufferLength >= 4) {
        if (isLittleEndian()) {
          lengthBytes[0] = buffer[3];
          lengthBytes[1] = buffer[2];
          lengthBytes[2] = buffer[1];
          lengthBytes[3] = static_cast<char>(buffer[0] & (~b4));
        } else {
          lengthBytes[0] = static_cast<char>(buffer[0] & (~b4));
          lengthBytes[1] = buffer[1];
          lengthBytes[2] = buffer[2];
          lengthBytes[3] = buffer[3];
        }
        length.lengthSize = 4;
      }
    } else if ((buffer[0] & b3) == b3) {
      if (bufferLength >= 3) {
        if (isLittleEndian()) {
          lengthBytes[0] = buffer[2];
          lengthBytes[1] = buffer[1];
          lengthBytes[2] = static_cast<char>(buffer[0] & (~b3));
        } else {
          lengthBytes[0] = static_cast<char>(buffer[0] & (~b3));
          lengthBytes[1] = buffer[1];
          lengthBytes[2] = buffer[2];
        }
        length.lengthSize = 3;
      }
    } else if ((buffer[0] & b2) == b2) {
      if (bufferLength >= 2) {
        if (isLittleEndian()) {
          lengthBytes[0] = buffer[1];
          lengthBytes[1] = static_cast<char>(buffer[0] & (~b2));
        } else {
          lengthBytes[0] = static_cast<char>(buffer[0] & (~b2));
          lengthBytes[1] = buffer[1];
        }
        length.lengthSize = 2;
      }
    } else {
      if (isLittleEndian()) {
        lengthBytes[0] = static_cast<char>(buffer[0]);
      } else {
        lengthBytes[1] = static_cast<char>(buffer[0]);
      }
      length.lengthSize = 1;
    }
  }

  return length;
}

} // namespace mikrotik
} // namespace channels
} // namespace devmand

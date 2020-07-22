/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string>
#include <google/protobuf/message.h>
#include "magma_logging.h"
#include "GrpcMagmaUtils.h"

#define MAGMA_PRINT_GRPC_PAYLOAD "MAGMA_PRINT_GRPC_PAYLOAD"

std::string grpcLoginLevel = get_env_var(MAGMA_PRINT_GRPC_PAYLOAD);

std::string get_env_var(std::string const& key) {
  MLOG(MINFO) << "Checking env var " << key;
  char* val;
  val                = getenv(key.c_str());
  std::string retval = "";
  if (val != NULL) {
    retval = val;
  }
  return std::string(retval);
}

void PrintGrpcMessage(const google::protobuf::Message& message) {
  if (grpcLoginLevel == "1") {
    const google::protobuf::Descriptor* desc = message.GetDescriptor();
    MLOG(MINFO) << "\n"
                << "  " << desc->full_name().c_str() << " {\n"
                << indentText(message.DebugString(), 6) << "  }";
  }
}

std::string indentText(std::string basicString, int indent) {
  std::stringstream iss(basicString);
  std::string blanks(indent, ' ');
  std::string result = "";
  while (iss.good()) {
    std::string SingleLine;
    getline(iss, SingleLine, '\n');
    // skip empty lines
    if (SingleLine == "") {
      continue;
    }
    result += blanks;
    result += SingleLine;
    // do not add \n on the last line
    result += "\n";
  }
  return result;
}

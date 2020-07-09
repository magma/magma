/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include "GRPCReceiver.h"

std::string get_env_var(std::string const& key);

void PrintGrpcMessage(const google::protobuf::Message& message);

std::string indentText(std::string basicString, int indent);

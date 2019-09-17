// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <devmand/channels/cnmaestro/Response.h>

namespace devmand {
namespace channels {
namespace cnmaestro {

Response::Response(const std::string& msg_) : Response(msg_, false) {}
std::string Response::get() const {
  return msg;
}

bool Response::isError() const {
  return isErr;
}

Response::Response(const std::string& msg_, bool isError_)
    : msg(msg_), isErr(isError_) {}

ErrorResponse::ErrorResponse(const std::string& msg_) : Response(msg_, true) {}

} // namespace cnmaestro
} // namespace channels
} // namespace devmand

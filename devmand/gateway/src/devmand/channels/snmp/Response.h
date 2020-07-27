/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <vector>

#include <folly/dynamic.h>

#include <devmand/channels/snmp/Oid.h>

namespace devmand {
namespace channels {
namespace snmp {

class Response {
 public:
  Response(const Oid& _oid, const folly::dynamic& _value);
  Response() = delete;
  ~Response() = default;
  Response(const Response&) = default;
  Response& operator=(const Response&) = default;
  Response(Response&&) = default;
  Response& operator=(Response&&) = default;

  bool isError() const;

  friend bool operator==(const Response& lhs, const Response& rhs) {
    return lhs.oid == rhs.oid and lhs.value == lhs.value;
  }

 public:
  Oid oid;
  folly::dynamic value;
};

class ErrorResponse final : public Response {
 public:
  ErrorResponse(const folly::dynamic& _value);
  ErrorResponse() = delete;
  ~ErrorResponse() = default;
  ErrorResponse(const ErrorResponse&) = default;
  ErrorResponse& operator=(const ErrorResponse&) = default;
  ErrorResponse(ErrorResponse&&) = default;
  ErrorResponse& operator=(ErrorResponse&&) = default;
};

using Responses = std::vector<Response>;

} // namespace snmp
} // namespace channels
} // namespace devmand

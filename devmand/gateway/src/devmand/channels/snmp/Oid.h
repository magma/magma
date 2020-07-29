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

#include <net-snmp/net-snmp-config.h>
#include <net-snmp/net-snmp-includes.h>
#undef FREE // snmp...
#undef READ
#undef WRITE

#include <folly/Conv.h>
#include <folly/GLog.h>

namespace devmand {
namespace channels {
namespace snmp {

class Oid final {
 public:
  Oid(oid* buf, size_t len);
  Oid(const std::string& repr);
  Oid() = delete;
  ~Oid() = default;
  Oid(const Oid&) = default;
  Oid& operator=(const Oid&) = delete;
  Oid(Oid&&) = default;
  Oid& operator=(Oid&&) = delete;

 public:
  const oid* const get() const;
  size_t getLength() const;
  std::string toString() const;

  friend bool operator==(const Oid& lhs, const Oid& rhs) {
    return snmp_oid_compare(lhs.buffer, lhs.length, rhs.buffer, lhs.length) ==
        0;
  }

  friend bool operator<(const Oid& lhs, const Oid& rhs) {
    return snmp_oid_compare(lhs.buffer, lhs.length, rhs.buffer, lhs.length) ==
        -1;
  }

  bool isDescendant(const Oid& tree) const {
    return tree.length <= length and
        snmp_oidtree_compare(tree.buffer, tree.length, buffer, length) == 0;
  }

 public:
  static const Oid error;

 private:
  size_t length{MAX_OID_LEN};
  oid buffer[MAX_OID_LEN];
};

} // namespace snmp
} // namespace channels
} // namespace devmand

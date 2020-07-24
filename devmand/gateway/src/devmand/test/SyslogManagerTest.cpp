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

#include <gtest/gtest.h>

#include <devmand/syslog/Manager.h>

namespace devmand {
namespace test {

TEST(SyslogManagerTest, addAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addRemAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.removeIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("2.2.2.2", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupDiffAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("2.2.2.2", "bar");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupIdenAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupIdenDiffAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("1.1.1.1", "bar");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

} // namespace test
} // namespace devmand

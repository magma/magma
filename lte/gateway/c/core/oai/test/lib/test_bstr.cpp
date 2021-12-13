/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <string.h>
#include <gtest/gtest.h>

#include "bstrlib.h"

namespace magma {
namespace lte {

class BstrTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}
};

TEST_F(BstrTest, TestBstrchrp) {
  const char hello[6] = "hello";
  const char h        = 'h';
  const char e        = 'e';

  // Create bstring for hello
  bstring hello_bstr = bfromcstr_with_str_len(hello, 5);

  // search for 'e' from the beginning
  EXPECT_EQ(bstrchrp(hello_bstr, e, 0), 1);

  // search for 'e' from the middle
  EXPECT_EQ(bstrchrp(hello_bstr, e, 2), BSTR_ERR);

  // search for 'e' after the end
  EXPECT_EQ(bstrchrp(hello_bstr, e, 6), BSTR_ERR);

  bdestroy(hello_bstr);
}

TEST_F(BstrTest, TestBtrimws) {
  const char hello[10] = "  hello  ";
  const char space[5]  = "    ";

  // Check trimming from start and end
  bstring hello_bstr = bfromcstr_with_str_len(hello, 9);
  EXPECT_EQ(btrimws(hello_bstr), BSTR_OK);
  EXPECT_EQ(hello_bstr->slen, 5);

  // Check trimming of string with all spaces
  bstring space_bstr = bfromcstr_with_str_len(space, 4);
  EXPECT_EQ(btrimws(space_bstr), BSTR_OK);
  EXPECT_EQ(space_bstr->slen, 0);

  // Check trimming of null string
  EXPECT_EQ(btrimws(nullptr), BSTR_ERR);
  bdestroy(space_bstr);
  bdestroy(hello_bstr);
}

TEST_F(BstrTest, TestBdelete) {
  const char habit[6] = "habit";
  const char bit[4]   = "bit";

  // null check
  EXPECT_EQ(bdelete(nullptr, 1, 2), BSTR_ERR);

  bstring habit_bstr = bfromcstr_with_str_len(habit, 5);
  EXPECT_EQ(bdelete(habit_bstr, 0, 2), BSTR_OK);
  EXPECT_EQ(memcmp(habit_bstr->data, bit, 3), 0);

  bdestroy(habit_bstr);
}

TEST_F(BstrTest, TestBiseqcstrcaseless) {
  const char hello_lower[6]        = "hello";
  const char hello_upper[6]        = "HELLO";
  const char world_lower[6]        = "world";
  const char hello_world_lower[12] = "hello world";

  // null check
  EXPECT_EQ(biseqcstrcaseless(nullptr, nullptr), BSTR_ERR);

  // Create bstring for hello
  const bstring hello_bstr = bfromcstr_with_str_len(hello_lower, 5);

  EXPECT_EQ(biseqcstrcaseless(hello_bstr, hello_lower), 1);

  EXPECT_EQ(biseqcstrcaseless(hello_bstr, hello_upper), 1);

  EXPECT_EQ(biseqcstrcaseless(hello_bstr, world_lower), BSTR_OK);

  EXPECT_EQ(biseqcstrcaseless(hello_bstr, hello_world_lower), BSTR_OK);
  bdestroy(hello_bstr);
}

TEST_F(BstrTest, TestBiseqcaselessblk) {
  const char hello_lower[6]        = "hello";
  const char hello_upper[6]        = "HELLO";
  const char world_lower[6]        = "world";
  const char hello_world_lower[12] = "hello world";

  // null check
  EXPECT_EQ(biseqcaselessblk(nullptr, nullptr, 0), BSTR_ERR);

  // Create bstring for hello
  const bstring hello_bstr = bfromcstr_with_str_len(hello_lower, 5);

  EXPECT_EQ(biseqcaselessblk(hello_bstr, hello_lower, 5), 1);

  EXPECT_EQ(biseqcaselessblk(hello_bstr, hello_upper, 5), 1);

  EXPECT_EQ(biseqcaselessblk(hello_bstr, world_lower, 5), BSTR_OK);

  EXPECT_EQ(biseqcaselessblk(hello_bstr, hello_world_lower, 11), BSTR_OK);
  bdestroy(hello_bstr);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace lte
}  // namespace magma

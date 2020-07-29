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

namespace devmand {

class FileUtils final {
 public:
  FileUtils() = delete;
  ~FileUtils() = delete;
  FileUtils(const FileUtils&) = delete;
  FileUtils& operator=(const FileUtils&) = delete;
  FileUtils(FileUtils&&) = delete;
  FileUtils& operator=(FileUtils&&) = delete;

 public:
  static bool touch(const std::string& filename);
  static void write(const std::string& filename, const std::string& contents);
  static std::string readContents(const std::string& filename);
  static bool mkdir(const std::string& dir);
};

} // namespace devmand

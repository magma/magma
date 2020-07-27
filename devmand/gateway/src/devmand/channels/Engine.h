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

#include <string>

namespace devmand {
namespace channels {

/*
 * This is the common interface all channel engines must implement. A channel
 * engine lives for the duration of the process and maintains state the is
 * independent of individual channels. Not all channels require engines: just
 * ones that maintain per process state.
 */
class Engine {
 public:
  Engine(const std::string& engineName_) : engineName(engineName_) {}

  Engine() = delete;
  virtual ~Engine() = default;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  unsigned long long getNumIterations() const {
    return iterations;
  }

  unsigned long long getNumRequests() const {
    return requests;
  }

  std::string getName() const {
    return engineName;
  }

  void incrementRequests() {
    ++requests;
  }

 protected:
  void incrementIterations() {
    ++iterations;
  }

 private:
  unsigned long long iterations{0};
  unsigned long long requests{0};
  std::string engineName;
};

} // namespace channels
} // namespace devmand

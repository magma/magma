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

#include <devmand/error/ErrorHandler.h>

//#include <folly/experimental/exception_tracer/ExceptionTracer.h>

namespace devmand {

void ErrorHandler::executeWithCatch(
    std::function<void()> runable,
    std::function<void()> onFailure) {
  try {
    runable();
  } catch (const std::exception& e) {
    onFailure();
    LOG(ERROR) << "Caught exception: " << e.what();
    trace();
  } catch (...) {
    onFailure();
    LOG(ERROR) << "Caught unknown exception";
    trace();
  }
}

std::string ErrorHandler::getErrorMsg(
    const std::string& device,
    const std::string& channel,
    const std::string& path,
    const std::string& context) {
  return folly::sformat(errorTemplate, device, channel, path, context);
}

void ErrorHandler::trace() {
  /* TODO libunwind has a valgrind error caused by this. debug later.
  auto exceptions = ::folly::exception_tracer::getCurrentExceptions();
  for (auto& e : exceptions) {
    LOG(ERROR) << e;
  }
  */
}

} // namespace devmand

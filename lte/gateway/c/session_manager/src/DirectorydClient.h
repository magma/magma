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
#pragma once

#include <orc8r/protos/common.pb.h>
#include <orc8r/protos/directoryd.grpc.pb.h>
#include <orc8r/protos/directoryd.pb.h>

#include <memory>
#include <mutex>
#include <unordered_map>

#include "includes/GRPCReceiver.h"
#include "SessionState.h"

namespace magma {
using namespace orc8r;

/**
 * DirectorydClient is the base class for managing interactions with DirectoryD.
 */
class DirectorydClient {
 public:
  virtual ~DirectorydClient() = default;
  /**
   * Update the DirectoryD record
   * @param update_request - request used to update the record
   */
  virtual void update_directoryd_record(
      const UpdateRecordRequest& request,
      std::function<void(Status status, Void)> callback) = 0;
  /**
   * Delete the DirectoryD record for the specified ID
   * @param delelete_request - request used to delete the record
   */
  virtual void delete_directoryd_record(
      const DeleteRecordRequest& request,
      std::function<void(Status status, Void)> callback) = 0;

  /**
   * Get all DirectoryD records
   */
  virtual void get_all_directoryd_records(
      std::function<void(Status status, AllDirectoryRecords)> callback) = 0;
};

/**
 * AsyncDirectorydClient sends asynchronous calls to DirectoryD to retrieve
 * UE information.
 */
class AsyncDirectorydClient : public GRPCReceiver, public DirectorydClient {
 public:
  AsyncDirectorydClient();

  AsyncDirectorydClient(std::shared_ptr<grpc::Channel> directoryd_channel);

  /**
   * Update the DirectoryD record
   * @param update_request - request used to update the record
   */
  void update_directoryd_record(
      const UpdateRecordRequest& request,
      std::function<void(Status status, Void)> callback);
  /**
   * Delete the DirectoryD record for the specified ID
   * @param delelete_request - request used to delete the record
   */
  void delete_directoryd_record(
      const DeleteRecordRequest& request,
      std::function<void(Status status, Void)> callback);

  /**
   * Get all DirectoryD records
   */
  void get_all_directoryd_records(
      std::function<void(Status status, AllDirectoryRecords)> callback);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<GatewayDirectoryService::Stub> stub_;
};

}  // namespace magma

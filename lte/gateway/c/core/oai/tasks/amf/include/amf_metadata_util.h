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

#include "include/amf_client_servicer.h"
#include <memory>
#include <string>

typedef struct amf_metadata_s {
  /* Amf Client Servicer Object */
  std::shared_ptr<magma5g::AmfClientServicer> amf_client_servicer_;

  /* Initialize the client servicer layer */
  void amf_client_servicer_init() {
    auto authentication_client =
        std::make_shared<magma5g::AsyncM5GAuthenticationServiceClient>();

    amf_client_servicer_ =
        std::make_shared<magma5g::AmfClientServicer>(authentication_client);
  }

  std::shared_ptr<magma5g::AmfClientServicer>& get_client_servicer_ref() {
    return amf_client_servicer_;
  }

} amf_metadata_t;

void amf_metadata_intialize(amf_metadata_t* metadata_p);

/*
 * This strucuture is used in itti_create_task
 * Container of amf structure/class at the task level.
 */
extern amf_metadata_t amf_metadata;

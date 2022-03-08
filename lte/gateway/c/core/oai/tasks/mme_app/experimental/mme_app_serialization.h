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

#ifdef __cplusplus
extern "C" {
#endif

// C includes --------------------------------------------------------------
#include "lte/gateway/c/core/oai/include/mme_app_desc.h"

void mme_app_schedule_test_protobuf_serialization(uint num_ues);
void mme_app_test_protobuf_serialization(mme_app_desc_t* mme_app_desc,
                                         uint num_ues);
#ifdef __cplusplus
}
#endif

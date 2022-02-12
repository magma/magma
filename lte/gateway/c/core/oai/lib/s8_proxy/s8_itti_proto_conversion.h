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
#include "feg/protos/s8_proxy.grpc.pb.h"
#include "lte/gateway/c/core/oai/include/s8_messages_types.h"

void get_pco_from_proto_msg(
    const magma::feg::ProtocolConfigurationOptions& proto_pco,
    protocol_configuration_options_t* s8_pco);

void get_qos_from_proto_msg(const magma::feg::QosInformation& proto_qos,
                            bearer_qos_t* bearer_qos);

void get_fteid_from_proto_msg(const magma::feg::Fteid& proto_fteid,
                              fteid_t* pgw_fteid);

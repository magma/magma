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

// "#pragma once" will not work for this file, because this file is included
// in include/messages_def.h, which is included more than once within enum
// and structure in the file intertask_interface_types.h
// See comment in "lte/gateway/c/core/oai/include/messages_def.h" for details

MESSAGE_DEF(AGW_OFFLOAD_REQ, ha_agw_offload_req_t, ha_agw_offload_req)

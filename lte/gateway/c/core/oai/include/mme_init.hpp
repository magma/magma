/**
 * Copyright 2022 The Magma Authors.
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

/*! \file mme_init.hpp
  \brief
  \author
  \company
  \email:
*/
/* TODO(rsarwad): mme_init.hpp is temporary file created to bridge between
 * main(), which is in oai_mme.c file and cpp version of individual tasks.
 * On final conversion, initialization shall be declared within task's
 * header file
 */
#pragma once

#include "lte/gateway/c/core/oai/include/mme_config.hpp"

/** \brief S1AP layer top init
 * @returns -1 in case of failure
 **/
#ifdef __cplusplus
extern "C" {
#endif
status_code_e s1ap_mme_init(const mme_config_t* mme_config);
#ifdef __cplusplus
}
#endif

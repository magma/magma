/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#ifndef FILE_EMM_HEARDERS_SEEN
#define FILE_EMM_HEADERS_SEEN

/*TODO: This file has temporary function declarations to
 * resolve undefined references. Delete
 * this file after moving all the files to c++
 * GH issue: https://github.com/magma/magma/issues/13096
 */
#include <sys/types.h>

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/nas/securityDef.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/queue.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/lib/hashtable/obj_hashtable.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_fsm.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsBearerContextStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsNetworkFeatureSupport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MobileStationClassmark2.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeNetworkCapability.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

#endif /* FILE_EMM_HEADERS_SEEN*/

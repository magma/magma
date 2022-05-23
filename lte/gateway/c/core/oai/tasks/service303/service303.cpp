/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include <stdio.h>
#include <assert.h>
#include <stdarg.h>

#include "lte/gateway/c/core/oai/include/service303.hpp"
#include "orc8r/gateway/c/common/service303/MagmaService.hpp"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "orc8r/protos/service303.pb.h"

using magma::service303::MagmaService;

static MagmaService* magma_service;

void start_service303_server(bstring name, bstring version) {
  magma_service = new MagmaService(bdata(name), bdata(version));
  magma_service->Start();
}

void stop_service303_server(void) {
  magma_service->Stop();
  magma_service->WaitForShutdown();
  delete magma_service;
  magma_service = NULL;
}

void service303_set_application_health(application_health_t health) {
  ServiceInfo::ApplicationHealth appHealthEnum;
  switch (health) {
    case APP_UNKNOWN: {
      appHealthEnum = ServiceInfo::APP_UNKNOWN;
      break;
    }
    case APP_HEALTHY: {
      appHealthEnum = ServiceInfo::APP_HEALTHY;
      break;
    }
    case APP_UNHEALTHY: {
      appHealthEnum = ServiceInfo::APP_UNHEALTHY;
      break;
    }
    default: {
      // invalid state
      assert(false);
    }
  }
  magma_service->setApplicationHealth(appHealthEnum);
}

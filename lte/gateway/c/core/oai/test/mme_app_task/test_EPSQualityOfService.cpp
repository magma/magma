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
#include <gtest/gtest.h>
#include <glog/logging.h>

extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/EpsQualityOfService.h"
}

TEST(test_qos_params_to_eps_qos_for_apnambr_test,
     qos_params_to_eps_qos_for_apnambr_test) {
  EpsQualityOfService eps_qos = {0};

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 28000, 28000, 28000, 28000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 28);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 28);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 28);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 28);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 63000, 63000, 63000, 63000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 63);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 63);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 63);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 63);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 72000, 72000, 72000, 72000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 65);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 65);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 65);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 65);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 568000, 568000, 568000, 568000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 127);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 127);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 127);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 127);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 575000, 575000, 575000, 575000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 127);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 127);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 127);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 127);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 640000, 640000, 640000, 640000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 129);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 129);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 129);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 129);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(
      qos_params_to_eps_qos(1, 704000, 704000, 704000, 704000, &eps_qos, false),
      RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 130);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 130);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 130);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 130);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 8640000, 8640000, 8640000, 8640000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 8800000, 8800000, 8800000, 8800000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 16000000, 16000000, 16000000, 16000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 17000000, 17000000, 17000000, 17000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 20000000, 20000000, 20000000, 20000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 128000000, 128000000, 128000000, 128000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 130000000, 130000000, 130000000, 130000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 150000000, 150000000, 150000000, 150000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 256000000, 256000000, 256000000, 256000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 0);

  ASSERT_EQ(qos_params_to_eps_qos(1, 260000000, 260000000, 260000000, 260000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 1);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 1);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 1);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 1);

  ASSERT_EQ(qos_params_to_eps_qos(1, 300000000, 300000000, 300000000, 300000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 11);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 11);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 11);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 11);

  ASSERT_EQ(qos_params_to_eps_qos(1, 500000000, 500000000, 500000000, 500000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 61);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 61);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 61);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 61);

  ASSERT_EQ(qos_params_to_eps_qos(1, 510000000, 510000000, 510000000, 510000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 62);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 62);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 62);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 62);

  ASSERT_EQ(qos_params_to_eps_qos(1, 600000000, 600000000, 600000000, 600000000,
                                  &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 71);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 71);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 71);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 71);

  ASSERT_EQ(qos_params_to_eps_qos(1, 1500000000, 1500000000, 1500000000,
                                  1500000000, &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 161);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 161);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 161);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 161);

  ASSERT_EQ(qos_params_to_eps_qos(1, 1600000000, 1600000000, 1600000000,
                                  1600000000, &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 162);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 162);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 162);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 162);

  ASSERT_EQ(qos_params_to_eps_qos(1, 2000000000, 2000000000, 2000000000,
                                  2000000000, &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 166);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 166);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 166);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 166);

  ASSERT_EQ(qos_params_to_eps_qos(1, 10000000000, 10000000000, 10000000000,
                                  10000000000, &eps_qos, false),
            RETURNok);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.guarBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.guarBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 246);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 246);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForUL), 246);
  ASSERT_EQ((eps_qos.bitRatesExt2.guarBitRateForDL), 246);
}

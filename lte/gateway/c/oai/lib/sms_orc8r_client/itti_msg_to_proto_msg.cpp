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

#include <iomanip>
#include <iostream>
#include <string>

#include "lte/protos/mconfig/mconfigs.pb.h"
#include "MConfigLoader.h"
#include "itti_msg_to_proto_msg.h"
#include "bstrlib.h"

extern "C" {
#include "ie_to_bytes.h"
}

#define MME_SERVICE "mme"
#define NUM_OF_MMEC_DIGITS 2
#define NUM_OF_MMEGID_DIGITS 4
#define DEFAULT_CSFB_MNC "01"
#define DEFAULT_CSFB_MCC "001"
#define DEFAULT_MME_CODE 1
#define DEFAULT_MME_GID 1

std::string int_to_hex_string(int input, int num_of_digit) {
  std::stringstream stream;
  stream << std::setfill('0') << std::setw(num_of_digit) << std::hex << input;
  return stream.str();
}

static magma::mconfig::MME get_default_mconfig() {
  magma::mconfig::MME mconfig;
  mconfig.set_csfb_mnc(DEFAULT_CSFB_MNC);
  mconfig.set_csfb_mcc(DEFAULT_CSFB_MCC);
  mconfig.set_mme_code(DEFAULT_MME_CODE);
  mconfig.set_mme_gid(DEFAULT_MME_GID);
  return mconfig;
}

static magma::mconfig::MME load_mconfig() {
  magma::mconfig::MME mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(MME_SERVICE, &mconfig)) {
    std::cout << "[ERROR] Unable to load mconfig for mme, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

static std::string make_mme_name() {
  auto mme_mconfig = load_mconfig();

  std::string mnc = mme_mconfig.csfb_mnc();
  std::string mcc = mme_mconfig.csfb_mcc();
  if (mnc.length() == 2) {
    mnc = "0" + mnc;
  }
  std::string mme_code =
      int_to_hex_string(mme_mconfig.mme_code(), NUM_OF_MMEC_DIGITS);
  std::string mme_gid =
      int_to_hex_string(mme_mconfig.mme_gid(), NUM_OF_MMEGID_DIGITS);

  std::string mme_name = ".mmec" + mme_code + ".mmegi" + mme_gid +
                         ".mme.epc.mnc" + mnc + ".mcc" + mcc +
                         ".3gppnetwork.org";

  return mme_name;
}

static std::string get_mme_name() {
  static std::string mme_name = make_mme_name();

  return mme_name;
}

namespace magma {
using namespace lte;

SMOUplinkUnitdata convert_itti_sgsap_uplink_unitdata_to_proto_msg(
    const itti_sgsap_uplink_unitdata_t* msg) {
  SMOUplinkUnitdata ret;
  ret.Clear();

  ret.set_imsi(msg->imsi, msg->imsi_length);

  ret.set_nas_message_container(
      bdata(msg->nas_msg_container), blength(msg->nas_msg_container));

  // optional fields
  if (msg->presencemask & UPLINK_UNITDATA_IMEISV_PARAMETER_PRESENT) {
    ret.set_imeisv(msg->opt_imeisv, msg->opt_imeisv_length);
  }
  if (msg->presencemask & UPLINK_UNITDATA_UE_TIMEZONE_PARAMETER_PRESENT) {
    char ue_time_zone = static_cast<char>(msg->opt_ue_time_zone);
    ret.set_ue_time_zone(&ue_time_zone, IE_LENGTH_UE_TIMEZONE);
  }
  if (msg->presencemask &
      UPLINK_UNITDATA_MOBILE_STATION_CLASSMARK_2_PARAMETER_PRESENT) {
    char mobile_station_classmark2[IE_LENGTH_MOBILE_STATION_CLASSMARK2];
    mobile_station_classmark2_to_bytes(
        &msg->opt_mobilestationclassmark2, mobile_station_classmark2);
    ret.set_mobile_station_classmark2(
        mobile_station_classmark2, IE_LENGTH_MOBILE_STATION_CLASSMARK2);
  }
  if (msg->presencemask & UPLINK_UNITDATA_TAI_PARAMETER_PRESENT) {
    char tai[IE_LENGTH_TAI];
    tai_to_bytes(&msg->opt_tai, tai);
    ret.set_tai(tai, IE_LENGTH_TAI);
  }
  if (msg->presencemask & UPLINK_UNITDATA_ECGI_PARAMETER_PRESENT) {
    char ecgi[IE_LENGTH_ECGI];
    ecgi_to_bytes(&msg->opt_ecgi, ecgi);
    ret.set_e_cgi(ecgi, IE_LENGTH_ECGI);
  }

  return ret;
}

}  // namespace magma

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

#include <thread>
#include <mutex>
#include <chrono>
#include <condition_variable>
#include "OpenflowController.h"
#include "ControllerMain.h"
extern "C" {
#include "log.h"
#include "common_defs.h"
}

using namespace fluid_base;
using namespace fluid_msg;

std::condition_variable cv;
std::mutex cv_mutex;

namespace openflow {

OpenflowController::OpenflowController(
    const char* address, const int port, const int n_workers, bool secure,
    std::shared_ptr<OpenflowMessenger> messenger)
    : OFServer(
          address, port, n_workers, secure,
          OFServerSettings()
              .supported_version(OF_13_VERSION)  // OF 1.3
              .use_hello_elements(true)          // bitmask version negotiation
              .keep_data_ownership(false)),
      running_(true),
      messenger_(messenger) {}

OpenflowController::OpenflowController(
    const char* address, const int port, const int n_workers, bool secure)
    : OpenflowController(
          address, port, n_workers, secure,
          std::shared_ptr<DefaultMessenger>(new DefaultMessenger())) {}

void OpenflowController::register_for_event(
    Application* app, ControllerEventType event_type) {
  event_listeners[event_type].push_back(app);
}

void OpenflowController::stop() {
  running_ = false;
  OFServer::stop();
}

void OpenflowController::message_callback(
    OFConnection* ofconn, uint8_t type, void* data, size_t len) {
  if (type == OFPT_PACKET_IN_TYPE) {
    OAILOG_DEBUG(LOG_GTPV1U, "Openflow controller got packet-in message\n");
    dispatch_event(PacketInEvent(ofconn, *this, data, len));
  } else if (type == OFPT_FEATURES_REPLY_TYPE) {
    OAILOG_INFO(LOG_GTPV1U, "Openflow controller connected to switch \n");
    // Save OF connection for external events
    latest_ofconn_ = ofconn;
    std::lock_guard<std::mutex> lck(cv_mutex);
    cv.notify_all();
    OAILOG_INFO(
        LOG_GTPV1U,
        "Send signal that Controller is connected to switch to all waiting "
        "threads \n");
    dispatch_event(SwitchUpEvent(ofconn, *this, data, len));
  } else if (type == OFPT_ERROR) {
    dispatch_event(
        ErrorEvent(ofconn, reinterpret_cast<struct ofp_error_msg*>(data)));
  } else {
    OAILOG_DEBUG(LOG_GTPV1U, "Openflow controller unknown callback %d\n", type);
  }
}
void OpenflowController::connection_callback(
    OFConnection* ofconn, OFConnection::Event type) {
  if (type == OFConnection::EVENT_CLOSED || type == OFConnection::EVENT_DEAD) {
    OAILOG_ERROR(LOG_GTPV1U, "Openflow controller lost connection to switch\n");
    dispatch_event(SwitchDownEvent(ofconn));
  }
}

void OpenflowController::dispatch_event(const ControllerEvent& ev) {
  if (not running_) {
    throw std::runtime_error(
        "Openflow controller needs to be running before handling an event\n");
    return;
  }
  std::vector<Application*> listeners = event_listeners[ev.get_type()];
  for (auto it = listeners.begin(); it != listeners.end(); it++) {
    ((Application*) (*it))->event_callback(ev, *messenger_);
  }
}

void OpenflowController::inject_external_event(
    std::shared_ptr<ExternalEvent> ev, void* (*cb)(std::shared_ptr<void>) ) {
  if (latest_ofconn_ == NULL) {
#define CONNECTION_EVENT_WAIT_TIME 5
    if (is_controller_connected_to_switch(CONNECTION_EVENT_WAIT_TIME) ==
        RETURNok) {
      OAILOG_INFO(LOG_GTPV1U, "Controller is now connected to switch \n");
    }
  }
  ev->set_of_connection(latest_ofconn_);
  latest_ofconn_->add_immediate_event(cb, ev);
}

bool OpenflowController::is_controller_connected_to_switch(int conn_timeout) {
  /* c++ provided conditional variable is added to wait for
   * conn_timeout seconds to make sure connection is established
   * between Controller and switch before inserting the OVS rules
   */

  OAILOG_INFO(
      LOG_GTPV1U,
      "Openflow controller is waiting for %d seconds to establish connection "
      "with Switch \n",
      conn_timeout);
  std::unique_lock<std::mutex> lck(cv_mutex);
  if (cv.wait_for(lck, std::chrono::seconds(conn_timeout)) ==
      std::cv_status::timeout) {
    OAILOG_CRITICAL(
        LOG_GTPV1U,
        "Failed to connect openflow controller to switch, waited for %d "
        "seconds \n",
        conn_timeout);
    OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNerror);
  };
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

}  // namespace openflow

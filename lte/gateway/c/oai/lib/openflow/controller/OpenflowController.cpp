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

#include "OpenflowController.h"
extern "C" {
#include "log.h"
}

using namespace fluid_base;
using namespace fluid_msg;

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
    OAILOG_DEBUG(LOG_GTPV1U, "Openflow controller connected to switch\n");
    // Save OF connection for external events
    latest_ofconn_ = ofconn;
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
        "Openflow controller needs to be running beforehandling an event\n");
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
    OAILOG_ERROR(
        LOG_GTPV1U, "Null connection on event type %d", ev->get_type());
    throw std::runtime_error("Controller not connected to switch:\n");
  }
  ev->set_of_connection(latest_ofconn_);
  latest_ofconn_->add_immediate_event(cb, ev);
}

}  // namespace openflow

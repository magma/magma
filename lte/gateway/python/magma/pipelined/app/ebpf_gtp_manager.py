"""
Copyright 2025 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Author: Nitin Rajput (coRAN LABS)

eBPF GTP Manager App for Pipelined

This app integrates the eBPF GTP manager into the pipelined service framework.
It's a lightweight wrapper that provides the eBPF GTP functionality as a
pipelined service.
"""

import logging
from ryu.base import app_manager
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.ofproto import ofproto_v1_4

from magma.pipelined.app.base import MagmaController

LOG = logging.getLogger("pipelined.app.ebpf_gtp_manager")


class EbpfGtpManagerController(MagmaController):
    
    APP_NAME = "ebpf_gtp_manager"
    APP_TYPE = "service"
    
    def __init__(self, service_manager, *args, **kwargs):
        super(EbpfGtpManagerController, self).__init__(service_manager, *args, **kwargs)
        self._ebpf_gtp_manager = None
        self._config = kwargs.get('config', {})
        
    def initialize_on_connect(self, datapath):

        LOG.info("eBPF GTP Manager Controller datapath connected")
        
        if self._ebpf_gtp_manager is None:
            ebpf_gtp_config = self._config.get('ebpf_gtp', {})
            if ebpf_gtp_config.get('enabled', False):
                try:
                    from magma.pipelined.ebpf.ebpf_gtp_manager import get_ebpf_gtp_manager

                    self._ebpf_gtp_manager = get_ebpf_gtp_manager(self._config)

                    if self._ebpf_gtp_manager:
                        self._service_manager.ebpf_gtp = self._ebpf_gtp_manager
                        LOG.info("EbpfGtpManagerController: Got eBPF GTP manager singleton")
                        LOG.info("eBPF GTP Manager is ready for operations")
                    else:
                        LOG.error("EbpfGtpManagerController: Failed to get eBPF GTP manager")
                        self._ebpf_gtp_manager = None
                except Exception as e:
                    LOG.error("EbpfGtpManagerController: Error creating eBPF GTP manager: %s", e)
            else:
                LOG.info("EbpfGtpManagerController: eBPF GTP is not enabled")
        else:
            LOG.info("eBPF GTP Manager already initialized and ready for operations")
        
        self.init_finished = True
        
        app_future = self._app_futures.get(self.APP_NAME)
        if app_future and not app_future.done():
            app_future.set_result(self)
    
    def cleanup_on_disconnect(self, datapath):
        LOG.info("eBPF GTP Manager Controller cleanup")
        
        if self._ebpf_gtp_manager:
            try:
                self._ebpf_gtp_manager.cleanup()
                LOG.info("eBPF GTP Manager cleaned up successfully")
            except Exception as e:
                LOG.error(f"Error cleaning up eBPF GTP Manager: {e}")
    
    def delete_all_flows(self, datapath):
        LOG.debug("eBPF GTP Manager: No OpenFlow flows to delete")
    
    @set_ev_cls(ofp_event.EventOFPSwitchFeatures, MAIN_DISPATCHER)
    def switch_features_handler(self, ev):
        LOG.debug("eBPF GTP Manager: Switch features event received")
    
    def get_stats(self):
        
        if self._ebpf_gtp_manager:
            try:
                return self._ebpf_gtp_manager.get_statistics()
            except Exception as e:
                LOG.error(f"Error getting eBPF GTP statistics: {e}")
                return {}
        else:
            return {}
    
    def is_ebpf_gtp_enabled(self):
        return (self._ebpf_gtp_manager is not None and
                self._ebpf_gtp_manager.is_ebpf_gtp_enabled())

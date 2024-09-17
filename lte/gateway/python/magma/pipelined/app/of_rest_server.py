"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from ryu import cfg
from ryu.app import wsgi
from ryu.lib import hub


def configure(pipelined_config):
    CONF = cfg.CONF
    CONF.wsapi_port = pipelined_config['of_server_port']


def start(app_manager):
    webapp = wsgi.start_service(app_manager)
    return hub.spawn(webapp)

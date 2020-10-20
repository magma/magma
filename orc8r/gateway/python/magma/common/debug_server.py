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
import wsgiserver
from flask import Flask, request

def debug(**kwargs):
    # pylint: disable=broad-except,eval-used
    try:
        expr = request.form.get('q')
        if kwargs['namespace']:
            namespace = {
                'manager' : kwargs['namespace']
            }
            response = str(eval(expr, namespace))
        else:
            response = str(eval(expr))
        return response
    except Exception as e:
        return "Exception %s " % str(e)

def run_debug_server(debug_sockpath, namespace=None):
    app = Flask(__name__)
    app.add_url_rule(
        '/debug', 'index', debug,
        defaults={'namespace': namespace},
        methods=['POST',])
    server = wsgiserver.WSGIServer(app)
    server.bind_addr = debug_sockpath
    try:
        server.start()
    finally:
        # When the flask server finishes running, do any other cleanup
        pass

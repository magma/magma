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
from flask import Flask, jsonify


def check_quota_response(**kwargs):
    response = kwargs['response']
    return jsonify({
        'version': 1,
        'payload': {
            'emptyWallet': not response
        }
    })


def setup_flask_server(json_response):
    app = Flask(__name__)

    app.add_url_rule(
        '/', 'index', check_quota_response,
        defaults={'response': json_response}
    )
    app.add_url_rule(
        '/<path:dummy>', 'index', check_quota_response,
        defaults={'response': json_response}
    )
    return app


def run_flask(ip, port, response, exit_callback):
    app = setup_flask_server(response)
    server = wsgiserver.WSGIServer(app, host=ip, port=port)
    try:
        server.start()
    finally:
        # When the flask server finishes running, do any other cleanup
        exit_callback()

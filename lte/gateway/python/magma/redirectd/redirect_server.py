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

import logging
from collections import namedtuple

import wsgiserver
from flask import Flask, redirect, render_template, request
from magma.redirectd.redirect_store import RedirectDict

# Use 404 when subscriber not found, 302 for 'Found' redirect
HTTP_NOT_FOUND = 404
HTTP_REDIRECT = 302

NOT_FOUND_HTML = '404.html'

RedirectInfo = namedtuple('RedirectInfo', ['subscriber_ip', 'server_response'])
ServerResponse = namedtuple(
    'ServerResponse', ['redirect_address', 'http_code'],
)


def flask_redirect(**kwargs):
    """ Check redis for src_ip, redirect if found and send 404 if not """
    response = kwargs['get_redirect_response'](request.remote_addr)
    redirect_info = RedirectInfo(request.remote_addr, response)

    logging.info(
        "Request from %s: sent http code %s - redirected to %s",
        redirect_info.subscriber_ip, response.http_code,
        response.redirect_address,
    )

    if response.http_code is HTTP_NOT_FOUND:
        return render_template(
            response.redirect_address,
            subscriber={'ip': redirect_info.subscriber_ip},
        ), HTTP_NOT_FOUND

    return redirect(response.redirect_address, code=response.http_code)


def setup_flask_server():
    app = Flask(__name__)
    url_dict = RedirectDict()

    def get_redirect_response(src_ip):
        """
        If addr type is IPv4/IPv6 prepend http, if url don't change
        TODO: not sure what to do with SIP_URI
        """
        if src_ip not in url_dict:
            return ServerResponse(NOT_FOUND_HTML, HTTP_NOT_FOUND)

        redirect_addr = url_dict[src_ip].server_address
        if url_dict[src_ip].address_type == url_dict[src_ip].IPv4:
            redirect_addr = 'http://' + redirect_addr + '/'
        elif url_dict[src_ip].address_type == url_dict[src_ip].IPv6:
            redirect_addr = 'http://[' + redirect_addr + ']/'

        return ServerResponse(redirect_addr, HTTP_REDIRECT)

    app.add_url_rule(
        '/',
        'index',
        flask_redirect,
        defaults={'get_redirect_response': get_redirect_response},
    )
    app.add_url_rule(
        '/<path:path>',
        'index',
        flask_redirect,
        defaults={'get_redirect_response': get_redirect_response},
    )
    return app


def run_flask(ip, port, exit_callback):
    """
    Run the flask server. this is a daemon, so it exits when redirectd exits
    """

    app = setup_flask_server()

    server = wsgiserver.WSGIServer(app, host=ip, port=port)
    try:
        server.start()
    finally:
        # When the flask server finishes running, do any other cleanup
        exit_callback()

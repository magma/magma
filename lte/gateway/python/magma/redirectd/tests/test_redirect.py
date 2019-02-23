"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from unittest.mock import MagicMock

from lte.protos.policydb_pb2 import RedirectInformation
from magma.redirectd.redirect_server import HTTP_NOT_FOUND, HTTP_REDIRECT, \
    NOT_FOUND_HTML, RedirectInfo, ServerResponse, setup_flask_server


class RedirectdTest(unittest.TestCase):
    def setUp(self):
        """
        Sets up a test version of the redirect server, mocks scribe/url_dict
        """
        self._scribe_client = MagicMock()
        app = setup_flask_server(self._scribe_client)
        app.config['TESTING'] = True

        test_dict = {
            '192.5.82.1':
                RedirectInformation(
                    support=1,
                    address_type=2,
                    server_address='http://www.example.com/'
                )
        }

        def get_resp(src_ip):
            if src_ip not in test_dict:
                return ServerResponse(NOT_FOUND_HTML, HTTP_NOT_FOUND)
            return ServerResponse(
                test_dict[src_ip].server_address, HTTP_REDIRECT
            )
        # Replaces all url_dict polls with a mocked dict (for all url rules)
        for rule in app.url_map._rules:
            if rule is not None and rule.defaults is not None:
                rule.defaults['get_redirect_response'] = get_resp
        self.client = app.test_client()

    def test_302_homepage(self):
        """
        Assert 302 http response, proper reponse headers with new dest url

        Correct scribe logging
        """
        resp = self.client.get('/', environ_base={'REMOTE_ADDR': '192.5.82.1'})

        self.assertEqual(resp.status_code, HTTP_REDIRECT)
        self.assertEqual(resp.headers['Location'], 'http://www.example.com/')

        self._scribe_client.log_to_scribe.assert_called_with(
            RedirectInfo(
                subscriber_ip='192.5.82.1',
                server_response=ServerResponse(
                    redirect_address='http://www.example.com/',
                    http_code=HTTP_REDIRECT
                )
            )
        )

    def test_302_with_path(self):
        """
        Assert 302 http response, proper reponse headers with new dest url

        Correct scribe logging
        """
        resp = self.client.get('/generate_204',
                               environ_base={'REMOTE_ADDR': '192.5.82.1'})

        self.assertEqual(resp.status_code, HTTP_REDIRECT)
        self.assertEqual(resp.headers['Location'], 'http://www.example.com/')

        self._scribe_client.log_to_scribe.assert_called_with(
            RedirectInfo(
                subscriber_ip='192.5.82.1',
                server_response=ServerResponse(
                    redirect_address='http://www.example.com/',
                    http_code=HTTP_REDIRECT
                )
            )
        )

    def test_404(self):
        """
        Assert 404 http response

        Correct scribe logging
        """
        resp = self.client.get('/', environ_base={'REMOTE_ADDR': '127.0.0.1'})

        self.assertEqual(resp.status_code, HTTP_NOT_FOUND)

        self._scribe_client.log_to_scribe.assert_called_with(
            RedirectInfo(
                subscriber_ip='127.0.0.1',
                server_response=ServerResponse(
                    redirect_address='404.html', http_code=404
                )
            )
        )

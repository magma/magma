from __future__ import absolute_import

from graphql.execution import ExecutionResult
from graphql.language.printer import print_ast

from .http import HTTPTransport
import json

class RequestsHTTPTransport(HTTPTransport):
    def __init__(self, session, url, auth=None, use_json=False, timeout=None, **kwargs):
        """
        :param session: The session
        :param auth: Auth tuple or callable to enable Basic/Digest/Custom HTTP Auth
        :param use_json: Send request body as JSON instead of form-urlencoded
        :param timeout: Specifies a default timeout for requests (Default: None)
        """
        super(RequestsHTTPTransport, self).__init__(url, **kwargs)
        self.session = session
        self.auth = auth
        self.default_timeout = timeout
        self.use_json = use_json

    def execute(self, document, variable_values=None, timeout=None):
        query_str = print_ast(document)
        payload = {
            'query': query_str,
            'variables': variable_values or {}
        }

        request = self.session.post(
            self.url,
            data=json.dumps(payload).encode('utf-8'),
            headers=self.headers)
        request.raise_for_status()

        result = request.json()

        extensions = {}
        if "x-correlation-id" in request.headers:
            extensions["trace_id"] = request.headers["x-correlation-id"]

        assert 'errors' in result or 'data' in result, \
            'Received non-compatible response "{}"'.format(result)
        return ExecutionResult(
            errors=result.get('errors'),
            data=result.get('data'),
            extensions=extensions
        )

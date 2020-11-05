#!/usr/bin/env python

import SimpleHTTPServer
import SocketServer

PORT = 80


class GetHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):

    def do_GET(self):
        self.send_head()
        for h in self.headers:
            self.send_header(h, self.headers[h])
        self.end_headers()
        self.send_response(200, "")


Handler = GetHandler
httpd = SocketServer.TCPServer(("", PORT), Handler)
print("starting server")
httpd.serve_forever()

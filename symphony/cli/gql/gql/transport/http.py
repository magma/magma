#!/usr/bin/env python3

from typing import Dict, Optional

from .transport import Transport


class HTTPTransport(Transport):
    def __init__(self, url: str, headers: Optional[Dict[str, str]] = None) -> None:
        self.url: str = url
        self.headers = headers

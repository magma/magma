#!/usr/bin/env python3
from typing import Union

import requests


class Client:
    def __init__(self, endpoint, headers=None):
        self.endpoint = endpoint

        headers = headers or {}
        self.__headers = {
            **headers,
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Encoding': 'gzip',
        }

    def call(self, query,
             variables=None,
             return_json=False) -> Union[dict, str]:

        headers = self.__headers.copy()

        payload = {
            'query': query
        }
        if variables:
            payload['variables'] = variables

        response = requests.post(self.endpoint, json=payload, headers=headers)
        response.raise_for_status()
        return response.json() if return_json else response.text

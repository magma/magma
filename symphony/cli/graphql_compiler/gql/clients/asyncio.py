from typing import Callable, Mapping, Union

import aiohttp


class AsyncIOClient:
    def __init__(self, endpoint, headers=None):
        self.endpoint = endpoint

        headers = headers or {}
        self.__headers = {
            **headers,
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Encoding': 'gzip',
        }
        self.__session = None

    @property
    def session(self):
        if not self.__session or self.__session.closed:
            self.__session = aiohttp.ClientSession(headers=self.__headers)

        return self.__session

    async def call(self, query,
                   variables=None,
                   return_json=False,
                   on_before_callback: Callable[[Mapping[str, str], Mapping[str, str]], None] = None) -> Union[dict, str]:

        headers = self.__headers.copy()

        payload = {
            'query': query
        }
        if variables:
            payload['variables'] = variables

        if on_before_callback:
            on_before_callback(payload, headers)

        async with self.session.post(self.endpoint, json=payload, headers=headers, raise_for_status=True) as resp:
            if return_json:
                return await resp.json()

            return await resp.text()

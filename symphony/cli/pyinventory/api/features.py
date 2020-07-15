#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List

from pysymphony import SymphonyClient

from ..common.constant import GET_FEATURES_URL, SET_FEATURE_URL
from ..exceptions import assert_ok


def get_enabled_features(client: SymphonyClient) -> List[str]:
    """Returns list of the enabled features that are accessible publicly

        Returns:
            list of feature strings

        Raises:
            AssertionError: error returned by server

        Example:
            ```
            features = client.get_enabled_features()
            ```
    """
    resp = client.get(GET_FEATURES_URL)
    assert_ok(resp)
    return list(map(str, resp.json()["features"]))


def set_feature(client: SymphonyClient, feature_id: str, enabled: bool) -> None:
    """Enable or disable given feature if the feature is publicly accessible

        Args:
            feature_id (str): the feature identifier to set
            enabled (bool): enabled or disabled

        Example:
            ```
            features = client.get_enabled_features()
            client.set_feature(feature[0], False)
            ```
    """
    resp = client.post(SET_FEATURE_URL.format(feature_id), {"enabled": enabled})
    assert_ok(resp)

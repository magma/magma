/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import {findIndex} from 'lodash';

// Get current tab position within tabList from url specified
export function GetCurrentTabPos(url: string, tabItems: Array<string>): number {
  const tabPos = findIndex(tabItems, route =>
    location.pathname.startsWith(url + '/' + route),
  );
  return tabPos != -1 ? tabPos : 0;
}

export const DetailTabItems = ['overview', 'event', 'config'];

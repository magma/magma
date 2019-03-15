/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import querystring from 'querystring';
import {format, parse} from 'url';

export function addQueryParamsToUrl(
  url: string,
  params: {[string]: any},
): string {
  const parsedUrl = parse(url, true /* parseQueryString */);
  if (params) {
    parsedUrl.search = querystring.stringify({
      ...parsedUrl.query,
      ...params,
    });
  }
  return format(parsedUrl);
}

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Entries} from './Tokenizer';

export default function TokenizerBasicPostFetchDecorator<TEntry>(
  response: Entries<TEntry>,
  queryString: string,
  currentTokens: Entries<TEntry>,
): Entries<TEntry> {
  return response.filter(
    entry =>
      entry.label.toLowerCase().includes(queryString.toLowerCase()) &&
      !currentTokens.some(token => token.key === entry.key),
  );
}

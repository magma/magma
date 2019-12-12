/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import grammar from './__generated__/PromQLGrammar.js';
import nearley from 'nearley';

export default function Parser() {
  return new nearley.Parser(nearley.Grammar.fromCompiled(grammar));
}

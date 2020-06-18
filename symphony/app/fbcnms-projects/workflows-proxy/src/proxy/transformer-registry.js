/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import logging from '@fbcnms/logging';

import type {
  TransformerCtx,
  TransformerEntry,
  TransformerRegistrationFun,
} from '../types';

const logger = logging.getLogger(module);

export default async function(
  registrationCtx: TransformerCtx,
  transformFx: Array<TransformerRegistrationFun>,
): Promise<Array<TransformerEntry>> {
  logger.debug(
    `Registering transformer modules: [${JSON.stringify(
      transformFx,
    )}] using context ${JSON.stringify(registrationCtx)}`,
  );

  const transformers: Array<TransformerEntry> = [];
  for (const registrationFun of transformFx) {
    const items: Array<TransformerEntry> = registrationFun(registrationCtx);
    transformers.push(...items);
  }
  logger.debug(`Returning ${JSON.stringify(transformers)}`);
  return transformers;
}

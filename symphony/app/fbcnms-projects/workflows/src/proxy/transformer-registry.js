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

const logger = logging.getLogger(module);

export default async function(registrationCtx) {
  // TODO populate from fs
  const transformerModules = ['./transformers/metadata-taskdef'];
  logger.debug(
    `Registering transformer modules: [${transformerModules}] using context ${JSON.stringify(
      registrationCtx,
    )}`,
  );

  const transformers = [];
  for (const file of transformerModules) {
    const transformerModule = await import(file);
    const registrationFun = transformerModule.default;
    const items = registrationFun(registrationCtx);
    logger.debug(`Registering ${file} with ${items.length} items`);
    transformers.push(...items);
  }
  logger.debug(`Returning ${JSON.stringify(transformers)}`);
  return transformers;
}

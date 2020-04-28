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

import bulk from './transformers/bulk';
import event from './transformers/event';
import metadataTaskdef from './transformers/metadata-taskdef';
import metadataWorkflowdef from './transformers/metadata-workflowdef';
import workflow from './transformers/workflow';

import type {
  TransformerCtx,
  TransformerEntry,
  TransformerRegistrationFun,
} from '../types';

const logger = logging.getLogger(module);

export default async function(
  registrationCtx: TransformerCtx,
): Promise<Array<TransformerEntry>> {
  // TODO populate from fs
  const transformerModules: Array<TransformerRegistrationFun> = [
    bulk,
    event,
    metadataTaskdef,
    metadataWorkflowdef,
    workflow,
  ];
  logger.debug(
    `Registering transformer modules: [${JSON.stringify(
      transformerModules,
    )}] using context ${JSON.stringify(registrationCtx)}`,
  );

  const transformers: Array<TransformerEntry> = [];
  for (const registrationFun of transformerModules) {
    const items: Array<TransformerEntry> = registrationFun(registrationCtx);
    transformers.push(...items);
  }
  logger.debug(`Returning ${JSON.stringify(transformers)}`);
  return transformers;
}

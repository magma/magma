/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {
  sanitizeWorkflowdefAfter,
  sanitizeWorkflowdefBefore,
} from '../metadata-workflowdef.js';

const tenant = 'FACEBOOK';

describe('Workflow def prefixing', () => {
  const testCases = [
    [
      'Decision with nested tasks',
      'nested_tasks_decision',
      'nested_tasks_decision_prefixed',
    ],
  ];

  test.each(testCases)('%s', (_, file, filePrefixed) => {
    const workflowDef = require(`./workflow_defs/${file}.json`);
    const workflowDefPrefixed = require(`./workflow_defs/${filePrefixed}.json`);
    const workflowDefTest = JSON.parse(JSON.stringify(workflowDef));

    sanitizeWorkflowdefBefore(tenant, workflowDefTest);
    expect(workflowDefTest).toStrictEqual(workflowDefPrefixed);

    sanitizeWorkflowdefAfter(tenant, workflowDefTest);
    expect(workflowDefTest).toStrictEqual(workflowDef);
  });
});

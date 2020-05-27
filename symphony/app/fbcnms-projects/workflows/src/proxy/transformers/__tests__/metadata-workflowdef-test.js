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
} from '../metadata-workflowdef';

const tenant = 'FACEBOOK';

describe('Workflow def prefixing', () => {
  const testCases = [
    [
      'Decision with nested tasks',
      require('./workflow_defs/nested_tasks_decision.json'),
      require('./workflow_defs/nested_tasks_decision_prefixed.json'),
    ],
  ];

  test.each(testCases)('%s', (_, workflowDef, workflowDefPrefixed) => {
    const workflowDefTest = JSON.parse(JSON.stringify(workflowDef));

    sanitizeWorkflowdefBefore(tenant, workflowDefTest);
    expect(workflowDefTest).toStrictEqual(workflowDefPrefixed);

    sanitizeWorkflowdefAfter(tenant, workflowDefTest);
    expect(workflowDefTest).toStrictEqual(workflowDef);
  });
});

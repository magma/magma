/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ActionID, OperatorID, TriggerFilterID, TriggerID} from './types';

// A set of valid filter options for any given trigger
type TriggerFilterOption = {
  id: TriggerFilterID,
  name: string,
  operatorIDs: OperatorID[],
};

type TriggerActionOption = {
  actionID: ActionID,
  name: string,
};

// All FilterTriggers for every Trigger.
const TRIGGER_FILTERS: {[TriggerID]: TriggerFilterOption[]} = {
  magma_alert: [
    {
      id: 'alert_name',
      name: "that alert's name",
      operatorIDs: ['containsAny'],
    },
    {
      id: 'alert_gatewayid',
      name: "that alert's gatewayID",
      operatorIDs: ['containsAny', 'containsAll'],
    },
    {
      id: 'alert_networkid',
      name: "that alert's networkID",
      operatorIDs: ['containsAny', 'containsAll'],
    },
  ],
};

const TRIGGER_ACTIONS = {
  magma_alert: [
    {
      actionID: 'magma_reboot_gateway',
      name: 'reboot the gateways',
    },
    {
      actionID: 'magma_silence_alert',
      name: 'silence the alerts for gateways',
    },
  ],
};

// Mapping of operatorID to display string
const OPERATORS: {[OperatorID]: string} = {
  containsAny: 'is any of',
  containsAll: 'is all of',
};

export function getFiltersForTrigger(
  triggerID: TriggerID,
): TriggerFilterOption[] {
  return TRIGGER_FILTERS[triggerID];
}

export function getActionsForTrigger(
  triggerID: TriggerID,
): TriggerActionOption[] {
  return TRIGGER_ACTIONS[triggerID];
}

export function getOperatorDisplayName(operatorID: OperatorID) {
  return OPERATORS[operatorID];
}

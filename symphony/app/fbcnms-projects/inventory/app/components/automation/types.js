/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type FilterData = string | string[];

export type ActionData = string | string[] | null;

// A user-configured TriggerFilter when creating a Rule
export type RuleTriggerFilter = {|
  triggerFilterID: TriggerFilterID,
  operatorID: OperatorID,
  data: FilterData,
|};

// A user-configured action when creating a Rule
export type RuleAction = {|
  actionID: ActionID,
  data: ActionData,
|};

// Set of valid triggers that can occur
export type TriggerID = 'magma_alert';

// Set of valid actions that can be executed
export type ActionID = 'magma_reboot_gateway' | 'magma_silence_alert';

// ID of any supported operator
export type OperatorID = 'containsAny' | 'containsAll';

export type TriggerFilterID =
  | 'alert_gatewayid'
  | 'alert_networkid'
  | 'alert_name';

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import type {ApiRequest} from '../AlarmsApi';
import type {ApiUtil} from '../AlarmsApi';
export type RuleEditorProps<TRule> = {
  apiUtil: ApiUtil,
  rule: ?GenericRule<TRule>,
  // invoked when rule is modified locally
  onRuleUpdated: (rule: GenericRule<TRule>) => void,
  onExit: () => void,
  isNew: boolean,
  // component used to swap rule types, used by AddEditAlert
  ruleTypeSelector?: ?React.Node,
};

export type RuleViewerProps<_TRule> = {};

/**
 * Rules should be mapped to the generic rule for rendering in tables and
 * passing to shared components. The raw rule should be passed to
 * ruletype-specific components such as PrometheusEditor.
 */
export type GenericRule<TRule> = {
  severity: string,
  name: string,
  description: string,
  period: string,
  expression: string,
  /**
   * Type of rule, used for selecting which rule editor to use. Must exist as
   * a key inside the RuleInterfaceMap passed to the Alarms component
   */
  ruleType: string,
  // The original rule which this generic rule is mapped from
  rawRule: TRule,
};

export type RuleInterface<TRule> = {|
  friendlyName: string,
  /**
   * Component to create and edit an alerting rule.
   */
  RuleEditor: React.ComponentType<RuleEditorProps<TRule>>,
  /**
   * Component to be rendered inside of the "view rule" modal. Use this to
   * display information specific to the rule type. If this is not provided,
   * the rule will be rendered as json in a pre tag
   */
  RuleViewer?: React.ComponentType<RuleViewerProps<TRule>>,
  /**
   * Retrieve all rules for this rule type
   */
  getRules: (req: ApiRequest) => Promise<Array<GenericRule<TRule>>>,
  deleteRule: (req: {ruleName: string} & ApiRequest) => Promise<void>,
|};

export type RuleInterfaceMap<TUnion> = {
  [string]: RuleInterface<TUnion>,
};

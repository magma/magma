/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  ActionRow_data,
  ActionRow_data$key,
} from './__generated__/ActionRow_data.graphql';
import type {RuleAction} from './types';

import ActionsAutoComplete from './ActionsAutoComplete';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import nullthrows from '@fbcnms/util/nullthrows';
import {graphql, useFragment} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  control: {
    display: 'inline-block',
    width: 100,
    fontWeight: 'bold',
    textAlign: 'right',
  },
  autoCompleteContainer: {
    width: 400,
  },
  spacing: {
    paddingLeft: 10,
  },
}));

const query = graphql`
  fragment ActionRow_data on ActionsTrigger {
    triggerID
    supportedActions {
      actionID
      dataType
      description
    }
  }
`;

type Props = {|
  trigger: ActionRow_data$key,
  ruleAction: ?RuleAction,
  onChange: RuleAction => void,
|};

export default function ActionRow(props: Props) {
  const {trigger, ruleAction, onChange} = props;
  const classes = useStyles();
  const data: ActionRow_data = useFragment<ActionRow_data>(query, trigger);

  const defaultRuleAction: RuleAction = {
    actionID: nullthrows(data.supportedActions[0]).actionID,
    data: [],
  };

  const thisRuleAction = ruleAction || defaultRuleAction;

  return (
    <>
      <Grid item xs={3} className={classes.control}>
        Then
      </Grid>
      <Grid item xs={9}>
        <Select
          value={thisRuleAction.actionID}
          onChange={({target}) => {
            onChange({
              ...thisRuleAction,
              actionID: target.value,
            });
          }}>
          {data.supportedActions.filter(Boolean).map(supportedAction => (
            <MenuItem
              key={supportedAction.actionID}
              value={supportedAction.actionID}>
              {supportedAction.description}
            </MenuItem>
          ))}
        </Select>
      </Grid>
      <Grid item xs={3} className={classes.control} />
      <Grid item xs={9}>
        <div className={classes.autoCompleteContainer}>
          <ActionsAutoComplete
            value={thisRuleAction.data}
            options={[
              '{{ gatewayID }}',
              '{{ networkID }}',
              'test_network1',
              'test_gateway1',
            ]}
            onChange={(evt, newValue) => {
              onChange({
                ...thisRuleAction,
                data: newValue,
              });
            }}
          />
        </div>
      </Grid>
    </>
  );
}

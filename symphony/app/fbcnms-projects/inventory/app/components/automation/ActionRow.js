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
import React from 'react';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';

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
  first: boolean,
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
  const actionsItems = {};
  data.supportedActions
    .filter(Boolean)
    .forEach(action => (actionsItems[action.actionID] = action.description));

  return (
    <>
      <Grid item xs={3} className={classes.control}>
        {props.first ? null : 'Then'}
      </Grid>
      <Grid item xs={9}>
        <TypedSelect
          value={thisRuleAction.actionID}
          items={actionsItems}
          onChange={actionID => {
            onChange({
              ...thisRuleAction,
              actionID,
            });
          }}
        />
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

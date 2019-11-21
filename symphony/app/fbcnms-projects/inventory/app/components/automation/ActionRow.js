/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ActionID, RuleAction, TriggerID} from './types';

import ActionsAutoComplete from './ActionsAutoComplete';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import {getActionsForTrigger} from './constants';

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

type Props = {
  triggerID: TriggerID,
  action: RuleAction,
  onChange: RuleAction => void,
};

export default function ActionRow(props: Props) {
  const {action} = props;
  const classes = useStyles();

  const validActions = getActionsForTrigger(props.triggerID);

  return (
    <>
      <Grid item xs={3} className={classes.control}>
        Then
      </Grid>
      <Grid item xs={9}>
        <Select
          value={action.actionID}
          onChange={({target}) => {
            props.onChange({
              ...action,
              /* eslint-disable-next-line flowtype/no-weak-types */
              actionID: ((target.value: any): ActionID),
            });
          }}>
          {validActions.map((validAction, i) => (
            <MenuItem key={i} value={validAction.actionID}>
              {validAction.name}
            </MenuItem>
          ))}
        </Select>
      </Grid>
      <Grid item xs={3} className={classes.control} />
      <Grid item xs={9}>
        <div className={classes.autoCompleteContainer}>
          <ActionsAutoComplete
            value={action.data || []}
            options={[
              '{{ gatewayID }}',
              '{{ networkID }}',
              'test_network1',
              'test_gateway1',
            ]}
            onChange={(evt, newValue) => {
              props.onChange({
                ...action,
                data: newValue,
              });
            }}
          />
        </div>
      </Grid>
    </>
  );
}

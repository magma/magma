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
  OperatorID,
  RuleTriggerFilter,
  TriggerFilterID,
  TriggerID,
} from './types';

import {getFiltersForTrigger, getOperatorDisplayName} from './constants';

import ActionsAutoComplete from './ActionsAutoComplete';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {find} from 'lodash';
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

type TriggerFilterProps = {
  triggerID: TriggerID,
  filter: RuleTriggerFilter,
  onChange: RuleTriggerFilter => void,
};

export default function TriggerFilterRow(props: TriggerFilterProps) {
  const classes = useStyles();
  const {filter} = props;

  const validFilters = getFiltersForTrigger(props.triggerID);

  const selectedFilter =
    find(validFilters, f => f.id === filter.triggerFilterID) || validFilters[0];

  const selectedOperatorID = filter.operatorID || selectedFilter.operatorIDs[0];
  const validOperators = selectedFilter.operatorIDs;

  return (
    <>
      <Grid item xs={3} className={classes.control}>
        If
      </Grid>
      <Grid item xs={9}>
        <span>
          <Select
            value={filter.triggerFilterID}
            onChange={({target}) =>
              props.onChange({
                ...filter,
                /* string is guaranteed as TriggerFilterID */
                /* eslint-disable-next-line flowtype/no-weak-types */
                triggerFilterID: ((target.value: any): TriggerFilterID),
              })
            }>
            {validFilters.map((filter, i) => (
              <MenuItem key={i} value={filter.id}>
                {filter.name}
              </MenuItem>
            ))}
          </Select>
        </span>
        <span className={classes.spacing}>
          <Select
            value={selectedOperatorID}
            onChange={({target}) => {
              props.onChange({
                ...filter,
                /* eslint-disable-next-line flowtype/no-weak-types */
                operatorID: ((target.value: any): OperatorID),
              });
            }}>
            {validOperators.map((operatorID, i) => (
              <MenuItem key={i} value={operatorID}>
                {getOperatorDisplayName(operatorID)}
              </MenuItem>
            ))}
          </Select>
        </span>
        <span>
          <div className={classes.autoCompleteContainer}>
            <ActionsAutoComplete
              value={filter.data}
              options={['test_network1', 'test_gateway1']}
              onChange={(evt, newValue) => {
                props.onChange({
                  ...filter,
                  data: newValue,
                });
              }}
            />
          </div>
        </span>
      </Grid>
    </>
  );
}

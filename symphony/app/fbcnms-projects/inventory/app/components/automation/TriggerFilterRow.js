/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {RuleFilter} from './types';
import type {
  TriggerFilterRow_data,
  TriggerFilterRow_data$key,
} from './__generated__/TriggerFilterRow_data.graphql';

import ActionsAutoComplete from './ActionsAutoComplete';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TriggerFilterOperator from './TriggerFilterOperator';

import nullthrows from '@fbcnms/util/nullthrows';
import {find} from 'lodash';
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
  fragment TriggerFilterRow_data on ActionsTrigger {
    triggerID
    supportedFilters {
      filterID
      description
      supportedOperators {
        operatorID
      }
      ...TriggerFilterOperator_data
    }
  }
`;

type TriggerFilterProps = {|
  trigger: TriggerFilterRow_data$key,
  ruleFilter: ?RuleFilter,
  onChange: RuleFilter => void,
|};

export default function TriggerFilterRow(props: TriggerFilterProps) {
  const {trigger, ruleFilter, onChange} = props;
  const classes = useStyles();

  const data: TriggerFilterRow_data = useFragment<TriggerFilterRow_data>(
    query,
    trigger,
  );
  const supportedFilters = data.supportedFilters;
  const defaultFilter = data.supportedFilters[0];
  const defaultOperator = defaultFilter?.supportedOperators[0];

  const defaultRuleFilter: RuleFilter = {
    filterID: defaultFilter?.filterID || '',
    data: [],
    operatorID: defaultOperator?.operatorID || '',
  };

  const thisRuleFilter = ruleFilter || defaultRuleFilter;

  const selectedFilter =
    find(
      supportedFilters,
      filter => filter?.filterID === ruleFilter?.filterID,
    ) || nullthrows(supportedFilters[0]);

  return (
    <>
      <Grid item xs={3} className={classes.control}>
        If
      </Grid>
      <Grid item xs={9}>
        <span>
          <Select
            value={thisRuleFilter.filterID}
            onChange={({target}) =>
              onChange({
                ...thisRuleFilter,
                filterID: target.value,
              })
            }>
            {supportedFilters.map((filter, i) => (
              <MenuItem key={i} value={filter?.filterID}>
                {filter?.description}
              </MenuItem>
            ))}
          </Select>
        </span>
        <span className={classes.spacing}>
          <TriggerFilterOperator
            selectedOperatorID={thisRuleFilter.operatorID || ''}
            filter={selectedFilter}
            onChange={operatorID => onChange({...thisRuleFilter, operatorID})}
          />
        </span>
        <span>
          <div className={classes.autoCompleteContainer}>
            <ActionsAutoComplete
              value={thisRuleFilter.data}
              options={['test_network1', 'test_gateway1']}
              onChange={(evt, newValue) => {
                onChange({
                  ...thisRuleFilter,
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

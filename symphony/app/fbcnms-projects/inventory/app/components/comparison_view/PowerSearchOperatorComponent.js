/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Operator} from './ComparisonViewTypes';

import * as React from 'react';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import {OperatorMap} from './ComparisonViewTypes';
import {dateValues, getOperatorLabel} from './FilterUtils';
import {makeStyles} from '@material-ui/styles';

import classNames from 'classnames';

const useStyles = makeStyles(_theme => ({
  pill: {
    backgroundColor: '#ebedf0',
    borderRadius: '4px',
    padding: '6px 8px',
  },
  operator: {
    margin: '0px 4px',
  },
  root: {
    padding: '4px 32px 4px 12px',
    display: 'flex',
    alignItems: 'center',
  },
  test: {
    padding: '12px',
    backgroundColor: 'red',
  },
}));

export const POWER_SEARCH_OPERATOR_ID = 'power_search_operator_select';

type Props = {
  operator: Operator,
  onOperatorChange?: Operator => void,
};

const PowerSearchOperatorComponent = (props: Props) => {
  const classes = useStyles();
  const {operator, onOperatorChange} = props;

  switch (operator) {
    case 'date_greater_than':
    case 'date_less_than':
      return (
        <div className={classes.operator}>
          <Select
            id={POWER_SEARCH_OPERATOR_ID}
            value={operator}
            margin="none"
            variant="outlined"
            classes={{
              root: classes.root,
            }}
            onChange={event => {
              const v = event.target.value;
              if (
                Object.keys(OperatorMap).indexOf(v) !== -1 &&
                onOperatorChange
              ) {
                // $FlowFixMe
                onOperatorChange(v);
              }
            }}>
            {dateValues.map(option => (
              <MenuItem key={option.value} value={option.value}>
                {option.label}
              </MenuItem>
            ))}
          </Select>
        </div>
      );
    default:
      return (
        <div className={classNames(classes.pill, classes.operator)}>
          <Text>{getOperatorLabel(operator)}</Text>
        </div>
      );
  }
};

export default PowerSearchOperatorComponent;

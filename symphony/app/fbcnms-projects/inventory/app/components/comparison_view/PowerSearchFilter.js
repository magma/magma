/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Operator} from './ComparisonViewTypes';

import * as React from 'react';
import ClearIcon from '@material-ui/icons/Clear';
import IconButton from '@material-ui/core/IconButton';
import PowerSearchOperatorComponent from './PowerSearchOperatorComponent';
import Text from '@fbcnms/ui/components/design-system/Text';
import {getOperatorLabel} from './FilterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    border: '1px solid #dadcde',
    borderRadius: '6px',
    display: 'inline-flex',
    padding: '4px',
  },
  pill: {
    backgroundColor: '#ebedf0',
    borderRadius: '4px',
    padding: '6px 8px',
  },
  filterName: {
    fontSize: '12px',
    lineHeight: '16px',
  },
  operator: {
    margin: '0px 4px',
  },
  filterDesc: {
    fontSize: '12px',
    lineHeight: '28px',
    margin: '0px 4px',
  },
  removeButton: {
    marginLeft: '2px',
    padding: '6px',
  },
  removeButtonIcon: {
    fontSize: '14px',
    color: theme.palette.grey.A200,
  },
}));

type Props = {
  name: string,
  operator: Operator,
  onOperatorChange?: Operator => void,
  value: ?string,
  editMode: boolean,
  input: React.Element<any>,
  onRemoveFilter: () => void,
};

const PowerSearchFilter = (props: Props) => {
  const classes = useStyles();
  const {
    input,
    value,
    editMode,
    name,
    operator,
    onRemoveFilter,
    onOperatorChange,
  } = props;

  if (!editMode && value) {
    return (
      <div className={classes.root}>
        <Text className={classes.filterDesc}>
          {name} {getOperatorLabel(operator)} <b>{value}</b>
        </Text>
        <IconButton onClick={onRemoveFilter} className={classes.removeButton}>
          <ClearIcon className={classes.removeButtonIcon} />
        </IconButton>
      </div>
    );
  }

  return (
    <div className={classes.root}>
      <div className={classes.pill}>
        <Text className={classes.filterName}>{name}</Text>
      </div>
      <PowerSearchOperatorComponent
        operator={operator}
        onOperatorChange={onOperatorChange}
      />
      <div>{input}</div>
    </div>
  );
};

export default PowerSearchFilter;

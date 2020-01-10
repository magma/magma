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
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
import {makeStyles} from '@material-ui/styles';
import type {RuleInterfaceMap} from './RuleInterface';

const useRuleTypeStyles = makeStyles(_theme => ({
  select: {
    textTransform: 'capitalize',
  },
  menuItem: {
    textTransform: 'capitalize',
  },
}));
export default function SelectRuleType<TRuleUnion>({
  ruleMap,
  value,
  onChange,
}: {
  ruleMap: RuleInterfaceMap<TRuleUnion>,
  onChange: string => void,
  value: string,
}) {
  const classes = useRuleTypeStyles();
  const ruleTypes = React.useMemo<Array<{type: string, friendlyName: string}>>(
    () =>
      Object.keys(ruleMap || {}).map(key => ({
        type: key,
        friendlyName: ruleMap[key].friendlyName || key,
      })),
    [ruleMap],
  );

  // if there's < 2 rule types, just stick with the default rule type
  if (ruleTypes.length < 2) {
    return null;
  }

  /**
   * Grid structure is chosen here to match the selected editor's width
   * and padding.
   */
  return (
    <TextField
      label="Rule Type"
      value={value}
      onChange={event => onChange(event.target.value)}
      classes={{root: classes.select}}
      select
      fullWidth>
      {ruleTypes.map(ruleType => (
        <MenuItem
          className={classes.menuItem}
          key={ruleType.type}
          value={ruleType.type}>
          {ruleType.friendlyName}
        </MenuItem>
      ))}
    </TextField>
  );
}

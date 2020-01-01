/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

import type {ApiUtil} from './AlarmsApi';
import type {GenericRule, RuleInterfaceMap} from './RuleInterface';

type Props<TRuleUnion> = {
  apiUtil: ApiUtil,
  ruleMap: RuleInterfaceMap<TRuleUnion>,
  onExit: () => void,
  //TODO rename?
  initialConfig: ?GenericRule<TRuleUnion>,
  isNew: boolean,
  thresholdEditorEnabled?: ?boolean,
  defaultRuleType?: string,
};

const useStyles = makeStyles(theme => ({
  gridContainer: {
    flexGrow: 1,
  },
  editingSpace: {
    height: '100%',
    padding: theme.spacing(3),
  },
}));

export default function AddEditAlert<TRuleUnion>(props: Props<TRuleUnion>) {
  const {isNew, apiUtil, ruleMap, onExit} = props;
  const [rule, setRule] = useState<?GenericRule<TRuleUnion>>(
    props.initialConfig,
  );

  const [selectedRuleType, setSelectedRuleType] = React.useState<string>(
    rule?.ruleType || props.defaultRuleType || 'prometheus',
  );

  const classes = useStyles();
  const {RuleEditor} = ruleMap[selectedRuleType];

  return (
    <Grid
      className={classes.gridContainer}
      container
      spacing={0}
      data-testid="add-edit-alert">
      <Grid className={classes.editingSpace} item xs>
        {isNew && (
          <SelectRuleType
            ruleMap={ruleMap}
            value={selectedRuleType}
            onChange={setSelectedRuleType}
          />
        )}
        <RuleEditor
          apiUtil={apiUtil}
          isNew={isNew}
          onExit={onExit}
          onRuleUpdated={setRule}
          rule={rule}
          //TODO remove this prop once context is created
          thresholdEditorEnabled={props.thresholdEditorEnabled}
        />
      </Grid>
    </Grid>
  );
}

const useRuleTypeStyles = makeStyles(_theme => ({
  select: {
    textTransform: 'capitalize',
  },
  menuItem: {
    textTransform: 'capitalize',
  },
}));
function SelectRuleType<TRuleUnion>({
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
    <Grid container spacing={3}>
      <Grid container item spacing={2}>
        <Grid item xs={12} sm={3}>
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
        </Grid>
      </Grid>
    </Grid>
  );
}

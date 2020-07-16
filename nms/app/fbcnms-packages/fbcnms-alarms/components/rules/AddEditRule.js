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
import React from 'react';
import RuleContext from './RuleContext';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useState} from 'react';

import type {GenericRule} from '../rules/RuleInterface';

type Props<TRuleUnion> = {
  onExit: () => void,
  //TODO rename?
  initialConfig: ?GenericRule<TRuleUnion>,
  isNew: boolean,
  defaultRuleType?: string,
};

const useStyles = makeStyles(_theme => ({
  gridContainer: {
    flexGrow: 1,
  },
}));

export default function AddEditRule<TRuleUnion>(props: Props<TRuleUnion>) {
  const {ruleMap} = useAlarmContext();
  const {isNew, onExit} = props;
  const classes = useStyles();
  const [rule, setRule] = useState<?GenericRule<TRuleUnion>>(
    props.initialConfig,
  );

  const [selectedRuleType, setSelectedRuleType] = React.useState<string>(
    rule?.ruleType || props.defaultRuleType || 'prometheus',
  );

  const {RuleEditor} = ruleMap[selectedRuleType];

  return (
    <RuleContext.Provider
      value={{
        ruleMap: ruleMap,
        ruleType: selectedRuleType,
        selectRuleType: setSelectedRuleType,
      }}>
      <Grid
        className={classes.gridContainer}
        container
        spacing={0}
        data-testid="add-edit-alert">
        <RuleEditor
          isNew={isNew}
          onExit={onExit}
          onRuleUpdated={setRule}
          rule={rule}
        />
      </Grid>
    </RuleContext.Provider>
  );
}

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {RuleAction, RuleTriggerFilter, TriggerID} from './types';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import TriggerFilterRow from './TriggerFilterRow';

import ActionRow from './ActionRow';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(_theme => ({
  control: {
    display: 'inline-block',
    width: 100,
    fontWeight: 'bold',
    textAlign: 'right',
  },
}));

type Props = {
  triggerID: TriggerID,
  onClose: () => void,
  onSave: () => void,
};

export default function ActionsAddDialog(props: Props) {
  const classes = useStyles();

  const [filters, setFilters] = useState<RuleTriggerFilter[]>([
    {
      triggerFilterID: 'alert_gatewayid',
      operatorID: 'containsAny',
      data: '',
    },
  ]);

  const [actions, setActions] = useState<RuleAction[]>([
    {
      actionID: 'magma_reboot_gateway',
      data: null,
    },
  ]);

  const onSave = () => {
    const payload = {
      triggerID: props.triggerID,
      filters,
      actions,
    };
    console.log('Save', payload);
    props.onSave();
  };

  const onChangeFilter = (newFilter, i) => {
    const newFilters = [...filters];
    newFilters[i] = newFilter;
    setFilters(newFilters);
  };

  const onChangeAction = (newAction, i) => {
    const newActions = [...actions];
    newActions[i] = newAction;
    setActions(newActions);
  };

  return (
    <Dialog open={true} onClose={() => props.onClose()} maxWidth="lg">
      <DialogTitle>Reboot a device when an alert is fired</DialogTitle>
      <DialogContent>
        <Grid container spacing={1}>
          <Grid item xs={3} className={classes.control}>
            Whenever
          </Grid>
          <Grid item xs={9}>
            we receive an alert
          </Grid>
          {filters.map((filter, i) => (
            <TriggerFilterRow
              key={i}
              triggerID={props.triggerID}
              filter={filter}
              onChange={newFilter => onChangeFilter(newFilter, i)}
            />
          ))}
          {actions.map((action, i) => (
            <ActionRow
              key={i}
              action={action}
              triggerID={props.triggerID}
              onChange={newAction => onChangeAction(newAction, i)}
            />
          ))}
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => props.onClose()} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

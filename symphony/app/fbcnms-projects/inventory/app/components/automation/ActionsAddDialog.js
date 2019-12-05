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
  ActionsAddDialog_triggerData,
  ActionsAddDialog_triggerData$key,
} from './__generated__/ActionsAddDialog_triggerData.graphql';
import type {RuleAction, RuleFilter} from './types';

import ActionRow from './ActionRow';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import TriggerFilterRow from './TriggerFilterRow';

import {graphql, useFragment} from 'react-relay/hooks';
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

type Props = {|
  trigger: ActionsAddDialog_triggerData$key,
  onClose: () => void,
  onSave: () => void,
|};

const query = graphql`
  fragment ActionsAddDialog_triggerData on ActionsTrigger {
    triggerID
    description
    ...ActionRow_data
    ...TriggerFilterRow_data
  }
`;

export default function ActionsAddDialog(props: Props) {
  const classes = useStyles();
  const data: ActionsAddDialog_triggerData = useFragment<ActionsAddDialog_triggerData>(
    query,
    props.trigger,
  );

  const [filters, setFilters] = useState<RuleFilter[]>([]);

  const [actions, setActions] = useState<RuleAction[]>([]);

  const onSave = () => {
    const payload = {
      triggerID: data.triggerID,
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
      <DialogTitle>Create new Action</DialogTitle>
      <DialogContent>
        <Grid container spacing={1}>
          <Grid item xs={3} className={classes.control}>
            Whenever
          </Grid>
          <Grid item xs={9}>
            {data.description}
          </Grid>
          {filters.map((filter, i) => (
            <TriggerFilterRow
              key={i}
              trigger={data}
              ruleFilter={filter}
              onChange={newFilter => onChangeFilter(newFilter, i)}
            />
          ))}
          {filters.length === 0 ? (
            <TriggerFilterRow
              trigger={data}
              ruleFilter={null}
              onChange={newFilter => onChangeFilter(newFilter, 0)}
            />
          ) : null}
          {actions.map((action, i) => (
            <ActionRow
              key={i}
              trigger={data}
              ruleAction={action}
              onChange={newAction => onChangeAction(newAction, i)}
            />
          ))}
          {actions.length === 0 ? (
            <ActionRow
              trigger={data}
              ruleAction={null}
              onChange={newAction => onChangeAction(newAction, 0)}
            />
          ) : null}
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

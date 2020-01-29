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
import AddActionsRuleMutation from '../../mutations/AddActionsRuleMutation';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import EditActionsRuleMutation from '../../mutations/EditActionsRuleMutation';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
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
  divider: {
    margin: '10px 0',
    width: '100%',
  },
  button: {
    float: 'right',
  },
}));

type Props = {|
  trigger: ActionsAddDialog_triggerData$key,
  onClose: () => void,
  onSave: () => void,
  rule?: {
    id: string,
    name: string,
    ruleFilters: $ReadOnlyArray<?RuleFilter>,
    ruleActions: $ReadOnlyArray<?RuleAction>,
  },
|};

const query = graphql`
  fragment ActionsAddDialog_triggerData on ActionsTrigger {
    triggerID
    description
    ...ActionRow_data
    ...TriggerFilterRow_data
  }
`;

const EMPTY_ITEM = null;

export default function ActionsAddDialog(props: Props) {
  const classes = useStyles();
  const data: ActionsAddDialog_triggerData = useFragment<ActionsAddDialog_triggerData>(
    query,
    props.trigger,
  );

  const rule = props.rule;
  const [name, setName] = useState<string>(rule?.name || '');
  const [filters, setFilters] = useState<(?RuleFilter)[]>(
    rule ? [...rule.ruleFilters] : [EMPTY_ITEM],
  );
  const [actions, setActions] = useState<(?RuleAction)[]>(
    rule ? [...rule.ruleActions] : [EMPTY_ITEM],
  );

  const onSave = () => {
    const input = {
      name: name,
      triggerID: data.triggerID,
      ruleActions: actions.filter(Boolean).map(action => ({
        ...action,
        data: JSON.stringify(action.data),
      })),
      ruleFilters: filters.filter(Boolean).map(filter => ({
        ...filter,
        data: JSON.stringify(filter.data),
      })),
    };
    if (rule) {
      EditActionsRuleMutation(
        {id: rule.id, input},
        {onCompleted: props.onSave},
      );
    } else {
      AddActionsRuleMutation({input}, {onCompleted: props.onSave});
    }
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
    <Dialog open={true} onClose={() => props.onClose()}>
      <DialogTitle>{rule ? 'Edit Rule' : 'Create New Rule'}</DialogTitle>
      <DialogContent>
        <Grid container spacing={1}>
          <Grid item xs={3} className={classes.control}>
            Name
          </Grid>
          <Grid item xs={9}>
            <TextInput
              value={name}
              onChange={({target}) => setName(target.value)}
            />
          </Grid>
          <Grid item xs={3} className={classes.control}>
            Whenever
          </Grid>
          <Grid item xs={9}>
            {data.description}
            <Button
              className={classes.button}
              size="small"
              variant="outlined"
              color="primary"
              onClick={() => setFilters([...filters, EMPTY_ITEM])}>
              Add Condition
            </Button>
          </Grid>
          {filters.map((filter, i) => (
            <TriggerFilterRow
              key={i}
              trigger={data}
              ruleFilter={filter}
              onChange={newFilter => onChangeFilter(newFilter, i)}
            />
          ))}
          <Divider className={classes.divider} />
          <Grid item xs={3} className={classes.control}>
            Do
          </Grid>
          <Grid item xs={9}>
            <Button
              className={classes.button}
              size="small"
              variant="outlined"
              color="primary"
              onClick={() => setActions([...actions, EMPTY_ITEM])}>
              Add Action
            </Button>
          </Grid>
          {actions.map((action, i) => (
            <ActionRow
              key={i}
              first={i === 0}
              trigger={data}
              ruleAction={action}
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

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
  ActionsListCard_actionsRule,
  ActionsListCard_actionsRule$key,
} from './__generated__/ActionsListCard_actionsRule.graphql';
import type {
  ActionsListCard_rulesQuery,
  ActionsListCard_rulesQueryResponse,
} from './__generated__/ActionsListCard_rulesQuery.graphql';

import * as React from 'react';
import ActionsAddDialog from './ActionsAddDialog';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage';
import RemoveActionsRuleMutation from '../../mutations/RemoveActionsRuleMutation';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import {Route} from 'react-router-dom';
import {graphql, useFragment, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRelativePath, useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_theme => ({
  paper: {
    padding: '24px 16px',
  },
}));

const query = graphql`
  query ActionsListCard_rulesQuery {
    actionsRules {
      results {
        ...ActionsListCard_actionsRule
      }
    }
  }
`;

const actionRuleFragment = graphql`
  fragment ActionsListCard_actionsRule on ActionsRule {
    id
    name
    trigger {
      description
      ...ActionsAddDialog_triggerData
    }
    ruleActions {
      actionID
      data
    }
    ruleFilters {
      filterID
      operatorID
      data
    }
  }
`;

export default function ActionsListCard() {
  const data: ActionsListCard_rulesQueryResponse = useLazyLoadQuery<ActionsListCard_rulesQuery>(
    query,
  );

  const rules = (data.actionsRules?.results || []).filter(Boolean);
  if (rules.length === 0) {
    return (
      <NoNetworksMessage>
        You currently do not have any actions configured
      </NoNetworksMessage>
    );
  }

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Name</TableCell>
          <TableCell>Description</TableCell>
          <TableCell />
        </TableRow>
      </TableHead>
      <TableBody>
        {rules.map((rule, i) => (
          <RuleRow key={i} rule={rule} />
        ))}
      </TableBody>
    </Table>
  );
}

function RuleRow(props: {rule: ActionsListCard_actionsRule$key}) {
  const relativeUrl = useRelativeUrl();
  const relativePath = useRelativePath();
  const {history} = useRouter();
  const classes = useStyles();
  const rule: ActionsListCard_actionsRule = useFragment<ActionsListCard_actionsRule>(
    actionRuleFragment,
    props.rule,
  );

  const enqueueSnackbar = useEnqueueSnackbar();

  const onDelete = rule => {
    RemoveActionsRuleMutation(
      {id: rule.id},
      {
        onCompleted: () => {
          enqueueSnackbar('Rule deleted successfully', {variant: 'success'});
        },
      },
      store => store.delete(rule.id),
    );
  };

  return (
    <>
      <TableRow className={classes.paper}>
        <TableCell>{rule.name}</TableCell>
        <TableCell>{rule.trigger.description}</TableCell>
        <TableCell>
          <IconButton
            color="primary"
            onClick={() => history.push(relativeUrl(`/edit/${rule.id}`))}>
            <EditIcon />
          </IconButton>
          <IconButton color="primary" onClick={() => onDelete(rule)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      <Route
        path={relativePath(`/edit/${rule.id}`)}
        render={() => (
          <ActionsAddDialog
            trigger={rule.trigger}
            rule={{
              id: rule.id,
              name: rule.name,
              ruleActions: rule.ruleActions
                .filter(Boolean)
                .map(a => ({...a, data: JSON.parse(a.data)})),
              ruleFilters: rule.ruleFilters
                .filter(Boolean)
                // T62071472
                // $FlowFixMe v0.118.0+ filter fields may be null
                .map(f => ({...f, data: JSON.parse(f.data)})),
            }}
            onClose={() => history.push(relativeUrl(''))}
            onSave={() => history.push(relativeUrl(''))}
          />
        )}
      />
    </>
  );
}

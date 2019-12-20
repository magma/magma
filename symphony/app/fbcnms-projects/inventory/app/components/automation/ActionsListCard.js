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
  ActionsListCard_rulesQuery,
  ActionsListCard_rulesQueryResponse,
} from './__generated__/ActionsListCard_rulesQuery.graphql';

import * as React from 'react';
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

import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  paper: {
    padding: theme.spacing(3, 2),
  },
}));

const query = graphql`
  query ActionsListCard_rulesQuery {
    actionsRules {
      results {
        id
        name
        trigger {
          description
        }
      }
    }
  }
`;

export default function ActionsListCard() {
  const classes = useStyles();
  const data: ActionsListCard_rulesQueryResponse = useLazyLoadQuery<ActionsListCard_rulesQuery>(
    query,
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
        {rules.map(rule => (
          <TableRow key={rule.id} className={classes.paper}>
            <TableCell>{rule.name}</TableCell>
            <TableCell>{rule.trigger.description}</TableCell>
            <TableCell>
              <IconButton color="primary">
                <EditIcon />
              </IconButton>
              <IconButton color="primary" onClick={() => onDelete(rule)}>
                <DeleteIcon />
              </IconButton>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

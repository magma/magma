/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ActionsAddDialog_triggerData$key} from './__generated__/ActionsAddDialog_triggerData.graphql';
import type {
  Actions_ActionsQuery,
  Actions_ActionsQueryResponse,
} from './__generated__/Actions_ActionsQuery.graphql';

import ActionsAddDialog from './ActionsAddDialog';
import ActionsCard from './ActionsCard';
import ActionsHead from './ActionsHead';
import Button from '@fbcnms/ui/components/design-system/Button';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
import SettingsApplicationsIcon from '@material-ui/icons/SettingsApplications';
import Text from '@fbcnms/ui/components/design-system/Text';

import useRouter from '@fbcnms/ui/hooks/useRouter';
import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  spacer: {
    flexGrow: 1,
  },
  main: {
    margin: theme.spacing(2),
  },
}));

const query = graphql`
  query Actions_ActionsQuery {
    actionsTriggers {
      results {
        triggerID
        description
        ...ActionsAddDialog_triggerData
      }
      count
    }
  }
`;

export default function Actions() {
  const classes = useStyles();
  const [
    triggerToAdd,
    setTriggerToAdd,
  ] = useState<?ActionsAddDialog_triggerData$key>(null);
  const {history} = useRouter();
  const data: Actions_ActionsQueryResponse = useLazyLoadQuery<Actions_ActionsQuery>(
    query,
  );
  const results = data.actionsTriggers?.results;
  if (!results) {
    return null;
  }

  return (
    <div className={classes.root}>
      <ActionsHead>
        <div className={classes.spacer} />
        <div>
          <Button
            variant="text"
            onClick={() => history.push('/automation/actions/list')}>
            View Rules
          </Button>
        </div>
      </ActionsHead>
      <div className={classes.main}>
        <Text variant="h3">Actions</Text>
        <div className={classes.spacer} />
        <Text variant="h6">Whenever...</Text>
        <Grid container spacing={3}>
          {results.map(trigger => (
            <Grid key={trigger?.triggerID} item xs={3}>
              <ActionsCard
                onClick={() => setTriggerToAdd(trigger)}
                icon={<SettingsApplicationsIcon />}
                message={trigger?.description || ''}
              />
            </Grid>
          ))}
        </Grid>
        {triggerToAdd ? (
          <ActionsAddDialog
            trigger={triggerToAdd}
            onClose={() => setTriggerToAdd(null)}
            onSave={() => setTriggerToAdd(null)}
          />
        ) : null}
      </div>
    </div>
  );
}

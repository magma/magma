/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TriggerID} from './types';

import ActionsCard from './ActionsCard';
import ActionsHead from './ActionsHead';
import Button from '@fbcnms/ui/components/design-system/Button';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
import SettingsApplicationsIcon from '@material-ui/icons/SettingsApplications';
import Text from '@fbcnms/ui/components/design-system/Text';

import ActionsAddDialog from './ActionsAddDialog';
import useRouter from '@fbcnms/ui/hooks/useRouter';
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

export default function Actions() {
  const classes = useStyles();
  const [triggerIDToAdd, setTriggerIDToAdd] = useState<?TriggerID>(null);
  const {history} = useRouter();

  return (
    <div className={classes.root}>
      <ActionsHead>
        <div className={classes.spacer} />
        <div>
          <Button
            variant="text"
            onClick={() => history.push('/automation/actions/list')}>
            Upload Locations
          </Button>
        </div>
      </ActionsHead>
      <div className={classes.main}>
        <Text variant="h6">Get Started</Text>
        <Grid container spacing={3}>
          <Grid item xs={3}>
            <ActionsCard
              onClick={() => setTriggerIDToAdd('magma_alert')}
              icon={<SettingsApplicationsIcon />}
              message="Reboot a device when it goes offline"
            />
          </Grid>
        </Grid>
        {triggerIDToAdd ? (
          <ActionsAddDialog
            triggerID={triggerIDToAdd}
            onClose={() => setTriggerIDToAdd(null)}
            onSave={() => setTriggerIDToAdd(null)}
          />
        ) : null}
      </div>
    </div>
  );
}

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
import LinkedDeviceInput from './LinkedDeviceInput';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

type Props = {
  deviceID: string,
  onChange: string => void,
};

const useStyles = makeStyles(theme => ({
  subheader: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  titleText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
}));

export default function LinkedDeviceAddEditSection(props: Props) {
  const classes = useStyles();
  const {deviceID: id, onChange} = props;
  const [deviceID, networkID] = id.split('.');

  return (
    <>
      <div className={classes.subheader}>
        <Text variant="subtitle1" className={classes.titleText}>
          Link to Orchestrator Device
        </Text>
      </div>
      <Grid container spacing={2}>
        <Grid key={'device-id'} item xs={12} sm={12} lg={6} xl={4}>
          <LinkedDeviceInput
            label="Device ID"
            value={deviceID}
            onChange={event =>
              onChange(event.target.value + '.' + (networkID ?? ''))
            }
          />
        </Grid>
        <Grid key={'network-id'} item xs={12} sm={12} lg={6} xl={4}>
          <LinkedDeviceInput
            label="Network ID"
            value={networkID ?? ''}
            onChange={event => onChange(deviceID + '.' + event.target.value)}
          />
        </Grid>
      </Grid>
    </>
  );
}

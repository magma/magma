/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DevicesGateway} from './DevicesUtils';

import AppBar from '@material-ui/core/AppBar';
import DevicesGatewayDevmandFields from './DevicesGatewayDevmandFields';
import Dialog from '@material-ui/core/Dialog';
import GatewayCommandFields from '@fbcnms/magmalte/app/components/GatewayCommandFields';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import nullthrows from '@fbcnms/util/nullthrows';

import {MagmaAPIUrls, fetchDevice} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEffect, useState} from 'react';

const useStyles = makeStyles({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
});

type Props = {
  onClose: () => void,
  onSave: (gateway: {[string]: any}) => void,
  gateway: DevicesGateway,
};

export default function DevicesEditControllerDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const [tab, setTab] = useState(0);
  const [devmandManagedDevices, setDevmandManagedDevices] = useState<string[]>(
    [],
  );

  const {isLoading, error, response} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.devicesDevmandConfigs(
      match,
      nullthrows(props.gateway).id,
    ),
  });

  useEffect(() => {
    if (response) {
      let managedDevices = response.data.managed_devices.filter(
        d => d.length > 0,
      );
      if (managedDevices.length == 0) {
        managedDevices = [''];
      }
      setDevmandManagedDevices(managedDevices);
    }
  }, [response]);

  if (error) {
    return <div>Error: {error}</div>;
  }

  const onTabChange = (event, tab) => setTab(tab);
  const onSave = gatewayID => {
    fetchDevice(match, gatewayID).then(props.onSave);
  };

  let content;

  switch (tab) {
    case 0:
      if (isLoading || devmandManagedDevices.length == 0) {
        content = <LoadingFiller />;
      } else {
        content = (
          <DevicesGatewayDevmandFields
            onClose={props.onClose}
            gateway={props.gateway}
            onSave={onSave}
            devmandManagedDevices={devmandManagedDevices}
            showRestartCommand={true}
            showRebootEnodebCommand={true}
            showPingCommand={true}
            showGenericCommand={true}
          />
        );
      }
      break;
    case 1:
      content = (
        <GatewayCommandFields
          onClose={props.onClose}
          gatewayID={props.gateway.id}
          showRestartCommand={true}
          showRebootEnodebCommand={true}
          showPingCommand={true}
          showGenericCommand={true}
        />
      );
      break;
  }

  return (
    <Dialog open={true} onClose={props.onClose} maxWidth="md" scroll="body">
      <AppBar position="static" className={classes.appBar}>
        <Tabs
          indicatorColor="primary"
          textColor="primary"
          value={tab}
          onChange={onTabChange}>
          <Tab label="Managed Devices" />
          <Tab label="Commands" />
        </Tabs>
      </AppBar>
      {content}
    </Dialog>
  );
}

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {DevicesAgent} from './DevicesUtils';
import type {symphony_agent} from '@fbcnms/magma-api';

import AppBar from '@material-ui/core/AppBar';
import DevicesAgentFields from './DevicesAgentFields';
import Dialog from '@material-ui/core/Dialog';
import FormLabel from '@material-ui/core/FormLabel';
import GatewayCommandFields from '@fbcnms/magmalte/app/components/GatewayCommandFields';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
}));

type Props = {
  onClose: () => void,
  onSave: symphony_agent => void,
  agent: DevicesAgent,
};

export default function DevicesEditAgentDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const [tab, setTab] = useState(0);
  const [devmandManagedDevices, setDevmandManagedDevices] = useState<string[]>(
    [],
  );

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdAgentsByAgentIdManagedDevices,
    {
      networkId: nullthrows(match.params.networkId),
      agentId: props.agent.id,
    },
    useCallback(response => {
      let managedDevices = [];
      if (response) {
        managedDevices = response.filter(d => d.length > 0);
      }
      if (managedDevices.length == 0) {
        managedDevices = [''];
      }
      setDevmandManagedDevices(managedDevices);
    }, []),
  );

  const onTabChange = (event, tab) => setTab(tab);
  const onSave = agentId => {
    MagmaV1API.getSymphonyByNetworkIdAgentsByAgentId({
      networkId: nullthrows(match.params.networkId),
      agentId,
    }).then(props.onSave);
  };

  let content;

  switch (tab) {
    case 0:
      if (isLoading) {
        content = <LoadingFiller />;
      } else {
        content = (
          <DevicesAgentFields
            onClose={props.onClose}
            agent={props.agent}
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
          gatewayID={props.agent.id}
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
      <>
        <FormLabel error>{(error && error.toString()) || ''}</FormLabel>
        {content}
      </>
    </Dialog>
  );
}

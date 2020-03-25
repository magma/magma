/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DevicesAgent} from './DevicesUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import DevicesEditAgentDialog from './DevicesEditAgentDialog';
import DevicesNewAgentDialog from './DevicesNewAgentDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useCallback, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import {map} from 'lodash';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {buildDevicesAgentFromPayload} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

type Props = WithAlert & {};

function DevicesAgents(props: Props) {
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [agents, setAgents] = useState<?Array<DevicesAgent>>(null);
  const [errorMessage, setErrorMessage] = useState<?string>(null);
  const [editingAgent, setEditingAgent] = useState<?DevicesAgent>(null);
  const classes = useStyles();

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdAgents,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => {
      if (response != null) {
        setAgents(
          map(response, (agent, _) =>
            buildDevicesAgentFromPayload(agent),
          ).sort((a, b) => a.id.localeCompare(b.id)),
        );
      }
    }, []),
  );

  if (error || isLoading || !agents) {
    return <LoadingFiller />;
  }

  const onSave = agentPayload => {
    const agent = buildDevicesAgentFromPayload(agentPayload);
    const newAgents = agents.slice(0);
    if (editingAgent) {
      newAgents[newAgents.indexOf(editingAgent)] = agent;
    } else {
      newAgents.push(agent);
    }
    setAgents(newAgents);
    setEditingAgent(null);
  };

  const deleteAgent = agent => {
    if (!agent.id) {
      setErrorMessage('Error: cannot delete because id is empty');
    } else {
      props
        .confirm(`Are you sure you want to delete ${agent.id}?`)
        .then(confirmed => {
          if (!confirmed) {
            return;
          }
          MagmaV1API.deleteSymphonyByNetworkIdAgentsByAgentId({
            networkId: nullthrows(match.params.networkId),
            agentId: agent.id,
          }).then(() => setAgents(agents.filter(a => a.id != agent.id)));
          setErrorMessage(null);
        });
    }
  };

  const rows = agents.map(agent => (
    <TableRow key={agent.id}>
      <TableCell>
        <DeviceStatusCircle
          isGrey={agent.status == null}
          isActive={!!agent.up}
        />
        {agent.id || 'Error: Missing ID'}
      </TableCell>
      <TableCell>
        {agent.hardware_id}
        {agent.devmand_config === undefined && (
          <Text color="error">missing devmand config</Text>
        )}
      </TableCell>
      <TableCell>
        <IconButton color="primary" onClick={() => setEditingAgent(agent)}>
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => deleteAgent(agent)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Configure Agents</Text>
        <NestedRouteLink to="/new">
          <Button>Add Agent</Button>
        </NestedRouteLink>
      </div>
      <Paper elevation={2}>
        {errorMessage && <div>{errorMessage.toString()}</div>}
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Hardware UUID</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      {editingAgent && (
        <DevicesEditAgentDialog
          key={editingAgent.id}
          agent={editingAgent}
          onClose={() => setEditingAgent(null)}
          onSave={onSave}
        />
      )}
      <Route
        path={relativePath('/new')}
        render={() => (
          <DevicesNewAgentDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={rawAgent => {
              setAgents([...agents, buildDevicesAgentFromPayload(rawAgent)]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </div>
  );
}

export default withAlert(DevicesAgents);

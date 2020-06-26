/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddNetworkDialog from './AddNetworkDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DialogWithConfirmationPhrase from '@fbcnms/ui/components/DialogWithConfirmationPhrase';
import EditIcon from '@material-ui/icons/Edit';
import EditNetworkDialog from './EditNetworkDialog';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
  noNetworks: {
    height: '70vh',
  },
}));

function Networks() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {relativePath, relativeUrl, history} = useRouter();
  const [networks, setNetworks] = useState(null);
  const [networkToDelete, setNetworkToDelete] = useState(null);

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworks,
    {},
    useCallback(res => setNetworks(sortBy(res, [n => n.toLowerCase()])), []),
  );

  if (error || isLoading || !networks) {
    return <LoadingFiller />;
  }

  const rows = networks.map(network => (
    <TableRow key={network}>
      <TableCell>{network}</TableCell>
      <TableCell>
        <IconButton
          onClick={() => history.push(relativeUrl(`/edit/${network}`))}>
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => setNetworkToDelete(network)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  const closeDialog = () => history.push(relativeUrl(''));
  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <div />
        <NestedRouteLink to="/new">
          <Button>Add Network</Button>
        </NestedRouteLink>
      </div>
      {rows.length === 0 ? (
        <div className={classes.noNetworks}>
          <NoNetworksMessage>
            You currently do not have any networks configured. Click "Add
            Network" to create a new network
          </NoNetworksMessage>
        </div>
      ) : (
        <Paper elevation={2}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Network ID</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </Paper>
      )}
      {networkToDelete && (
        <DialogWithConfirmationPhrase
          title="Warning!"
          message={
            'Deleting a network is a serious action and cannot be ' +
            'un-done. Please type the Network ID below if you are confident ' +
            'about this action.'
          }
          label="Network ID"
          confirmationPhrase={networkToDelete}
          onClose={() => setNetworkToDelete(null)}
          onConfirm={async () => {
            const payload = {
              networkID: networkToDelete,
            };
            axios
              .post('/nms/network/delete', payload)
              .then(response => {
                if (!response.data.success) {
                  enqueueSnackbar('Network delete failed', {
                    variant: 'error',
                  });
                }
              })
              .catch(() => {
                enqueueSnackbar('Network delete failed', {
                  variant: 'error',
                });
              });
            setNetworks(networks.filter(n => n != networkToDelete));
            setNetworkToDelete(null);
          }}
        />
      )}
      <Route
        path={relativePath('/new')}
        render={() => (
          <AddNetworkDialog
            onClose={closeDialog}
            onSave={networkID => {
              setNetworks([...networks, networkID]);
              enqueueSnackbar('Network created successfully', {
                variant: 'success',
              });
              closeDialog();
            }}
          />
        )}
      />
      <Route
        path={relativePath('/edit/:networkID')}
        render={() => (
          <EditNetworkDialog
            onClose={closeDialog}
            onSave={_ => {
              enqueueSnackbar('Network updated successfully', {
                variant: 'success',
              });
              closeDialog();
            }}
          />
        )}
      />
    </div>
  );
}

export default Networks;

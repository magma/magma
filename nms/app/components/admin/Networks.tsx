/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import AddNetworkDialog from './AddNetworkDialog';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import EditIcon from '@material-ui/icons/Edit';
import EditNetworkDialog from './EditNetworkDialog';
import EmptyState from '../EmptyState';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../LoadingFiller';
import NestedRouteLink from '../NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useCallback, useContext, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import MagmaAPI from '../../api/MagmaAPI';
import NetworkContext from '../context/NetworkContext';
import useMagmaAPI from '../../api/useMagmaAPI';
import {Route, Routes, useNavigate} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'flex-end',
  },
  paper: {
    margin: '10px',
  },
  noNetworks: {
    height: '70vh',
  },
}));

type DialogConfirmationProps = {
  title: string;
  message: string;
  confirmationPhrase: string;
  label: string;
  onClose: () => void;
  onConfirm: () => void | Promise<void>;
};

function DialogWithConfirmationPhrase(props: DialogConfirmationProps) {
  const [confirmationPhrase, setConfirmationPhrase] = useState('');
  const {title, message, label, onClose, onConfirm} = props;

  return (
    <Dialog
      open={true}
      onClose={onClose}
      TransitionProps={{onExited: onClose}}
      maxWidth="sm">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          {message}
          <TextField
            label={label}
            value={confirmationPhrase}
            onChange={({target}) => setConfirmationPhrase(target.value)}
          />
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => void onConfirm()}
          disabled={confirmationPhrase !== props.confirmationPhrase}>
          Confirm
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function Networks() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const navigate = useNavigate();
  const [networks, setNetworks] = useState<Array<string> | null>(null);
  const [networkToDelete, setNetworkToDelete] = useState<string | null>(null);
  const {networkId: selectedNetworkId} = useContext(NetworkContext);

  const {error, isLoading} = useMagmaAPI(
    MagmaAPI.networks.networksGet,
    {},
    useCallback(
      (res: Array<string>) => setNetworks(sortBy(res, [n => n.toLowerCase()])),
      [],
    ),
  );

  if (error || isLoading || !networks) {
    return <LoadingFiller />;
  }

  const rows = networks.map(network => (
    <TableRow key={network}>
      <TableCell>{network}</TableCell>
      <TableCell>
        <IconButton onClick={() => navigate(`edit/${network}`)}>
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => setNetworkToDelete(network)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  const closeDialog = () => navigate('');

  const cardActions = {
    buttonText: 'Add Network',
    onClick: () => navigate('new'),
    linkText: 'Learn more about Networks',
    link: 'https://docs.magmacore.org/docs/next/nms/network',
  };
  return (
    <div className={classes.paper}>
      {rows.length === 0 ? (
        <Grid container justify="space-between" spacing={3}>
          <EmptyState
            title={'Set up a Network'}
            instructions={
              'Add your first network from the Network Selector.  If you already have multiple networks, use the Network Selector to switch networks.'
            }
            cardActions={cardActions}
            overviewTitle={'Network Overview'}
            overviewDescription={
              'In Magma, a Network is formed by a group of subscribers that have access to a group of access-gateways. On the Network dashboard, you can configure your network information, and view the summary of the network components.'
            }
          />
        </Grid>
      ) : (
        <>
          <div className={classes.header}>
            <NestedRouteLink to="new">
              <Button variant="contained" color="primary">
                Add Network
              </Button>
            </NestedRouteLink>
          </div>
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
        </>
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
          onConfirm={() => {
            const payload = {
              networkID: networkToDelete,
            };
            axios
              .post<{success: boolean}>('/nms/network/delete', payload)
              .then(response => {
                if (!response.data.success) {
                  enqueueSnackbar('Network delete failed', {
                    variant: 'error',
                  });
                } else {
                  setNetworks(networks.filter(n => n != networkToDelete));
                  setNetworkToDelete(null);
                  if (selectedNetworkId === networkToDelete) {
                    window.location.replace('/nms');
                  }
                }
              })
              .catch(() => {
                enqueueSnackbar('Network delete failed', {
                  variant: 'error',
                });
              });
          }}
        />
      )}
      <Routes>
        <Route
          path="/new"
          element={
            <AddNetworkDialog
              onClose={closeDialog}
              onSave={networkID => {
                setNetworks([...networks, networkID]);
                enqueueSnackbar('Network created successfully', {
                  variant: 'success',
                });
                closeDialog();
                if (!selectedNetworkId) {
                  window.location.replace(`/nms/${networkID}/admin/networks`);
                }
              }}
            />
          }
        />
        <Route
          path="/edit/:networkID"
          element={
            <EditNetworkDialog
              onClose={closeDialog}
              onSave={() => {
                enqueueSnackbar('Network updated successfully', {
                  variant: 'success',
                });
                closeDialog();
              }}
            />
          }
        />
      </Routes>
    </div>
  );
}

export default Networks;

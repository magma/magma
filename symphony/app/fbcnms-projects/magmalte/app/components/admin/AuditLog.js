/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import JSONTree from 'react-json-tree';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MoreVert from '@material-ui/icons/MoreVert';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../../theme/design-system/Text';
import {colors} from '../../theme/default';

import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '@fbcnms/ui/hooks';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  actionsCell: {
    textAlign: 'right',
  },
  actionsColumn: {
    width: '160px',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
  iconButton: {
    color: colors.primary.brightGray,
    padding: '5px',
  },
}));

const THEME = {
  scheme: 'monokai',
  base00: '#272822',
  base01: '#383830',
  base02: '#49483e',
  base03: '#75715e',
  base04: '#a59f85',
  base05: '#f8f8f2',
  base06: '#f5f4f1',
  base07: '#f9f8f5',
  base08: '#f92672',
  base09: '#fd971f',
  base0A: '#f4bf75',
  base0B: '#a6e22e',
  base0C: '#a1efe4',
  base0D: '#66d9ef',
  base0E: '#ae81ff',
  base0F: '#cc6633',
};

export default function() {
  const classes = useStyles();
  const {history, relativeUrl, relativePath} = useRouter();
  const {response, error, isLoading} = useAxios({
    url: '/admin/auditlog/async',
    method: 'get',
  });

  if (error || isLoading || !response || !response.data) {
    return <LoadingFiller />;
  }

  const rows = response.data.map(log => (
    <TableRow key={log.id}>
      <TableCell>{log.actingUserId}</TableCell>
      <TableCell>
        <DeviceStatusCircle
          isGrey={false}
          isActive={log.status === 'SUCCESS'}
        />
        {log.mutationType}
      </TableCell>
      <TableCell>{log.objectType}</TableCell>
      <TableCell>{log.objectId}</TableCell>
      <TableCell className={classes.actionsCell}>
        <NestedRouteLink to={`/details/${log.id}`}>
          <IconButton>
            <MoreVert />
          </IconButton>
        </NestedRouteLink>
      </TableCell>
    </TableRow>
  ));

  return (
    <>
      <div className={classes.paper}>
        <div className={classes.header}>
          <Text variant="h5">Audit Log</Text>
        </div>
        <Paper elevation={2}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>User</TableCell>
                <TableCell>Action</TableCell>
                <TableCell>Object Type</TableCell>
                <TableCell>Object ID</TableCell>
                <TableCell className={classes.actionsColumn} />
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </Paper>
      </div>
      <Route
        path={relativePath('/details/:id')}
        render={routeProps => (
          <JSONDialog
            json={
              response.data.find(log => log.id == routeProps.match.params.id)
                .mutationData
            }
            onClose={() => history.push(relativeUrl(''))}
          />
        )}
      />
    </>
  );
}

const JSONDialog = props => {
  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Details</DialogTitle>
      <DialogContent>
        <JSONTree data={props.json} theme={THEME} />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
};

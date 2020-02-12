/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {OrganizationRawType} from '@fbcnms/sequelize-models/models/organization';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditUserDialog from '@fbcnms/ui/components/auth/EditUserDialog';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import OrganizationDialog from './OrganizationDialog';
import Paper from '@material-ui/core/Paper';
import PersonAdd from '@material-ui/icons/PersonAdd';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';

import renderList from '@fbcnms/util/renderList';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Link, Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

export type Organization = OrganizationRawType & {id: number};

const useStyles = makeStyles(_ => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
}));

type Props = {...WithAlert};

function Organizations(props: Props) {
  const classes = useStyles();
  const {relativePath, relativeUrl, history} = useRouter();
  const [organizations, setOrganizations] = useState<?(Organization[])>(null);
  const [addingUserFor, setAddingUserFor] = useState<?Organization>(null);
  const enqueueSnackbar = useEnqueueSnackbar();
  const {error, isLoading} = useAxios({
    url: '/master/organization/async',
    onResponse: useCallback(
      res => setOrganizations(res.data.organizations),
      [],
    ),
  });

  if (error || isLoading || !organizations) {
    return <LoadingFiller />;
  }

  const onDelete = org => {
    props
      .confirm('Are you sure you want to delete this org?')
      .then(async confirm => {
        if (!confirm) return;
        await axios.delete(`/master/organization/async/${org.id}`);
        const newOrganizations = [...organizations];
        newOrganizations.splice(organizations.indexOf(org), 1);
        setOrganizations(newOrganizations);
      });
  };

  const rows = organizations
    .sort((org1, org2) => (org1.name < org2.name ? -1 : 1))
    .map(row => (
      <TableRow key={row.id}>
        <TableCell>
          <Link to={relativePath(`/detail/${row.name}`)}>{row.name}</Link>
        </TableCell>
        <TableCell>{renderList(row.networkIDs)}</TableCell>
        <TableCell>{row.tabs && renderList(row.tabs)}</TableCell>
        <TableCell>
          <IconButton onClick={() => onDelete(row)}>
            <DeleteIcon />
          </IconButton>
          <IconButton onClick={() => setAddingUserFor(row)}>
            <PersonAdd />
          </IconButton>
        </TableCell>
      </TableRow>
    ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <div />
        <NestedRouteLink to="/new">
          <Button>Add Organization</Button>
        </NestedRouteLink>
      </div>
      <Paper className={classes.tableRoot} elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Network IDs</TableCell>
              <TableCell>Tabs</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      <Route
        path={relativePath('/new')}
        render={() => (
          <OrganizationDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={org => {
              const newOrganizations = [...organizations];
              newOrganizations.push(org);
              setOrganizations(newOrganizations);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
      {addingUserFor && (
        <EditUserDialog
          open={true}
          ssoEnabled={!!addingUserFor.ssoEntrypoint}
          onClose={() => setAddingUserFor(null)}
          onEditUser={() => {}}
          editingUser={null}
          allNetworkIDs={addingUserFor.networkIDs}
          onCreateUser={user => {
            axios
              .post(
                `/master/organization/async/${addingUserFor.name}/add_user`,
                user,
              )
              .then(() => {
                enqueueSnackbar('User added successfully', {
                  variant: 'success',
                });
                setAddingUserFor(null);
              })
              .catch(error =>
                enqueueSnackbar(error?.response?.data?.error || error, {
                  variant: 'error',
                }),
              );
          }}
        />
      )}
    </div>
  );
}

export default withAlert(Organizations);

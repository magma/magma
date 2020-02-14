/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FeatureFlag} from './FeatureFlagsDialog';

import EditIcon from '@material-ui/icons/Edit';
import FeatureFlagsDialog from './FeatureFlagsDialog';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import renderList from '@fbcnms/util/renderList';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

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

export default function Features() {
  const classes = useStyles();
  const {relativePath, relativeUrl, history} = useRouter();
  const [featureFlags, setFeatureFlags] = useState<?(FeatureFlag[])>(null);
  useEffect(() => {
    axios.get('/master/feature/async').then(({data}) => setFeatureFlags(data));
  }, []);

  if (!featureFlags) {
    return <LoadingFiller />;
  }

  const rows = featureFlags.map(row => (
    <TableRow key={row.id}>
      <TableCell>{row.title}</TableCell>
      <TableCell>{row.enabledByDefault ? 'Yes' : 'No'}</TableCell>
      <TableCell>
        {renderList(
          Object.keys(row.config).filter(org => row.config[org].enabled),
        )}
      </TableCell>
      <TableCell>
        {renderList(
          Object.keys(row.config).filter(org => !row.config[org].enabled),
        )}
      </TableCell>
      <TableCell>
        <IconButton
          onClick={() => history.push(relativeUrl(`/edit/${row.id}`))}>
          <EditIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.paper}>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Feature</TableCell>
              <TableCell>Enabled By Default</TableCell>
              <TableCell>Enabled For</TableCell>
              <TableCell>Disabled For</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      <Route
        path={relativePath('/edit/:id')}
        render={({match}) => (
          <FeatureFlagsDialog
            featureFlag={nullthrows(
              featureFlags.find(f => f.id === match.params.id),
            )}
            onClose={() => history.push(relativeUrl(''))}
            onSave={flag => {
              const newFeatureFlags = [...featureFlags];
              for (let i = 0; i < newFeatureFlags.length; i++) {
                if (newFeatureFlags[i].id === flag.id) {
                  newFeatureFlags[i] = flag;
                  break;
                }
              }
              setFeatureFlags(newFeatureFlags);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </div>
  );
}

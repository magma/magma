/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
    backgroundColor: symphony.palette.white,
  },
}));

export default function UsersView() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Table
        dataRowsSeparator="border"
        showSelection={true}
        data={[]}
        columns={[
          {
            key: 'name',
            title: <fbt desc="Name column header in users table">Name</fbt>,
            render: _row => 'No Value',
          },
          {
            key: 'role',
            title: <fbt desc="Role column header in users table">Role</fbt>,
            render: _row => 'No Value',
          },
          {
            key: 'job_title',
            title: (
              <fbt desc="Job Title column header in users table">Job Title</fbt>
            ),
            render: _row => 'No Value',
          },
          {
            key: 'groups',
            title: <fbt desc="Groups column header in users table">Groups</fbt>,
            render: _row => 'No Value',
          },
          {
            key: 'status',
            title: <fbt desc="Status column header in users table">Status</fbt>,
            render: _row => 'No Value',
          },
        ]}
      />
    </div>
  );
}

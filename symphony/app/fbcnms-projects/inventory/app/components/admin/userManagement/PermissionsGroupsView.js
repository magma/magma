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
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
  },
}));

export default function UsersView() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Table
        data={[]}
        columns={[
          {
            key: 'name',
            title: (
              <fbt desc="Group Name column header in permission groups table">
                Group Name
              </fbt>
            ),
            render: _row => 'No Value',
          },
          {
            key: 'members',
            title: (
              <fbt desc="Members column header in permission groups table">
                Members
              </fbt>
            ),
            render: _row => 'No Value',
          },
          {
            key: 'status',
            title: (
              <fbt desc="Status column header in permission groups table">
                Status
              </fbt>
            ),
            render: _row => 'No Value',
          },
        ]}
      />
    </div>
  );
}

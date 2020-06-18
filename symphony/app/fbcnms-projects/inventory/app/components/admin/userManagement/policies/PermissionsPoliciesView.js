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
import PermissionsPoliciesTable from './PermissionsPoliciesTable';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {usePermissionsPolicies} from '../data/PermissionsPolicies';
import {useRouter} from '@fbcnms/ui/hooks';

export const PERMISSION_POLICIES_VIEW_NAME = fbt(
  'Polices',
  'Header for view showing system permissions policies settings',
);

const useStyles = makeStyles(() => ({
  root: {
    maxHeight: '100%',
  },
}));

export default function PermissionsPoliciesView() {
  const classes = useStyles();
  const {history} = useRouter();

  const policies = usePermissionsPolicies();

  return (
    <div className={classes.root}>
      <PermissionsPoliciesTable
        policies={policies}
        onPolicySelected={policyId => {
          history.push(`policy/${policyId}`);
        }}
      />
    </div>
  );
}

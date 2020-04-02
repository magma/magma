/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {UserPermissionsGroup} from './UserManagementUtils';

import * as React from 'react';
import ViewHeader from '@fbcnms/ui/components/design-system/View/ViewHeader';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
  },
}));

type Props = {
  group: UserPermissionsGroup,
  className?: ?string,
};

export default function PermissionsGroupPoliciesPane({
  group,
  className,
}: Props) {
  const classes = useStyles();

  return (
    <div className={classNames(classes.root, className)}>
      <ViewHeader title={<fbt desc="">Policies</fbt>} />
      {group.name} Policies
    </div>
  );
}

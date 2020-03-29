/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../../common/InventoryErrorBoundary';
import PermissionsGroupDetailsPane from './PermissionsGroupDetailsPane';
import PermissionsGroupMembersPane from './PermissionsGroupMembersPane';
import PermissionsGroupPoliciesPane from './PermissionsGroupPoliciesPane';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {PERMISSION_GROUPS_VIEW_NAME} from './PermissionsGroupsView';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useUserManagement} from './UserManagementContext';

const useStyles = makeStyles(() => ({
  detailsPane: {
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    '&:not(:first-child)': {
      marginTop: '16px',
    },
  },
}));

type Props = {
  redirectToGroupsView?: () => void,
};

export default function PermissionsGroupCard({redirectToGroupsView}: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const {groups} = useUserManagement();

  const groupId = match.params.id;
  const group = groups.find(group => group.id === groupId);
  if (group == null) {
    if (redirectToGroupsView != null) {
      redirectToGroupsView();
    }
    return null;
  }
  const breadcrumbs = [
    {
      id: 'groups',
      name: `${PERMISSION_GROUPS_VIEW_NAME}`,
      onClick: redirectToGroupsView,
    },
    {
      id: 'groupName',
      name: group.name,
    },
  ];
  const header = {
    title: <Breadcrumbs breadcrumbs={breadcrumbs} />,
    subtitle: fbt('Manage group details, members and policies', ''),
  };
  return (
    <InventoryErrorBoundary>
      <ViewContainer header={header} useBodyScrollingEffect={false}>
        <Grid container spacing={2}>
          <Grid item xs={8} sm={8} lg={8} xl={8}>
            <PermissionsGroupDetailsPane
              group={group}
              className={classes.detailsPane}
            />
            <PermissionsGroupPoliciesPane
              group={group}
              className={classes.detailsPane}
            />
          </Grid>
          <Grid item xs={4} sm={4} lg={4} xl={4}>
            <PermissionsGroupMembersPane
              group={group}
              className={classes.detailsPane}
            />
          </Grid>
        </Grid>
      </ViewContainer>
    </InventoryErrorBoundary>
  );
}

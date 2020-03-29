/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {UserPermissionsGroup} from './TempTypes';

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../../common/InventoryErrorBoundary';
import PermissionsGroupDetailsPane from './PermissionsGroupDetailsPane';
import PermissionsGroupMembersPane from './PermissionsGroupMembersPane';
import PermissionsGroupPoliciesPane from './PermissionsGroupPoliciesPane';
import Strings from '../../../common/CommonStrings';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {PERMISSION_GROUPS_VIEW_NAME} from './PermissionsGroupsView';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useMemo, useState} from 'react';
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
  redirectToGroupsView: () => void,
  onClose: () => void,
};

export default function PermissionsGroupCard({
  redirectToGroupsView,
  onClose,
}: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const {groups, editGroup} = useUserManagement();
  const [group, setGroup] = useState<?UserPermissionsGroup>(null);

  const groupId = match.params.id;
  useEffect(() => {
    const requestedGroup = groups.find(group => group.id === groupId);
    if (requestedGroup == null) {
      redirectToGroupsView();
    }
    setGroup(requestedGroup);
  }, [groupId, groups, redirectToGroupsView]);

  const header = useMemo(() => {
    const breadcrumbs = [
      {
        id: 'groups',
        name: `${PERMISSION_GROUPS_VIEW_NAME}`,
        onClick: redirectToGroupsView,
      },
      {
        id: 'groupName',
        name: group?.name || '',
      },
    ];
    const actions = [
      {
        title: Strings.common.cancelButton,
        action: onClose,
        skin: 'regular',
      },
      {
        title: Strings.common.saveButton,
        action: () => {
          if (group == null) {
            return;
          }
          editGroup(group).then(onClose);
        },
      },
    ];
    return {
      title: <Breadcrumbs breadcrumbs={breadcrumbs} />,
      subtitle: fbt('Manage group details, members and policies', ''),
      actionButtons: actions,
    };
  }, [editGroup, group, onClose, redirectToGroupsView]);

  if (group == null) {
    return null;
  }
  return (
    <InventoryErrorBoundary>
      <ViewContainer header={header} useBodyScrollingEffect={false}>
        <Grid container spacing={2}>
          <Grid item xs={8} sm={8} lg={8} xl={8}>
            <PermissionsGroupDetailsPane
              group={group}
              onChange={setGroup}
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

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {UserPermissionsGroup} from '../utils/UserManagementUtils';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../../../common/InventoryErrorBoundary';
import PermissionsGroupDetailsPane from './PermissionsGroupDetailsPane';
import PermissionsGroupMembersPane from './PermissionsGroupMembersPane';
import PermissionsGroupPoliciesPane from './PermissionsGroupPoliciesPane';
import Strings from '../../../../common/CommonStrings';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {
  GROUP_STATUSES,
  NEW_GROUP_DIALOG_PARAM,
} from '../utils/UserManagementUtils';
import {PERMISSION_GROUPS_VIEW_NAME} from './PermissionsGroupsView';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useMemo, useState} from 'react';
import {useRouteMatch} from 'react-router-dom';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  detailsPane: {
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    '&:not(:first-child)': {
      marginTop: '16px',
    },
  },
  container: {
    maxHeight: '100%',
  },
}));

type Props = {
  redirectToGroupsView: () => void,
  onClose: () => void,
};

const initialNewGroup: UserPermissionsGroup = {
  id: '',
  name: '',
  description: '',
  status: GROUP_STATUSES.ACTIVE.key,
  members: [],
  memberUsers: [],
};

export default function PermissionsGroupCard({
  redirectToGroupsView,
  onClose,
}: Props) {
  const classes = useStyles();
  const match = useRouteMatch();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');
  const {groups, editGroup, addGroup} = useUserManagement();
  const groupId = match.params.id;
  const isOnNewGroup = groupId === NEW_GROUP_DIALOG_PARAM;
  const [group, setGroup] = useState<?UserPermissionsGroup>(
    isOnNewGroup ? {...initialNewGroup} : null,
  );

  useEffect(() => {
    if (isOnNewGroup) {
      return;
    }
    const requestedGroup = groups.find(group => group.id === groupId);
    if (requestedGroup == null) {
      redirectToGroupsView();
    }
    setGroup(requestedGroup);
  }, [groupId, groups, isOnNewGroup, redirectToGroupsView]);

  const header = useMemo(() => {
    const breadcrumbs = [
      {
        id: 'groups',
        name: `${PERMISSION_GROUPS_VIEW_NAME}`,
        onClick: redirectToGroupsView,
      },
      {
        id: 'groupName',
        name: isOnNewGroup ? `${fbt('New Group', '')}` : group?.name || '',
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
          const saveAction = isOnNewGroup ? addGroup : editGroup;
          saveAction(group).then(onClose);
        },
      },
    ];
    return {
      title: <Breadcrumbs breadcrumbs={breadcrumbs} />,
      subtitle: fbt('Manage group details, members and policies', ''),
      actionButtons: actions,
    };
  }, [addGroup, editGroup, group, isOnNewGroup, onClose, redirectToGroupsView]);

  if (group == null) {
    return null;
  }
  return (
    <InventoryErrorBoundary>
      <ViewContainer header={header} useBodyScrollingEffect={false}>
        <Grid container spacing={2} className={classes.container}>
          <Grid item xs={8} sm={8} lg={8} xl={8} className={classes.container}>
            <PermissionsGroupDetailsPane
              group={group}
              onChange={setGroup}
              className={classes.detailsPane}
            />
            {userManagementDevMode ? (
              <PermissionsGroupPoliciesPane
                group={group}
                className={classes.detailsPane}
              />
            ) : null}
          </Grid>
          <Grid item xs={4} sm={4} lg={4} xl={4} className={classes.container}>
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

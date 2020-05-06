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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../../../common/InventoryErrorBoundary';
import PermissionsGroupDetailsPane from './PermissionsGroupDetailsPane';
import PermissionsGroupMembersPane from './PermissionsGroupMembersPane';
import PermissionsGroupPoliciesPane from './PermissionsGroupPoliciesPane';
import Strings from '@fbcnms/strings/Strings';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {
  ButtonAction,
  IconAction,
} from '@fbcnms/ui/components/design-system/View/ViewHeaderActions';
import {GROUP_STATUSES, NEW_DIALOG_PARAM} from '../utils/UserManagementUtils';
import {PERMISSION_GROUPS_VIEW_NAME} from './PermissionsGroupsView';
import {generateTempId} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouteMatch} from 'react-router-dom';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  detailsPane: {
    display: 'flex',
    flexDirection: 'column',
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

type Props = $ReadOnly<{|
  redirectToGroupsView: () => void,
  onClose: () => void,
  ...WithAlert,
|}>;

const initialNewGroup: UserPermissionsGroup = {
  id: generateTempId(),
  name: '',
  description: '',
  status: GROUP_STATUSES.ACTIVE.key,
  members: [],
  memberUsers: [],
};

function PermissionsGroupCard(props: Props) {
  const {redirectToGroupsView, onClose} = props;
  const classes = useStyles();
  const match = useRouteMatch();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');
  const {groups, editGroup, addGroup, deleteGroup} = useUserManagement();
  const groupId = match.params.id;
  const isOnNewGroup = groupId === NEW_DIALOG_PARAM;
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

  const enqueueSnackbar = useEnqueueSnackbar();
  const handleError = useCallback(
    (error: string) => {
      enqueueSnackbar(error, {variant: 'error'});
    },
    [enqueueSnackbar],
  );

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
      <ButtonAction skin="regular" action={onClose}>
        {Strings.common.cancelButton}
      </ButtonAction>,
      <ButtonAction
        disableOnFromError={true}
        action={() => {
          if (group == null) {
            return;
          }
          const saveAction = isOnNewGroup ? addGroup : editGroup;
          saveAction(group)
            .then(onClose)
            .catch(handleError);
        }}>
        {Strings.common.saveButton}
      </ButtonAction>,
    ];
    if (!isOnNewGroup && userManagementDevMode) {
      actions.unshift(
        <IconAction
          skin="gray"
          icon={DeleteIcon}
          action={() => {
            if (group == null) {
              return;
            }
            props
              .confirm(
                <fbt desc="">Are you sure you want to delete this group?</fbt>,
              )
              .then(confirm => {
                if (!confirm) {
                  return;
                }
                return deleteGroup(group.id).then(onClose);
              })
              .catch(handleError);
          }}
        />,
      );
    }
    return {
      title: <Breadcrumbs breadcrumbs={breadcrumbs} />,
      subtitle: fbt('Manage group details, members and policies', ''),
      actionButtons: actions,
    };
  }, [
    addGroup,
    deleteGroup,
    editGroup,
    group,
    handleError,
    isOnNewGroup,
    onClose,
    props,
    redirectToGroupsView,
    userManagementDevMode,
  ]);

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

export default withAlert(PermissionsGroupCard);

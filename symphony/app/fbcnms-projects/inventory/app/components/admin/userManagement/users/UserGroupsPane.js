/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from '../utils/UserManagementUtils';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import GroupSearchBar from '../utils/search/GroupSearchBar';
import GroupsSearchEmptyState from '../utils/search/GroupsSearchEmptyState';
import NoGroupsEmptyState from '../utils/NoGroupsEmptyState';
import UserGroupsList from '../utils/UserGroupsList';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {
  GroupSearchContextProvider,
  useGroupSearchContext,
} from '../utils/search/GroupSearchContext';
import {PERMISSION_GROUPS_VIEW_NAME} from '../groups/PermissionsGroupsView';
import {TOGGLE_BUTTON_DISPLAY} from '../utils/ListItem';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    padding: '0 16px',
  },
  groupsSearchBox: {
    paddingTop: '8px',
    maxWidth: '320px',
  },
}));

type Props = $ReadOnly<{|
  user: User,
|}>;

function GroupsList(props: Props) {
  const {user} = props;
  const groupSearch = useGroupSearchContext();

  return (
    <UserGroupsList
      user={user}
      groups={
        groupSearch.isEmptySearchTerm
          ? user.groups.filter(Boolean)
          : groupSearch.results
      }
      assigmentButton={
        groupSearch.isEmptySearchTerm ? undefined : TOGGLE_BUTTON_DISPLAY.always
      }
      emptyState={
        <GroupsSearchEmptyState noSearchEmptyState={<NoGroupsEmptyState />} />
      }
    />
  );
}

export default function UserGroupsPane(props: Props) {
  const {user} = props;
  const classes = useStyles();

  return (
    <GroupSearchContextProvider>
      <ViewContainer
        className={classes.root}
        header={{
          title: PERMISSION_GROUPS_VIEW_NAME,
          subtitle: (
            <fbt desc="">
              Groups determine which policies apply on users. To edit and create
              more groups, go to the
              <fbt:param name="link to groups">
                <Button
                  variant="text"
                  onClick={() => window.open('/admin/user_management/groups')}>
                  {PERMISSION_GROUPS_VIEW_NAME}
                </Button>
              </fbt:param>
              tab in the left navigation.
            </fbt>
          ),
          searchBar: (
            <GroupSearchBar
              staticShownGroups={user.groups.filter(Boolean)}
              className={classes.groupsSearchBox}
            />
          ),
        }}>
        <GroupsList {...props} />
      </ViewContainer>
    </GroupSearchContextProvider>
  );
}

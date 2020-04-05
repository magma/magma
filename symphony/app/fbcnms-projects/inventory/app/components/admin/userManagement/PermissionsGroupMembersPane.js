/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {GroupMember} from './GroupMemberViewer';
import type {PermissionsGroupMembersPaneUserSearchQuery} from './__generated__/PermissionsGroupMembersPaneUserSearchQuery.graphql';
import type {UserPermissionsGroup} from './UserManagementUtils';

import * as React from 'react';
import GroupMemberViewer, {ASSIGNMENT_BUTTON_VIEWS} from './GroupMemberViewer';
import InputAffix from '@fbcnms/ui/components/design-system/Input/InputAffix';
import RelayEnvironment from '../../../common/RelayEnvironment';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {
  CloseIcon,
  ProfileIcon,
} from '@fbcnms/ui/components/design-system/Icons/';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';
import {userResponse2User} from './UserManagementUtils';

const userSearchQuery = graphql`
  query PermissionsGroupMembersPaneUserSearchQuery(
    $filters: [UserFilterInput!]!
  ) {
    userSearch(filters: $filters) {
      users {
        id
        authID
        firstName
        lastName
        email
        status
        role
        groups {
          id
          name
        }
        profilePhoto {
          id
          fileName
          storeKey
        }
      }
    }
  }
`;

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    height: '100%',
  },
  header: {
    paddingBottom: '5px',
  },
  title: {
    marginBottom: '16px',
    display: 'flex',
    alignItems: 'center',
  },
  titleIcon: {
    marginRight: '4px',
  },
  userSearch: {
    position: 'relative',
    overflow: 'hidden',
    borderRadius: '4px',
    marginTop: '8px',
  },
  searchProgress: {
    position: 'absolute',
    borderBottom: `3px solid transparent`,
    bottom: 0,
    left: '0%',
  },
  runSearch: {
    borderBottomColor: symphony.palette.primary,
    animation: '$progress 2s infinite',
  },
  clearSearchIcon: {},
  usersListHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    marginTop: '12px',
    marginBottom: '-3px',
  },
  usersList: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
  },
  user: {
    borderBottom: `1px solid ${symphony.palette.separatorLight}`,
  },
  '@keyframes progress': {
    '0%': {
      right: '100%',
      left: '0%',
    },
    '50%': {
      left: '0%',
    },
    '100%': {
      right: '0%',
      left: '100%',
    },
  },
}));

type Props = $ReadOnly<{|
  group: UserPermissionsGroup,
  className?: ?string,
|}>;

const NO_SEARCH_VALUE = '';

export default function PermissionsGroupMembersPane(props: Props) {
  const {group, className} = props;
  const classes = useStyles();
  const [searchIsInProgress, setSearchIsInProgress] = useState(false);
  const [userSearchValue, setUserSearchValue] = useState(NO_SEARCH_VALUE);
  const [usersSearchResult, setUsersSearchResult] = useState<
    Array<GroupMember>,
  >([]);

  const queryUsers = useCallback(
    debounce((searchTerm: string) => {
      setSearchIsInProgress(true);
      fetchQuery<PermissionsGroupMembersPaneUserSearchQuery>(
        RelayEnvironment,
        userSearchQuery,
        {
          filters: [
            {
              filterType: 'USER_NAME',
              operator: 'CONTAINS',
              stringValue: searchTerm,
            },
          ],
        },
      )
        .then(response => {
          if (!response?.userSearch) {
            return;
          }
          setUsersSearchResult(
            response.userSearch.users.filter(Boolean).map(userNode => {
              const userData = userResponse2User(userNode);
              return {
                user: userData,
                isMember:
                  userData.groups.find(
                    userGroup => userGroup?.id == group.id,
                  ) != null,
              };
            }),
          );
        })
        .finally(() => setSearchIsInProgress(false));
    }, 200),
    [group],
  );

  const updateSearch = useCallback(
    newSearchValue => {
      setUserSearchValue(newSearchValue);
      if (newSearchValue.trim() == NO_SEARCH_VALUE) {
        setUsersSearchResult([]);
        setSearchIsInProgress(false);
      } else {
        queryUsers(newSearchValue);
      }
    },
    [queryUsers],
  );

  const title = useMemo(
    () => (
      <div className={classes.title}>
        <ProfileIcon className={classes.titleIcon} />
        <fbt desc="">Members</fbt>
      </div>
    ),
    [classes.title, classes.titleIcon],
  );

  const isOnSearchMode = userSearchValue.trim() != NO_SEARCH_VALUE;

  const subtitle = useMemo(
    () => (
      <>
        <Text variant="caption" color="gray" useEllipsis={true}>
          <fbt desc="">Add users to apply policies.</fbt>
        </Text>
        <Text variant="caption" color="gray" useEllipsis={true}>
          <fbt desc="">Users can be members in multiple groups.</fbt>
        </Text>
      </>
    ),
    [],
  );

  const searchBar = useMemo(
    () => (
      <>
        <div className={classes.userSearch}>
          <div
            className={classNames(classes.searchProgress, {
              [classes.runSearch]: searchIsInProgress,
            })}
          />
          <TextInput
            type="string"
            variant="outlined"
            placeholder={`${fbt('Search users...', '')}`}
            fullWidth={true}
            value={userSearchValue}
            onChange={e => updateSearch(e.target.value)}
            suffix={
              isOnSearchMode ? (
                <InputAffix onClick={() => updateSearch(NO_SEARCH_VALUE)}>
                  <CloseIcon className={classes.clearSearchIcon} color="gray" />
                </InputAffix>
              ) : null
            }
          />
        </div>
        {isOnSearchMode ? null : (
          <div className={classes.usersListHeader}>
            <Text variant="subtitle2" useEllipsis={true}>
              <fbt desc="">
                <fbt:plural count={group.members.length} showCount="yes">
                  Member
                </fbt:plural>
              </fbt>
            </Text>
          </div>
        )}
      </>
    ),
    [
      classes.clearSearchIcon,
      classes.runSearch,
      classes.searchProgress,
      classes.userSearch,
      classes.usersListHeader,
      group.members.length,
      isOnSearchMode,
      searchIsInProgress,
      updateSearch,
      userSearchValue,
    ],
  );

  const header = useMemo(
    () => ({
      title,
      subtitle,
      searchBar,
      className: classes.header,
    }),
    [classes.header, searchBar, subtitle, title],
  );

  const memberUsers: $ReadOnlyArray<GroupMember> = useMemo(
    () =>
      isOnSearchMode
        ? usersSearchResult
        : group.memberUsers.map(user => ({
            user: user,
            isMember: true,
          })),
    [isOnSearchMode, usersSearchResult, group.memberUsers],
  );

  return (
    <div className={classNames(classes.root, className)}>
      <ViewContainer header={header}>
        <div className={classes.usersList}>
          {memberUsers.map(member => (
            <GroupMemberViewer
              className={classes.user}
              member={member}
              assigmentButton={
                isOnSearchMode
                  ? ASSIGNMENT_BUTTON_VIEWS.always
                  : ASSIGNMENT_BUTTON_VIEWS.onHover
              }
              group={group}
            />
          ))}
        </div>
      </ViewContainer>
    </div>
  );
}

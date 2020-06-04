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
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import PermissionsGroupMembersList from './PermissionsGroupMembersList';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserSearchBox from '../utils/search/UserSearchBox';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {ProfileIcon} from '@fbcnms/ui/components/design-system/Icons/';
import {
  UserSearchContextProvider,
  useUserSearchContext,
} from '../utils/search/UserSearchContext';
import {makeStyles} from '@material-ui/styles';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
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
    marginTop: '8px',
  },
  usersListHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    marginTop: '12px',
    marginBottom: '-3px',
  },
  noMembers: {
    width: '124px',
    paddingTop: '50%',
    alignSelf: 'center',
  },
  noSearchResults: {
    paddingTop: '50%',
    alignSelf: 'center',
    textAlign: 'center',
  },
  clearSearchWrapper: {
    marginTop: '16px',
  },
  clearSearch: {
    ...symphony.typography.subtitle1,
  },
}));

type Props = $ReadOnly<{|
  group: UserPermissionsGroup,
  onChange: UserPermissionsGroup => void,
  className?: ?string,
|}>;

function SearchBar(
  props: $ReadOnly<{|
    group: UserPermissionsGroup,
  |}>,
) {
  const {group} = props;
  const classes = useStyles();
  const userSearch = useUserSearchContext();

  return (
    <>
      <div className={classes.userSearch}>
        <UserSearchBox />
      </div>
      {!userSearch.isEmptySearchTerm ? null : (
        <div className={classes.usersListHeader}>
          {group.members.length > 0 ? (
            <Text variant="subtitle2" useEllipsis={true}>
              <fbt desc="">
                <fbt:plural count={group.members.length} showCount="yes">
                  Member
                </fbt:plural>
              </fbt>
            </Text>
          ) : null}
        </div>
      )}
    </>
  );
}

export default function PermissionsGroupMembersPane(props: Props) {
  const {group, onChange, className} = props;
  const classes = useStyles();

  const title = useMemo(
    () => (
      <div className={classes.title}>
        <ProfileIcon className={classes.titleIcon} />
        <fbt desc="">Members</fbt>
      </div>
    ),
    [classes.title, classes.titleIcon],
  );

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

  const searchBar = useMemo(() => <SearchBar group={group} />, [group]);

  const header = useMemo(
    () => ({
      title,
      subtitle,
      searchBar,
      className: classes.header,
    }),
    [classes.header, searchBar, subtitle, title],
  );

  return (
    <Card className={classNames(classes.root, className)} margins="none">
      <UserSearchContextProvider>
        <ViewContainer header={header}>
          <PermissionsGroupMembersList group={group} onChange={onChange} />
        </ViewContainer>
      </UserSearchContextProvider>
    </Card>
  );
}

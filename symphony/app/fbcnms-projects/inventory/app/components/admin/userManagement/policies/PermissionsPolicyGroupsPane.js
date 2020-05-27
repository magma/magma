/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../utils/UserManagementUtils';

import * as React from 'react';
import GroupSearchBox from '../utils/search/GroupSearchBox';
import PermissionsPolicyGroupsList from './PermissionsPolicyGroupsList';
import Text from '@fbcnms/ui/components/design-system/Text';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {GroupIcon} from '@fbcnms/ui/components/design-system/Icons/';
import {
  GroupSearchContextProvider,
  useGroupSearchContext,
} from '../utils/search/GroupSearchContext';
import {makeStyles} from '@material-ui/styles';
import {useMemo} from 'react';

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
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  className?: ?string,
|}>;

function SearchBar(
  props: $ReadOnly<{|
    policy: PermissionsPolicy,
  |}>,
) {
  const {policy} = props;
  const classes = useStyles();
  const userSearch = useGroupSearchContext();

  return (
    <>
      <div className={classes.userSearch}>
        <GroupSearchBox />
      </div>
      {!userSearch.isEmptySearchTerm ? null : (
        <div className={classes.usersListHeader}>
          {policy.groups.length > 0 ? (
            <Text variant="subtitle2" useEllipsis={true}>
              <fbt desc="">
                <fbt:plural count={policy.groups.length} showCount="yes">
                  Group
                </fbt:plural>
              </fbt>
            </Text>
          ) : null}
        </div>
      )}
    </>
  );
}

export default function PermissionsPolicyGroupsPane(props: Props) {
  const {policy, onChange, className} = props;
  const classes = useStyles();

  const title = useMemo(
    () => (
      <div className={classes.title}>
        <GroupIcon className={classes.titleIcon} />
        <fbt desc="">Groups</fbt>
      </div>
    ),
    [classes.title, classes.titleIcon],
  );

  const subtitle = useMemo(
    () => (
      <Text variant="caption" color="gray" useEllipsis={true}>
        <fbt desc="">
          Add this policy to groups to apply it on their members.
        </fbt>
      </Text>
    ),
    [],
  );

  const searchBar = useMemo(() => <SearchBar policy={policy} />, [policy]);

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
    <div className={classNames(classes.root, className)}>
      <GroupSearchContextProvider>
        <ViewContainer header={header}>
          <PermissionsPolicyGroupsList policy={policy} onChange={onChange} />
        </ViewContainer>
      </GroupSearchContextProvider>
    </div>
  );
}

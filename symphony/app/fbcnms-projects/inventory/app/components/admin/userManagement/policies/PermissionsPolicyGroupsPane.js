/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../data/PermissionsPolicies';

import * as React from 'react';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import GroupSearchBox from '../utils/search/GroupSearchBox';
import InfoTinyIcon from '@fbcnms/ui/components/design-system/Icons/Indications/InfoTinyIcon';
import PermissionsPolicyGroupsList from './PermissionsPolicyGroupsList';
import Switch from '@fbcnms/ui/components/design-system/switch/Switch';
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
    '&:not($isGlobal)': {
      height: '100%',
    },
  },
  isGlobal: {},
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
  body: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
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
  globalGroupContainer: {
    display: 'flex',
    margin: '16px 0',
  },
  captionContainer: {
    display: 'flex',
  },
  switchWrapper: {
    flexGrow: 1,
  },
  switch: {
    float: 'right',
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

  const searchBar = useMemo(
    () => (policy.isGlobal ? null : <SearchBar policy={policy} />),
    [policy],
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

  return (
    <Card
      className={classNames(classes.root, className, {
        [classes.isGlobal]: policy.isGlobal,
      })}
      margins="none">
      <GroupSearchContextProvider>
        <ViewContainer header={header}>
          <div className={classes.body}>
            {policy.isGlobal !== true ? (
              <PermissionsPolicyGroupsList
                policy={policy}
                onChange={onChange}
              />
            ) : null}
            <div className={classes.globalGroupContainer}>
              <div
                className={classes.captionContainer}
                title={`${fbt(
                  'When turned on, all current and future users will have this policy applied.',
                  '',
                )}`}>
                <Text>
                  <fbt desc="">Apply this policy on all users</fbt>
                </Text>
                <InfoTinyIcon />
              </div>
              <div className={classes.switchWrapper}>
                <Switch
                  title=""
                  checked={policy.isGlobal}
                  onChange={isGlobal => {
                    onChange({
                      ...policy,
                      isGlobal,
                    });
                  }}
                  className={classes.switch}
                />
              </div>
            </div>
          </div>
        </ViewContainer>
      </GroupSearchContextProvider>
    </Card>
  );
}

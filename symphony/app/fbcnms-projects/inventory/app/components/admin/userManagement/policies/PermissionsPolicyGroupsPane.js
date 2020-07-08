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
import GroupSearchBar from '../utils/search/GroupSearchBar';
import InfoTinyIcon from '@fbcnms/ui/components/design-system/Icons/Indications/InfoTinyIcon';
import PermissionsPolicyGroupsList from './PermissionsPolicyGroupsList';
import Switch from '@fbcnms/ui/components/design-system/switch/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {GroupIcon} from '@fbcnms/ui/components/design-system/Icons/';
import {GroupSearchContextProvider} from '../utils/search/GroupSearchContext';
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
    () => <fbt desc="">Search groups to apply this policy to them.</fbt>,
    [],
  );

  const searchBar = useMemo(
    () =>
      policy.isGlobal ? null : (
        <GroupSearchBar staticShownGroups={policy.groups} />
      ),
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
                  'When selected, this policy will apply to all users, including those not in this group.',
                  '',
                )}`}>
                <Text>
                  <fbt desc="">Apply this policy to all users</fbt>
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

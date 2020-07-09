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
import GroupsSearchEmptyState from '../utils/search/GroupsSearchEmptyState';
import NoGroupsEmptyState from '../utils/NoGroupsEmptyState';
import PolicyGroupsList from '../utils/PolicyGroupsList';
import {TOGGLE_BUTTON_DISPLAY} from '../utils/ListItem';
import {makeStyles} from '@material-ui/styles';
import {useGroupSearchContext} from '../utils/search/GroupSearchContext';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: '1',
  },
}));

type Props = $ReadOnly<{|
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
|}>;

export default function PermissionsPolicyGroupsList(props: Props) {
  const {policy, onChange} = props;
  const classes = useStyles();
  const groupSearch = useGroupSearchContext();

  return (
    <PolicyGroupsList
      groups={
        groupSearch.isEmptySearchTerm ? policy.groups : groupSearch.results
      }
      policy={policy}
      onChange={onChange}
      assigmentButton={
        groupSearch.isEmptySearchTerm
          ? TOGGLE_BUTTON_DISPLAY.onHover
          : TOGGLE_BUTTON_DISPLAY.always
      }
      emptyState={
        <GroupsSearchEmptyState noSearchEmptyState={<NoGroupsEmptyState />} />
      }
      className={classes.root}
    />
  );
}

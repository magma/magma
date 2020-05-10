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
import Button from '@fbcnms/ui/components/design-system/Button';
import PolicyGroupsList from '../utils/PolicyGroupsList';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {TOGGLE_BUTTON_DISPLAY} from '../utils/ListItem';
import {makeStyles} from '@material-ui/styles';
import {useGroupSearchContext} from '../utils/search/GroupSearchContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {},
  noGroups: {
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
}));

type Props = $ReadOnly<{|
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
|}>;

export default function PermissionsPolicyGroupsList(props: Props) {
  const {policy, onChange} = props;
  const classes = useStyles();
  const groupSearch = useGroupSearchContext();

  const policyGroupsEmptyState = useMemo(
    () => (
      <img
        className={classes.noGroups}
        src="/inventory/static/images/noGroups.png"
      />
    ),
    [classes.noGroups],
  );

  const emptyState = useMemo(() => {
    if (groupSearch.isEmptySearchTerm) {
      return policyGroupsEmptyState;
    }

    if (groupSearch.isSearchInProgress) {
      return null;
    }

    return (
      <div className={classes.noSearchResults}>
        <Text variant="h6" color="gray">
          <fbt desc="">
            No groups found for '<fbt:param name="given search term">
              {groupSearch.searchTerm}
            </fbt:param>'
          </fbt>
        </Text>
        <div className={classes.clearSearchWrapper}>
          <Button variant="text" skin="gray" onClick={groupSearch.clearSearch}>
            <Text variant="subtitle2">
              <fbt desc="">Clear Search</fbt>
            </Text>
          </Button>
        </div>
      </div>
    );
  }, [
    classes.clearSearchWrapper,
    classes.noSearchResults,
    groupSearch.clearSearch,
    groupSearch.isEmptySearchTerm,
    groupSearch.isSearchInProgress,
    groupSearch.searchTerm,
    policyGroupsEmptyState,
  ]);

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
      emptyState={emptyState}
    />
  );
}

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useGroupSearchContext} from './GroupSearchContext';

const useStyles = makeStyles(() => ({
  noSearchResults: {
    position: 'absolute',
    top: '33%',
    right: '50%',
    transform: 'translate(50%, -50%)',
    textAlign: 'center',
  },
  clearSearchWrapper: {
    marginTop: '16px',
  },
}));

type GroupsSearchEmptyStateProps = $ReadOnly<{|
  noSearchEmptyState: React.Node,
|}>;

export default function GroupsSearchEmptyState(
  props: GroupsSearchEmptyStateProps,
) {
  const {noSearchEmptyState} = props;
  const groupSearch = useGroupSearchContext();
  const classes = useStyles();

  if (groupSearch.isEmptySearchTerm) {
    return noSearchEmptyState;
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
}

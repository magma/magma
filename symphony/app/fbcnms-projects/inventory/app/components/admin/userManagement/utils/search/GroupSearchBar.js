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
import GroupSearchBox from './GroupSearchBox';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useGroupSearchContext} from './GroupSearchContext';
import type {UsersGroup} from '../../data/UsersGroups';

const useStyles = makeStyles(() => ({
  groupsSearch: {
    marginTop: '8px',
  },
  groupsListHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    marginTop: '12px',
    marginBottom: '-3px',
  },
}));

type GroupSearchBarProps = $ReadOnly<{|
  staticShownGroups: $ReadOnlyArray<UsersGroup>,
  className?: ?string,
|}>;

export default function GroupSearchBar(props: GroupSearchBarProps) {
  const {staticShownGroups, className} = props;
  const classes = useStyles();
  const groupSearch = useGroupSearchContext();

  return (
    <div className={className}>
      <div className={classes.groupsSearch}>
        <GroupSearchBox />
      </div>
      {!groupSearch.isEmptySearchTerm ? null : (
        <div className={classes.groupsListHeader}>
          {staticShownGroups.length > 0 ? (
            <Text variant="subtitle2" useEllipsis={true}>
              <fbt desc="">
                <fbt:plural count={staticShownGroups.length} showCount="yes">
                  Group
                </fbt:plural>
              </fbt>
            </Text>
          ) : null}
        </div>
      )}
    </div>
  );
}

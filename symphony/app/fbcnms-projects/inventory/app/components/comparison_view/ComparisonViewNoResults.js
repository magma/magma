/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {SearchIcon} from '@fbcnms/ui/components/design-system/Icons';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
    width: '100%',
  },
  searchIcon: {
    marginRight: '8px',
  },
}));

const ComparisonViewNoResults = () => {
  const classes = useStyles();
  return (
    <div className={classes.noResultsRoot}>
      <SearchIcon className={classes.searchIcon} color="gray" />
      <Text variant="h6" color="gray">
        No results found
      </Text>
    </div>
  );
};

export default ComparisonViewNoResults;

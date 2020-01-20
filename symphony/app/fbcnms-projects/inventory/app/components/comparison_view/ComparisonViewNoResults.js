/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import SearchIcon from '@material-ui/icons/Search';
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: theme.palette.grey[600],
  },
  searchIcon: {
    color: theme.palette.grey[600],
    marginBottom: '6px',
    fontSize: '36px',
  },
}));

const ComparisonViewNoResults = () => {
  const classes = useStyles();
  return (
    <div className={classes.noResultsRoot}>
      <SearchIcon className={classes.searchIcon} />
      <Text variant="h6" className={classes.noResultsLabel}>
        No results found
      </Text>
    </div>
  );
};

export default ComparisonViewNoResults;

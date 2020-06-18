/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';

import AppContext from '@fbcnms/ui/context/AppContext';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  version: {
    bottom: '10px',
    cursor: 'pointer',
    paddingLeft: '90px',
    position: 'absolute',
    textDecoration: 'underline',
  },
}));

export default function VersionTooltip() {
  const classes = useStyles();
  return (
    <AppContext.Consumer>
      {({version}) =>
        version && (
          <Tooltip
            title={version}
            placement="top"
            onClick={() =>
              navigator.clipboard.writeText(version.split(' ')[0])
            }>
            <Text className={classes.version} variant="caption">
              Version
            </Text>
          </Tooltip>
        )
      }
    </AppContext.Consumer>
  );
}

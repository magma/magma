/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
    flexGrow: 1,
  },
  dashboardsIframe: {
    width: '100%',
    border: 0,
  },
}));

type Props = {
  grafanaURL: string,
};

export default function GrafanaDashboards(props: Props) {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <iframe src={props.grafanaURL} className={classes.dashboardsIframe} />
    </div>
  );
}

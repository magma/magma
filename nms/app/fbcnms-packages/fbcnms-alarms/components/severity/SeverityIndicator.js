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
import Typography from '@material-ui/core/Typography';
import classnames from 'classnames';
import {SEVERITY} from './Severity';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  // the circle
  indicator: {
    display: 'inline-block',
    height: '10px',
    width: '10px',
    borderRadius: '50%',
  },
  text: {
    marginLeft: theme.spacing(1),
    textTransform: 'capitalize',
  },
  critical: {
    backgroundColor: SEVERITY.CRITICAL.color,
  },
  major: {
    backgroundColor: SEVERITY.MAJOR.color,
  },
  minor: {
    backgroundColor: SEVERITY.MINOR.color,
  },
  warning: {
    backgroundColor: SEVERITY.WARNING.color,
  },
  info: {
    backgroundColor: SEVERITY.INFO.color,
  },
  notice: {
    backgroundColor: SEVERITY.NOTICE.color,
  },
  unknown: {
    backgroundColor: SEVERITY.NOTICE.color,
  },
}));

type Props = {
  severity: string,
};

export default function SeverityIndicator({severity}: Props) {
  const value =
    severity && severity.trim() !== '' ? severity.toLowerCase() : 'unknown';
  const classes = useStyles();

  const colorClassname = React.useMemo(
    () => classnames(classes.indicator, classes[value] ?? classes.unknown),
    [value, classes],
  );

  return (
    <Typography noWrap>
      <span className={colorClassname} />
      <span className={classes.text}>{value}</span>
    </Typography>
  );
}

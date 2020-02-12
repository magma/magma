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
import grey from '@material-ui/core/colors/grey';
import orange from '@material-ui/core/colors/orange';
import red from '@material-ui/core/colors/red';
import yellow from '@material-ui/core/colors/yellow';
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
  red: {
    backgroundColor: red.A400,
  },
  orange: {
    backgroundColor: orange.A400,
  },
  yellow: {
    backgroundColor: yellow.A400,
  },
  grey: {
    backgroundColor: grey[500],
  },
}));

const SEVERITY_CLASS_MAP = {
  critical: {className: 'red'},
  major: {className: 'orange'},
  minor: {className: 'yellow'},
  warning: {className: 'yellow'},
  info: {className: 'grey'},
  notice: {className: 'grey'},
  unknown: {className: 'grey'},
};

type Props = {
  severity: string,
};

export default function SeverityIndicator({severity}: Props) {
  const value =
    severity && severity.trim() !== '' ? severity.toLowerCase() : 'unknown';
  const classes = useStyles();

  const colorClassname = React.useMemo(
    () =>
      classnames(
        classes.indicator,
        classes[(SEVERITY_CLASS_MAP[value]?.className)],
      ),
    [value, classes],
  );

  return (
    <Typography noWrap>
      <span className={colorClassname} />
      <span className={classes.text}>{value}</span>
    </Typography>
  );
}

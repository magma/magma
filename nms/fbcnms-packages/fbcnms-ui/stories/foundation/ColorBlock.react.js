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
import Theme from '../../theme/symphony';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

type Props = {
  color: string,
  name: string,
  code?: string,
  className?: string,
};

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    width: '160px',
  },
  block: {
    height: '160px',
    padding: '9px 12px',
    boxSizing: 'border-box',
  },
  nameContainer: {
    display: 'flex',
    flexDirection: 'row',
    padding: '8px',
    backgroundColor: 'white',
    alignItems: 'center',
  },
  name: {
    color: Theme.palette.D900,
    fontSize: '16px',
    lineHeight: '19px',
    fontWeight: 'bold',
    flexGrow: 1,
  },
  hex: {
    color: Theme.palette.D900,
    fontSize: '16px',
    lineHeight: '21px',
    fontWeight: 500,
  },
}));

const ColorBlock = (props: Props) => {
  const {color, name, code, className} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <div className={classes.block} style={{backgroundColor: color}} />
      <div className={classes.nameContainer}>
        <Typography className={classes.name}>{name}</Typography>
        <Typography className={classes.hex}>{code ?? color}</Typography>
      </div>
    </div>
  );
};

export default ColorBlock;

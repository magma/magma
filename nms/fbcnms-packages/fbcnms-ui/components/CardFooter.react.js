/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import {withStyles} from '@material-ui/core/styles';
import * as React from 'react';
import CardActions from '@material-ui/core/CardActions';

const styles = {
  cardRoot: {
    boxShadow: '0px -1px 4px 0px rgba(0, 0, 0, 0.11)',
    padding: '12px',
    display: 'flex',
  },
  placeholder: {
    flexGrow: 1,
  },
};

export type FooterAlign = 'left' | 'right';

type Props = {
  alignItems: FooterAlign,
  children: Array<React.Node>,
} & WithStyles;

const CardFooter = (props: Props) => {
  const {alignItems, children, classes} = props;
  return (
    <CardActions classes={{root: classes.cardRoot}}>
      {alignItems === 'right' ? <div className={classes.placeholder} /> : null}
      {children}
    </CardActions>
  );
};

CardFooter.defaultProps = {
  alignItems: 'right',
};

export default withStyles(styles)(CardFooter);

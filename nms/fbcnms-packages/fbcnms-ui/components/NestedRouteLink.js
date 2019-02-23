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
import {Link} from 'react-router-dom';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../hooks';

const useStyles = makeStyles({
  link: {
    textDecoration: 'none',
  },
});

type Props = {
  children: any,
  to: string,
};

export default function NestedRouteLink(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const {children, to, ...childProps} = props;
  // remove trailing/leading slashes
  const base = match.url.replace(/\/$/, '');
  const url = to.replace(/^\//, '');
  return (
    <Link className={classes.link} to={`${base}/${url}`} {...childProps}>
      {children}
    </Link>
  );
}

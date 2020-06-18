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
import {Link} from 'react-router-dom';

import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../hooks';

const useStyles = makeStyles(() => ({
  link: {
    textDecoration: 'none',
  },
}));

type Props = {
  children: React.Node,
  to: string,
  className?: string,
};

function NestedRouteLink(props: Props, ref: React.Ref<*>) {
  const classes = useStyles();
  const {match} = useRouter();
  const {children, to, className: childClassName, ...childProps} = props;
  // remove trailing/leading slashes
  const base = match.url.replace(/\/$/, '');
  const url = to.replace(/^\//, '');
  return (
    <Link
      {...childProps}
      innerRef={ref}
      className={classNames(classes.link, childClassName)}
      to={`${base}/${url}`}>
      {children}
    </Link>
  );
}

export default React.forwardRef<Props, typeof NestedRouteLink>(NestedRouteLink);

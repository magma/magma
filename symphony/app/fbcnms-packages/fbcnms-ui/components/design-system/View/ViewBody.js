/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CombinedRefs} from '../hooks/useCombinedRefs';

import * as React from 'react';
import classNames from 'classnames';
import useCombinedRefs from '../hooks/useCombinedRefs';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useRef, useState} from 'react';

const paddingRight = 24;
const scrollWidth = 12;
const useStyles = makeStyles(() => ({
  viewWrapper: {
    flexGrow: 1,
    overflowX: 'hidden',
    overflowY: 'auto',
    padding: `8px ${paddingRight}px 4px 24px`,
    display: 'flex',
  },
  withScrollY: {
    paddingRight: `${paddingRight - scrollWidth}px`,
  },
}));

type Props = {
  children: React.Node,
};

const ViewBody = React.forwardRef<Props, HTMLElement>((props, ref) => {
  const classes = useStyles();
  const refs: CombinedRefs = [useRef(null), ref];
  const combinedRef = useCombinedRefs(refs);
  const [hasScrollY, setHasScrollY] = useState(false);

  useEffect(() => {
    window.requestAnimationFrame(() => {
      const element = combinedRef.current;
      setHasScrollY(element && element.scrollHeight > element.offsetHeight);
    });
  });

  return (
    <div
      ref={combinedRef}
      className={classNames(
        {
          [classes.withScrollY]: hasScrollY,
        },
        classes.viewWrapper,
      )}>
      {props.children}
    </div>
  );
});

export default ViewBody;

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FullViewHeaderProps} from './ViewHeader';

import * as React from 'react';
import ViewBody from './ViewBody';
import ViewHeader from './ViewHeader';
import useVerticalScrollingEffect from '../hooks/useVerticalScrollingEffect';
import {makeStyles} from '@material-ui/styles';
import {useRef, useState} from 'react';

const useStyles = makeStyles(() => ({
  viewPanel: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    maxHeight: '100%',
  },
}));

export type ViewContainerProps = {
  header?: ?FullViewHeaderProps,
  useBodyScrollingEffect?: ?boolean,
  children: React.Node,
};

export default function ViewContainer(props: ViewContainerProps) {
  const {header, useBodyScrollingEffect = true, children} = props;
  const classes = useStyles();
  const bodyElement = useRef(null);
  const [bodyIsScrolled, setBodyIsScrolled] = useState(false);

  const handleBodyScroll = verticalScrollValues => {
    setBodyIsScrolled(verticalScrollValues.startIsHidden);
  };

  useVerticalScrollingEffect(
    bodyElement,
    handleBodyScroll.bind(this),
    !!useBodyScrollingEffect,
  );

  return (
    <div className={classes.viewPanel}>
      {!!header && <ViewHeader {...header} showMinimal={bodyIsScrolled} />}
      <ViewBody ref={bodyElement}>{children}</ViewBody>
    </div>
  );
}

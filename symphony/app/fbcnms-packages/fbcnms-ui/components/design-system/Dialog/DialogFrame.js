/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Portal from '../Core/Portal';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  anchor: {
    alignItems: 'flex-start',
    display: 'flex',
    justifyContent: 'center',
    pointerEvents: 'none',
  },
  dialog: {
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
    pointerEvents: 'all',
    position: 'relative',
    zIndex: 0,
    backgroundColor: symphony.palette.white,
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP3,
  },
  root: {
    alignItems: 'stretch',
    boxSizing: 'border-box',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    minHeight: '100vh',
    zIndex: 1,
    position: 'fixed',
    left: 0,
    right: 0,
    top: 0,
    bottom: 0,
  },
  hidden: {
    visibility: 'hidden',
  },
  mask: {
    backgroundColor: symphony.palette.overlay,
    position: 'fixed',
    bottom: 0,
    right: 0,
    left: 0,
    top: 0,
  },
}));

type Props = $ReadOnly<{|
  children: React.Node,
  hidden?: boolean,
  onClose?: () => void,
  className?: string,
|}>;

function DialogFrame(props: Props) {
  const {children, className, hidden = false, onClose} = props;
  const classes = useStyles();
  return (
    <Portal target={document.body}>
      <div className={classNames(classes.root, {[classes.hidden]: hidden})}>
        <div className={classes.mask} onClick={() => onClose && onClose()} />
        <div className={classes.anchor}>
          <div className={classNames(classes.dialog, className)}>
            {children}
          </div>
        </div>
      </div>
    </Portal>
  );
}

export default DialogFrame;

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
import {useRef} from 'react';

const SIDE_PANEL_WIDTH = '474px';

const useStyles = makeStyles(() => ({
  dialog: {
    zIndex: 2,
    position: 'fixed',
    display: 'flex',
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP3,
    '&$hidden': {
      visibility: 'hidden',
    },
  },
  center: {
    transform: 'translate(-50%, -50%)',
    left: '50%',
    top: '50%',
    '&:not($hidden)': {
      animation: '$fadeIn 200ms forwards',
    },
    '&$hidden': {
      animation: '$fadeOut 200ms forwards',
    },
  },
  right: {
    bottom: 0,
    top: 0,
    '&:not($hidden)': {
      animation: '$slideIn 500ms forwards',
    },
    '&$hidden': {
      animation: '$slideOut 500ms forwards',
    },
  },
  hidden: {},
  mask: {
    zIndex: 1,
    backgroundColor: symphony.palette.overlay,
    position: 'fixed',
    bottom: 0,
    right: 0,
    left: 0,
    top: 0,
    '&:not($hidden)': {
      animation: '$fadeIn 500ms forwards',
    },
    '&$hidden': {
      animation: '$fadeOut 500ms forwards',
    },
  },
  '@keyframes fadeIn': {
    from: {
      opacity: 0,
      visibility: 'hidden',
    },
    to: {
      opacity: 1,
      visibility: 'visible',
    },
  },
  '@keyframes fadeOut': {
    from: {
      opacity: 1,
      visibility: 'visible',
    },
    to: {
      opacity: 0,
      visibility: 'hidden',
    },
  },
  '@keyframes slideIn': {
    from: {
      right: `-${SIDE_PANEL_WIDTH}`,
      left: '100%',
      visibility: 'hidden',
    },
    to: {
      right: 0,
      left: `calc(100% - ${SIDE_PANEL_WIDTH})`,
      visibility: 'visible',
    },
  },
  '@keyframes slideOut': {
    from: {
      right: 0,
      left: `calc(100% - ${SIDE_PANEL_WIDTH})`,
      visibility: 'visible',
    },
    to: {
      right: `-${SIDE_PANEL_WIDTH}`,
      left: '100%',
      visibility: 'hidden',
    },
  },
}));

export const POSITION = {
  center: 'center',
  right: 'right',
};
export type DialogPosition = $Keys<typeof POSITION>;

type Props = $ReadOnly<{|
  children: React.Node,
  position?: ?DialogPosition,
  isModal?: ?boolean,
  hidden?: boolean,
  onClose?: () => void,
  className?: string,
|}>;

function DialogFrame(props: Props) {
  const {
    children,
    className,
    hidden = false,
    position: positionProp,
    isModal: isModalProp,
    onClose,
  } = props;
  const classes = useStyles();

  const renderedOnce = useRef(false);

  if (hidden && renderedOnce.current === false) {
    return null;
  }
  renderedOnce.current = true;

  const position = positionProp ?? POSITION.center;
  const isModal = isModalProp !== false;

  return (
    <Portal target={document.body}>
      {isModal && (
        <div
          className={classNames(classes.mask, {[classes.hidden]: hidden})}
          onClick={onClose}
        />
      )}
      <div
        className={classNames(
          classes.dialog,
          classes[position],
          {[classes.hidden]: hidden},
          className,
        )}>
        {children}
      </div>
    </Portal>
  );
}

export default DialogFrame;

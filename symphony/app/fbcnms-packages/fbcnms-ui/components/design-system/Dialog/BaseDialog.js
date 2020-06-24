/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DialogPosition} from './DialogFrame';

import * as React from 'react';
import DialogFrame from './DialogFrame';
import IconButton from '../IconButton';
import Text from '../Text';
import ViewContainer from '../View/ViewContainer';
import classNames from 'classnames';
import {CloseIcon} from '../Icons';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    paddingBottom: '1px',
  },
  titleContainer: {
    display: 'flex',
    flexDirection: 'row',
    marginBottom: '16px',
  },
  titleText: {
    flexGrow: 1,
    maxWidth: '560px',
    overflow: 'hidden',
    marginRight: '16px',
  },
  content: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    height: 0,
    overflowY: 'auto',
  },
}));

// TODO(T68894541)
// Consider Moving the stick-right option to separate component
export type BaseDialogProps = $ReadOnly<{|
  position?: ?DialogPosition,
  isModal?: ?boolean,
  title: React.Node,
  children: React.Node,
  showCloseButton?: ?boolean,
  onClose?: ?() => void,
|}>;

export type BaseDialogComponentProps = $ReadOnly<{|
  ...BaseDialogProps,
  hidden?: boolean,
  className?: ?string,
|}>;

function BaseDialog(props: BaseDialogComponentProps) {
  const {
    className,
    title,
    children,
    onClose,
    showCloseButton,
    ...rootProps
  } = props;
  const classes = useStyles();

  const callOnClose = onClose ?? undefined;

  return (
    <DialogFrame
      className={classNames(classes.root, className)}
      onClose={callOnClose}
      {...rootProps}>
      <ViewContainer
        header={{
          title: (
            <div className={classes.titleContainer}>
              <Text className={classes.titleText} weight="medium">
                {title}
              </Text>
              {showCloseButton != false && (
                <IconButton
                  skin="gray"
                  icon={CloseIcon}
                  onClick={callOnClose}
                />
              )}
            </div>
          ),
        }}>
        {children}
      </ViewContainer>
    </DialogFrame>
  );
}

export default BaseDialog;

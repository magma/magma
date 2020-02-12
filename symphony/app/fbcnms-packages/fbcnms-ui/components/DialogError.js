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
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import Text from './design-system/Text';
import classNames from 'classnames';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    display: 'flex',
    flexDirection: 'row',
    borderWidth: '1px',
    borderRadius: '4px',
    borderStyle: 'solid',
    margin: '4px 8px',
    padding: '4px 10px',
    alignItems: 'center',
  },
  errorBorder: {
    borderColor: symphony.palette.R600,
  },
  warningBorder: {
    borderColor: symphony.palette.Y600,
  },
  errorIcon: {
    margin: '2px 8px',
    color: symphony.palette.R600,
  },
  warningIcon: {
    margin: '2px 8px',
    color: symphony.palette.Y600,
  },
});

type Props = {
  message: ?string,
  color?: 'error' | 'warning',
};

const DialogError = ({message, color = 'error'}: Props) => {
  const classes = useStyles();
  return (
    <div
      className={classNames({
        [classes.root]: true,
        [classes.errorBorder]: color == 'error',
        [classes.warningBorder]: color == 'warning',
      })}>
      <ErrorOutlineIcon
        className={classNames({
          [classes.errorIcon]: color == 'error',
          [classes.warningIcon]: color == 'warning',
        })}
      />
      <Text variant="subtitle2" color={color}>
        {message}
      </Text>
    </div>
  );
};

export default DialogError;

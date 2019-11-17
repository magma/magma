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
import SymphonyTheme from '../../../theme/symphony';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useInput} from './InputContext';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    color: SymphonyTheme.palette.D400,
  },
  hasValue: {},
  disabled: {
    '&:not($hasValue) $text': {
      color: SymphonyTheme.palette.disabled,
    },
  },
  text: {
    color: SymphonyTheme.palette.D400,
  },
  clickable: {
    cursor: 'pointer',
  },
}));

type Props = {
  className?: string,
  children: React.Node,
  onClick?: () => void,
};

const InputAffix = (props: Props) => {
  const {children, className, onClick} = props;
  const classes = useStyles();
  const {disabled, value} = useInput();
  return (
    <div
      className={classNames(
        classes.root,
        {
          [classes.disabled]: disabled,
          [classes.hasValue]: Boolean(value),
          [classes.clickable]: onClick !== undefined,
        },
        className,
      )}
      onClick={onClick}>
      {typeof children === 'string' ? (
        <Text className={classes.text} variant="body2">
          {children}
        </Text>
      ) : (
        children
      )}
    </div>
  );
};

export default InputAffix;

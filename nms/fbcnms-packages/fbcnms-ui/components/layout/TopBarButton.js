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
import Checkbox from '@material-ui/core/Checkbox';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  /* Styles applied to the root element. */
  root: {
    color: '#fff',
  },
  checked: {
    color: '#fff',
  },
  colorPrimary: {
    '&$checked': {
      color: '#fff',
    },
    '&$disabled': {
      color: '#fff',
    },
  },
  colorSecondary: {
    '&$checked': {
      color: '#fff',
    },
    '&$disabled': {
      color: '#fff',
    },
  },
});

type Props = {|
  checked?: boolean | string,
  checkedIcon?: React$Node,
  className?: string,
  defaultChecked?: boolean,
  disableRipple?: boolean,
  disabled?: boolean,
  icon?: React$Node,
  indeterminate?: boolean,
  indeterminateIcon?: React$Node,
  inputProps?: Object,
  inputRef?: Function,
  name?: string,
  onChange?: Function,
  tabIndex?: number | string,
  value?: string,
|};

export default function TopBarButton(
  props: Props,
): React$Element<typeof Checkbox> {
  const classes = useStyles();
  return <Checkbox classes={classes} {...props} />;
}

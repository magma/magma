/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Checkbox from '@material-ui/core/Checkbox';
import React from 'react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  /* Styles applied to the root element. */
  root: {
    color: theme.palette.primaryText,
  },
  checked: {
    color: theme.palette.primaryText,
  },
  colorPrimary: {
    '&$checked': {
      color: theme.palette.primaryText,
    },
    '&$disabled': {
      color: theme.palette.primaryText,
    },
  },
  colorSecondary: {
    '&$checked': {
      color: theme.palette.primaryText,
    },
    '&$disabled': {
      color: theme.palette.primaryText,
    },
  },
}));

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

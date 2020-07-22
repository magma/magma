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
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CloseIcon from '@material-ui/icons/Close';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';
import Typography from '@material-ui/core/Typography';
import WarningIcon from '@material-ui/icons/Warning';
import classNames from 'classnames';
import {blue60, gray4, green, red, yellow} from '../theme/colors';
import {makeStyles} from '@material-ui/styles';
import {useSnackbar} from 'notistack';
import {withForwardRef} from './ForwardRef';
import type {ForwardRef} from './ForwardRef';
import type {Variants} from 'notistack';

const useStyles = makeStyles(theme => ({
  root: {
    backgroundColor: theme.palette.common.white,
    boxShadow: '0 0 0 1px #ccd0d5, 0 4px 8px 1px rgba(0,0,0,0.15)',
    borderRadius: '2px',
    display: 'flex',
    width: '420px',
  },
  bar: {
    borderLeft: '6px solid',
    borderRadius: '2px 0px 0px 2px',
  },
  errorBar: {
    borderColor: red,
  },
  successBar: {
    borderColor: green,
  },
  warningBar: {
    borderColor: yellow,
  },
  defaultBar: {
    borderColor: gray4,
  },
  infoBar: {
    borderColor: blue60,
  },
  content: {
    marginLeft: '6px',
    display: 'flex',
    flexDirection: 'row',
    padding: '12px 12px 12px 0px',
    alignItems: 'center',
    flexGrow: 1,
  },
  message: {
    marginLeft: '6px',
    fontSize: '13px',
    lineHeight: '17px',
    flexGrow: 1,
  },
  icon: {
    '&&': {fontSize: '20px'},
    marginRight: '12px',
  },
  errorIcon: {
    '&&': {fill: red},
  },
  successIcon: {
    '&&': {fill: green},
  },
  warningIcon: {
    '&&': {fill: yellow},
  },
  defaultIcon: {
    '&&': {fill: gray4},
  },
  infoIcon: {
    '&&': {fill: blue60},
  },
  closeButton: {
    marginLeft: '16px',
    color: theme.palette.grey[400],
    '&:hover': {
      color: theme.palette.grey[600],
    },
    cursor: 'pointer',
  },
}));

type Props = {
  id: number | string,
  message: string,
  variant: Variants,
} & ForwardRef;

const IconVariants: {[Variants]: React.ComponentType<*>} = {
  error: ErrorIcon,
  success: CheckCircleIcon,
  warning: WarningIcon,
  default: InfoIcon,
  info: InfoIcon,
};

const SnackbarItem = withForwardRef((props: Props) => {
  const {id, message, variant, fwdRef} = props;
  const classes = useStyles();
  const {closeSnackbar} = useSnackbar();
  const Icon = IconVariants[variant];
  return (
    <div className={classes.root} ref={fwdRef}>
      <div className={classNames(classes.bar, classes[variant + 'Bar'])} />
      <div className={classes.content}>
        <Icon className={classNames(classes.icon, classes[variant + 'Icon'])} />
        <Typography className={classes.message}>{message}</Typography>
        <CloseIcon
          className={classes.closeButton}
          onClick={() => closeSnackbar(id)}
        />
      </div>
    </div>
  );
});

export default SnackbarItem;

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
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CloseIcon from '@material-ui/icons/Close';
import ErrorIcon from '@material-ui/icons/Error';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {green, red} from '../theme/colors';
import {makeStyles} from '@material-ui/styles';
import {useSnackbar} from 'notistack';

export type WithSnackbarProps = {
  enqueueSnackbar: (message: string | React.Node, options?: Object) => null,
};

const useStyles = makeStyles(theme => ({
  root: {
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
  id: number,
  message: string,
  variant: 'success' | 'error',
};

const IconVariants = {
  error: ErrorIcon,
  success: CheckCircleIcon,
};

const SnackbarItem = (props: Props) => {
  const {id, message, variant} = props;
  const classes = useStyles();
  const {closeSnackbar} = useSnackbar();
  const Icon = IconVariants[variant];
  return (
    <div className={classes.root}>
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
};

export default SnackbarItem;

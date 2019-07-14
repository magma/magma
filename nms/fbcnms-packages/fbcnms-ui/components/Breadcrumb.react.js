/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {makeStyles} from '@material-ui/styles';
import * as React from 'react';
import classNames from 'classnames';
import KeyboardArrowRightIcon from '@material-ui/icons/KeyboardArrowRight';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  upperSection: {
    display: 'flex',
    flexDirection: 'column',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
  },
  arrowIcon: {
    color: theme.palette.grey[600],
  },
  breadcrumbName: {
    cursor: 'pointer',
    fontWeight: 'bold',
    '&:hover': {
      textDecoration: 'underline',
    },
  },
  parentBreadcrumb: {
    color: theme.palette.grey.A200,
  },
}));

export type BreadcrumbData = {
  id: string,
  name: string,
  subtext?: ?string | React.Node,
  onClick?: ?(id: string) => void,
};

type Props = {
  data: BreadcrumbData,
  isLastBreadcrumb: boolean,
  size?: 'default' | 'small' | 'large',
};

const Breadcrumb = (props: Props) => {
  const {data, isLastBreadcrumb, size} = props;
  const {id, name, subtext, onClick} = data;
  const classes = useStyles();
  return (
    <div key={id} className={classes.root}>
      <div className={classes.upperSection}>
        <Typography
          className={classNames({
            [classes.breadcrumbName]: true,
            [classes.parentBreadcrumb]: !isLastBreadcrumb,
          })}
          variant={size === 'large' || size === 'default' ? 'h6' : 'body2'}
          onClick={() => onClick && onClick(id)}>
          {name}
        </Typography>
        {typeof subtext === 'string' ? (
          <Typography className={classes.subtext}>{subtext}</Typography>
        ) : (
          subtext
        )}
      </div>
      {!isLastBreadcrumb && (
        <KeyboardArrowRightIcon
          className={classes.arrowIcon}
          fontSize={size}
          style={{
            height: size === 'large' || size === 'default' ? '32px' : '21px',
          }}
        />
      )}
    </div>
  );
};

Breadcrumb.defaultProps = {
  size: 'default',
};

export default Breadcrumb;

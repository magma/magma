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
import SymphonyTheme from '../theme/symphony';
import Text from './design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  upperSection: {
    display: 'flex',
    flexDirection: 'column',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(13),
  },
  slash: {
    color: SymphonyTheme.palette.D400,
    margin: '0 6px',
  },
  breadcrumbName: {
    whiteSpace: 'nowrap',
    fontWeight: 500,
    color: theme.palette.blueGrayDark,
  },
  parentBreadcrumb: {
    color: SymphonyTheme.palette.D400,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    display: 'inline-block',
  },
  hover: {
    '&:hover': {
      color: theme.palette.primary.main,
    },
    cursor: 'pointer',
  },
  largeText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  smallText: {
    fontSize: '14px',
    lineHeight: '24px',
    fontWeight: 500,
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
        <Tooltip
          placement="top"
          title={
            typeof subtext === 'string' ? (
              <Text className={classes.subtext} color="light">
                {subtext}
              </Text>
            ) : (
              subtext ?? ''
            )
          }>
          <Text
            className={classNames({
              [classes.breadcrumbName]: true,
              [classes.parentBreadcrumb]: !isLastBreadcrumb,
              [classes.hover]: !!onClick,
              [classes.smallText]: size === 'small',
              [classes.largeText]: size !== 'small',
            })}
            onClick={() => onClick && onClick(id)}>
            {name}
          </Text>
        </Tooltip>
      </div>
      {!isLastBreadcrumb && (
        <Text
          className={classNames({
            [classes.slash]: true,
            [classes.smallText]: size === 'small',
            [classes.largeText]: size !== 'small',
          })}>
          {'/'}
        </Text>
      )}
    </div>
  );
};

Breadcrumb.defaultProps = {
  size: 'default',
};

export default Breadcrumb;

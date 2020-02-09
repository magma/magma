/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TextVariant} from '../theme/symphony';

import * as React from 'react';
import SymphonyTheme from '../theme/symphony';
import Text from './design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {typographyStyles} from './design-system/Text';

const useStyles = makeStyles(theme => ({
  root: {
    overflowX: 'hidden',
    textOverflow: 'ellipsis',
    display: 'flex',
  },
  notShrinkable: {
    flexShrink: '0',
  },
  tooltipAnchor: {
    overflowX: 'hidden',
    textOverflow: 'ellipsis',
    color: theme.palette.blueGrayDark,
  },
  parentBreadcrumb: {
    color: SymphonyTheme.palette.D400,
  },
  hover: {
    '&:hover': {
      color: theme.palette.primary.main,
    },
    cursor: 'pointer',
  },
  slash: {
    color: SymphonyTheme.palette.D400,
    margin: '0 6px',
  },
  breadcrumbName: {
    color: 'unset',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
}));

export type BreadcrumbData = {
  id: string,
  name: string,
  subtext?: ?string | React.Node,
  onClick?: ?(id: string, target: HTMLElement) => void,
};

type Props = {
  data: BreadcrumbData,
  isLastBreadcrumb: boolean,
  useEllipsis?: boolean,
  size?: 'default' | 'small' | 'large',
  variant?: TextVariant,
};

const Breadcrumb = (props: Props) => {
  const {data, isLastBreadcrumb, size, variant, useEllipsis} = props;
  const {id, name, subtext, onClick} = data;
  const classes = useStyles();
  const typographyClasses = typographyStyles();
  const textVariant = variant ? variant : size === 'small' ? 'subtitle2' : 'h6';

  return (
    <div
      key={id}
      className={classNames(
        {
          [classes.notShrinkable]: !useEllipsis,
        },
        classes.root,
      )}>
      <Tooltip
        arrow
        interactive
        placement="top"
        title={
          typeof subtext === 'string' ? (
            <Text className={classes.subtext} variant="caption" color="light">
              {subtext}
            </Text>
          ) : (
            subtext ?? ''
          )
        }>
        <div
          className={classNames(
            {
              [classes.parentBreadcrumb]: !isLastBreadcrumb,
              [classes.hover]: !!onClick,
            },
            classes.tooltipAnchor,
            typographyClasses[textVariant],
          )}>
          <Text
            variant={textVariant}
            className={classes.breadcrumbName}
            onClick={e => onClick && onClick(id, e.currentTarget)}>
            {name}
          </Text>
        </div>
      </Tooltip>
      {!isLastBreadcrumb && (
        <Text variant={textVariant} className={classes.slash}>
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

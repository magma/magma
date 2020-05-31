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
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    padding: '3px',
    '&:not(:first-child)': {
      marginTop: '16px',
    },
  },
  cardContainer: {
    width: '100%',
    maxWidth: '100%',
    height: '100%',
    overflow: 'hidden',
    boxSizing: 'border-box',
    display: 'flex',
    flexDirection: 'column',
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    backgroundColor: symphony.palette.white,
  },
  standardMargins: {
    padding: '24px',
  },
}));

export const CARD_MARGINS = {
  none: 'none',
  standard: 'standard',
};
type Margins = $Keys<typeof CARD_MARGINS>;

type Props = $ReadOnly<{|
  className?: ?string,
  contentClassName?: ?string,
  margins?: ?Margins,
  children: React.Node,
|}>;

const Card = (props: Props) => {
  const {children, margins: marginsProp, className, contentClassName} = props;
  const classes = useStyles();
  const margins: string & Margins = marginsProp || CARD_MARGINS.standard;

  return (
    <div className={classNames(classes.root, className)}>
      <div
        className={classNames(
          classes.cardContainer,
          classes[`${margins}Margins`],
          contentClassName,
        )}>
        {children}
      </div>
    </div>
  );
};

export default Card;

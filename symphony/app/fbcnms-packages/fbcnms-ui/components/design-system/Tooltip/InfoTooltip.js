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
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import InfoTinyIcon from '../Icons/Indications/InfoTinyIcon';
import Text from '../Text';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  tooltip: {
    borderRadius: '2px',
    backgroundColor: symphony.palette.secondary,
    padding: '8px 10px',
    width: '181px',
  },
  iconContainer: {
    display: 'inline-flex',
    '&:hover $icon': {
      fill: symphony.palette.primary,
    },
  },
  icon: {},
  arrow: {
    position: 'absolute',
    bottom: '0',
    left: '8px',
    '&:before': {
      borderBottom: '4px solid transparent',
      borderLeft: '4px solid transparent',
      borderRight: '4px solid transparent',
      borderTop: `4px solid ${symphony.palette.secondary}`,
      botom: '-8px',
      content: '""',
      position: 'absolute',
      zIndex: 10,
    },
  },
});

type Props = {
  description: React.Node,
};

const InfoTooltip = ({description}: Props) => {
  const classes = useStyles();
  return (
    <BasePopoverTrigger
      position="above"
      popover={
        <div className={classes.tooltip}>
          <Text color="light" variant="caption">
            {description}
          </Text>
          <span className={classes.arrow} />
        </div>
      }>
      {(onShow, onHide, contextRef) => (
        <div
          className={classes.iconContainer}
          ref={contextRef}
          onMouseOver={onShow}
          onMouseOut={onHide}>
          <InfoTinyIcon color="gray" className={classes.icon} />
        </div>
      )}
    </BasePopoverTrigger>
  );
};

export default InfoTooltip;

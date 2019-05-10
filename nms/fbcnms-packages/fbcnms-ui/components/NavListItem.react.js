/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Link} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import classNames from 'classnames';
import React, {useState, useCallback} from 'react';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles(theme => ({
  icon: {
    color: theme.palette.gray50,
  },
  link: {
    width: '100%',
  },
  root: {
    display: 'flex',
    justifyContent: 'center',
    width: '100%',
    padding: '15px 0px',
    '&:hover $icon, &$selected $icon': {
      color: theme.palette.common.white,
    },
  },
  selected: {
    backgroundColor: theme.palette.primary.main,
  },
  tooltip: {
    '&&': {
      padding: '8px 12px',
      backgroundColor: theme.palette.primary.dark,
    },
  },
  arrow: {
    position: 'absolute',
    left: '-8px',
    '&:before': {
      borderBottom: '4px solid transparent',
      borderLeft: '4px solid transparent',
      borderRight: `4px solid ${theme.palette.primary.dark}`,
      borderTop: '4px solid transparent',
      top: '-5px',
      content: '""',
      position: 'absolute',
      zIndex: 10,
    },
  },
  bootstrapPlacementLeft: {
    margin: '0 8px',
  },
  tooltipLabel: {
    '&&': {
      fontSize: '12px',
      lineHeight: '16px',
      color: theme.palette.common.white,
      fontWeight: 'bold',
    },
  },
}));

type Props = {
  path: string,
  label: string,
  icon: any,
  hidden: boolean,
};

export default function NavListItem(props: Props) {
  const {hidden} = props;
  const classes = useStyles();
  const router = useRouter();
  const [arrowArrow, setArrowRef] = useState(null);
  const handleArrowRef = useCallback(node => {
    if (node !== null) {
      setArrowRef(node);
    }
  }, []);

  if (hidden) {
    return null;
  }

  const isSelected = router.location.pathname.includes(props.path);

  return (
    <Link to={props.path} className={classes.link}>
      <Tooltip
        placement="right"
        title={
          <>
            <Typography className={classes.tooltipLabel}>
              {props.label}
            </Typography>
            <span className={classes.arrow} ref={handleArrowRef} />
          </>
        }
        classes={{
          tooltip: classes.tooltip,
          popper: classes.arrowPopper,
          tooltipPlacementLeft: classes.bootstrapPlacementLeft,
        }}
        PopperProps={{
          popperOptions: {
            modifiers: {
              arrow: {
                enabled: Boolean(arrowArrow),
                element: arrowArrow,
              },
            },
          },
        }}>
        <div
          className={classNames({
            [classes.root]: true,
            [classes.selected]: isSelected,
          })}>
          <div className={classes.icon}>{props.icon}</div>
        </div>
      </Tooltip>
    </Link>
  );
}

NavListItem.defaultProps = {
  hidden: false,
};

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
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  tabs: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'row',
    backgroundColor: 'white',
    overflowX: 'auto',
    overflowY: 'hidden',
    '&$spread $tab': {
      flexShrink: 1,
      flexBasis: '250px',
    },
  },
  divider: {
    position: 'absolute',
    left: 0,
    right: 0,
    bottom: 0,
    borderBottom: `1px solid ${symphony.palette.D50}`,
  },
  standard: {
    minHeight: '48px',
    height: '48px',
    padding: '0px 16px',
    '& $tab': {
      margin: '0px 8px',
      paddingLeft: '8px',
      paddingRight: '8px',
    },
  },
  large: {
    minHeight: '56px',
    height: '56px',
    padding: '0px 20px',
    '& $tab': {
      margin: '0px 8px',
      paddingLeft: '12px',
      paddingRight: '12px',
    },
  },
  tab: {
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    whiteSpace: 'nowrap',
    cursor: 'pointer',
    '&:hover:not($disabledTab) $tabName': {
      color: symphony.palette.primary,
    },
  },
  tabName: {
    color: symphony.palette.D400,
  },
  selectedTab: {
    '& $tabName': {
      color: symphony.palette.primary,
    },
  },
  disabledTab: {
    cursor: 'default',
    '& $tabName': {
      color: symphony.palette.disabled,
    },
  },
  selectedTabIndicator: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: symphony.palette.primary,
    height: '2px',
    zIndex: 2,
  },
  small: {
    minHeight: '28px',
    height: '28px',
    '& $tab': {
      margin: '0px 12px',
      '&:first-child': {
        marginLeft: '0px',
      },
      '&:last-child': {
        marginRight: '0px',
      },
      paddingLeft: '4px',
      paddingRight: '4px',
    },
    '& $divider': {
      display: 'none',
    },
  },
  spread: {},
}));

export type TabProps = {|
  label: string,
  className?: ?string,
  disabled?: boolean,
|};

export type Props = {
  tabs: Array<TabProps>,
  activeTabIndex: number,
  onChange?: number => void,
  size?: 'small' | 'standard' | 'large',
  spread?: ?boolean,
  className?: ?string,
};

export default function TabsBar(props: Props) {
  const {
    spread = true,
    activeTabIndex = 0,
    onChange,
    tabs,
    className,
    size = 'standard',
  } = props;
  const classes = useStyles();
  return (
    <div
      className={classNames(
        classes.tabs,
        {[classes.spread]: spread},
        classes[size],
        className,
      )}>
      {tabs.map((tab, ind) => (
        <div
          key={`tab${ind}`}
          className={classNames(
            classes.tab,
            {
              [classes.selectedTab]: activeTabIndex === ind,
              [classes.disabledTab]: tab.disabled === true,
            },
            tab.className,
          )}
          onClick={() => tab.disabled !== true && onChange && onChange(ind)}>
          <Text className={classes.tabName} variant="body1" weight="medium">
            {tab.label}
          </Text>
          {activeTabIndex === ind && (
            <div className={classes.selectedTabIndicator} />
          )}
        </div>
      ))}
      <div className={classes.divider} />
    </div>
  );
}

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
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  tabsContainer: {
    height: '100%',
  },
  standard: {},
  tabs: {
    backgroundColor: 'white',
    borderBottom: `1px ${symphony.palette.separatorLight} solid`,
    overflowX: 'auto',
    overflowY: 'hidden',
    '&$spread $tab': {
      flexShrink: 1,
      flexBasis: '250px',
    },
  },
  large: {
    minHeight: '60px',
  },
  tab: {
    textTransform: 'unset',
    minWidth: 'unset',
    padding: '0 24px',
    whiteSpace: 'nowrap',
  },
  spread: {},
}));

export type TabProps = {|
  label: string,
  className?: ?string,
|};

export type Props = {
  tabs: Array<TabProps>,
  activeTabIndex: number,
  onChange?: number => void,
  size?: 'standard' | 'large',
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
    <Tabs
      className={classNames(
        classes.tabs,
        {[classes.spread]: spread},
        classes[size],
        className,
      )}
      classes={{flexContainer: classes.tabsContainer}}
      value={activeTabIndex}
      onChange={
        onChange
          ? (_e, newActiveTabIndex) => onChange(newActiveTabIndex)
          : undefined
      }
      indicatorColor="primary"
      textColor="primary">
      {tabs.map((tab, ind) => (
        <Tab
          key={`tab${ind}`}
          classes={{root: classNames(classes.tab, tab.className)}}
          label={tab.label}
        />
      ))}
    </Tabs>
  );
}

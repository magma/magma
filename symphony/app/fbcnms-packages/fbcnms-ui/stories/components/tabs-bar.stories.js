/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TabProps} from '../../components/design-system/Tabs/TabsBar';

import React, {useState} from 'react';
import TabsBar from '../../components/design-system/Tabs/TabsBar';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
  },
  tabsContainer: {
    marginTop: '16px',
  },
}));

const tabs: Array<TabProps> = [
  {
    label: 'Option 1',
  },
  {
    label: 'Option 2',
  },
  {
    label: 'Option 3',
  },
  {
    label: 'Option 4',
    disabled: true,
  },
];

const TabsBarRoot = () => {
  const classes = useStyles();
  const [activeTab, setActiveTab] = useState(0);

  return (
    <div className={classes.root}>
      <TabsBar
        className={classes.tabsContainer}
        tabs={tabs}
        size="large"
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={false}
      />
      <TabsBar
        className={classes.tabsContainer}
        tabs={tabs}
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={false}
      />
      <TabsBar
        className={classes.tabsContainer}
        tabs={tabs}
        size="large"
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={true}
      />
      <TabsBar
        className={classes.tabsContainer}
        tabs={tabs}
        size="small"
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={false}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('TabsBar', () => (
  <TabsBarRoot />
));

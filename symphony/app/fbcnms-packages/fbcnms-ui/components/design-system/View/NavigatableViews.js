/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {MenuItem} from '@fbcnms/ui/components/design-system/Menu/SideMenu';
import type {ViewContainerProps} from '@fbcnms/ui/components/design-system/View/ViewContainer';

import * as React from 'react';
import SideMenu from '@fbcnms/ui/components/design-system/Menu/SideMenu';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    maxHeight: '100%',
    minHeight: '100%',
    width: '100%',
    overflow: 'hidden',
  },
}));

export const NAVIGATION_TYPES = {
  side: 'side',
};

export type NavigatableView = {
  navigation: MenuItem,
} & ViewContainerProps;

type Props = {
  navigation?: $Keys<typeof NAVIGATION_TYPES>,
  header?: ?React.Node,
  views: Array<NavigatableView>,
};

export default function NavigatableViews(props: Props) {
  const {header, navigation = NAVIGATION_TYPES.side, views} = props;
  const classes = useStyles();
  const [activeView, setActiveView] = useState(0);

  return (
    <div className={classes.root}>
      {navigation === NAVIGATION_TYPES.side && (
        <SideMenu
          header={header}
          items={views.map(view => view.navigation)}
          activeItemIndex={activeView}
          onActiveItemChanged={(_item, index) => setActiveView(index)}
        />
      )}
      <ViewContainer {...views[activeView]} />
    </div>
  );
}

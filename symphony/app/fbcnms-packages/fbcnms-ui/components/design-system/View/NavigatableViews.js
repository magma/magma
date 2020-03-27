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
import {
  Redirect,
  Route,
  Switch,
  useHistory,
  useLocation,
} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    maxHeight: '100%',
    minHeight: '100%',
    width: '100%',
    overflow: 'hidden',
  },
}));

export const NAVIGATION_VARIANTS = {
  side: 'side',
};

export type NavigatableView_MenuItemOnly = {|
  menuItem: MenuItem,
|};
export type NavigatableView_MenuItemWithRelatedComponent = {|
  menuItem: MenuItem,
  component: ViewContainerProps,
|};
export type NavigatableView_MenuItemRoutingToGivenComponent = {|
  menuItem: MenuItem,
  component: ViewContainerProps,
  routingPath: string,
|};
export type NavigatableView_MenuItemRouting = {|
  menuItem: MenuItem,
  routingPath: string,
|};
export type NavigatableView_PossibleRoutingToGivenComponent = {|
  component: ViewContainerProps,
  routingPath: string,
|};
export type NavigatableView =
  | NavigatableView_MenuItemOnly
  | NavigatableView_MenuItemWithRelatedComponent
  | NavigatableView_MenuItemRoutingToGivenComponent
  | NavigatableView_MenuItemRouting
  | NavigatableView_PossibleRoutingToGivenComponent;

type Props = {
  variant?: $Keys<typeof NAVIGATION_VARIANTS>,
  header?: ?React.Node,
  views: Array<NavigatableView>,
  routingBasePath?: string,
};

export default function NavigatableViews(props: Props) {
  const {
    header,
    variant = NAVIGATION_VARIANTS.side,
    views,
    routingBasePath,
  } = props;
  const classes = useStyles();
  const history = useHistory();
  const location = useLocation();
  const [activeViewIndex, setActiveView] = useState(0);

  const onNavigation = useCallback(
    navigatedViewIndex => {
      if (routingBasePath == null) {
        setActiveView(navigatedViewIndex);
        return;
      }
      const navigatedView = views[navigatedViewIndex];
      if (navigatedView.routingPath == null) {
        return;
      }
      history.push(`${routingBasePath}/${navigatedView.routingPath}`);
    },
    [history, routingBasePath, views],
  );

  useEffect(() => {
    if (routingBasePath == null) {
      return;
    }
    const activePath = location.pathname.substring(routingBasePath.length + 1);
    const relatedActiveViewIndex = views.findIndex(
      view => view.routingPath != null && view.routingPath === activePath,
    );
    if (
      relatedActiveViewIndex !== -1 &&
      relatedActiveViewIndex !== activeViewIndex
    ) {
      setActiveView(relatedActiveViewIndex);
    }
  }, [views, activeViewIndex, routingBasePath, location.pathname]);

  const menuItems = useMemo(() => {
    const arr: Array<MenuItem> = [];
    views.forEach(view => {
      if (view.menuItem == null) {
        return;
      }
      arr.push(view.menuItem);
      // Why with 'forEach' and not filter&map - good question!
      // Flow doesn't allow this :
      //  views
      //    .filter(view => view.menuItem != null)
      //    .map(view => view.menuItem);
    });
    return arr;
  }, [views]);
  const routableViews = useMemo(() => {
    const arr: Array<{path: string, component: ViewContainerProps}> = [];
    views.forEach(view => {
      if (view.routingPath == null || view.component == null) {
        return;
      }
      arr.push({path: view.routingPath, component: view.component});
    });
    return arr;
  }, [views]);

  if (views.length === 0) {
    return null;
  }

  const activeView = views[activeViewIndex];
  return (
    <div className={classes.root}>
      {menuItems.length > 0 && variant === NAVIGATION_VARIANTS.side && (
        <SideMenu
          header={header}
          items={menuItems}
          activeItemIndex={activeViewIndex}
          onActiveItemChanged={(_item, index) => onNavigation(index)}
        />
      )}
      {activeView.component != null &&
      (routingBasePath == null || activeView.routingPath == null) ? (
        <ViewContainer {...activeView.component} />
      ) : null}
      {routingBasePath != null && routableViews.length > 0 ? (
        <Switch>
          {routableViews.map(routableView =>
            routableView.component.header != null ? (
              <Route
                key={routableView.path}
                path={`${routingBasePath}/${routableView.path}`}
                render={() => <ViewContainer {...routableView.component} />}
              />
            ) : (
              <Route
                key={routableView.path}
                path={`${routingBasePath}/${routableView.path}`}
                children={routableView.component.children}
              />
            ),
          )}
          <Redirect
            from={`${routingBasePath}/`}
            to={`${routingBasePath}/${routableViews[0].path}`}
          />
        </Switch>
      ) : null}
    </div>
  );
}

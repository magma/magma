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
  relatedMenuItemIndex: ?number,
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

const PATH_DELIMITER = '/';
const PATH_PARAM_PREFIX = ':';
const getPathParts: string => Array<string> = path => {
  const parts = path.split(PATH_DELIMITER);
  if (parts[0] == '') {
    return parts.slice(1);
  }
  return parts;
};

const pathNameFitsDefinition = (definition: string, pathName: string) => {
  const definitionParts = getPathParts(definition);
  const pathParts = getPathParts(pathName);

  if (definitionParts.length != pathParts.length) {
    return false;
  }

  let i = 0;
  while (i < definitionParts.length) {
    if (
      !definitionParts[i].startsWith(PATH_PARAM_PREFIX) &&
      pathParts[i] != definitionParts[i]
    ) {
      return false;
    }
    i++;
  }

  return true;
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
      view =>
        view.routingPath != null &&
        pathNameFitsDefinition(view.routingPath, activePath),
    );
    if (
      relatedActiveViewIndex !== -1 &&
      relatedActiveViewIndex !== activeViewIndex
    ) {
      setActiveView(relatedActiveViewIndex);
    }
  }, [views, activeViewIndex, routingBasePath, location.pathname]);

  const menuItemViews = useMemo(
    () =>
      views
        .map((view, ind) => {
          if (view.menuItem != null) {
            return {menuItem: view.menuItem, viewIndex: ind};
          }
          return null;
        })
        .filter(Boolean),
    [views],
  );
  const routableViews = useMemo(
    () =>
      views
        .map(view => {
          if (view.routingPath != null && view.component != null) {
            return {path: view.routingPath, component: view.component};
          }
          return null;
        })
        .filter(Boolean),
    [views],
  );

  if (views.length === 0) {
    return null;
  }

  const activeView = views[activeViewIndex];
  return (
    <div className={classes.root}>
      {menuItemViews.length > 0 && variant === NAVIGATION_VARIANTS.side && (
        <SideMenu
          header={header}
          items={menuItemViews.map(item => item.menuItem)}
          activeItemIndex={
            activeView.relatedMenuItemIndex != null
              ? activeView.relatedMenuItemIndex
              : activeViewIndex
          }
          onActiveItemChanged={(_item, index) =>
            onNavigation(menuItemViews[index].viewIndex)
          }
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

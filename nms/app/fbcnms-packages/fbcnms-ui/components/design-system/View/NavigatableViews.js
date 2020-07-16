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
  targetPath?: string,
|};
export type NavigatableView_MenuItemWithRelatedComponent = {|
  ...NavigatableView_MenuItemOnly,
  component: ViewContainerProps,
|};
export type NavigatableView_MenuItemRoutingToGivenComponent = {|
  ...NavigatableView_MenuItemWithRelatedComponent,
  routingPath: string,
|};
export type NavigatableView_MenuItemRouting = {|
  ...NavigatableView_MenuItemOnly,
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

function getViewTargetPath(view: NavigatableView) {
  return view.targetPath != null
    ? view.targetPath
    : view.routingPath != null
    ? view.routingPath
    : null;
}

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
      const viewTargetPath = getViewTargetPath(navigatedView);
      if (viewTargetPath == null) {
        return;
      }
      history.push(`${routingBasePath}/${viewTargetPath}`);
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

  const menu = useMemo(
    () =>
      views
        .map((view, ind) => {
          if (view.menuItem != null) {
            return {item: view.menuItem, viewIndex: ind};
          }
          return null;
        })
        .filter(Boolean),
    [views],
  );

  const routes = useMemo(
    () =>
      views
        .map(view => {
          if (view.routingPath != null && view.component != null) {
            return {
              path: view.routingPath,
              component: view.component,
            };
          }
          return null;
        })
        .filter(Boolean),
    [views],
  );

  const firstRoutablePath = useMemo(() => {
    let routablePath;
    let viewIndex = 0;
    while (routablePath == null && viewIndex < views.length) {
      routablePath = getViewTargetPath(views[viewIndex]);
      viewIndex++;
    }
    return routablePath;
  }, [views]);

  const activeView = views[activeViewIndex];
  const menuActiveItemIndex = useMemo(
    () =>
      menu.findIndex(
        m =>
          m.viewIndex ==
          (activeView.relatedMenuItemIndex != null
            ? activeView.relatedMenuItemIndex
            : activeViewIndex),
      ),
    [activeView, activeViewIndex, menu],
  );

  if (views.length === 0) {
    return null;
  }
  return (
    <div className={classes.root}>
      {menu.length > 0 && variant === NAVIGATION_VARIANTS.side && (
        <SideMenu
          header={header}
          items={menu.map(m => m.item)}
          activeItemIndex={menuActiveItemIndex}
          onActiveItemChanged={(_item, index) =>
            onNavigation(menu[index].viewIndex)
          }
        />
      )}
      {activeView.component != null &&
      (routingBasePath == null || activeView.routingPath == null) ? (
        <ViewContainer {...activeView.component} />
      ) : null}
      {routingBasePath != null && routes.length > 0 ? (
        <Switch>
          {routes.map(route =>
            route.component.header != null ? (
              <Route
                key={route.path}
                path={`${routingBasePath}/${route.path}`}
                render={() => <ViewContainer {...route.component} />}
              />
            ) : (
              <Route
                key={route.path}
                path={`${routingBasePath}/${route.path}`}
                children={route.component.children}
              />
            ),
          )}
          {firstRoutablePath != null && (
            <Redirect
              from={`${routingBasePath}/`}
              to={`${routingBasePath}/${firstRoutablePath}`}
            />
          )}
        </Switch>
      ) : null}
    </div>
  );
}

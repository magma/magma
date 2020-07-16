/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import React from 'react';
import Settings from '../Settings';

import useSections from './useSections';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

export default function SectionRoutes() {
  const {relativePath, match} = useRouter();
  const [landingPath, sections] = useSections();

  return (
    <Switch>
      {sections.map(section => (
        <Route
          key={section.path}
          path={relativePath(`/${section.path}`)}
          component={section.component}
        />
      ))}
      <Route
        key="settings"
        path={relativePath(`/settings`)}
        component={Settings}
      />
      {landingPath && (
        <Route
          path={relativePath('')}
          render={() => (
            <Redirect to={`/nms/${match.params.networkId}/${landingPath}`} />
          )}
        />
      )}
      <LoadingFiller />
    </Switch>
  );
}

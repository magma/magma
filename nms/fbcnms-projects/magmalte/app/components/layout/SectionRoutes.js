/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import LoadingFiller from '../LoadingFiller';
import MagmaTopBar from '../MagmaTopBar';
import Settings from '../Settings';

import {Route, Switch} from 'react-router-dom';
import useSections from './useSections';
import {useRouter} from '@fbcnms/ui/hooks';

export default function SectionRoutes() {
  const {match} = useRouter();
  const sections = useSections();

  if (!sections.length) {
    return (
      <>
        <MagmaTopBar />
        <LoadingFiller />
      </>
    );
  }

  return (
    <Switch>
      {sections.map(section => (
        <Route
          key={section.path}
          path={`${match.path}/${section.path}`}
          component={section.component}
        />
      ))}
      <Route
        key="settings"
        path={`${match.path}/settings`}
        component={Settings}
      />
    </Switch>
  );
}

/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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

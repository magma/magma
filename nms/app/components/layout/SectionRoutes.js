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

import AccountSettings from '../AccountSettings';
import Admin from '../admin/Admin';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AppContext from '../../../app/components/context/AppContext';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import React, {useContext} from 'react';
import useSections from './useSections';
import {Navigate, Route, Routes} from 'react-router-dom';

export default function SectionRoutes() {
  const [landingPath, sections] = useSections();
  const {user, ssoEnabled} = useContext(AppContext);

  return (
    <Routes>
      {sections.map(section => (
        <Route
          key={section.path}
          path={`${section.path}/*`}
          element={<section.component />}
        />
      ))}
      {user.isSuperUser && (
        <Route key="admin" path="admin/*" element={<Admin />} />
      )}
      {!ssoEnabled && <Route path="settings" element={<AccountSettings />} />}
      {landingPath && (
        <Route index element={<Navigate to={landingPath} replace />} />
      )}
      <Route element={<LoadingFiller />} />
    </Routes>
  );
}

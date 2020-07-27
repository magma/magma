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

import * as React from 'react';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import NetworkContext from '../context/NetworkContext';

import useSections from './useSections';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export default function SectionLinks(): React.Node {
  const [_landingPath, sections] = useSections();
  const {relativeUrl} = useRouter();
  const {networkId} = useContext(NetworkContext);

  if (!sections) {
    return <LoadingFiller />;
  }

  if (!networkId) {
    return null;
  }

  return (
    <>
      {sections.map(section => (
        <NavListItem
          key={section.label}
          label={section.label}
          path={relativeUrl(`/${section.path}`)}
          icon={section.icon}
        />
      ))}
    </>
  );
}

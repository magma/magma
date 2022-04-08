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

import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import NavListItem from '../../../fbc_js_core/ui/components/NavListItem';
import NetworkContext from '../context/NetworkContext';
import type {SectionsConfigs} from './Section';

import useSections from './useSections';
import {useContext} from 'react';
import {useRouter} from '../../../fbc_js_core/ui/hooks';

type Props = {
  sections: ?SectionsConfigs,
};

/**
 * SectionLinks is the vertical navigation panel for a network
 */
export default function SectionLinks(props: Props): React.Node {
  let [_landingPath, sections] = useSections();
  const {relativeUrl} = useRouter();
  const {networkId} = useContext(NetworkContext);
  if (props.sections) {
    sections = props.sections;
  }

  if (!sections) {
    return <LoadingFiller />;
  }

  if (!networkId) {
    return null;
  }

  console.log('SectionLinks, sections: ', sections);

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

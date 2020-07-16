/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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

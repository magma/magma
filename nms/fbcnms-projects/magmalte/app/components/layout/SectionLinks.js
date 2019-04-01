/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useContext} from 'react';

import Divider from '@material-ui/core/Divider';

import List from '@material-ui/core/List';
import LoadingFiller from '../LoadingFiller';
import NavListItem from '@fbcnms/ui/components/NavListItem.react';
import SettingsIcon from '@material-ui/icons/Settings';

import {useRouter} from '@fbcnms/ui/hooks';
import useSections from './useSections';
import NetworkContext from '../context/NetworkContext';

export default function SectionLinks() {
  const sections = useSections();
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
      <List>
        {sections.map(section => (
          <NavListItem
            key={section.label}
            label={section.label}
            path={relativeUrl(`/${section.path}`)}
            icon={section.icon}
          />
        ))}
      </List>
      <Divider />
      <List>
        <NavListItem
          label="Settings"
          path={relativeUrl('/settings/security/')}
          icon={<SettingsIcon />}
        />
      </List>
    </>
  );
}

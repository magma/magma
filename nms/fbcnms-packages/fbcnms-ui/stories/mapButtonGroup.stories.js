/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
declare module '@material-ui/icons/AccessAlarms' {
  declare module.exports: React$ComponentType<SvgIconExports>;
}

// $FlowFixMe: ListAlt exists
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapIcon from '@material-ui/icons/Map';
import React from 'react';
import {storiesOf} from '@storybook/react';

const AddMapButtonGroup = () => {
  return (
    <MapButtonGroup
      onIconClicked={() => {}}
      icons={[<ListAltIcon />, <MapIcon />]}
    />
  );
};
const AddThreeMapButton = () => {
  return (
    <MapButtonGroup
      onIconClicked={() => {}}
      icons={[<ListAltIcon />, <MapIcon />, <MapIcon />]}
    />
  );
};

storiesOf('MapButtonGroup', module).add('two', () => {
  return <AddMapButtonGroup />;
});
storiesOf('MapButtonGroup', module).add('three', () => {
  return <AddThreeMapButton />;
});

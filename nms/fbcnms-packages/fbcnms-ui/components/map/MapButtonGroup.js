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
import MapButton from '@fbcnms/ui/components/map/MapButton';
import MapToggleButtonGroup from '@fbcnms/ui/components/map/MapToggleButtonGroup';
import MapToggleContainer from '@fbcnms/ui/components/map/MapToggleContainer';
import {useState} from 'react';

type Props = {
  icons: Array<React.Node>,
  onIconClicked: (id: number) => void,
};

const MapButtonGroup = (props: Props) => {
  const [selectedButtonId, setSelectedButtonId] = useState(0);
  const {onIconClicked, icons} = props;
  return (
    <MapToggleContainer>
      <MapToggleButtonGroup>
        <>
          {icons.map((icon, i) => {
            return (
              <MapButton
                onClick={() => {
                  setSelectedButtonId(i);
                  onIconClicked(i);
                }}
                icon={icon}
                isSelected={selectedButtonId === i}
              />
            );
          })}
        </>
      </MapToggleButtonGroup>
    </MapToggleContainer>
  );
};

export default MapButtonGroup;

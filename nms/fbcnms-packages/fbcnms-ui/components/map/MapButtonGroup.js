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

type ButtonItem = {
  item: React.Node,
  id: string,
};

type Props = {
  buttons: Array<ButtonItem>,
  onIconClicked: (id: string) => void,
};

const MapButtonGroup = (props: Props) => {
  const [selectedButtonId, setSelectedButtonId] = useState(0);
  const {onIconClicked, buttons} = props;
  return (
    <MapToggleContainer>
      <MapToggleButtonGroup>
        <>
          {buttons.map((button, i) => {
            return (
              <MapButton
                key={button.id}
                onClick={() => {
                  setSelectedButtonId(i);
                  onIconClicked(button.id);
                }}
                icon={button.item}
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

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import MapButton from '@fbcnms/ui/components/map/MapButton';
import MapToggleButtonGroup from '@fbcnms/ui/components/map/MapToggleButtonGroup';
import MapToggleContainer from '@fbcnms/ui/components/map/MapToggleContainer';
import Text from '../design-system/Text';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

type ButtonItem = {
  item: React.Node | string,
  id: string,
};

type Props = {
  buttons: Array<ButtonItem>,
  onIconClicked: (id: string) => void,
  initiallySelectedButton?: number,
};

const useStyles = makeStyles(() => ({
  text: {
    fontSize: '12px',
    LineHeight: '14px',
  },
}));

const MapButtonGroup = (props: Props) => {
  const {onIconClicked, buttons} = props;
  const [selectedButtonId, setSelectedButtonId] = useState(
    props.initiallySelectedButton,
  );
  const classes = useStyles();
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
                icon={
                  typeof button.item === 'string' ? (
                    <Text className={classes.text}>{button.item}</Text>
                  ) : (
                    button.item
                  )
                }
                isSelected={selectedButtonId === i}
              />
            );
          })}
        </>
      </MapToggleButtonGroup>
    </MapToggleContainer>
  );
};

MapButtonGroup.defaultProps = {
  initiallySelectedButton: 0,
};

export default MapButtonGroup;

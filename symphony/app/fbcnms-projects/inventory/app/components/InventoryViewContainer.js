/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ViewContainerProps} from '@fbcnms/ui/components/design-system/View/ViewContainer';

import * as React from 'react';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapIcon from '@material-ui/icons/Map';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import {useState} from 'react';

export const DisplayOptions = {
  table: 'table',
  map: 'map',
};
export type DisplayOptionTypes = $Keys<typeof DisplayOptions>;

type ViewToggleProps = {
  onViewToggleClicked?: (id: string) => void,
};

type Props = ViewContainerProps & ViewToggleProps;

const InventoryView = (props: Props) => {
  const {onViewToggleClicked, ...restProps} = props;
  const viewProps: ViewContainerProps = {
    ...restProps,
  };
  const [selectedDisplayOption, setSelectedDisplayOption] = useState(
    DisplayOptions.table,
  );
  if (viewProps.header && onViewToggleClicked) {
    const onViewOptionClicked = displayOptionId => {
      setSelectedDisplayOption(displayOptionId);
      if (onViewToggleClicked) {
        onViewToggleClicked(displayOptionId);
      }
    };
    viewProps.header.viewOptions = {
      onItemClicked: onViewOptionClicked,
      selectedButtonId: selectedDisplayOption,
      buttons: [
        {
          item: <ListAltIcon />,
          id: DisplayOptions.table,
        },
        {
          item: <MapIcon />,
          id: DisplayOptions.map,
        },
      ],
    };
  }
  return <ViewContainer {...viewProps} />;
};

export default InventoryView;

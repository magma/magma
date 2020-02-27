/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  FullViewHeaderProps,
  ViewHeaderActionsProps,
  ViewHeaderProps,
} from '@fbcnms/ui/components/design-system/View/ViewHeader';

import * as React from 'react';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapIcon from '@material-ui/icons/Map';
import ViewHeader from '@fbcnms/ui/components/design-system/View/ViewHeader';
import {useState} from 'react';

export const DisplayOptions = {
  table: 'table',
  map: 'map',
};
export type DisplayOptionTypes = $Keys<typeof DisplayOptions>;

type ViewToggleProps = {
  onViewToggleClicked?: (id: string) => void,
};

type Props = ViewHeaderProps & ViewHeaderActionsProps & ViewToggleProps;

const InventoryViewHeader = (props: Props) => {
  const {onViewToggleClicked, ...restProps} = props;
  const viewHeaderProps: FullViewHeaderProps = {
    ...restProps,
  };
  const [selectedDisplayOption, setSelectedDisplayOption] = useState(
    DisplayOptions.table,
  );
  const onViewOptionClicked = displayOptionId => {
    setSelectedDisplayOption(displayOptionId);
    if (onViewToggleClicked) {
      onViewToggleClicked(displayOptionId);
    }
  };
  viewHeaderProps.viewOptions = {
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
  return <ViewHeader {...viewHeaderProps} />;
};

export default InventoryViewHeader;

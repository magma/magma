/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import React from 'react';

type Props = {
  onClick: (id: string) => void,
};

const WorkOrderMapButtons = (props: Props) => {
  const {onClick} = props;
  return (
    <MapButtonGroup
      onIconClicked={onClick}
      buttons={[
        {
          item: 'Status',
          id: 'status',
        },
        {
          item: 'Technician',
          id: 'technician',
        },
      ]}
    />
  );
};
export default WorkOrderMapButtons;

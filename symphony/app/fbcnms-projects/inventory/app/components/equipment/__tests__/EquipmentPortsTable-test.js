/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import 'jest-dom/extend-expect';
import EquipmentPortsTable from '../EquipmentPortsTable';
import React from 'react';
import TestWrapper from '../../../common/TestWrapper';
import emptyFunction from '@fbcnms/util/emptyFunction';
import nullthrows from '@fbcnms/util/nullthrows';
import shortid from 'shortid';
import {cleanup, render, wait} from '@testing-library/react';

jest.mock('react-relay', () => ({
  createFragmentContainer: component => component,
}));

afterEach(cleanup);

// eslint-disable-next-line flowtype/no-weak-types
const createMockEquipment = (name: string): Object => {
  const portDefinitions = [
    {
      id: shortid.generate(),
      name: 'Port1',
    },
    {
      id: shortid.generate(),
      name: 'Port2',
    },
  ];
  const positionDefinitions = [
    {
      id: shortid.generate(),
      name: 'Position 1',
    },
    {
      id: shortid.generate(),
      name: 'Position 2',
    },
  ];
  return {
    id: shortid.generate(),
    name,
    equipmentType: {
      id: shortid.generate(),
      name: 'Eq Type',
      portDefinitions,
      positionDefinitions,
      propertyTypes: [],
    },
    ports: [
      {
        id: shortid.generate(),
        definition: portDefinitions[0],
      },
    ],
    positions: [
      {
        attachedEquipment: null,
        id: shortid.generate(),
        definition: positionDefinitions[0],
      },
    ],
  };
};

describe('<EquipmentPortsTable />', () => {
  /**
   * 3 levels of ports hierarchy:
   * equipment_1 - 2 ports, 1 position ->
   *  (on the position) equipment_2 - 2 ports, 1 position ->
   *    (on the position) equipment_3 - 2 port, 1 position
   */
  const rootEq = createMockEquipment('Root Eq');
  const childEq = createMockEquipment('Child Eq');
  const grandChildEq = createMockEquipment('Grandchild Eq');

  grandChildEq.positions[0].attachedEquipment = null;
  childEq.positions[0].attachedEquipment = grandChildEq;
  rootEq.positions[0].attachedEquipment = childEq;

  it('renders Edit', async () => {
    const {container} = render(
      <TestWrapper>
        <EquipmentPortsTable
          equipment={rootEq}
          workOrderId={emptyFunction}
          onPortEquipmentClicked={emptyFunction}
          onParentLocationClicked={emptyFunction}
          onWorkOrderSelected={emptyFunction}
        />
      </TestWrapper>,
    );

    expect(nullthrows(container.querySelector('tbody')).children).toHaveLength(
      6,
    );

    await wait();
  });
});

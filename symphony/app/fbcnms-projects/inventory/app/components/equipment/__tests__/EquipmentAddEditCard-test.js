/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('../../../common/RelayEnvironment');

import 'jest-dom/extend-expect';
import EquipmentAddEditCard from '../EquipmentAddEditCard';
import React from 'react';
import RelayEnvironment from '../../../common/RelayEnvironment';
import TestWrapper from '../../../common/TestWrapper';
import {MockPayloadGenerator} from 'relay-test-utils';

import {act, cleanup, fireEvent, render, wait} from '@testing-library/react';

afterEach(cleanup);

describe('<EquipmentAddEditCard />', () => {
  it('renders Edit', async () => {
    const {getByText, getByDisplayValue} = render(
      <TestWrapper>
        <EquipmentAddEditCard
          editingEquipmentId={'36b58c43-b9a8-dd75-c4ff-c3837a541b98'}
          locationId={null}
          equipmentPosition={null}
          workOrderId={null}
          type={null}
          onCancel={() => {}}
          onSave={() => {}}
        />
      </TestWrapper>,
    );

    act(() => {
      // $FlowFixMe (T62907961) Relay flow types
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, {
          EquipmentType() {
            return {
              name: 'Nexus 3048',
            };
          },
          Equipment() {
            return {
              name: 'test_equipment_1',
            };
          },
        }),
      );
    });

    await wait();

    expect(getByText('Nexus 3048')).toBeInTheDocument();
    expect(getByDisplayValue('test_equipment_1')).toBeInTheDocument();
  });

  it('renders Add', async () => {
    const {getByText} = render(
      <TestWrapper>
        <EquipmentAddEditCard
          editingEquipmentId={null}
          locationId={null}
          equipmentPosition={null}
          workOrderId={null}
          // $FlowFixMe for test only
          type={{id: 'nexus', name: 'Nexus 3048'}}
          onCancel={() => {}}
          onSave={() => {}}
        />
      </TestWrapper>,
    );

    act(() => {
      // $FlowFixMe (T62907961) Relay flow types
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, {
          EquipmentType() {
            return {
              name: 'Nexus 3048',
            };
          },
        }),
      );
    });

    await wait(() => {
      expect(getByText('Nexus 3048')).toBeInTheDocument();
    });

    act(() => {
      fireEvent.click(getByText('Cancel'));
    });
  });
});

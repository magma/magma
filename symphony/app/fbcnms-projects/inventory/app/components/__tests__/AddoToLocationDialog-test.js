/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('../../common/RelayEnvironment');

import 'jest-dom/extend-expect';
import AddToLocationDialog from '../AddToLocationDialog';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import defaultTheme from '@fbcnms/ui/theme/default';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {MockPayloadGenerator} from 'relay-test-utils';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {act, cleanup, fireEvent, render, wait} from '@testing-library/react';

global.CONFIG = {
  appData: {enabledFeatures: []},
};

const Wrapper = props => (
  <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <AppContextProvider>{props.children}</AppContextProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<AddToLocationDialog />', () => {
  describe('location', () => {
    it('renders', async () => {
      const {getByText, queryByText} = render(
        <Wrapper>
          <AddToLocationDialog
            open={true}
            show="location"
            onClose={() => {}}
            onEquipmentTypeSelected={() => {}}
            onLocationTypeSelected={() => {}}
          />
        </Wrapper>,
      );

      expect(getByText('Add')).toBeInTheDocument();
      expect(getByText('Select a location type')).toBeInTheDocument();

      expect(queryByText('Upload Exported Service')).not.toBeInTheDocument();
      act(() => {
        fireEvent.click(getByText('Bulk Upload'));
      });
      expect(getByText('Upload Exported Service')).toBeInTheDocument();
    });

    it('saves and cancels', async () => {
      const onSave = jest.fn(() => {});
      const onCancel = jest.fn(() => {});

      const {getByText} = render(
        <Wrapper>
          <AddToLocationDialog
            open={true}
            show="location"
            onClose={onCancel}
            onEquipmentTypeSelected={() => {}}
            onLocationTypeSelected={onSave}
          />
        </Wrapper>,
      );

      act(() => {
        RelayEnvironment.mock.resolveMostRecentOperation(operation =>
          MockPayloadGenerator.generate(operation, {
            LocationType() {
              return {
                name: 'Building',
              };
            },
          }),
        );
      });

      act(() => {
        fireEvent.click(getByText('Add'));
      });

      await wait(() => {
        expect(onSave).not.toBeCalled();
      });

      act(() => {
        fireEvent.click(getByText('Building'));
      });
      act(() => {
        fireEvent.click(getByText('Add'));
      });

      expect(onSave).toBeCalled();

      expect(onCancel).not.toBeCalled();
      fireEvent.click(getByText('Cancel'));
      expect(onCancel).toBeCalled();
    });
  });

  describe('equipment', () => {
    it('renders', async () => {
      const {getByText} = render(
        <Wrapper>
          <AddToLocationDialog
            open={true}
            show="equipment"
            onClose={() => {}}
            onEquipmentTypeSelected={() => {}}
            onLocationTypeSelected={() => {}}
          />
        </Wrapper>,
      );

      act(() => {
        RelayEnvironment.mock.resolveMostRecentOperation(operation =>
          MockPayloadGenerator.generate(operation),
        );
      });

      await wait(() => {
        expect(getByText('Select an equipment type')).toBeInTheDocument();
      });
    });

    it('saves and cancels', async () => {
      const onSave = jest.fn(() => {});
      const onCancel = jest.fn(() => {});

      const {getByText} = render(
        <Wrapper>
          <AddToLocationDialog
            open={true}
            show="equipment"
            onClose={onCancel}
            onEquipmentTypeSelected={onSave}
            onLocationTypeSelected={() => {}}
          />
        </Wrapper>,
      );

      act(() => {
        RelayEnvironment.mock.resolveMostRecentOperation(operation =>
          MockPayloadGenerator.generate(operation, {
            EquipmentType() {
              return {
                name: 'Ubiquiti NanoBeam M5',
              };
            },
          }),
        );
      });

      act(() => {
        fireEvent.click(getByText('Add'));
      });

      await wait(() => {
        expect(onSave).not.toBeCalled();
      });

      act(() => {
        fireEvent.click(getByText('Ubiquiti NanoBeam M5'));
      });
      act(() => {
        fireEvent.click(getByText('Add'));
      });
      expect(onSave).toBeCalled();

      act(() => {
        fireEvent.click(getByText('Cancel'));
      });
      expect(onCancel).toBeCalled();
    });
  });

  describe('upload', () => {
    it('renders', async () => {
      const {getByText} = render(
        <Wrapper>
          <AddToLocationDialog
            open={true}
            show="upload"
            onClose={() => {}}
            onEquipmentTypeSelected={() => {}}
            onLocationTypeSelected={() => {}}
          />
        </Wrapper>,
      );

      expect(getByText('Upload Exported Equipment')).toBeInTheDocument();
      expect(getByText('Upload Exported Ports')).toBeInTheDocument();
      expect(getByText('Upload Exported Links')).toBeInTheDocument();
      expect(getByText('Upload Locations')).toBeInTheDocument();
      expect(getByText('Upload Exported Service')).toBeInTheDocument();
    });
  });
});

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

import LocationsTree from '../LocationsTree';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import {MemoryRouter} from 'react-router-dom';
import {MockPayloadGenerator} from 'relay-test-utils';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';

import 'jest-dom/extend-expect';
import defaultTheme from '@fbcnms/ui/theme/default';

import {
  act,
  cleanup,
  fireEvent,
  render,
  waitForElement,
} from '@testing-library/react';

const Wrapper = props => (
  <MemoryRouter initialEntries={['/inventory']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>{props.children}</SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

const MOCK_RESOLVER = {
  Location(ctx) {
    switch (ctx.path.join('.')) {
      case 'location':
      case 'locations.edges.node':
        // root location
        return {
          id: 'usa',
          externalId: 'USA',
          name: 'United States',
          numChildren: 1,
        };
      case 'location.children':
        // children from root
        return {
          id: 'california',
          externalId: 'cali',
          name: 'California',
          numChildren: 0,
          siteSurveyNeeded: false,
        };
      default:
        throw new Error('Invalid ctx path for Location');
    }
  },
  LocationType(ctx) {
    switch (ctx.path.join('.')) {
      case 'location.locationType':
      case 'locations.edges.node.locationType':
        // root location
        return {
          id: 'country',
          name: 'Country',
        };
      case 'location.children.locationType':
        // children
        return {
          id: 'state',
          name: 'State',
        };
      default:
        throw new Error('Invalid ctx path for LocationType');
    }
  },
};

describe('<LocationsTree />', () => {
  it('renders', async () => {
    const {getByText} = render(
      <Wrapper>
        <LocationsTree
          selectedLocationId={null}
          onAddLocation={() => {}}
          onSelect={() => {}}
        />
      </Wrapper>,
    );

    act(() => {
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
      );
    });

    expect(getByText('Locations')).toBeInTheDocument();
    expect(getByText('Add top-level location')).toBeInTheDocument();
  });

  it('renders tree with locations', async () => {
    const {getByText} = render(
      <Wrapper>
        <LocationsTree
          selectedLocationId={null}
          onAddLocation={() => {}}
          onSelect={() => {}}
        />
      </Wrapper>,
    );

    act(() => {
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
      );
    });

    expect(getByText('United States')).toBeInTheDocument();
  });

  it('handles selecting location', async () => {
    const onSelect = jest.fn(() => {});

    const {getByText} = render(
      <Wrapper>
        <LocationsTree
          selectedLocationId={null}
          onAddLocation={() => {}}
          onSelect={onSelect}
        />
      </Wrapper>,
    );

    act(() => {
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
      );
    });

    act(() => {
      fireEvent.click(getByText('United States'));
    });

    expect(onSelect.mock.calls.length).toBe(1);
  });

  it('handles selected location change', async () => {
    const {getByText, getByTestId} = render(
      <Wrapper>
        <LocationsTree
          selectedLocationId={null}
          onAddLocation={() => {}}
          onSelect={() => {}}
        />
      </Wrapper>,
    );

    act(() => {
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
      );
    });

    expect(getByText('United States')).toBeInTheDocument();

    act(() => {
      fireEvent.click(getByTestId('inventory-expand-usa'));
    });

    act(() => {
      RelayEnvironment.mock.resolveMostRecentOperation(operation =>
        MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
      );
    });

    const loadedCalifornia = await waitForElement(() =>
      getByText('California'),
    );

    expect(getByText('United States')).toBeInTheDocument();
    expect(loadedCalifornia).toBeInTheDocument();
  });
});

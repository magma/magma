/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import DataGrid from '../DataGrid';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../theme/default';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render} from '@testing-library/react';
import type {DataRows} from '../DataGrid';

const data: DataRows[] = [
  [
    {
      category: 'Total',
      value: 'eNodeBs',
      tooltip: 'Tooltip text',
    },
    {
      category: 'Severe Events',
      value: 'Value used as a tooltip',
    },
    {
      category: 'Max Latency',
      value: 100,
      unit: 'ms',
    },
  ],
];

const Wrapper = () => {
  return (
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <DataGrid data={data} />
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  );
};

describe('<DataGrid />', () => {
  it('displays the passed tooltip', async () => {
    const {getByText} = render(<Wrapper />);

    const dataEntryElement = getByText('eNodeBs');
    expect(dataEntryElement).toHaveAttribute('title', 'Tooltip text');
  });

  it('defaults to the data entry value when the tooltip prop in not passed', async () => {
    const {getByText} = render(<Wrapper />);

    const dataEntryElement = getByText('Value used as a tooltip');
    expect(dataEntryElement).toHaveAttribute(
      'title',
      'Value used as a tooltip',
    );
  });

  it('displays the data unit along with data value as the tooltip when unit prop is passed', async () => {
    const {getByText} = render(<Wrapper />);

    const dataEntryElement = getByText('100ms');
    expect(dataEntryElement).toHaveAttribute('title', '100ms');
  });
});

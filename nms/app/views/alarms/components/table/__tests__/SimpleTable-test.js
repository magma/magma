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

import * as React from 'react';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SimpleTable, {LabelsCell} from '../SimpleTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../../../theme/default';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render} from '@testing-library/react';

// replace the default chip with a more easily queryable version
jest.mock('@material-ui/core/Chip', () => ({label, ...props}) => (
  <div data-chip {...props} children={label} />
));

function Wrapper(props: {route?: string, children: React.Node}) {
  return (
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        {props.children}
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  );
}

test('renders with required default props', () => {
  const {getByText} = render(
    <Wrapper>
      <SimpleTable columnStruct={mockColumns()} tableData={[]} />
    </Wrapper>,
  );
  expect(getByText('name')).toBeInTheDocument();
  expect(getByText('age')).toBeInTheDocument();
});

function mockColumns() {
  return [
    {title: 'name', field: 'name'},
    {title: 'age', field: 'age'},
  ];
}

test('rendered row is transformed by path expression', () => {
  const rows = [
    {
      name: 'bob',
      labels: {
        description: 'bob description',
      },
    },
    {
      name: 'mary',
      labels: {
        description: 'mary description',
      },
    },
  ];

  const {getByText, ..._result} = render(
    <Wrapper>
      <SimpleTable
        columnStruct={[
          {title: 'name', field: 'name'},
          {title: 'description', field: 'labels.description'},
        ]}
        tableData={rows}
      />
    </Wrapper>,
  );

  expect(getByText('name')).toBeInTheDocument();
  expect(getByText('description')).toBeInTheDocument();
  expect(getByText('bob')).toBeInTheDocument();
  expect(getByText('bob description')).toBeInTheDocument();
  expect(getByText('mary')).toBeInTheDocument();
  expect(getByText('mary description')).toBeInTheDocument();
});

test('if menuItems is passed, actions menu is rendered', () => {
  const {getAllByTitle} = render(
    <Wrapper>
      <SimpleTable
        columnStruct={mockColumns()}
        tableData={[{name: 'name', age: 'age'}]}
        menuItems={[
          {
            name: 'View',
          },
        ]}
      />
    </Wrapper>,
  );
  expect(getAllByTitle('Actions')[0]).toBeInTheDocument();
});

describe('column renderers', () => {
  test('if cell value is an object, renders label chips', () => {
    const {container} = render(
      <Wrapper>
        <SimpleTable
          columnStruct={[
            {
              title: 'labels',
              field: 'labels',
              render: row => <LabelsCell value={row.labels} />,
            },
          ]}
          tableData={[{labels: {name: 'name', age: 'age'}}]}
        />
      </Wrapper>,
    );
    /**
     * Replace the default material-ui chip with one which passes an
     * easily queryable identifier to ensure that a chip is rendered.
     */
    const chips = container.querySelectorAll('[data-chip]');
    // ensure that 2 chips are rendered
    expect(chips.length).toBe(2);
    /**
     * chip text is broken up by multiple elements. textContent combines text
     * from all children so we can check for that instead.
     */
    const textContent = [].map.call(chips, chip => chip.textContent);
    expect(textContent).toContain('name=name');
    expect(textContent).toContain('age=age');
  });
});

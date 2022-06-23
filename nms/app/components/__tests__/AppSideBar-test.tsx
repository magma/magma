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
 */

import AppSideBar from '../AppSideBar';
import DashboardIcon from '@material-ui/icons/Dashboard';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';
import {MemoryRouter} from 'react-router-dom';
import {fireEvent, render} from '@testing-library/react';
import type {MemoryRouterProps} from 'react-router-dom';

const Wrapper = (props: MemoryRouterProps) => {
  return (
    <MemoryRouter initialEntries={['/nms']} initialIndex={0}>
      {props.children}
    </MemoryRouter>
  );
};

describe('AppSideBar', () => {
  it('renders without items', () => {
    const {queryAllByRole} = render(
      <Wrapper>
        <AppSideBar items={[]} />
      </Wrapper>,
    );

    expect(queryAllByRole('link')).toHaveLength(0);
  });

  const items = [
    {
      path: 'dashboard',
      label: 'Dashboard',
      icon: <DashboardIcon />,
    },
    {
      path: 'equipment',
      label: 'Equipment',
      icon: <RouterIcon />,
    },
  ];

  it('renders items as links', () => {
    const {queryAllByRole} = render(
      <Wrapper>
        <AppSideBar items={items} />
      </Wrapper>,
    );

    expect(queryAllByRole('link')).toHaveLength(2);
  });

  it('expands on hover', () => {
    const {queryAllByRole, getByTestId, queryByText} = render(
      <Wrapper>
        <AppSideBar items={items} />
      </Wrapper>,
    );
    const links = queryAllByRole('link');

    fireEvent.mouseOver(links[0]);
    expect(links.map(link => link.textContent)).toEqual([
      'Dashboard',
      'Equipment',
    ]);
    expect(queryByText('Account & Settings')).toBeInTheDocument();

    fireEvent.mouseLeave(getByTestId('app-sidebar'));
    expect(links.map(link => link.textContent)).toEqual(['', '']);
    expect(queryByText('Account & Settings')).not.toBeInTheDocument();
  });
});

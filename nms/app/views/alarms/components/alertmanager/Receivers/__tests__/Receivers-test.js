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
 * @flow
 * @format
 */

import * as React from 'react';
import Receivers from '../Receivers';
import {
  act,
  fireEvent,
  render,
  wait,
  waitForElement,
} from '@testing-library/react';
import {alarmTestUtil} from '../../../../test/testHelpers';
import type {AlarmsWrapperProps} from '../../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ApiUtil} from '../../../AlarmsApi';

describe('Receivers', () => {
  let AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>;
  let apiUtil: ApiUtil;

  beforeEach(() => {
    ({AlarmsWrapper, apiUtil} = alarmTestUtil());
  });

  it('renders', () => {
    render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );
  });

  it('clicking the View button on a row shows the view dialog', async () => {
    jest.spyOn(apiUtil, 'useAlarmsApi').mockReturnValueOnce({
      response: [
        {
          name: 'test_receiver',
          slack_configs: [
            {
              api_url: 'test.com',
              channel: '#test',
              text: '{{text}}',
              title: '{{title}}',
            },
          ],
        },
      ],
    });
    const {getByText, getAllByText, queryByText, getAllByTitle} = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );
    const actionMenu = getAllByTitle('Actions');
    expect(actionMenu[0]).toBeInTheDocument();
    act(() => {
      fireEvent.click(actionMenu[0]);
    });
    act(() => {
      fireEvent.click(getAllByText('View')[0]);
    });
    // clicking View should open the dialog
    await waitForElement(() => getByText(/View Receiver/i));
    expect(getByText(/View Receiver/i)).toBeInTheDocument();

    // clicking Close should close the dialog
    act(() => {
      fireEvent.click(getByText(/close/i));
    });
    await wait(() => {
      expect(queryByText(/View Receiver/i)).not.toBeInTheDocument();
    });
  });

  it('clicking edit button should show AddEditReceiver in edit mode', () => {
    jest.spyOn(apiUtil, 'useAlarmsApi').mockReturnValueOnce({
      response: [
        {
          name: 'test_receiver',
          slack_configs: [
            {
              api_url: 'test.com',
              channel: '#test',
              text: '{{text}}',
              title: '{{title}}',
            },
          ],
        },
      ],
    });
    const {getAllByText, getByTestId, queryByTestId, getAllByTitle} = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );

    const actionMenu = getAllByTitle('Actions');
    expect(actionMenu[0]).toBeInTheDocument();
    act(() => {
      fireEvent.click(actionMenu[0]);
    });
    expect(queryByTestId('add-edit-receiver')).not.toBeInTheDocument();
    act(() => {
      fireEvent.click(getAllByText('Edit')[0]);
    });
    expect(getByTestId('add-edit-receiver')).toBeInTheDocument();
  });

  it('clicking add button should show AddEditReceiver', () => {
    jest.spyOn(apiUtil, 'useAlarmsApi').mockReturnValueOnce({
      response: [],
    });
    const {getByTestId, queryByTestId} = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );

    expect(queryByTestId('add-edit-receiver')).not.toBeInTheDocument();
    act(() => {
      fireEvent.click(getByTestId('add-receiver-button'));
    });
    expect(getByTestId('add-edit-receiver')).toBeInTheDocument();
  });
});

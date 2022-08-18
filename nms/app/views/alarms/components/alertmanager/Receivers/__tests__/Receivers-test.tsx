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

import * as React from 'react';
import Receivers from '../Receivers';
import {MockApiUtil, alarmTestUtil} from '../../../../test/testHelpers';
import {act, fireEvent, waitFor} from '@testing-library/react';
import {render} from '../../../../../../util/TestingLibrary';
import type {AlarmsWrapperProps} from '../../../../test/testHelpers';

describe('Receivers', () => {
  let AlarmsWrapper: React.ComponentType<Partial<AlarmsWrapperProps>>;
  let apiUtil: MockApiUtil;

  const defaultResponse = {
    error: undefined,
    isLoading: false,
  };

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
      ...defaultResponse,
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
    const {getByText, getAllByText, queryByText, openActionsTableMenu} = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );
    await openActionsTableMenu(0);
    act(() => {
      fireEvent.click(getAllByText('View')[0]);
    });
    // clicking View should open the dialog
    expect(queryByText('View Receiver')).toBeInTheDocument(),
      // clicking Close should close the dialog
      act(() => {
        fireEvent.click(getByText(/close/i));
      });
    await waitFor(() =>
      expect(queryByText('View Receiver')).not.toBeInTheDocument(),
    );
  });

  it('clicking edit button should show AddEditReceiver in edit mode', async () => {
    jest.spyOn(apiUtil, 'useAlarmsApi').mockReturnValueOnce({
      ...defaultResponse,
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
    const {
      getAllByText,
      getByTestId,
      queryByTestId,
      openActionsTableMenu,
    } = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );

    await openActionsTableMenu(0);
    expect(queryByTestId('add-edit-receiver')).not.toBeInTheDocument();
    act(() => {
      fireEvent.click(getAllByText('Edit')[0]);
    });
    expect(getByTestId('add-edit-receiver')).toBeInTheDocument();
  });

  it('clicking add button should show AddEditReceiver', () => {
    jest.spyOn(apiUtil, 'useAlarmsApi').mockReturnValueOnce({
      ...defaultResponse,
      response: [],
    });
    const {getAllByTestId, getByTestId, queryByTestId} = render(
      <AlarmsWrapper>
        <Receivers />
      </AlarmsWrapper>,
    );

    expect(queryByTestId('add-edit-receiver')).not.toBeInTheDocument();
    act(() => {
      fireEvent.click(getAllByTestId('add-receiver-button')[0]);
    });
    expect(getByTestId('add-edit-receiver')).toBeInTheDocument();
  });
});

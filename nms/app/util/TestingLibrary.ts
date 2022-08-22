/**
 * Copyright 2022 The Magma Authors.
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

import {
  RenderOptions,
  fireEvent,
  render as rtlRender,
  waitFor,
} from '@testing-library/react';

export function render(
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'queries'>,
) {
  const result = rtlRender(ui, options);

  async function openActionsTableMenu(index: number) {
    const actionMenuButtons = (
      await result.findAllByRole('button', {name: 'Actions'})
    ).filter(elem => elem instanceof HTMLButtonElement); // exclude header

    const actionsMenu = result.getByTestId('actions-menu');

    expect(actionsMenu).not.toBeVisible();
    fireEvent.click(actionMenuButtons[index]);
    await waitFor(() => expect(actionsMenu).toBeVisible());

    return actionsMenu;
  }

  return {...result, openActionsTableMenu};
}

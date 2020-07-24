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

import fbt from 'fbt';

const Strings = {
  helpers: {
    productNameParam: 'product name',
  },
  common: {
    productName: `${fbt('Symphony', '')}`,
    emptyField: `${fbt(
      'None',
      'Text to be displayed in case a user input field has no value',
    )}`,
    unassignedItem: `${fbt(
      'Unassigned',
      'Text to be displayed in case an assignable item was not assigned yet',
    )}`,
    closeButton: `${fbt('Close', 'Text for button closing message dialog')}`,
    okButton: `${fbt(
      'OK',
      'Text for button approving message or dialog content',
    )}`,
    doneButton: `${fbt('Done', '')}`,
    updateButton: `${fbt('Update', 'Text for button that applies changes')}`,
    saveButton: `${fbt(
      'Save',
      'Text for button that saves current view changes',
    )}`,
    cancelButton: `${fbt(
      'Cancel',
      'Text for button that cancels current operation',
    )}`,
    deleteButton: `${fbt(
      'Delete',
      'Text for button that will cause a delete operation',
    )}`,
    createButton: `${fbt(
      'Create',
      'Text for button that creates a new instance',
    )}`,
    nextButton: `${fbt('Next', 'Text for button that go to next operation')}`,
    backButton: `${fbt(
      'Back',
      'Text for button that go to previous operation',
    )}`,
    addButton: `${fbt('Add', 'Text for button that adds an item')}`,
    removeButton: `${fbt('Remove', 'Text for button that removes an item')}`,
    yesButton: `${fbt('Yes', 'Yes response to a question')}`,
    noButton: `${fbt('No', 'No response to a question')}`,
    fields: {
      url: {
        label: 'URL',
        placeholder: `${fbt(
          'https://example.com/',
          'Example text for URL input field',
        )}`,
      },
    },
  },
  admin: {
    users: {
      viewHeader: `${fbt(
        'User Management',
        'Header for view showing and managing all system user and permissions settings',
      )}`,
    },
  },
};

export default Strings;

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import fbt from 'fbt';

const Strings = {
  common: {
    emptyField: `${fbt(
      'None',
      'Text to be displayed in case a user input field has no value',
    )}`,
    unassignedItem: `${fbt(
      'Unassigned',
      'Text to be displayed in case an assignable item was not assigned yet',
    )}`,
    okButton: `${fbt(
      'OK',
      'Text for button approving message or dialog content',
    )}`,
    cancelButton: `${fbt(
      'Cancel',
      'Text for button that cancels current operation',
    )}`,
    nextButton: `${fbt('Next', 'Text for button that go to next operation')}`,
    backButton: `${fbt(
      'Back',
      'Text for button that go to previous operation',
    )}`,
    addButton: `${fbt('Add', 'Text for button that adds an item')}`,
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
  documents: {
    viewHeader: `${fbt(
      'Documents',
      'Header text for a view showing documents',
    )}`,
    uploadButton: `${fbt('Upload File', 'Upload files button caption')}`,
    addLinkButton: `${fbt('Add URL', 'Open Add URL dialog button caption')}`,
    categories: [
      'Archivos de Estudios Pre-instalación',
      'Archivos de Contratos',
      'Archivos de TSS',
      'DataFills',
      'ATP',
      'Topología',
      'Archivos Simulación',
      'Reportes de Mantenimiento',
      'Fotos',
    ],
  },
};

export default Strings;

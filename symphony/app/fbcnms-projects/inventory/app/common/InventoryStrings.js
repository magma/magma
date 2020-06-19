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
      'Site Folder',
    ],
  },
};

export default Strings;

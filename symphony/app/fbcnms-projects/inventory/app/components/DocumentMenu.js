/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FileAttachmentType} from '../common/FileAttachment';

import OptionsPopoverButton from './OptionsPopoverButton';
import React from 'react';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../common/DocumentAPI';

type Props = {
  document: FileAttachmentType,
  onDocumentDeleted: (document: FileAttachmentType) => void,
  onDialogOpen: () => void,
  popoverMenuClassName?: ?string,
  onVisibilityChange?: (isVisible: boolean) => void,
};

class DocumentMenu extends React.Component<Props> {
  downloadFileRef: {
    current: null | HTMLAnchorElement,
  } = React.createRef<HTMLAnchorElement>();

  handleDownload = () => {
    if (this.downloadFileRef.current != null) {
      this.downloadFileRef.current.click();
    }
  };

  handlePreview = () => {
    this.props.onDialogOpen();
  };

  handleDelete = () => {
    this.props.onDocumentDeleted(this.props.document);
  };

  render() {
    const {document, popoverMenuClassName, onVisibilityChange} = this.props;
    const storeKey = nullthrows(document.storeKey);
    const menuOptions = [
      {
        onClick: this.handlePreview.bind(this),
        caption: fbt(
          'Preview',
          'Caption for menu option for showing image in preview mode',
        ),
        ignorePermissions: true,
      },
      {
        onClick: this.handleDelete.bind(this),
        caption: fbt(
          'Delete',
          'Caption for menu option for deleting a file from files table',
        ),
      },
      {
        onClick: this.handleDownload.bind(this),
        caption: fbt(
          'Download',
          'Caption for menu option for downloading a file from files table',
        ),
        ignorePermissions: true,
      },
    ];
    return (
      <>
        <OptionsPopoverButton
          options={menuOptions}
          popoverMenuClassName={popoverMenuClassName}
          onVisibilityChange={onVisibilityChange}
        />
        <a
          href={DocumentAPIUrls.download_url(storeKey, document.fileName)}
          ref={this.downloadFileRef}
          style={{display: 'none'}}
          download
        />
      </>
    );
  }
}

export default DocumentMenu;

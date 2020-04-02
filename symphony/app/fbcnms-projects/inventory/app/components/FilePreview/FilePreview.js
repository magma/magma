/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FileAttachmentType} from '../../common/FileAttachment';

import * as React from 'react';
import DocumentMenu from '../DocumentMenu';
import ImageDialog from '../ImageDialog';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../../common/DocumentAPI';
import {
  WIDE_DIMENSION_HEIGHT_PX,
  WIDE_DIMENSION_WIDTH_PX,
} from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    height: WIDE_DIMENSION_HEIGHT_PX,
    width: WIDE_DIMENSION_WIDTH_PX,
    position: 'relative',
    borderRadius: '4px',
    '&:hover $popoverMenu, $popoverMenu$visiblePopoverMenu': {
      visibility: 'visible',
    },
  },
  image: {
    objectFit: 'cover',
    objectPosition: '50% 50%',
    height: '100%',
    width: '100%',
    borderRadius: '4px',
  },
  overlay: {
    background: 'linear-gradient(to bottom, rgba(255, 255, 255, 0), black)',
    position: 'absolute',
    left: 0,
    right: 0,
    top: 0,
    bottom: 0,
    borderRadius: '4px',
    zIndex: 1,
  },
  name: {
    position: 'absolute',
    left: 8,
    right: 8,
    bottom: 8,
    zIndex: 2,
  },
  popoverMenu: {
    position: 'absolute',
    right: 4,
    top: 4,
    zIndex: 2,
    visibility: 'hidden',
  },
  visiblePopoverMenu: {},
  moreIcon: {
    padding: '4px',
    backgroundColor: 'white',
    borderRadius: '100%',
    cursor: 'pointer',
  },
}));

type Props = {
  file: FileAttachmentType,
  onFileDeleted: (file: FileAttachmentType) => void,
  className?: string,
};

const FilePreview = ({file, onFileDeleted, className}: Props): React.Node => {
  const classes = useStyles();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isPreviewDialogOpen, setIsPreviewDialogOpen] = useState(false);

  return (
    <div className={classNames(classes.root, className)}>
      <img
        className={classes.image}
        src={DocumentAPIUrls.get_url(nullthrows(file.storeKey))}
      />
      <div className={classes.name}>
        <Text
          variant="caption"
          color="light"
          weight="medium"
          useEllipsis={true}>
          {file.fileName}
        </Text>
      </div>
      <div className={classes.overlay} />
      <DocumentMenu
        document={file}
        onDialogOpen={() => setIsPreviewDialogOpen(true)}
        onDocumentDeleted={() => onFileDeleted(file)}
        popoverMenuClassName={classNames(classes.popoverMenu, {
          [classes.visiblePopoverMenu]: isMenuOpen,
        })}
        onVisibilityChange={isVisible => setIsMenuOpen(isVisible)}
      />
      <ImageDialog
        onClose={() => setIsPreviewDialogOpen(false)}
        open={isPreviewDialogOpen}
        img={file}
      />
    </div>
  );
};

export default FilePreview;

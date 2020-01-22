/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {FileAttachment_file} from './__generated__/FileAttachment_file.graphql';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import DateTimeFormat from '../common/DateTimeFormat.js';
import DocumentMenu from './DocumentMenu';
import ImageDialog from './ImageDialog';
import InsertDriveFileIcon from '@material-ui/icons/InsertDriveFile';
import React from 'react';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import nullthrows from 'nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {formatFileSize} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
  nameCell: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
  thumbnail: {
    marginRight: '20px',
    display: 'flex',
    alignItems: 'center',
  },
  icon: {
    fontSize: '24px',
    width: '32px',
  },
  img: {
    width: '32px',
    height: '32px',
  },
  fileName: {
    ...symphony.typography.caption,
  },
  secondaryCell: {
    color: symphony.palette.D400,
  },
  cell: {
    height: '48px',
    ...symphony.typography.caption,
  },
  secondaryCell: {
    color: symphony.palette.D400,
  },
  moreIcon: {
    fill: symphony.palette.D400,
  },
});

type Props = {
  file: FileAttachment_file,
  onDocumentDeleted: (file: FileAttachment_file) => void,
} & WithStyles<typeof styles>;

type State = {
  isImageDialogOpen: boolean,
};

class FileAttachment extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;

  downloadFileRef: {
    current: null | HTMLAnchorElement,
  } = React.createRef<HTMLAnchorElement>();

  constructor(props: Props) {
    super(props);
    this.state = {
      isImageDialogOpen: false,
    };
  }

  handleDownload = () => {
    if (this.downloadFileRef.current != null) {
      this.downloadFileRef.current.click();
    }
  };

  handleDelete = async () => {
    this.props.onDocumentDeleted(this.props.file);
  };

  render() {
    const {classes, file} = this.props;
    if (file === null) {
      return null;
    }

    const categoriesEnabled = this.context.isFeatureEnabled('file_categories');

    return (
      <TableRow key={file.id} hover={false}>
        {categoriesEnabled && (
          <TableCell padding="none" component="th" scope="row">
            {file.category}
          </TableCell>
        )}
        <TableCell padding="none" component="th" scope="row">
          <div className={classes.nameCell}>
            <div className={classes.thumbnail}>
              {file.fileType === 'IMAGE' ? (
                <img
                  className={classes.img}
                  src={DocumentAPIUrls.get_url(nullthrows(file.storeKey))}
                />
              ) : (
                <InsertDriveFileIcon color="primary" className={classes.icon} />
              )}
            </div>
            <Text className={classes.fileName}>{file.fileName}</Text>
          </div>
        </TableCell>
        <TableCell
          padding="none"
          className={classNames(classes.cell, classes.secondaryCell)}
          component="th"
          scope="row">
          {file.fileName
            .split('.')
            .pop()
            .toUpperCase()}
        </TableCell>
        <TableCell
          padding="none"
          className={classNames(classes.cell, classes.secondaryCell)}
          component="th"
          scope="row">
          {file.sizeInBytes != null && formatFileSize(file.sizeInBytes)}
        </TableCell>
        <TableCell
          padding="none"
          className={classNames(classes.cell, classes.secondaryCell)}
          component="th"
          scope="row">
          {file.uploaded && DateTimeFormat.dateOnly(file.uploaded)}
        </TableCell>
        <TableCell
          padding="none"
          className={classNames(classes.cell, classes.secondaryCell)}
          component="th"
          scope="row"
          align="right">
          <DocumentMenu
            document={file}
            onDocumentDeleted={this.handleDelete}
            onDialogOpen={() => this.setState({isImageDialogOpen: true})}
          />
          {file.fileType === 'IMAGE' && (
            <ImageDialog
              onClose={() => this.setState({isImageDialogOpen: false})}
              open={this.state.isImageDialogOpen}
              img={file}
            />
          )}
        </TableCell>
      </TableRow>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(FileAttachment, {
    file: graphql`
      fragment FileAttachment_file on File {
        id
        fileName
        sizeInBytes
        uploaded
        fileType
        storeKey
        category
        ...DocumentMenu_document
        ...ImageDialog_img
      }
    `,
  }),
);

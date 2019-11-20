/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DocumentMenu_document} from './__generated__/DocumentMenu_document.graphql';
import type {WithStyles} from '@material-ui/core';

import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
  icon: {
    color: symphony.palette.D400,
    cursor: 'pointer',
  },
});

type Props = {
  document: DocumentMenu_document,
  onDocumentDeleted: (document: DocumentMenu_document) => void,
  onDialogOpen: () => void,
} & WithStyles<typeof styles>;

type State = {
  anchorElement: ?HTMLElement,
};

class DocumentMenu extends React.Component<Props, State> {
  state = {
    anchorElement: null,
  };

  downloadFileRef: {
    current: null | HTMLAnchorElement,
  } = React.createRef<HTMLAnchorElement>();

  handleClick = event => {
    this.setState({anchorElement: event.currentTarget});
  };

  hide = () => {
    this.setState({anchorElement: null});
  };

  handleDownload = () => {
    if (this.downloadFileRef.current != null) {
      this.downloadFileRef.current.click();
    }
    this.hide();
  };

  handlePreview = async _ => {
    this.props.onDialogOpen();
    this.hide();
  };

  handleDelete = async _ => {
    this.props.onDocumentDeleted(this.props.document);
    this.hide();
  };

  handleClose = () => {
    this.hide();
  };

  render() {
    const {classes, document} = this.props;
    const {anchorElement} = this.state;
    const storeKey = nullthrows(document.storeKey);
    return (
      <div>
        <MoreVertIcon className={classes.icon} onClick={this.handleClick} />
        <a
          href={DocumentAPIUrls.download_url(storeKey, document.fileName)}
          ref={this.downloadFileRef}
          style={{display: 'none'}}
          download
        />
        <Menu
          id="simple-menu"
          anchorEl={anchorElement}
          open={!!anchorElement}
          onClose={this.handleClose}>
          {document.fileType === 'IMAGE' && (
            <MenuItem onClick={this.handlePreview}>Preview</MenuItem>
          )}
          <MenuItem onClick={this.handleDelete}>Delete</MenuItem>
          <MenuItem onClick={this.handleDownload}>Download</MenuItem>
        </Menu>
      </div>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(DocumentMenu, {
    document: graphql`
      fragment DocumentMenu_document on File {
        id
        fileName
        storeKey
        fileType
      }
    `,
  }),
);

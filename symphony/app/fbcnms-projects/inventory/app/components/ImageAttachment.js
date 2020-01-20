/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {ImageAttachment_img} from './__generated__/ImageAttachment_img.graphql';
import type {WithStyles} from '@material-ui/core';

import DocumentMenu from './DocumentMenu';
import ImageDialog from './ImageDialog';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  img: ImageAttachment_img,
  onImageDeleted: (img: ImageAttachment_img) => void,
} & WithStyles<typeof styles>;

type State = {
  hovered: boolean,
  dialogOpen: boolean,
};

const styles = _ => ({
  root: {
    cursor: 'pointer',
    '&:hover': {
      '& $menu': {
        display: 'block',
      },
    },
  },
  img: {
    width: '100%',
    maxHeight: '100%',
  },
  menu: {
    display: 'none',
    position: 'absolute',
    top: '2px',
    right: '2px',
  },
});

class ImageAttachment extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {hovered: false, dialogOpen: false};
  }
  onMouseHover = () => {
    this.setState({hovered: true});
  };
  onMouseLeave = () => {
    this.setState({hovered: false});
  };
  onClick = () => {
    this.setState({dialogOpen: true});
  };
  onDialogClose = () => {
    this.setState({dialogOpen: false});
  };
  onDialogOpen = () => {
    this.setState({dialogOpen: true});
  };
  render() {
    const {classes} = this.props;
    return (
      <div className={this.props.classes.root}>
        <img
          className={this.props.classes.img}
          src={DocumentAPIUrls.get_url(nullthrows(this.props.img.storeKey))}
          onClick={this.onClick}
        />
        <div className={classes.menu}>
          <DocumentMenu
            document={this.props.img}
            onDocumentDeleted={this.props.onImageDeleted}
            onDialogOpen={this.onDialogOpen}
          />
          <ImageDialog
            onClose={this.onDialogClose}
            open={this.state.dialogOpen}
            img={this.props.img}
          />
        </div>
      </div>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(ImageAttachment, {
    img: graphql`
      fragment ImageAttachment_img on File @relay(mask: false) {
        id
        storeKey
        ...DocumentMenu_document
        ...ImageDialog_img
      }
    `,
  }),
);

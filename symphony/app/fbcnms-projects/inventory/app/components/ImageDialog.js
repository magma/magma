/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ImageDialog_img} from './__generated__/ImageDialog_img.graphql';
import type {WithStyles} from '@material-ui/core';

import CloseIcon from '@material-ui/icons/Close';
import Dialog from '@material-ui/core/Dialog';
import IconButton from '@material-ui/core/IconButton';
import MuiDialogContent from '@material-ui/core/DialogContent';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  closeButton: {
    float: 'right',
    fontSize: '16px',
    color: 'white',
    borderRadius: '4px',
  },
  dialog: {
    padding: '100px 100px 74px 100px',
  },
  paper: {
    backgroundColor: 'transparent',
    boxShadow: 'none',
    position: 'initial',
    margin: 0,
    width: '100%',
    height: '100%',
  },
  dialogTitle: {
    position: 'absolute',
    top: '100px',
    right: '100px',
  },
  dialogContent: {
    padding: 0,
    width: '100%',
    height: '100%',
  },
  closeIcon: {
    marginLeft: '6px',
  },
  div: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    width: '100%',
    height: '100%',
  },
  img: {
    display: 'block',
  },
  imgName: {
    marginTop: '26px',
    color: 'white',
  },
};

const IMAGE_TITLE_HEIGHT = 45;

type Props = {
  onClose: () => void,
  open: boolean,
  img: ImageDialog_img,
} & WithStyles<typeof styles>;

type State = {
  clientHeight: number,
  clientWidth: number,
};

class ImageDialog extends React.Component<Props, State> {
  constructor(props) {
    super(props);

    this.state = {clientHeight: 0, clientWidth: 0};
  }

  setImageElement = element => {
    if (element != null) {
      this.setState({
        clientHeight: element.clientHeight,
        clientWidth: element.clientWidth,
      });
    }
  };

  onDialogClick = _ => {
    this.props.onClose();
  };

  onImageClick = event => {
    event.stopPropagation();
  };

  render() {
    const {classes, onClose, open} = this.props;
    const img = this.props.img;
    return (
      <Dialog
        open={open}
        onClose={onClose}
        className={classes.dialog}
        classes={{paper: classes.paper}}
        maxWidth={false}
        onClick={this.onDialogClick}>
        <MuiDialogTitle className={classes.dialogTitle}>
          <IconButton
            aria-label="Close"
            onClick={onClose}
            className={classes.closeButton}>
            CLOSE
            <CloseIcon className={classes.closeIcon} />
          </IconButton>
        </MuiDialogTitle>
        <MuiDialogContent className={classes.dialogContent}>
          <div ref={this.setImageElement} className={classes.div}>
            <img
              style={{
                maxHeight: this.state.clientHeight - IMAGE_TITLE_HEIGHT,
                maxWidth: this.state.clientWidth,
              }}
              className={classes.img}
              src={DocumentAPIUrls.get_url(nullthrows(img.storeKey))}
              onClick={this.onImageClick}
            />
            <div className={classes.imgName}>{img.fileName}</div>
          </div>
        </MuiDialogContent>
      </Dialog>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(ImageDialog, {
    img: graphql`
      fragment ImageDialog_img on File {
        storeKey
        fileName
      }
    `,
  }),
);

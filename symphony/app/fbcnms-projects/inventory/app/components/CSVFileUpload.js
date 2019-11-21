/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import axios from 'axios';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  hiddenInput: {
    width: '0px',
    height: '0px',
    opacity: 0,
    overflow: 'hidden',
    position: 'absolute',
    zIndex: -1,
  },
  fileButton: {
    cursor: 'pointer',
    width: 'fit-content',
  },
};

type Props = {
  button: React.Element<any>,
  onProgress: () => void,
  onFileUploaded: () => void,
  onUploadFailed?: any => void,
  uploadPath: string,
} & WithStyles<typeof styles>;

class CSVFileUpload extends React.Component<Props> {
  input = null;

  constructor(props: Props) {
    super(props);
    this.input = null;
  }

  render() {
    const {button, classes} = this.props;
    return (
      <>
        <input
          className={classes.hiddenInput}
          type="file"
          accept=".csv"
          onChange={this.onFileChanged}
          multiple
          ref={ref => {
            this.input = ref;
          }}
        />
        <div className={classes.fileButton} onClick={this.openFileDialog}>
          {button}
        </div>
      </>
    );
  }

  openFileDialog = () => {
    if (this.input) {
      this.input.click();
    }
  };

  onFileChanged = async e => {
    const config = {
      onUploadProgress: _progressEvent => this.props.onProgress(),
      headers: {
        //setting custom CSV file charset from tenant (default utf-8)
        'X-Mime-Charset': window.CONFIG.appData.csvCharset,
      },
    };
    const files = e.target.files;
    if (!files || files.length === 0) {
      return;
    }

    const formData = new FormData();

    Array.from(files).forEach((file, idx) => {
      const name = 'file_' + idx;
      formData.append(name, file);
      idx++;
    });
    try {
      await axios.post(this.props.uploadPath, formData, config);
      this.props.onFileUploaded();
    } catch (error) {
      const message = error.response?.data;
      this.props.onUploadFailed && this.props.onUploadFailed(message);
    }
  };
}

export default withStyles(styles)(CSVFileUpload);

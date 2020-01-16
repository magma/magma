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
import fbt from 'fbt';
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

const SUCCESS_RESPONSE = 0;
const ERROR_RESPONSE = 1;

type Props = {
  button: React.Element<any>,
  onProgress: () => void,
  onFileUploaded: string => void,
  onUploadFailed?: any => void,
  uploadPath: string,
  entity?: string,
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
    formData.append('skip_lines', JSON.stringify([]));
    try {
      const response = await axios.post(
        this.props.uploadPath,
        formData,
        config,
      );
      let msg = '';
      let errorLines = fbt('', 'empty string');

      const entity = this.props.entity ? this.props.entity : '';
      const responseData = response.data;
      const summary = responseData.summary;

      if (summary.messageCode == SUCCESS_RESPONSE) {
        if (responseData.errors != null) {
          const lines = responseData.errors.map(e => '#' + e.line);
          errorLines = fbt(
            'Problematic lines are ' +
              fbt.param('list of rows', lines.toString()),
            'list of rows',
          );
        }
        msg = fbt(
          'Successfully uploaded ' +
            fbt.param('number of saved lines', summary.successLines) +
            ' of ' +
            fbt.param('number of all lines', summary.allLines) +
            ' ' +
            fbt.param('type that was saved', entity) +
            ' items. ' +
            fbt.param('error lines', errorLines),
          'message for a successful import',
        );
        this.props.onFileUploaded(msg);
      } else if (
        summary.messageCode == ERROR_RESPONSE &&
        responseData.errors != null
      ) {
        const tmpErrMessage =
          'Row ' +
          responseData.errors[0].line +
          ': ' +
          responseData.errors[0].message;
        msg = fbt(
          ' Uploaded Failed. ' +
            fbt.param('error message from server', tmpErrMessage),
          'message for a failed import',
        );
        this.props.onUploadFailed && this.props.onUploadFailed(msg);
      }
    } catch (error) {
      const message = error.response?.data;
      this.props.onUploadFailed && this.props.onUploadFailed(message);
    }
  };
}

export default withStyles(styles)(CSVFileUpload);

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import axios from 'axios';
import {makeStyles} from '@material-ui/styles';
import {useRef} from 'react';

const useStyles = makeStyles({
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
    display: 'flex',
    width: '100%',
  },
});

export const FileUploadButton = (props: {
  button: React.Node,
  onFileChanged: (SyntheticEvent<HTMLInputElement>) => void | Promise<void>,
}) => {
  const classes = useStyles();
  const inputRef = useRef();
  return (
    <>
      <input
        className={classes.hiddenInput}
        type="file"
        onChange={props.onFileChanged}
        ref={inputRef}
      />
      <div
        className={classes.fileButton}
        onClick={() => inputRef.current && inputRef.current.click()}>
        {props.button}
      </div>
    </>
  );
};

type Props = {
  button: React.Node,
  onProgress?: (progress: number) => void,
  onFileUploaded: (file: File, key: string) => void,
};

class FileUpload extends React.Component<Props> {
  input = null;

  constructor(props: Props) {
    super(props);

    this.input = null;
  }

  render() {
    const {button} = this.props;
    return (
      <FileUploadButton button={button} onFileChanged={this.onFileChanged} />
    );
  }

  openFileDialog = () => {
    if (this.input) {
      this.input.click();
    }
  };

  onFileChanged = async (e: SyntheticEvent<HTMLInputElement>) => {
    const files = e.currentTarget.files;
    if (!files || files.length === 0) {
      return;
    }

    const file = files[0];
    uploadFile(file, this.props.onFileUploaded, this.props.onProgress);
  };
}

export async function uploadFile(
  file: File,
  onUpload: (File, string) => void,
  onProgress?: number => void,
) {
  const signingResponse = await axios.get('/store/put', {
    params: {
      contentType: file.type,
    },
  });

  const config = {
    headers: {
      'Content-Type': file.type,
    },
    onUploadProgress: function(progressEvent) {
      const percentCompleted = Math.round(
        (progressEvent.loaded * 100) / progressEvent.total,
      );
      onProgress && onProgress(percentCompleted);
    },
  };
  await axios.put(signingResponse.data.URL, file, config);

  onUpload(file, signingResponse.data.key);
}

export default FileUpload;

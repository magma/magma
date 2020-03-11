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
import UploadIcon from '@fbcnms/ui/components/design-system/Icons/Actions/UploadIcon';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useRef, useState} from 'react';

const useStyles = makeStyles(() => ({
  photoUploadContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '112px',
    width: '112px',
    backgroundColor: symphony.palette.D10,
    border: `1px dashed ${symphony.palette.D100}`,
    borderRadius: '4px',
    '&:hover': {
      borderColor: symphony.palette.primary,
      cursor: 'pointer',
    },
    '&:active': {
      backgroundColor: symphony.palette.D50,
    },
  },
  hiddenInput: {
    width: '0px',
    height: '0px',
    opacity: 0,
    overflow: 'hidden',
    position: 'absolute',
    zIndex: -1,
  },
}));

export type FileUploadAreaProps = {
  onFileChanged: FileList => void,
  fileTypes?: string,
};

const FileUploadArea = (props: FileUploadAreaProps) => {
  const {fileTypes = 'file', onFileChanged} = props;
  const classes = useStyles();
  const [hoversUploadPhoto, setHoversUploadPhoto] = useState(false);
  const inputRef = useRef();
  const showFileDialog = () => inputRef?.current?.click();

  const onFileSelected: (SyntheticInputEvent<HTMLInputElement>) => void = e =>
    onFileChanged(e.currentTarget.files);
  return (
    <>
      <div
        onMouseEnter={() => setHoversUploadPhoto(true)}
        onMouseLeave={() => setHoversUploadPhoto(false)}
        onClick={showFileDialog}
        className={classes.photoUploadContainer}>
        <UploadIcon color={hoversUploadPhoto ? 'primary' : 'gray'} />
      </div>

      <input
        className={classes.hiddenInput}
        type={fileTypes}
        onChange={onFileSelected}
        ref={inputRef}
      />
    </>
  );
};

export default FileUploadArea;

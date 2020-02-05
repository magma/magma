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
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CSVFileUpload from './CSVFileUpload';
import CircularProgress from '@material-ui/core/CircularProgress';
import DialogConfirm from '@fbcnms/ui/components/DialogConfirm';
import DialogError from '@fbcnms/ui/components/DialogError';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {UploadAPIUrls} from '../common/UploadAPI';
import {useContext, useState} from 'react';
import {withStyles} from '@material-ui/core/styles';

type Props = {} & WithStyles<typeof styles>;

const styles = _ => ({
  uploadContent: {
    padding: '20px',
  },
  link: {
    paddingLeft: '28px',
  },
});

const deprecatedUploadsParams = [
  {
    text: 'Upload Equipment',
    uploadPath: UploadAPIUrls.equipment(),
    entity: null,
  },
  {
    text: 'Upload Position Def',
    uploadPath: UploadAPIUrls.position_definition(),
    entity: null,
  },
  {
    text: 'Upload Port Def',
    uploadPath: UploadAPIUrls.port_definition(),
    entity: null,
  },
  {
    text: 'Upload Port Connections',
    uploadPath: UploadAPIUrls.port_connect(),
    entity: null,
  },
];

const uploadParams = [
  {
    text: 'Upload Locations',
    uploadPath: UploadAPIUrls.locations(),
    entity: null,
  },
  {
    text: 'Upload Exported Equipment',
    uploadPath: UploadAPIUrls.exported_equipment(),
    entity: 'equipment',
  },
  {
    text: 'Upload Exported Ports',
    uploadPath: UploadAPIUrls.exported_ports(),
    entity: 'port',
  },
  {
    text: 'Upload Exported Links',
    uploadPath: UploadAPIUrls.exported_links(),
    entity: 'link',
  },
  {
    text: 'Upload Exported Service',
    uploadPath: UploadAPIUrls.exported_service(),
    entity: 'service',
  },
];

const CSVUploadDialog = (props: Props) => {
  const {classes} = props;

  const documentsLink = (path: string) => (
    <Text
      className={classes.link}
      onClick={() =>
        ServerLogger.info(
          LogEvents.DOCUMENTATION_LINK_CLICKED_FROM_EXPORT_DIALOG,
        )
      }>
      <a href={'/docs/docs/' + path}>Go to documentation page</a>
    </Text>
  );

  const [errorMessage, setErrorMessage] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const onFileUploaded = (msg: string) => {
    setIsLoading(false);
    setSuccessMessage(msg);
    setErrorMessage(null);
  };

  const onUploadFailed = (msg: string) => {
    setIsLoading(false);
    setErrorMessage(msg);
    setSuccessMessage(null);
  };

  const appContext = useContext(AppContext);

  const deprecatedImportsEnabled = appContext.isFeatureEnabled(
    'deprecated_imports',
  );
  const importScripts = deprecatedImportsEnabled
    ? deprecatedUploadsParams.concat(uploadParams)
    : uploadParams;
  return isLoading ? (
    <CircularProgress />
  ) : (
    <>
      {errorMessage && <DialogError message={errorMessage} />}
      {successMessage && <DialogConfirm message={successMessage} />}

      {documentsLink('csv-upload.html')}
      <div className={classes.uploadContent}>
        {importScripts.map(entity => (
          <CSVFileUpload
            key={entity.uploadPath}
            button={<Button variant="text">{entity.text}</Button>}
            onProgress={() => {
              setIsLoading(true);
              setErrorMessage(null);
            }}
            onFileUploaded={msg => onFileUploaded(msg)}
            uploadPath={entity.uploadPath}
            onUploadFailed={msg => onUploadFailed(msg)}
          />
        ))}
      </div>
    </>
  );
};

export default withStyles(styles)(CSVUploadDialog);

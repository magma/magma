/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  DeleteImageMutationResponse,
  DeleteImageMutationVariables,
} from '../../mutations/__generated__/DeleteImageMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveSiteSurveyMutationResponse,
  RemoveSiteSurveyMutationVariables,
} from '../../mutations/__generated__/RemoveSiteSurveyMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import DeleteIcon from '@material-ui/icons/Delete';
import DeleteImageMutation from '../../mutations/DeleteImageMutation';
import IconButton from '@material-ui/core/IconButton';
import Link from '@fbcnms/ui/components/Link';
import React from 'react';
import RemoveSiteSurveyMutation from '../../mutations/RemoveSiteSurveyMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import axios from 'axios';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {DocumentAPIUrls} from '../../common/DocumentAPI';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useRef} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

export type LocationSiteSurveyEntry = {
  +id: string,
  +completionTimestamp: number,
  +name: string,
  +ownerName: ?string,
  +sourceFile: ?{
    +id: string,
    +fileName: string,
    +storeKey: ?string,
  },
};

type Props = {
  survey: LocationSiteSurveyEntry,
  onSurveySelected: () => void,
} & WithAlert;

const useStyles = makeStyles(_theme => ({
  statusText: {
    fontWeight: 'bold',
  },
}));

function LocationSiteSurvey(props: Props) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {survey, onSurveySelected, confirm} = props;
  const storeKey = survey.sourceFile?.storeKey;
  const fileName = survey.sourceFile?.fileName;
  const fileId = survey.sourceFile?.id;
  const downloadFileRef = useRef(null);
  const handleDownload = useCallback(() => {
    if (downloadFileRef.current != null) {
      downloadFileRef.current.click();
    }
  }, [downloadFileRef]);

  const onSurveySourceFileDeleted = useCallback(() => {
    if (fileId == null) {
      return;
    }
    const variables: DeleteImageMutationVariables = {
      entityType: 'SITE_SURVEY',
      entityId: survey.id,
      id: fileId,
    };

    const callbacks: MutationCallbacks<DeleteImageMutationResponse> = {
      onCompleted: (_, errors) => {
        if (errors && errors[0]) {
          enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        }
      },
      onError: () => {},
    };

    DeleteImageMutation(variables, callbacks);
  }, [enqueueSnackbar, survey.id, fileId]);

  const handleDelete = useCallback(async () => {
    if (storeKey != null && fileId != null) {
      await axios.delete(`/store/delete?key=${storeKey}`);
      onSurveySourceFileDeleted();
    }
    const variables: RemoveSiteSurveyMutationVariables = {
      id: survey.id,
    };

    const callbacks: MutationCallbacks<RemoveSiteSurveyMutationResponse> = {
      onCompleted: (_, errors) => {
        if (errors && errors[0]) {
          enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        }
      },
      onError: () => {},
    };

    RemoveSiteSurveyMutation(variables, callbacks, store => {
      store.delete(survey.id);
    });
  }, [survey.id, storeKey, fileId, onSurveySourceFileDeleted, enqueueSnackbar]);

  const handleDeleteConfirmation = useCallback(() => {
    confirm(
      `Are you sure you want to delete "${survey.name}" site survey ?`,
    ).then(confirmed => {
      if (confirmed) {
        handleDelete();
      }
    });
  }, [confirm, survey.name, handleDelete]);

  return (
    <TableRow key={survey.id}>
      <TableCell>{survey.name}</TableCell>
      <TableCell>{survey.ownerName}</TableCell>
      <TableCell>
        <div>
          <Text className={classes.statusText}>Done:</Text>
        </div>
        <div>
          <Text>
            {new Intl.DateTimeFormat('en-US').format(
              new Date(survey.completionTimestamp * 1000),
            )}
          </Text>
        </div>
      </TableCell>
      <TableCell>
        <Link onClick={() => onSurveySelected()}>View Results</Link>
      </TableCell>
      <TableCell>
        <span>
          <IconButton onClick={() => handleDeleteConfirmation()}>
            <DeleteIcon />
          </IconButton>
          {storeKey && fileName && (
            <a
              href={DocumentAPIUrls.download_url(storeKey, fileName)}
              ref={downloadFileRef}
              style={{display: 'none'}}
              download
            />
          )}
          {survey.sourceFile?.storeKey && (
            <IconButton onClick={handleDownload}>
              <CloudDownloadIcon />
            </IconButton>
          )}
        </span>
      </TableCell>
    </TableRow>
  );
}

export default withAlert(LocationSiteSurvey);

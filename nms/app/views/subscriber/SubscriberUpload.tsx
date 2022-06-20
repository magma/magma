/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Alert from '@material-ui/lab/Alert';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import React from 'react';
import Text from '../../theme/design-system/Text';
import {
  CoreNetworkTypes,
  SUBSCRIBER_ACTION_TYPE,
  SubscriberInfo,
} from './SubscriberUtils';
import {DropzoneArea} from 'material-ui-dropzone';
import {SubscriberForbiddenNetworkTypesEnum} from '../../../generated-ts';
import {
  SubscribersDialogDetailProps,
  validateSubscribers,
} from './SubscriberUtils';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  uploadDialog: {
    width: '800px',
  },
  uploadInstructions: {
    marginTop: '16px',
    color: colors.primary.comet,
  },
}));
const forbiddenNetworkTypes = Object.values(CoreNetworkTypes);
const SUB_NAME_OFFSET = 0;
const SUB_IMSI_OFFSET = 1;
const SUB_AUTH_KEY_OFFSET = 2;
const SUB_AUTH_OPC_OFFSET = 3;
const SUB_STATE_OFFSET = 4;
const SUB_FORBIDDEN_NETWORK_TYPE_OFFSET = 5;
const SUB_DATAPLAN_OFFSET = 6;
const SUB_APN_OFFSET = 7;
const SUB_POLICY_OFFSET = 8;
const SUB_MAX_FIELDS = 9;
const MAX_UPLOAD_FILE_SZ_BYTES = 10 * 1024 * 1024;
const UPLOAD_DOC_LINK =
  'https://docs.magmacore.org/docs/nms/subscriber#uploading-a-subscriber-csv-file';
const ADD_INSTRUCTIONS =
  'You can download this template that automatically maps the fields. Find more instruction in ';
const DELETE_INSTRUCTIONS =
  'You can export all subscribers and select the subscribers you want to delete. Find more instruction in ';
const EDIT_INSTRUCTIONS =
  'You can export all subscribers to edit and upload the file. Find more instruction in ';

function parseSubscriber(line: string): SubscriberInfo {
  const items = line.split(',').map(item => item.trim());

  if (items.length > SUB_MAX_FIELDS) {
    throw new Error(
      `Too many fields to parse, expected ${SUB_MAX_FIELDS} fields, received ${items.length} fields`,
    );
  }

  return {
    name: items[SUB_NAME_OFFSET],
    imsi: items[SUB_IMSI_OFFSET],
    authKey: items[SUB_AUTH_KEY_OFFSET],
    authOpc: items[SUB_AUTH_OPC_OFFSET],
    state: items[SUB_STATE_OFFSET] === 'ACTIVE' ? 'ACTIVE' : 'INACTIVE',
    forbiddenNetworkTypes: forbiddenNetworkTypes.filter((value: string) =>
      items[SUB_FORBIDDEN_NETWORK_TYPE_OFFSET]?.split('|')
        .map(item => item.trim())
        .filter(Boolean)
        .includes(value),
    ) as Array<SubscriberForbiddenNetworkTypesEnum>,
    dataPlan: items[SUB_DATAPLAN_OFFSET],
    apns: items[SUB_APN_OFFSET]?.split('|')
      .map(item => item.trim())
      .filter(Boolean),
    policies: items?.[SUB_POLICY_OFFSET]?.split('|')
      .map(item => item.trim())
      .filter(Boolean),
  };
}

function parseSubscriberFile(fileObj: File): Promise<Array<SubscriberInfo>> {
  const reader = new FileReader();
  const subscribers: Array<SubscriberInfo> = [];
  return new Promise((resolve, reject) => {
    if (fileObj.size > MAX_UPLOAD_FILE_SZ_BYTES) {
      reject(
        'file size exceeds max upload size of 10MB, please upload smaller file',
      );
      return;
    }

    reader.onload = e => {
      try {
        if (!(e.target instanceof FileReader)) {
          reject('invalid target type');
          return;
        }

        const text = e.target.result;

        if (typeof text !== 'string') {
          reject('invalid file content');
          return;
        }

        for (const line of text
          .split('\n')
          .map(item => item.trim())
          .filter(Boolean)) {
          subscribers.push(parseSubscriber(line));
        }
      } catch (e) {
        reject(
          `Failed parsing the file ${fileObj.name}. ${getErrorMessage(e)}`,
        );
        return;
      }

      resolve(subscribers);
    };

    reader.readAsText(fileObj);
  });
}

export function SubscriberDetailsUpload(props: SubscribersDialogDetailProps) {
  const {
    setSubscribers,
    setAddError,
    setUpload,
    upload,
    subscribers,
    subscriberAction,
  } = props;
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [fileName, setFileName] = useState('');

  const DropzoneText = () => (
    <div>
      Drag and drop or <Link>browse files</Link>
    </div>
  );

  return (
    <>
      <DialogContent
        classes={{
          root: classes.uploadDialog,
        }}>
        <CardTitleRow label={'Upload CSV'} />
        <Grid container>
          <Grid item xs={12}>
            {subscriberAction !== SUBSCRIBER_ACTION_TYPE.EDIT && (
              <Alert severity="warning">
                This will replace the subscribers you entered on the previous
                page.
              </Alert>
            )}
          </Grid>
          {!fileName ? (
            <Grid item xs={12}>
              <DropzoneArea
                dropzoneText={((<DropzoneText />) as unknown) as string}
                useChipsForPreview
                showPreviewsInDropzone={false}
                filesLimit={1}
                showAlerts={false}
                // eslint-disable-next-line @typescript-eslint/no-misused-promises
                onChange={async files => {
                  if (files.length) {
                    try {
                      const newSubscribers: Array<SubscriberInfo> = await parseSubscriberFile(
                        files[0],
                      );

                      if (newSubscribers) {
                        setSubscribers([...newSubscribers]);
                        const errors = validateSubscribers(
                          newSubscribers,
                          subscriberAction,
                        );
                        setFileName(files[0].name);

                        if (!(subscriberAction === 'delete')) {
                          setUpload(false);
                          setAddError(errors);
                        }
                      }
                    } catch (e) {
                      enqueueSnackbar(e as string, {
                        variant: 'error',
                      });
                    }
                  }
                }}
              />
              <UploadInstructions action={subscriberAction} />
            </Grid>
          ) : (
            <Grid item xs={12}>
              <Alert severity="success">{`${fileName} is uploaded`}</Alert>
            </Grid>
          )}
        </Grid>
      </DialogContent>
      <DialogActions>
        <Grid container justifyContent="space-between">
          <Grid item>
            {upload && (
              <Button
                onClick={() => {
                  setUpload(false);

                  if (subscriberAction === 'delete' && subscribers.length > 0) {
                    setSubscribers([]);
                  }
                }}>
                Back
              </Button>
            )}
          </Grid>
          <Grid item>
            <Button onClick={props.onClose}> Cancel </Button>
            <Button
              data-testid="saveSubscriber"
              variant="contained"
              color="primary"
              onClick={() => {
                props.onSave?.(subscribers);
              }}>
              {subscriberAction === 'delete'
                ? 'Delete Subcribers'
                : subscriberAction === 'edit'
                ? 'Update Subscribers'
                : 'Save and Add Subscribers'}
            </Button>
          </Grid>
        </Grid>
      </DialogActions>
    </>
  );
}

function UploadInstructions({action}: {action: string}) {
  const classes = useStyles();
  const instructions = useMemo(() => {
    switch (action) {
      case 'delete':
        return DELETE_INSTRUCTIONS;

      case 'edit':
        return EDIT_INSTRUCTIONS;

      case 'add':
        return ADD_INSTRUCTIONS;

      default:
        return '';
    }
  }, [action]);
  return (
    <Text variant="body2" className={classes.uploadInstructions}>
      {`Accepted file type: .csv (<10 MB).  ${instructions}`}
      <Link href={UPLOAD_DOC_LINK}>documentation</Link>
    </Text>
  );
}

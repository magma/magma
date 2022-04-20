/**
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
 *
 * @flow
 * @format
 */

import * as React from 'react';
import BaseDialog from './BaseDialog';
import Button from '../Button';
import Checkbox from '../Checkbox/Checkbox';
import Text from '../Text';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    minWidth: '480px',
    minHeight: '210px',
    maxWidth: '600px',
    maxHeight: '600px',
  },
  checkboxContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    flexGrow: 1,
  },
  content: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  body: {
    flexGrow: 1,
    marginBottom: '40px',
  },
  footer: {
    paddingBottom: '8px',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

export type DialogSkin = 'primary' | 'red';

export type MessageDialogProps = $ReadOnly<{|
  title: React.Node,
  message: React.Node,
  verificationCheckbox?: ?{
    label: React.Node,
    isMandatory?: ?boolean,
  },
  cancelLabel?: React.Node,
  confirmLabel?: React.Node,
  skin?: DialogSkin,
  onCancel?: () => void,
  onClose: () => void,
  onConfirm?: (?boolean) => void,
|}>;

export type MessageDialogComponentProps = $ReadOnly<{|
  ...MessageDialogProps,
  hidden?: boolean,
|}>;

const MessageDialog = ({
  title,
  message,
  onClose,
  verificationCheckbox,
  cancelLabel = 'Cancel',
  confirmLabel = 'OK',
  onCancel,
  onConfirm,
  hidden,
  skin = 'primary',
}: MessageDialogComponentProps) => {
  const classes = useStyles();
  const [checkboxChecked, setCheckboxChecked] = useState(false);

  useEffect(() => {
    setCheckboxChecked(false);
  }, [hidden]);

  return (
    <BaseDialog
      className={classes.root}
      title={title}
      onClose={onClose}
      hidden={hidden}>
      <div className={classes.content}>
        <div className={classes.body}>
          <Text>{message}</Text>
        </div>
        <div className={classes.footer}>
          {verificationCheckbox && (
            <div className={classes.checkboxContainer}>
              <Checkbox
                checked={checkboxChecked}
                title={verificationCheckbox.label}
                onChange={selection =>
                  setCheckboxChecked(selection === 'checked' ? true : false)
                }
              />
            </div>
          )}
          {cancelLabel && (
            <Button
              skin="gray"
              onClick={onCancel}
              className={classes.cancelButton}>
              {cancelLabel}
            </Button>
          )}
          {confirmLabel && (
            <Button
              onClick={() =>
                onConfirm &&
                onConfirm(verificationCheckbox == null ? null : checkboxChecked)
              }
              skin={skin}
              disabled={
                verificationCheckbox?.isMandatory === true && !checkboxChecked
              }>
              {confirmLabel}
            </Button>
          )}
        </div>
      </div>
    </BaseDialog>
  );
};

export default MessageDialog;

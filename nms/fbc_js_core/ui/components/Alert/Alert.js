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

import type {Node} from 'react';

import Button from '../design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import Text from '../../../../app/theme/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  paper: {
    minWidth: `${theme.breakpoints.values.sm / 2}px`,
  },
}));

export type AlertSkin = 'primary' | 'red';

type Props = {|
  cancelLabel?: Node,
  confirmLabel?: Node,
  message: Node,
  skin?: AlertSkin,
  onCancel?: () => void,
  onClose?: () => void,
  onConfirm?: () => void,
  title?: ?Node,
  open?: boolean,
|};

const Alert = ({
  cancelLabel,
  confirmLabel,
  message,
  onCancel,
  onClose,
  onConfirm,
  title,
  open,
  skin = 'primary',
}: Props) => {
  const classes = useStyles();
  const hasActions = cancelLabel != null || confirmLabel != null;

  return (
    <Dialog
      classes={{paper: classes.paper}}
      open={open}
      onClose={onCancel}
      TransitionProps={{onExited: onClose}}
      maxWidth="sm">
      {title && <DialogTitle>{title}</DialogTitle>}
      <DialogContent>
        <Text>{message}</Text>
      </DialogContent>
      {hasActions && (
        <DialogActions>
          {cancelLabel && (
            <Button skin="regular" onClick={onCancel}>
              {cancelLabel}
            </Button>
          )}
          {confirmLabel && (
            <Button onClick={onConfirm} skin={skin}>
              {confirmLabel}
            </Button>
          )}
        </DialogActions>
      )}
    </Dialog>
  );
};

export default Alert;

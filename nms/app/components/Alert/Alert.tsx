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
 */

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import Text from '../../theme/design-system/Text';
import {Theme} from '@material-ui/core/styles';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles<Theme>(theme => ({
  paper: {
    minWidth: `${theme.breakpoints.values.sm / 2}px`,
  },
}));

type Props = {
  cancelLabel?: React.ReactNode;
  confirmLabel?: React.ReactNode;
  message: React.ReactNode;
  onCancel?: () => void;
  onClose?: () => void;
  onConfirm?: () => void;
  title?: React.ReactNode | null;
  open?: boolean;
};

const Alert = ({
  cancelLabel,
  confirmLabel,
  message,
  onCancel,
  onClose,
  onConfirm,
  title,
  open,
}: Props) => {
  const classes = useStyles();
  const hasActions = cancelLabel != null || confirmLabel != null;

  return (
    <Dialog
      classes={{paper: classes.paper}}
      open={!!open}
      onClose={onCancel}
      TransitionProps={{onExited: onClose}}
      maxWidth="sm">
      {title && <DialogTitle>{title}</DialogTitle>}
      <DialogContent>
        <Text>{message}</Text>
      </DialogContent>
      {hasActions && (
        <DialogActions>
          {cancelLabel && <Button onClick={onCancel}>{cancelLabel}</Button>}
          {confirmLabel && (
            <Button onClick={onConfirm} variant="contained" color="primary">
              {confirmLabel}
            </Button>
          )}
        </DialogActions>
      )}
    </Dialog>
  );
};

export default Alert;

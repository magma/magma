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

import Button from '../../../../../fbc_js_core/ui/components/design-system/Button';
import FormAction from './FormAction';
import React from 'react';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  cancelButton: {
    marginRight: '8px',
  },
}));

type Props = {
  isDisabled?: boolean,
  disabledMessage?: string,
  onSave: () => void,
  onCancel: () => void,
  classes?: {
    cancelButton?: string,
    saveButton?: string,
  },
  captions?: {
    cancelButton?: string,
    saveButton?: string,
  },
};

const FormSaveCancelPanel = (props: Props) => {
  const {
    captions,
    classes: propsClasses,
    isDisabled,
    disabledMessage,
    onCancel,
    onSave,
  } = props;
  const classes = useStyles();
  return (
    <div title={isDisabled && disabledMessage}>
      <FormAction
        ignorePermissions={true}
        disabled={isDisabled}
        tooltip={isDisabled ? disabledMessage : undefined}>
        <Button
          className={classNames(
            classes.cancelButton,
            propsClasses?.cancelButton,
          )}
          onClick={onCancel}
          skin="regular">
          {captions?.cancelButton || 'Cancel'}
        </Button>
      </FormAction>
      <FormAction
        disableOnFromError={true}
        disabled={isDisabled}
        tooltip={isDisabled ? disabledMessage : undefined}>
        <Button className={propsClasses?.saveButton} onClick={onSave}>
          {captions?.saveButton || 'Save'}
        </Button>
      </FormAction>
    </div>
  );
};

export default FormSaveCancelPanel;

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
 * @flow strict-local
 * @format
 */

import type {BaseDialogComponentProps, BaseDialogProps} from './BaseDialog';

import * as React from 'react';
import BaseDialog from './BaseDialog';
import emptyFunction from '../../../../../fbc_js_core/util/emptyFunction';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

export type DialogShowingContextValue = $ReadOnly<{|
  showDialog: (props: BaseDialogProps, replaceExisting?: ?boolean) => void,
  hideDialog: () => void,
|}>;

const DialogShowingContext = createContext<DialogShowingContextValue>({
  showDialog: emptyFunction,
  hideDialog: emptyFunction,
});

export function useDialogShowingContext() {
  return useContext(DialogShowingContext);
}

type Props = $ReadOnly<{|
  children: React.Node,
|}>;

const EMPTY_HIDDEN_DIALOG_PROPS: BaseDialogComponentProps = {
  title: null,
  children: null,
  onClose: emptyFunction,
  hidden: true,
};

export function DialogShowingContextProvider(props: Props) {
  const [
    baseDialogProps,
    setBaseDialogProps,
  ] = useState<BaseDialogComponentProps>(EMPTY_HIDDEN_DIALOG_PROPS);

  const showDialog = useCallback(
    (props: BaseDialogProps, replaceExisting?: ?boolean) => {
      if (baseDialogProps.hidden !== true && replaceExisting !== true) {
        return;
      }
      setBaseDialogProps(props);
    },
    [baseDialogProps],
  );
  const hideDialog = useCallback(() => {
    setBaseDialogProps(currentDialogProps => ({
      ...currentDialogProps,
      hidden: true,
    }));
  }, []);

  const value = useMemo(
    () => ({
      showDialog,
      hideDialog,
    }),
    [hideDialog, showDialog],
  );

  return (
    <DialogShowingContext.Provider value={value}>
      {props.children}
      <BaseDialog {...baseDialogProps} />
    </DialogShowingContext.Provider>
  );
}

export default DialogShowingContext;

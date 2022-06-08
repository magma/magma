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
 * @flow strict-local strict-local
 * @format
 */

'use strict';

import type {ComponentType, ElementConfig, Node} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Alert from './Alert';
import React from 'react';

type State = {|
  dialogs: Map<DialogMapKey, boolean>,
|};

export type DialogProps = {|
  cancelLabel?: Node,
  confirmLabel?: Node,
  message: Node,
  title?: ?Node,
|};

export type DialogMapKey = {|
  key: number,
  onCancel: () => void,
  onClose: () => void,
  onConfirm: () => void,
  ...DialogProps,
|};

export type WithAlert = {|
  alert: (Node | Error, ?Node) => Promise<boolean>,
  confirm: (Node | DialogProps) => Promise<boolean>,
|};

function withAlert<Props: WithAlert, TComponent: ComponentType<Props>>(
  Component: TComponent,
): ComponentType<$Diff<ElementConfig<TComponent>, WithAlert>> {
  return class extends React.Component<
    $Diff<ElementConfig<TComponent>, WithAlert>,
    State,
  > {
    state = {
      dialogs: new Map<DialogMapKey, boolean>(),
    };

    lastKey = 0;

    removeDialog(alert) {
      const nextAlerts = new Map<DialogMapKey, boolean>(this.state.dialogs);
      nextAlerts.delete(alert);
      this.setState({dialogs: nextAlerts});
    }

    closeDialog(alert) {
      this.setState({dialogs: this.state.dialogs.set(alert, false)});
    }

    addDialog(props: DialogProps): Promise<boolean> {
      let dialog: DialogMapKey;
      this.lastKey = this.lastKey + 1;
      return new Promise<boolean>(resolve => {
        dialog = {
          ...props,
          key: this.lastKey,
          onCancel: () => {
            this.closeDialog(dialog);
            resolve(false);
          },
          onConfirm: () => {
            this.closeDialog(dialog);
            resolve(true);
          },
          onClose: () => this.removeDialog(dialog),
        };
        this.setState({
          dialogs: this.state.dialogs.set(dialog, true),
        });
      });
    }

    alert = (
      message: Node | Error,
      confirmLabel?: Node = 'OK',
    ): Promise<boolean> => {
      return this.addDialog({
        message: message instanceof Error ? String(message) : message,
        confirmLabel,
      }).catch(() => {
        /* always resolve */
        return false;
      });
    };

    confirm = (messageOrProps: DialogProps | Node): Promise<boolean> => {
      let dialogProps: DialogProps;
      const confirmLabel = 'Confirm';
      const cancelLabel = 'Cancel';

      if (
        typeof messageOrProps === 'string' ||
        React.isValidElement(messageOrProps)
      ) {
        dialogProps = {
          confirmLabel,
          cancelLabel,
          // $FlowFixMe - we validated props is a Node
          message: (messageOrProps: Node),
        };
      } else {
        dialogProps = {
          confirmLabel,
          cancelLabel,
          // $FlowFixMe - we validated props is DialogProps
          ...(messageOrProps: DialogProps),
        };
      }
      return this.addDialog(dialogProps);
    };

    render() {
      return (
        <>
          <Component
            {...this.props}
            alert={this.alert}
            confirm={this.confirm}
          />
          {[...this.state.dialogs.entries()].map(([dialog, open]) => (
            <Alert {...dialog} open={open} />
          ))}
        </>
      );
    }
  };
}

export default withAlert;

/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';

import type {ComponentType, Node} from 'react';

import Alert from './Alert';
import React from 'react';

type State = {|
  dialogs: Map<DialogMapKey, boolean>,
|};

export type DialogProps = {
  cancelLabel?: Node,
  confirmLabel?: Node,
  message: Node,
  title?: ?Node,
};

export type DialogMapKey = DialogProps & {
  key: number,
  onCancel: () => void,
  onClose: () => void,
  onConfirm: () => void,
};

export type WithAlert = {|
  alert: (Node | Error, ?Node) => Promise<*>,
  confirm: (Node | DialogProps) => Promise<*>,
|};

function withAlert<Props: {}>(
  Component: ComponentType<Props>,
): ComponentType<Props> {
  return class extends React.Component<Props, State> {
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

    addDialog(props: DialogProps): Promise<*> {
      let dialog: DialogMapKey;
      this.lastKey = this.lastKey + 1;
      return new Promise<*>(resolve => {
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

    alert = (message: Node | Error, confirmLabel?: Node = 'Ok'): Promise<*> => {
      return this.addDialog({
        message: message instanceof Error ? String(message) : message,
        confirmLabel,
      }).catch(() => {
        /* always resolve */
      });
    };

    confirm = (messageOrProps: DialogProps | Node): Promise<*> => {
      let dialogProps: DialogProps;
      const confirmLabel = <>Confirm</>;
      const cancelLabel = <>Cancel</>;

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

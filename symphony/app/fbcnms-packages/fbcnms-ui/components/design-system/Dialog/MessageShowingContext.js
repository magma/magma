/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  MessageDialogComponentProps,
  MessageDialogProps,
} from './MessageDialog';

import * as React from 'react';
import MessageDialog from './MessageDialog';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

export type MessageShowingContextValue = $ReadOnly<{|
  showMessage: MessageDialogProps => void,
  hideMessage: () => void,
|}>;

const MessageShowingContext = createContext<MessageShowingContextValue>({
  showMessage: emptyFunction,
  hideMessage: emptyFunction,
});

export function useMessageShowingContext() {
  return useContext(MessageShowingContext);
}

type Props = $ReadOnly<{|
  children: React.Node,
|}>;

const EMPTY_HIDDEN_MESSAGE_PROPS: MessageDialogComponentProps = {
  title: null,
  message: null,
  onClose: emptyFunction,
  hidden: true,
};

export function MessageShowingContextProvider(props: Props) {
  const [
    messageDialogProps,
    setMessageDialogProps,
  ] = useState<MessageDialogComponentProps>(EMPTY_HIDDEN_MESSAGE_PROPS);

  const showMessage = useCallback(
    (props: MessageDialogProps) => {
      if (messageDialogProps.hidden !== true) {
        return;
      }
      setMessageDialogProps(props);
    },
    [messageDialogProps],
  );
  const hideMessage = useCallback(() => {
    setMessageDialogProps(currentMessageProps => ({
      ...currentMessageProps,
      hidden: true,
    }));
  }, []);

  const value = useMemo(
    () => ({
      showMessage,
      hideMessage,
    }),
    [hideMessage, showMessage],
  );

  return (
    <MessageShowingContext.Provider value={value}>
      {props.children}
      <MessageDialog {...messageDialogProps} />
    </MessageShowingContext.Provider>
  );
}

export default MessageShowingContext;

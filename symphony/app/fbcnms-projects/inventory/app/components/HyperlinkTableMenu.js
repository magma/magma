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
  DeleteHyperlinkMutationResponse,
  DeleteHyperlinkMutationVariables,
} from '../mutations/__generated__/DeleteHyperlinkMutation.graphql';
import type {HyperlinkTableMenu_hyperlink} from './__generated__/HyperlinkTableMenu_hyperlink.graphql';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import DeleteHyperlinkMutation from '../mutations/DeleteHyperlinkMutation';
import OptionsPopoverButton from './OptionsPopoverButton';
import React, {useCallback} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import fbt from 'fbt';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {withSnackbar} from 'notistack';

type Props = {
  entityId: string,
  hyperlink: HyperlinkTableMenu_hyperlink,
} & WithAlert &
  WithSnackbarProps;

const HyperlinkTableMenu = (props: Props) => {
  const {entityId, confirm, enqueueSnackbar} = props;
  const hyperlink = props.hyperlink;
  const deleteHandler = useCallback(() => {
    confirm(
      fbt(
        'Are you sure you want to delete link ' +
          fbt.param(
            'Hyperlink display text',
            !!hyperlink.displayName ? `'${hyperlink.displayName}' ` : '',
          ) +
          `to '` +
          fbt.param('Hyperlink target URL', hyperlink.url) +
          `'?`,
        'Hyperlink delete confirmation message',
      ),
    ).then(confirmed => {
      if (!confirmed) {
        return;
      }

      const variables: DeleteHyperlinkMutationVariables = {
        id: hyperlink.id,
      };

      const updater = store => {
        // $FlowFixMe (T62907961) Relay flow types
        const deletedNode = store.getRootField('deleteHyperlink');
        // $FlowFixMe (T62907961) Relay flow types
        const proxy = store.get(entityId);
        // $FlowFixMe (T62907961) Relay flow types
        const currNodes = proxy.getLinkedRecords('hyperlinks');
        // $FlowFixMe (T62907961) Relay flow types
        const nodesToKeep = currNodes.filter(hyperlinkNode => {
          return hyperlinkNode != deletedNode;
        });
        // $FlowFixMe (T62907961) Relay flow types
        proxy.setLinkedRecords(nodesToKeep, 'hyperlinks');
        // $FlowFixMe (T62907961) Relay flow types
        store.delete(hyperlink.id);
      };

      const errorMessageHandling = errorMessage => {
        enqueueSnackbar(errorMessage, {
          children: key => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        });
      };

      const cbs: MutationCallbacks<DeleteHyperlinkMutationResponse> = {
        onCompleted: (_, errors) => {
          if (errors && errors[0]) {
            errorMessageHandling(errors[0].message);
          }
        },
        onError: error => {
          errorMessageHandling(error.message);
        },
      };

      DeleteHyperlinkMutation(variables, cbs, updater);
    });
  }, [
    confirm,
    enqueueSnackbar,
    entityId,
    hyperlink.displayName,
    hyperlink.id,
    hyperlink.url,
  ]);
  const menuOptions = [
    {
      onClick: deleteHandler,
      caption: fbt(
        'Delete',
        'Caption for menu option for deleting a url from hyperlinks table',
      ),
    },
    ,
  ];
  return <OptionsPopoverButton options={menuOptions} />;
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(HyperlinkTableMenu, {
      hyperlink: graphql`
        fragment HyperlinkTableMenu_hyperlink on Hyperlink {
          id
          displayName
          url
          ...HyperlinkTableRow_hyperlink
        }
      `,
    }),
  ),
);

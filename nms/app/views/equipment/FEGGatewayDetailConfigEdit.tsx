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

import type {FederationGateway} from '../../../generated-ts';
import type {TabOption} from '../../components/feg/FEGGatewayDialog';

import Button from '@material-ui/core/Button';
import FEGGatewayDialog from '../../components/feg/FEGGatewayDialog';
import React from 'react';

import {useState} from 'react';

type ButtonProps = {
  editingGateway: FederationGateway;
  tabOption?: TabOption;
  title: string;
};

/**
 * Return a button which allows a user to edit the federation
 * gateway. It displays the FEGGatewayDialog component when it
 * is clicked / open.
 * @param {FederationGateway} editingGateway The federation gateway being edited.
 * @param {TabOption} tabOption The Tab that is being looked at.
 * @param {string} title Title of the button.
 */
export default function EditGatewayButton(props: ButtonProps) {
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      {open && (
        <FEGGatewayDialog
          editingGateway={props.editingGateway}
          tabOption={props.tabOption}
          onClose={handleClose}
          onSave={() => handleClose()}
        />
      )}
      <Button
        data-testid={(props.tabOption ?? '') + 'EditButton'}
        component="button"
        variant="text"
        onClick={handleClickOpen}>
        {props.title}
      </Button>
    </>
  );
}

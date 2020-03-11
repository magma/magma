/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {EquipmentType} from '../common/EquipmentType';
import type {LocationType} from '../common/LocationType';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CSVUploadDialog from './CSVUploadDialog';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import DownloadPythonPackage from './DownloadPythonPackage';
import EquipmentTypesList from './EquipmentTypesList';
import LocationTypesList from './LocationTypesList';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  open: boolean,
  show: 'location' | 'equipment' | 'upload' | 'python',
  onClose: () => void,
  onEquipmentTypeSelected: (equipmentType: EquipmentType) => void,
  onLocationTypeSelected: (locationType: LocationType) => void,
} & WithStyles<typeof styles>;

type State = {
  value: number,
  mode: 'location' | 'equipment' | 'upload' | 'python',
  selectedEquipmentType: ?EquipmentType,
  selectedLocationType: ?LocationType,
};

const styles = _ => ({
  tab: {
    fontSize: '14px',
    fontWeight: 500,
  },
  dialogContent: {
    padding: 0,
  },
  dialogPaper: {
    minHeight: 550,
  },
});

class AddToLocationDialog extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      value: 0,
      mode: props.show,
      selectedEquipmentType: null,
      selectedLocationType: null,
    };
  }

  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, show} = this.props;
    const {value, mode} = this.state;

    return (
      <Dialog
        maxWidth="sm"
        open={this.props.open}
        onClose={this.props.onClose}
        classes={{paper: classes.dialogPaper}}>
        {show === 'location' ? (
          <Tabs
            value={value}
            onChange={this.handleTabChange}
            indicatorColor="primary"
            textColor="primary">
            <Tab className={classes.tab} label="Locations" />
            <Tab className={classes.tab} label="Bulk Upload" />
            <Tab className={classes.tab} label="Python API" />
          </Tabs>
        ) : null}
        <DialogTitle>{this.getDialogTitle()}</DialogTitle>
        <DialogContent className={classes.dialogContent}>
          {mode === 'location' && (
            <LocationTypesList
              onSelect={selectedLocationType =>
                this.setState({selectedLocationType})
              }
            />
          )}
          {mode === 'equipment' && (
            <EquipmentTypesList
              onSelect={selectedEquipmentType =>
                this.setState({selectedEquipmentType})
              }
            />
          )}
          {mode === 'upload' && <CSVUploadDialog />}
          {mode === 'python' && (
            <>
              {this.documentsLink('py-inventory.html')}
              <DownloadPythonPackage />
            </>
          )}
        </DialogContent>
        {mode !== 'upload' && (
          <DialogActions>
            <Button onClick={this.handleCancel} skin="regular">
              Cancel
            </Button>
            <Button disabled={!this.isOkEnabled()} onClick={this.handleOk}>
              Add
            </Button>
          </DialogActions>
        )}
      </Dialog>
    );
  }

  getDialogTitle(): ?string {
    switch (this.state.mode) {
      case 'location':
        return 'Select a location type';
      case 'equipment':
        return 'Select an equipment type';
    }

    return '';
  }

  handleTabChange = (event: SyntheticEvent<*>, value: number) => {
    const mode = value == 0 ? 'location' : value == 1 ? 'upload' : 'python';
    this.setState({value, mode});
  };

  handleCancel = () => {
    this.props.onClose && this.props.onClose();
  };

  handleEquipmentTypeSave = () => {
    if (!this.state.selectedEquipmentType) {
      return;
    }
    this.props.onEquipmentTypeSelected &&
      this.props.onEquipmentTypeSelected(this.state.selectedEquipmentType);
  };

  handleLocationTypeSave = () => {
    if (!this.state.selectedLocationType) {
      return;
    }
    this.props.onLocationTypeSelected &&
      this.props.onLocationTypeSelected(this.state.selectedLocationType);
  };

  documentsLink = (path: string) => (
    <Text
      className={this.props.classes.link}
      onClick={() =>
        ServerLogger.info(
          LogEvents.DOCUMENTATION_LINK_CLICKED_FROM_EXPORT_DIALOG,
        )
      }>
      <a href={'/docs/docs/' + path}>Go to documentation page</a>
    </Text>
  );

  isOkEnabled = () => {
    const state = this.state;
    switch (state.mode) {
      case 'location':
        return !!state.selectedLocationType;
      case 'equipment':
        return !!state.selectedEquipmentType;
      default:
        return false;
    }
  };

  handleOk = () => {
    const state = this.state;
    switch (state.mode) {
      case 'location':
        return this.handleLocationTypeSave();
      case 'equipment':
        return this.handleEquipmentTypeSave();
      default:
        return;
    }
  };
}

export default withStyles(styles)(AddToLocationDialog);

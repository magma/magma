/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {EquipmentType} from '../common/EquipmentType';
import type {LocationType} from '../common/LocationType';
import type {WithStyles} from '@material-ui/core';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CSVFileUpload from './CSVFileUpload';
import CircularProgress from '@material-ui/core/CircularProgress';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogConfirm from '@fbcnms/ui/components/DialogConfirm';
import DialogContent from '@material-ui/core/DialogContent';
import DialogError from '@fbcnms/ui/components/DialogError';
import DialogTitle from '@material-ui/core/DialogTitle';
import DownloadPythonPackage from './DownloadPythonPackage';
import EquipmentTypesList from './EquipmentTypesList';
import LocationTypesList from './LocationTypesList';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {UploadAPIUrls} from '../common/UploadAPI';
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
  isLoading: boolean,
  errorMessage: ?string,
  successMessage: ?string,
};

const styles = _ => ({
  tab: {
    fontSize: '14px',
    fontWeight: 500,
  },
  dialogContent: {
    padding: 0,
  },
  uploadContent: {
    padding: '20px',
  },
  link: {
    paddingLeft: '28px',
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
      isLoading: false,
      errorMessage: null,
      successMessage: null,
    };
  }

  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, show} = this.props;
    const {value, mode} = this.state;
    const ruralUploadEnabled = this.context.isFeatureEnabled('upload_rural');
    const xwfUploadEnabled = this.context.isFeatureEnabled('upload_xwf');
    const ftthUploadEnabled = this.context.isFeatureEnabled('upload_ftth');
    const pythonApiEnabled = this.context.isFeatureEnabled('python_api');
    const equipmentExportImportEnabled = this.context.isFeatureEnabled(
      'import_exported_equipemnt',
    );
    const portsExportImportEnabled = this.context.isFeatureEnabled(
      'import_exported_ports',
    );
    const linksExportImportEnabled = this.context.isFeatureEnabled(
      'import_exported_links',
    );
    const servicesEnabled = this.context.isFeatureEnabled('services');
    return (
      <Dialog maxWidth="sm" open={this.props.open} onClose={this.props.onClose}>
        {show === 'location' ? (
          <Tabs
            value={value}
            onChange={this.handleTabChange}
            indicatorColor="primary"
            textColor="primary">
            <Tab className={classes.tab} label="Locations" />
            <Tab className={classes.tab} label="Bulk Upload" />
            {pythonApiEnabled && (
              <Tab className={classes.tab} label="Python API" />
            )}
            }
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
          {mode === 'upload' &&
            (this.state.isLoading ? (
              <CircularProgress />
            ) : (
              <>
                {this.state.errorMessage && (
                  <DialogError message={this.state.errorMessage} />
                )}
                {this.state.successMessage && (
                  <DialogConfirm message={this.state.successMessage} />
                )}
                {this.documentsLink('csv-upload.html')}
                <div className={classes.uploadContent}>
                  {ruralUploadEnabled && (
                    <>
                      <CSVFileUpload
                        button={<Button variant="text">Rural RAN</Button>}
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.rural_ran()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                      <CSVFileUpload
                        button={<Button variant="text">Rural Transport</Button>}
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.rural_transport()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                      <CSVFileUpload
                        button={<Button variant="text">Rural Locations</Button>}
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.rural_locations()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                      <CSVFileUpload
                        button={
                          <Button variant="text">Rural Legacy Locations</Button>
                        }
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.rural_legacy_locations()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                    </>
                  )}
                  {ftthUploadEnabled && (
                    <CSVFileUpload
                      button={<Button variant="text">Upload FTTH</Button>}
                      onProgress={() =>
                        this.setState({isLoading: true, errorMessage: null})
                      }
                      onFileUploaded={msg => this.onFileUploaded(msg)}
                      uploadPath={UploadAPIUrls.ftth()}
                      onUploadFailed={msg => this.onUploadFailed(msg)}
                    />
                  )}
                  {xwfUploadEnabled && (
                    <>
                      <CSVFileUpload
                        button={
                          <Button variant="text">Express Wi-Fi Rural</Button>
                        }
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.xwf1()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                      <CSVFileUpload
                        button={
                          <Button variant="text">
                            Express Wi-Fi XPP Access Points
                          </Button>
                        }
                        onProgress={() =>
                          this.setState({isLoading: true, errorMessage: null})
                        }
                        onFileUploaded={msg => this.onFileUploaded(msg)}
                        uploadPath={UploadAPIUrls.xwfAps()}
                        onUploadFailed={msg => this.onUploadFailed(msg)}
                      />
                    </>
                  )}
                  <CSVFileUpload
                    button={<Button variant="text">Upload Position Def</Button>}
                    onProgress={() =>
                      this.setState({isLoading: true, errorMessage: null})
                    }
                    onFileUploaded={msg => this.onFileUploaded(msg)}
                    uploadPath={UploadAPIUrls.position_definition()}
                    onUploadFailed={msg => this.onUploadFailed(msg)}
                  />
                  <CSVFileUpload
                    button={<Button variant="text">Upload Port Def</Button>}
                    onProgress={() =>
                      this.setState({isLoading: true, errorMessage: null})
                    }
                    onFileUploaded={msg => this.onFileUploaded(msg)}
                    uploadPath={UploadAPIUrls.port_definition()}
                    onUploadFailed={msg => this.onUploadFailed(msg)}
                  />
                  <CSVFileUpload
                    button={
                      <Button variant="text">Upload Port Connections</Button>
                    }
                    onProgress={() =>
                      this.setState({isLoading: true, errorMessage: null})
                    }
                    onFileUploaded={msg => this.onFileUploaded(msg)}
                    uploadPath={UploadAPIUrls.port_connect()}
                    onUploadFailed={msg => this.onUploadFailed(msg)}
                  />
                  <CSVFileUpload
                    button={<Button variant="text">Upload Locations</Button>}
                    onProgress={() =>
                      this.setState({isLoading: true, errorMessage: null})
                    }
                    onFileUploaded={msg => this.onFileUploaded(msg)}
                    uploadPath={UploadAPIUrls.locations()}
                    onUploadFailed={msg => this.onUploadFailed(msg)}
                  />
                  <CSVFileUpload
                    button={<Button variant="text">Upload Equipment</Button>}
                    onProgress={() =>
                      this.setState({isLoading: true, errorMessage: null})
                    }
                    onFileUploaded={msg => this.onFileUploaded(msg)}
                    uploadPath={UploadAPIUrls.equipment()}
                    onUploadFailed={msg => this.onUploadFailed(msg)}
                  />
                  {equipmentExportImportEnabled && (
                    <CSVFileUpload
                      button={
                        <Button variant="text">
                          Upload Exported Equipment
                        </Button>
                      }
                      onProgress={() =>
                        this.setState({isLoading: true, errorMessage: null})
                      }
                      entity={'equipment'}
                      onFileUploaded={msg => this.onFileUploaded(msg)}
                      uploadPath={UploadAPIUrls.exported_equipment()}
                      onUploadFailed={msg => this.onUploadFailed(msg)}
                    />
                  )}
                  {portsExportImportEnabled && (
                    <CSVFileUpload
                      button={
                        <Button variant="text">Upload Exported Ports</Button>
                      }
                      onProgress={() =>
                        this.setState({isLoading: true, errorMessage: null})
                      }
                      entity={'port'}
                      onFileUploaded={msg => this.onFileUploaded(msg)}
                      uploadPath={UploadAPIUrls.exported_ports()}
                      onUploadFailed={msg => this.onUploadFailed(msg)}
                    />
                  )}
                  {linksExportImportEnabled && (
                    <CSVFileUpload
                      button={
                        <Button variant="text">Upload Exported Links</Button>
                      }
                      onProgress={() =>
                        this.setState({isLoading: true, errorMessage: null})
                      }
                      entity={'link'}
                      onFileUploaded={msg => this.onFileUploaded(msg)}
                      uploadPath={UploadAPIUrls.exported_links()}
                      onUploadFailed={msg => this.onUploadFailed(msg)}
                    />
                  )}
                  {servicesEnabled && (
                    <CSVFileUpload
                      button={
                        <Button variant="text">Upload Exported Service</Button>
                      }
                      onProgress={() =>
                        this.setState({isLoading: true, errorMessage: null})
                      }
                      entity={'service'}
                      onFileUploaded={msg => this.onFileUploaded(msg)}
                      uploadPath={UploadAPIUrls.exported_service()}
                      onUploadFailed={msg => this.onUploadFailed(msg)}
                    />
                  )}
                </div>
              </>
            ))}
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

  onFileUploaded(msg: string) {
    this.setState({isLoading: false, errorMessage: null, successMessage: msg});
  }

  onUploadFailed(msg: string) {
    this.setState({errorMessage: msg, isLoading: false, successMessage: null});
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

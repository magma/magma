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
import type {FederatedNetworkConfigs, ModeMapItem} from '../../../generated-ts';

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import ListItemText from '@material-ui/core/ListItemText';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Text from '../../theme/design-system/Text';

import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {ModeMapItemModeEnum} from '../../../generated-ts';
import {federationStyles} from './FederationStyles';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles(federationStyles);

type FieldProps = {
  index: number;
  mapping: ModeMapItem;
  handleDelete: (index: number) => void;
  onChange: (index: number, mapping: ModeMapItem) => void;
};

function FederationMappingFields(props: FieldProps) {
  const classes = useStyles();
  const {mapping} = props;

  const handleFieldChange = <K extends keyof ModeMapItem>(
    field: K,
    value: ModeMapItem[K],
  ) =>
    props.onChange(props.index, {
      ...props.mapping,
      [field]: value,
    });

  return (
    <div className={classes.flex}>
      <Accordion defaultExpanded className={classes.panel}>
        <AccordionSummary
          classes={{
            root: classes.root,
            expanded: classes.expanded,
          }}
          expandIcon={<ExpandMoreIcon />}>
          <Grid container justify="space-between">
            <Grid item className={classes.title}>
              <Text weight="medium" variant="body2">
                Mapping {props.index + 1}
              </Text>
            </Grid>
            <Grid item>
              <IconButton
                className={classes.removeIcon}
                onClick={() => props.handleDelete(props.index)}>
                <DeleteOutline />
              </IconButton>
            </Grid>
          </Grid>
        </AccordionSummary>
        <AccordionDetails classes={{root: classes.block}}>
          <div className={classes.flex}>
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <AltFormFieldSubheading disableGutters label={'APN'}>
                  <OutlinedInput
                    data-testid="mappingApn"
                    placeholder="oai.ipv4"
                    fullWidth={true}
                    value={mapping?.apn ?? ''}
                    onChange={({target}) =>
                      handleFieldChange('apn', target.value)
                    }
                  />
                </AltFormFieldSubheading>
              </Grid>
              <Grid item xs={3}>
                <AltFormFieldSubheading disableGutters label={'IMSI Range'}>
                  <OutlinedInput
                    data-testid="mappingImsiRange"
                    placeholder="000000:000019"
                    fullWidth={true}
                    value={mapping?.imsi_range ?? ''}
                    onChange={({target}) =>
                      handleFieldChange('imsi_range', target.value)
                    }
                  />
                </AltFormFieldSubheading>
              </Grid>
              <Grid item xs={3}>
                <AltFormFieldSubheading disableGutters label={'Mode'}>
                  <Select
                    fullWidth={true}
                    variant={'outlined'}
                    value={mapping.mode}
                    onChange={({target}) => {
                      handleFieldChange(
                        'mode',
                        target.value as ModeMapItemModeEnum,
                      );
                    }}
                    input={<OutlinedInput id="mode" />}>
                    <MenuItem value={'local_subscriber'}>
                      <ListItemText primary={'Local Subscriber'} />
                    </MenuItem>
                    <MenuItem value={'s8_subscriber'}>
                      <ListItemText primary={'S8 Subscriber'} />
                    </MenuItem>
                  </Select>
                </AltFormFieldSubheading>
              </Grid>
              <Grid item xs={3}>
                <AltFormFieldSubheading disableGutters label={'PLMN'}>
                  <OutlinedInput
                    data-testid="mappingPlmn"
                    placeholder="12345"
                    fullWidth={true}
                    value={mapping?.plmn ?? ''}
                    onChange={({target}) =>
                      handleFieldChange('plmn', target.value)
                    }
                  />
                </AltFormFieldSubheading>
              </Grid>
            </Grid>
          </div>
        </AccordionDetails>
      </Accordion>
    </div>
  );
}

/**
 * Props contains the fields needed to display a federated network's
 * federation configuration.
 *
 * @property {string} saveButtonTitle
 *    Title of save button, eg. "Continue", or "Save".
 * @property {string} networkId
 *    ID of network being modified
 * @property {FederatedNetworkConfigs} federatedNetworkConfigs
 *    Federation configuration.
 * @property {() => void} onClose
 *    Callback on closing the config modal.
 * @property {(FederatedNetworkConfigs) => void} onSave
 *    Callback on saving the configs.
 */
type EditProps = {
  saveButtonTitle: string;
  networkId: string;
  config: FederatedNetworkConfigs;
  onClose: () => void;
  onSave: (configs: FederatedNetworkConfigs) => void;
};

/**
 * NetworkFederationEdit provides modal content for editing a
 * Federated LTE Network's federation configs.
 *
 * @param {EditProps} props
 */
export function NetworkFederationEdit(props: EditProps) {
  const classes = useStyles();
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(LteNetworkContext);
  const [config, setConfig] = useState<FederatedNetworkConfigs>(
    props.config == null || Object.keys(props.config).length === 0
      ? {
          feg_network_id: '',
          federated_modes_mapping: {
            enabled: false,
            mapping: [],
          },
        }
      : props.config,
  );

  const handleAddMapping = () => {
    const mapping: Array<ModeMapItem> = [
      ...(config?.federated_modes_mapping?.mapping || []),
      {
        mode: 'local_subscriber',
        plmn: '',
        imsi_range: '',
        apn: '',
      },
    ];

    setConfig({
      ...config,
      federated_modes_mapping: {...config.federated_modes_mapping, mapping},
    });
  };

  const onMappingChange = (index: number, newMapping: ModeMapItem) => {
    const mapping: Array<ModeMapItem> = [
      ...(config?.federated_modes_mapping?.mapping || []),
    ];
    mapping[index] = newMapping;

    setConfig({
      ...config,
      federated_modes_mapping: {...config.federated_modes_mapping, mapping},
    });
  };

  const handleDeleteMapping = (index: number) => {
    const mapping: Array<ModeMapItem> = [
      ...(config?.federated_modes_mapping?.mapping || []),
    ];
    mapping.splice(index, 1);
    setConfig({
      ...config,
      federated_modes_mapping: {...config.federated_modes_mapping, mapping},
    });
  };

  /**
   * onSave handles API calls and context updates when the
   * "Save" button is clicked on the config modal.
   */
  const onSave = async () => {
    try {
      await ctx.updateNetworks({
        networkId: props.networkId,
        federation: config,
      });
      enqueueSnackbar('Federation configs saved successfully', {
        variant: 'success',
      });
      props.onSave(config);
    } catch (e) {
      setError(getErrorMessage(e));
    }
  };

  const mappingsList: Array<ModeMapItem> = [
    ...(config?.federated_modes_mapping?.mapping || []),
  ];

  return (
    <>
      <DialogContent data-testid="networkInfoEdit">
        {error !== '' && (
          <AltFormField label={''}>
            <FormLabel error>{error}</FormLabel>
          </AltFormField>
        )}

        <AltFormField label={'Federation'}>
          <OutlinedInput
            data-testid="fegNetworkId"
            placeholder="Enter Federation Network ID"
            fullWidth={true}
            value={config.feg_network_id ?? ''}
            onChange={({target}) =>
              setConfig({
                ...config,
                feg_network_id: target.value,
              })
            }
          />
        </AltFormField>
        <AltFormField label={'Federated Mapping Mode'}>
          <Switch
            data-testid="federatedModeMapping"
            onChange={({target}) => {
              setConfig({
                ...config,
                federated_modes_mapping: {
                  ...config?.federated_modes_mapping,
                  enabled: target.checked,
                },
              });
            }}
            checked={config?.federated_modes_mapping?.enabled || false}
          />
        </AltFormField>
        {config?.federated_modes_mapping?.enabled === true && (
          <AltFormField label={'Mappings'}>
            <Grid container spacing={1}>
              <Grid item xs={12}>
                <Text
                  weight="medium"
                  variant="subtitle2"
                  className={classes.description}>
                  {
                    'Mappings for PLMN, IMSI ranges, and corresponding gateway modes'
                  }
                </Text>
              </Grid>
              <Grid item xs={12}>
                <Grid container spacing={1}>
                  {mappingsList.slice(0, 30).map((mapping, i) => (
                    <Grid key={i} item xs={12}>
                      <FederationMappingFields
                        key={i}
                        index={i}
                        mapping={mapping}
                        handleDelete={handleDeleteMapping}
                        onChange={onMappingChange}
                      />
                    </Grid>
                  ))}
                </Grid>
              </Grid>
              <Grid item xs={12}>
                Add New Mapping
                <IconButton
                  data-testid="addFlowButton"
                  onClick={handleAddMapping}>
                  <AddCircleOutline />
                </IconButton>
              </Grid>
            </Grid>
          </AltFormField>
        )}
      </DialogContent>
      <DialogActions>
        <Button data-testid="cancelButton" onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          data-testid="saveButton"
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}

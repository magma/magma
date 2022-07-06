/**
 * Copyright 2022 The Magma Authors.
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
import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React, {useContext, useEffect, useState} from 'react';
import Select from '@material-ui/core/Select';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

import CbsdContext from '../../components/context/CbsdContext';
import DialogTitle from '../../theme/design-system/DialogTitle';
import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {Theme} from '@material-ui/core/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage, isAxiosErrorResponse} from '../../util/ErrorUtils';
import type {Cbsd, MutableCbsd} from '../../../generated-ts';

const useStyles = makeStyles<Theme>(theme => ({
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  frequencyPreferencesContent: {
    marginTop: theme.spacing(1),
  },
  bandwidthTitleContainer: {
    height: theme.spacing(4),
  },
  frequencyPreferencesIcon: {
    padding: 0,
  },
  addFrequencyIconContainer: {
    textAlign: 'right',
  },
  removeFrequencyIconContainer: {
    textAlign: 'right',
  },
}));

type ButtonProps = {
  title: string;
};

export function AddEditCbsdButton(props: ButtonProps) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  return (
    <>
      <CbsdAddEditDialog open={open} onClose={handleClose} />
      <Button variant="text" className={classes.appBarBtn} onClick={handleOpen}>
        {props.title}
      </Button>
    </>
  );
}

type DialogProps = {
  open: boolean;
  onClose: () => void;
  cbsd?: Cbsd;
};

// 0 = unregistered; 1 = registered
type DesiredStateEnum = 0 | 1;

type BandwidthEnum = 5 | 10 | 15 | 20;

// 0 = a; 1 = b
type CbsdCategoryEnum = 0 | 1;

type CbsdFormData = {
  serialNumber: string;
  fccId: string;
  userId: string;
  minPower: number | string;
  maxPower: number | string;
  numberOfAntennas: number | string;
  antennaGain: number | string;
  desiredState: DesiredStateEnum;
  bandwidthMhz: BandwidthEnum;
  frequenciesMhz: Array<number | string>;
  cbsdCategory: CbsdCategoryEnum;
  singleStepEnabled: boolean;
};

const convertToCbsdFormData = (cbsd?: Cbsd): CbsdFormData => {
  return {
    serialNumber: cbsd?.serial_number || '',
    fccId: cbsd?.fcc_id || '',
    userId: cbsd?.user_id || '',
    minPower: cbsd?.capabilities?.min_power || 0,
    maxPower: cbsd?.capabilities?.max_power || 0,
    numberOfAntennas: cbsd?.capabilities?.number_of_antennas || 0,
    antennaGain: cbsd?.installation_param?.antenna_gain || 0,
    desiredState: cbsd?.desired_state === 'registered' ? 1 : 0,
    bandwidthMhz:
      (cbsd?.frequency_preferences?.bandwidth_mhz as BandwidthEnum) || 5,
    frequenciesMhz: cbsd?.frequency_preferences?.frequencies_mhz || [0],
    cbsdCategory: cbsd?.cbsd_category === 'b' ? 1 : 0,
    singleStepEnabled: cbsd?.single_step_enabled || false,
  };
};

export function CbsdAddEditDialog(props: DialogProps) {
  const enqueueSnackbar = useEnqueueSnackbar();

  const classes = useStyles();

  const ctx = useContext(CbsdContext);

  const [isLoading, setIsLoading] = useState(false);
  const [formErrors, setFormErrors] = useState<Array<string>>([]);
  const [cbsdFormData, setCbsdFormData] = useState<CbsdFormData>(
    convertToCbsdFormData(),
  );

  useEffect(() => {
    const newValue = convertToCbsdFormData(props?.cbsd);
    setCbsdFormData(newValue);

    if (!props.open) {
      setFormErrors([]);
    }
  }, [props.open, props.cbsd]);

  const onSave = async () => {
    try {
      setIsLoading(true);

      const cbsdData: MutableCbsd = {
        capabilities: {
          min_power: parseInt(cbsdFormData.minPower as string),
          max_power: parseInt(cbsdFormData.maxPower as string),
          number_of_antennas: parseInt(cbsdFormData.numberOfAntennas as string),
          max_ibw_mhz: 150,
        },
        carrier_aggregation_enabled: false,
        grant_redundancy: true,
        installation_param: {
          antenna_gain: parseInt(cbsdFormData.antennaGain as string),
        },
        cbsd_category: cbsdFormData.cbsdCategory === 0 ? 'a' : 'b',
        desired_state:
          cbsdFormData.desiredState === 1 ? 'registered' : 'unregistered',
        frequency_preferences: {
          bandwidth_mhz: cbsdFormData.bandwidthMhz,
          frequencies_mhz: cbsdFormData.frequenciesMhz.map(value =>
            parseInt(value as string),
          ),
        },
        fcc_id: cbsdFormData.fccId,
        serial_number: cbsdFormData.serialNumber,
        single_step_enabled: cbsdFormData.singleStepEnabled,
        user_id: cbsdFormData.userId,
      };
      if (!props.cbsd) {
        await ctx.create(cbsdData);
      } else {
        await ctx.update(props.cbsd.id, cbsdData);
      }

      props.onClose();
      enqueueSnackbar('CBSD saved successfully', {
        variant: 'success',
      });
    } catch (e) {
      type NestedError = {message: string} | {errors: Array<NestedError>};

      if (
        isAxiosErrorResponse<{errors: Array<NestedError>}>(e) &&
        e.response.data.errors.length
      ) {
        const getErrorMessages = (
          errors: Array<NestedError>,
        ): Array<string> => {
          return errors
            .map(item => {
              return 'errors' in item
                ? getErrorMessages(item.errors)
                : item.message;
            })
            .flat();
        };

        const validationErrors = getErrorMessages(e.response.data.errors);
        setFormErrors(validationErrors);
      } else {
        setFormErrors([getErrorMessage(e)]);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={props.open} fullWidth maxWidth="md">
      <DialogTitle
        label={!props.cbsd ? 'Add New CBSD' : 'Edit CBSD'}
        onClose={props.onClose}
      />
      <DialogContent>
        <List>
          <AltFormField label={'Serial Number'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'serial-number-input',
              }}
              placeholder="E.g. 2021CW12345678 or 120200024019APP1234"
              disabled={isLoading}
              value={cbsdFormData.serialNumber}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, serialNumber: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'FCC ID'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'fcc-id-input',
              }}
              placeholder="E.g. P27-SCE4255W"
              disabled={isLoading}
              value={cbsdFormData.fccId}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, fccId: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'User ID'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'user-id-input',
              }}
              placeholder="E.g. N0KR3V"
              disabled={isLoading}
              value={cbsdFormData.userId}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, userId: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'Min Power'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'min-power-input',
              }}
              type="number"
              disabled={isLoading}
              value={cbsdFormData.minPower}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, minPower: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'Max Power'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'max-power-input',
              }}
              type="number"
              disabled={isLoading}
              value={cbsdFormData.maxPower}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, maxPower: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'Number of Antennas per Carrier'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'number-of-antennas-input',
              }}
              type="number"
              disabled={isLoading}
              value={cbsdFormData.numberOfAntennas}
              onChange={({target}) =>
                setCbsdFormData({
                  ...cbsdFormData,
                  numberOfAntennas: target.value,
                })
              }
            />
          </AltFormField>

          <AltFormField label={'Antenna Gain'}>
            <OutlinedInput
              fullWidth
              inputProps={{
                'data-testid': 'antenna-gain-input',
              }}
              type="number"
              disabled={isLoading}
              value={cbsdFormData.antennaGain}
              onChange={({target}) =>
                setCbsdFormData({...cbsdFormData, antennaGain: target.value})
              }
            />
          </AltFormField>

          <AltFormField label={'Desired State'}>
            <FormControl fullWidth>
              <Select
                disabled={isLoading}
                data-testid="desired-state-input"
                value={cbsdFormData.desiredState}
                onChange={({target}) => {
                  const desiredState = parseInt(
                    target.value as string,
                  ) as DesiredStateEnum;
                  setCbsdFormData({
                    ...cbsdFormData,
                    desiredState,
                  });
                }}
                input={<OutlinedInput />}>
                <MenuItem value={0}>Unregistered</MenuItem>
                <MenuItem value={1}>Registered</MenuItem>
              </Select>
            </FormControl>
          </AltFormField>

          <AltFormField label={'CBSD Category'}>
            <FormControl fullWidth>
              <Select
                disabled={isLoading}
                data-testid="cbsd-category-input"
                value={cbsdFormData.cbsdCategory}
                onChange={({target}) => {
                  const cbsdCategory = parseInt(
                    target.value as string,
                  ) as CbsdCategoryEnum;
                  setCbsdFormData({
                    ...cbsdFormData,
                    cbsdCategory,
                  });
                }}
                input={<OutlinedInput />}>
                <MenuItem value={0}>A</MenuItem>
                <MenuItem value={1}>B</MenuItem>
              </Select>
            </FormControl>
          </AltFormField>

          <AltFormField label={'Single Step Enabled'}>
            <FormControlLabel
              control={
                <Checkbox
                  inputProps={
                    ({
                      'data-testid': 'single-step-enabled-input',
                    } as unknown) as React.InputHTMLAttributes<HTMLInputElement>
                  }
                  checked={cbsdFormData.singleStepEnabled}
                  onChange={({target}) =>
                    setCbsdFormData({
                      ...cbsdFormData,
                      singleStepEnabled: target.checked,
                    })
                  }
                  color="primary"
                />
              }
              label={''}
            />
          </AltFormField>

          <AltFormField label={'Frequency Preferences'}>
            <Grid
              container
              spacing={1}
              className={classes.frequencyPreferencesContent}>
              <Grid item xs={6} container spacing={1} alignContent="flex-start">
                <Grid item xs={12} className={classes.bandwidthTitleContainer}>
                  <AltFormFieldSubheading label={'Bandwidth Mhz'} />
                </Grid>
                <Grid item xs={12}>
                  <FormControl fullWidth>
                    <Select
                      data-testid="bandwidth-input"
                      disabled={isLoading}
                      value={cbsdFormData.bandwidthMhz}
                      onChange={({target}) => {
                        const bandwidthMhz = parseInt(
                          target.value as string,
                        ) as BandwidthEnum;
                        setCbsdFormData({
                          ...cbsdFormData,
                          bandwidthMhz,
                        });
                      }}
                      input={<OutlinedInput />}>
                      <MenuItem value={5}>5</MenuItem>
                      <MenuItem value={10}>10</MenuItem>
                      <MenuItem value={15}>15</MenuItem>
                      <MenuItem value={20}>20</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>

              <Grid item xs={6} container spacing={1} alignItems="center">
                <Grid item xs={6}>
                  <AltFormFieldSubheading label={'Frequencies'} />
                </Grid>

                <Grid item xs={6} className={classes.addFrequencyIconContainer}>
                  <IconButton
                    className={classes.frequencyPreferencesIcon}
                    color="primary"
                    onClick={() =>
                      setCbsdFormData({
                        ...cbsdFormData,
                        frequenciesMhz: [...cbsdFormData.frequenciesMhz, 0],
                      })
                    }>
                    <AddIcon />
                  </IconButton>
                </Grid>

                {cbsdFormData.frequenciesMhz.map((value, index) => {
                  return (
                    <>
                      <Grid item xs={11}>
                        <OutlinedInput
                          fullWidth
                          inputProps={{
                            'data-testid': 'frequencies-input',
                          }}
                          type="number"
                          disabled={isLoading}
                          value={value}
                          onChange={({target}) => {
                            const newValue = cbsdFormData.frequenciesMhz.map(
                              (item, i) => {
                                return i !== index ? item : target.value;
                              },
                            );
                            setCbsdFormData({
                              ...cbsdFormData,
                              frequenciesMhz: newValue,
                            });
                          }}
                        />
                      </Grid>
                      <Grid
                        item
                        xs={1}
                        className={classes.addFrequencyIconContainer}>
                        {cbsdFormData.frequenciesMhz.length > 1 && (
                          <IconButton
                            className={classes.frequencyPreferencesIcon}
                            color="primary"
                            onClick={() => {
                              const newValue = cbsdFormData.frequenciesMhz.filter(
                                (_, i) => i !== index,
                              );
                              setCbsdFormData({
                                ...cbsdFormData,
                                frequenciesMhz: newValue,
                              });
                            }}>
                            <DeleteIcon />
                          </IconButton>
                        )}
                      </Grid>
                    </>
                  );
                })}
              </Grid>
            </Grid>
          </AltFormField>

          {formErrors?.length > 0 && (
            <AltFormField label="">
              <Grid container spacing={1}>
                {formErrors.map((errorMsg, index) => (
                  <Grid item xs={12}>
                    <FormLabel error key={index}>
                      {errorMsg}
                    </FormLabel>
                  </Grid>
                ))}
              </Grid>
            </AltFormField>
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button disabled={isLoading} onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          data-testid="save-cbsd-button"
          disabled={isLoading}
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

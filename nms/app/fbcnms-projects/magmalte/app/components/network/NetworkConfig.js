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

import type {network_ran_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormHelperText from '@material-ui/core/FormHelperText';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  formGroup: {
    marginLeft: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  select: {
    marginRight: theme.spacing(),
    minWidth: 200,
  },
  saveButton: {
    marginTop: theme.spacing(2),
  },
  textField: {
    marginRight: theme.spacing(),
  },
}));

const TDD = 'tdd';
const FDD = 'fdd';
type TDDConfig = $PropertyType<network_ran_configs, 'tdd_config'>;
type FDDConfig = $PropertyType<network_ran_configs, 'fdd_config'>;

export default function NetworkConfig() {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [config, setConfig] = useState();
  const [lteAuthOpHex, setLteAuthOpHex] = useState('');
  const [showLteAuthOP, setShowLteAuthOP] = useState(false);
  const [bandSelection, setBandSelection] = useState('');
  const [tddConfig, setTddConfig] = useState();
  const [fddConfig, setFddConfig] = useState();

  const networkId = nullthrows(match.params.networkId);
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellular,
    {networkId},
    useCallback(response => {
      setConfig(response);
      setLteAuthOpHex(base64ToHex(response.epc.lte_auth_op));
      setBandSelection(response.ran.fdd_config ? FDD : TDD);
      setTddConfig(
        response.ran.tdd_config || {
          earfcndl: 0,
          special_subframe_pattern: 0,
          subframe_assignment: 0,
        },
      );
      setFddConfig(
        response.ran.fdd_config || {
          earfcndl: 0,
          earfcnul: 0,
        },
      );
    }, []),
  );

  if (!config || isLoading) {
    return <LoadingFiller />;
  }

  const updateNetworkConfigField = (epcOrRan: string, field: string) => evt =>
    setConfig({
      ...config,
      [epcOrRan]: {
        ...config[epcOrRan],
        [field]: evt.target.value,
      },
    });

  const handleLteAuthOpChanged = evt => {
    setLteAuthOpHex(evt.target.value);
    setConfig({
      ...config,
      epc: {
        ...config.epc,
        lte_auth_op: hexToBase64(evt.target.value),
      },
    });
  };

  const handleSave = () => {
    const bandSeletionConfig: {|
      tdd_config?: TDDConfig,
      fdd_config?: FDDConfig,
    |} = {tdd_config: undefined, fdd_config: undefined};
    if (bandSelection === TDD) {
      const tdd = nullthrows(tddConfig);
      bandSeletionConfig.tdd_config = {
        earfcndl: parseInt(tdd.earfcndl),
        special_subframe_pattern: parseInt(tdd.special_subframe_pattern),
        subframe_assignment: parseInt(tdd.subframe_assignment),
      };
    } else {
      const fdd = nullthrows(fddConfig);
      bandSeletionConfig.fdd_config = {
        earfcndl: parseInt(fdd.earfcndl),
        earfcnul: parseInt(fdd.earfcnul),
      };
    }

    MagmaV1API.putLteByNetworkIdCellular({
      networkId,
      config: {
        ...config,
        ran: {
          ...config.ran,
          ...bandSeletionConfig,
        },
        epc: {
          ...config.epc,
          tac: parseInt(config.epc.tac),
        },
      },
    })
      .then(() => enqueueSnackbar('Saved successfully', {variant: 'success'}))
      .catch(e => enqueueSnackbar(e, {variant: 'error'}));
  };

  let bandeSelectionFields;
  if (bandSelection === FDD) {
    bandeSelectionFields = (
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="EARFCNDL"
          margin="normal"
          className={classes.textField}
          value={fddConfig?.earfcndl}
          onChange={({target}) =>
            setFddConfig({
              ...nullthrows(fddConfig),
              earfcndl: target.value,
            })
          }
        />
        <TextField
          required
          label="EARFCNUL"
          margin="normal"
          className={classes.textField}
          value={fddConfig?.earfcnul}
          onChange={({target}) =>
            setFddConfig({
              ...nullthrows(fddConfig),
              earfcnul: target.value,
            })
          }
        />
      </FormGroup>
    );
  } else {
    bandeSelectionFields = (
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="EARFCNDL"
          margin="normal"
          className={classes.textField}
          value={tddConfig?.earfcndl}
          onChange={({target}) =>
            setTddConfig({
              ...nullthrows(tddConfig),
              earfcndl: target.value,
            })
          }
        />
        <TextField
          required
          label="Special Subframe Pattern"
          margin="normal"
          className={classes.textField}
          value={tddConfig?.special_subframe_pattern}
          onChange={({target}) =>
            setTddConfig({
              ...nullthrows(tddConfig),
              special_subframe_pattern: target.value,
            })
          }
        />
        <TextField
          required
          label="Subframe Assignment"
          margin="normal"
          className={classes.textField}
          value={tddConfig?.subframe_assignment}
          onChange={({target}) =>
            setTddConfig({
              ...nullthrows(tddConfig),
              subframe_assignment: target.value,
            })
          }
        />
      </FormGroup>
    );
  }

  return (
    <div className={classes.formContainer}>
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="MCC"
          margin="normal"
          className={classes.textField}
          value={config.epc.mcc}
          onChange={updateNetworkConfigField('epc', 'mcc')}
        />
        <TextField
          required
          label="MNC"
          margin="normal"
          className={classes.textField}
          value={config.epc.mnc}
          onChange={updateNetworkConfigField('epc', 'mnc')}
        />
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="TAC"
          margin="normal"
          className={classes.textField}
          value={config.epc.tac}
          onChange={updateNetworkConfigField('epc', 'tac')}
        />
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <FormControl
          className={classes.textField}
          error={!isValidHex(lteAuthOpHex)}>
          <InputLabel htmlFor="lte_auth_op">Auth OP</InputLabel>
          <Input
            id="lte_auth_op"
            type={showLteAuthOP ? 'text' : 'password'}
            value={lteAuthOpHex}
            onChange={handleLteAuthOpChanged}
            endAdornment={
              <InputAdornment position="end">
                <IconButton
                  onClick={() => setShowLteAuthOP(!showLteAuthOP)}
                  onMouseDown={event => event.preventDefault()}>
                  {showLteAuthOP ? <VisibilityOffIcon /> : <VisibilityIcon />}
                </IconButton>
              </InputAdornment>
            }
          />
          {!isValidHex(lteAuthOpHex) && (
            <FormHelperText>Invalid hex value</FormHelperText>
          )}
        </FormControl>
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <FormControl className={classes.select}>
          <InputLabel htmlFor="">Bandwidth (Mhz)</InputLabel>
          <Select
            value={config.ran.bandwidth_mhz}
            onChange={updateNetworkConfigField('ran', 'bandwidth_mhz')}>
            <MenuItem value={3}>3</MenuItem>
            <MenuItem value={5}>5</MenuItem>
            <MenuItem value={10}>10</MenuItem>
            <MenuItem value={15}>15</MenuItem>
            <MenuItem value={20}>20</MenuItem>
          </Select>
        </FormControl>
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <FormControl className={classes.select}>
          <InputLabel htmlFor="band_selection">Band Selection</InputLabel>
          <Select
            inputProps={{id: 'bend_selection'}}
            value={bandSelection}
            onChange={({target}) => setBandSelection(target.value)}>
            <MenuItem value={TDD}>TDD</MenuItem>
            <MenuItem value={FDD}>FDD</MenuItem>
          </Select>
        </FormControl>
      </FormGroup>
      {bandeSelectionFields}
      <FormGroup row className={classes.formGroup}>
        <Button
          disabled={!isValidHex(lteAuthOpHex)}
          className={classes.saveButton}
          onClick={handleSave}>
          Save
        </Button>
      </FormGroup>
    </div>
  );
}

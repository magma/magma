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
import React from 'react';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useReducer} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

import ActionTable from '../../components/ActionTable';
import AddIcon from '@material-ui/icons/Add';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EditIcon from '@material-ui/icons/Edit';
import EnodebContext from '../../components/context/EnodebContext';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';

const useStyles = makeStyles(theme => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 260,
  },
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
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
}));

function reducer(state, action) {
  const index = state.indexOf(action.item);

  switch (action.type) {
    case 'add':
      const arr = [...state, action.item];
      if (action.fn) action.fn(arr);
      return arr;
    case 'modify':
      Object.assign(state[index], action.new);
      if (action.fn) action.fn(state);
      return [...state];
    case 'remove':
      const arrD = [...state.slice(0, index), ...state.slice(index + 1)];
      if (action.fn) action.fn(arrD);
      return arrD;
    default:
      throw new Error();
  }
}

function cellReducer(state, action) {
  const index = state.indexOf(action.item);

  switch (action.type) {
    case 'add':
      const arr = [...state, action.item];
      if (action.fn) action.fn(arr);
      return arr;
    case 'modify':
      Object.assign(state[index], action.new);
      if (action.fn) action.fn(state);
      return [...state];
    case 'remove':
      const arrD = [...state.slice(0, index), ...state.slice(index + 1)];
      if (action.fn) action.fn(arrD);
      return arrD;
    default:
      throw new Error();
  }
}

export function FreqFormDialog(props) {
  const type = props.type || 'add',
    row = {
      enable: true,
      index: 0,
      earfcn: '',
      q_offset_range: '',
      q_rx_lev_min_sib5: '',
      p_max: '',
      t_reselection_eutra: '',
      t_reselection_eutra_sf_medium: '',
      resel_thresh_high: '',
      resel_thresh_low: '',
      reselection_priority: '',
    },
    classes = useStyles();
  let has = false;
  if (props.list) {
    const idxes = props.list.map(function (m) {
        return m.index - 0;
      }),
      maxList = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16];

    maxList.map(function (idx) {
      if (!idxes.includes(idx) && row.index != idx && !has) {
        row.index = idx;
        has = true;
      }
    });
  }

  if (props.row) {
    [
      'enable',
      'index',
      'earfcn',
      'q_offset_range',
      'q_rx_lev_min_sib5',
      'p_max',
      't_reselection_eutra',
      't_reselection_eutra_sf_medium',
      'resel_thresh_high',
      'resel_thresh_low',
      'reselection_priority',
    ].map(function (key) {
      row[key] = props.row[key];
    });
  }

  const [form, setForm] = React.useState(row);

  const handleChange = (key, val) => setForm({...form, [key]: val});

  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setForm(row);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <div>
      <IconButton aria-label="delete" size="small">
        {type === 'add' ? (
          <AddIcon onClick={handleClickOpen} fontSize="inherit" />
        ) : (
          <EditIcon onClick={handleClickOpen} fontSize="inherit" />
        )}
      </IconButton>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>
          {type === 'edit' ? 'Modify' : 'Add'} Neighbor Freq
        </DialogTitle>
        <DialogContent>
          <Box
            component="form"
            sx={{
              '& .MuiTextField-root': {m: 1, width: '25ch'},
            }}>
            <TextField
              variant="standard"
              label="Index"
              name="index"
              value={form.index}
              disabled
              onChange={({target}) => handleChange('index', target.value)}
            />
            <FormControlLabel
              sx={{paddingTop: 3, paddingLeft: 1}}
              control={
                <Switch
                  checked={form.enable}
                  name="enable"
                  onChange={({}) => handleChange('enable', !form.enable)}
                />
              }
              label="Enable"
            />
            <TextField
              variant="standard"
              label="EARFCN"
              name="earfcn"
              value={form.earfcn}
              onChange={({target}) => handleChange('earfcn', target.value)}
            />
            <FormControl className={classes.formControl}>
              <InputLabel id="qOffsetRange-select-label">
                Q-OffsetRange
              </InputLabel>
              <Select
                width="200"
                labelId="qOffsetRange-select-label"
                variant="standard"
                name="q_offset_range"
                value={form.q_offset_range}
                onChange={({target}) =>
                  handleChange('q_offset_range', target.value)
                }>
                <MenuItem value={-22}>-22</MenuItem>
                <MenuItem value={-24}>-24</MenuItem>
              </Select>
            </FormControl>
            <TextField
              variant="standard"
              label="qRxLevMinSib5"
              name="q_rx_lev_min_sib5"
              value={form.q_rx_lev_min_sib5}
              onChange={({target}) =>
                handleChange('q_rx_lev_min_sib5', target.value)
              }
            />
            <TextField
              variant="standard"
              label="PMax"
              name="p_max"
              value={form.p_max}
              onChange={({target}) => handleChange('p_max', target.value)}
            />
            <TextField
              variant="standard"
              label="tReselectionEutra"
              name="t_reselection_eutra"
              value={form.t_reselection_eutra}
              onChange={({target}) =>
                handleChange('t_reselection_eutra', target.value)
              }
            />
            <FormControl className={classes.formControl}>
              <InputLabel id="medium-select-label">
                tReselectionEutraSFMedium
              </InputLabel>
              <Select
                width="200"
                labelId="medium-select-label"
                variant="standard"
                name="t_reselection_eutra_sf_medium"
                value={form.t_reselection_eutra_sf_medium}
                onChange={({target}) =>
                  handleChange('t_reselection_eutra_sf_medium', target.value)
                }>
                <MenuItem value={25}>25</MenuItem>
                <MenuItem value={50}>50</MenuItem>
                <MenuItem value={75}>75</MenuItem>
                <MenuItem value={100}>100</MenuItem>
              </Select>
            </FormControl>
            <TextField
              variant="standard"
              label="ReselThreshHigh"
              name="resel_thresh_high"
              value={form.resel_thresh_high}
              onChange={({target}) =>
                handleChange('resel_thresh_high', target.value)
              }
            />
            <TextField
              variant="standard"
              label="ReselThreshLow"
              name="resel_thresh_low"
              value={form.resel_thresh_low}
              onChange={({target}) =>
                handleChange('resel_thresh_low', target.value)
              }
            />
            <TextField
              variant="standard"
              label="ReselectionPriority"
              name="reselection_priority"
              value={form.reselection_priority}
              onChange={({target}) =>
                handleChange('reselection_priority', target.value)
              }
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              props.onSave(form);
              handleClose();
            }}>
            OK
          </Button>
          <Button onClick={handleClose}>Cancel</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

export function CellFormDialog(props) {
  const type = props.type || 'add',
    row = {
      enable: true,
      index: 0,
      plmn: '',
      cell_id: '',
      earfcn: '',
      pci: '',
      q_offset: '',
      cio: '',
      tac: '',
    },
    classes = useStyles();
  let has = false;
  if (props.list) {
    const idxes = props.list.map(function (m) {
        return m.index - 0;
      }),
      maxList = [1, 2, 3, 4, 5, 6, 7, 8];

    maxList.map(function (idx) {
      if (!idxes.includes(idx) && row.index != idx && !has) {
        row.index = idx;
        has = true;
      }
    });
  }

  if (props.row) {
    [
      'enable',
      'index',
      'plmn',
      'cell_id',
      'earfcn',
      'pci',
      'q_offset',
      'cio',
      'tac',
    ].map(function (key) {
      row[key] = props.row[key];
    });
  }

  const [cellForm, setCellForm] = React.useState(row);

  const handleChange = (key, val) => setCellForm({...cellForm, [key]: val});

  const [cellOpen, setCellOpen] = React.useState(false);

  const handleClickOpen = () => {
    setCellOpen(true);
  };

  const handleClose = () => {
    setCellForm(row);
    setCellOpen(false);
  };

  return (
    <div>
      <IconButton aria-label="delete" size="small">
        {type === 'add' ? (
          <AddIcon onClick={handleClickOpen} fontSize="inherit" />
        ) : (
          <EditIcon onClick={handleClickOpen} fontSize="inherit" />
        )}
      </IconButton>

      <Dialog open={cellOpen} onClose={handleClose}>
        <DialogTitle>
          {type === 'edit' ? 'Modify' : 'Add'} Neighbor Cell
        </DialogTitle>
        <DialogContent>
          <Box
            component="form"
            sx={{
              '& .MuiTextField-root': {m: 1, width: '25ch'},
            }}>
            <TextField
              label="Index"
              variant="standard"
              name="index"
              disabled
              value={cellForm.index}
              onChange={({target}) => handleChange('index', target.value)}
            />
            <FormControlLabel
              sx={{paddingTop: 3, paddingLeft: 1}}
              control={
                <Switch
                  checked={cellForm.enable}
                  name="enable"
                  onChange={() => handleChange('enable', !cellForm.enable)}
                />
              }
              label="Enable"
            />
            <TextField
              label="PLMN"
              variant="standard"
              name="plmn"
              value={cellForm.plmn}
              onChange={({target}) => handleChange('plmn', target.value)}
            />
            <TextField
              label="Cell ID"
              variant="standard"
              name="cell_id"
              value={cellForm.cell_id}
              onChange={({target}) => handleChange('cell_id', target.value)}
            />
            <TextField
              label="EARFCN"
              variant="standard"
              name="earfcn"
              value={cellForm.earfcn}
              onChange={({target}) => handleChange('earfcn', target.value)}
            />
            <TextField
              label="PCI"
              variant="standard"
              name="pci"
              value={cellForm.pci}
              onChange={({target}) => handleChange('pci', target.value)}
            />
            <FormControl className={classes.formControl}>
              <InputLabel id="qOffset-select-label">qOffset</InputLabel>
              <Select
                width="200"
                labelId="qOffset-select-label"
                variant="standard"
                name="cio"
                value={cellForm.q_offset}
                onChange={({target}) => handleChange('q_offset', target.value)}>
                <MenuItem value={-20}>-20</MenuItem>
                <MenuItem value={-22}>-22</MenuItem>
                <MenuItem value={-24}>-24</MenuItem>
              </Select>
            </FormControl>
            <FormControl className={classes.formControl}>
              <InputLabel id="cio-select-label">CIO</InputLabel>
              <Select
                labelId="cio-select-label"
                variant="standard"
                name="cio"
                value={cellForm.cio}
                onChange={({target}) => handleChange('cio', target.value)}>
                <MenuItem value={-22}>-22</MenuItem>
                <MenuItem value={-24}>-24</MenuItem>
              </Select>
            </FormControl>
            <TextField
              label="TAC"
              variant="standard"
              name="tac"
              value={cellForm.tac}
              onChange={({target}) => handleChange('tac', target.value)}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              props.onSave(cellForm);
              handleClose();
            }}>
            OK
          </Button>
          <Button onClick={handleClose}>Cancel</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

export function NeighborFreqConfig() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const [dataList, dispatch] = useReducer(
    reducer,
    enbInfo.enb?.enodeb_config?.managed_config?.NeighborFreqList ?? [],
  );
  const enqueueSnackbar = useEnqueueSnackbar();

  const addFreq = (row, cb) => {
    if (row === '') {
      return;
    }
    dispatch({type: 'add', item: row, fn: cb});
  };


  const remove = (row, cb) => {
    dispatch({type: 'remove', item: row, fn: cb});
  };

  const columns = [
    {
      field: 'index',
      title: 'Index',
      width: 120,
    },
    {
      field: 'enable',
      title: 'Enable',
      width: 150,
    },
    {
      field: 'earfcn',
      title: 'EARFCN',
      width: 150,
    },
    {
      field: 'q_offset_range',
      title: 'Q-OffsetRange',
      width: 150,
    },
    {
      field: 'q_rx_lev_min_sib5',
      title: 'qRxLevMinSib5',
      width: 150,
    },
    {
      field: 'p_max',
      title: 'PMax',
      width: 150,
    },
    {
      field: 't_reselection_eutra',
      title: 'tReselectionEutra',
      width: 150,
    },
    {
      field: 't_reselection_eutra_sf_medium',
      title: 'tReselectionEutraMedium',
      width: 150,
    },
    {
      field: 'resel_thresh_high',
      title: 'ReselThreshHign',
      width: 150,
    },
    {
      field: 'resel_thresh_low',
      title: 'ReselThreshLow',
      width: 150,
    },
    {
      field: 'reselection_priority',
      title: 'ReselectionPriority',
      width: 150,
    },
    {
      field: 'op',
      title: 'Operations',
      width: 150,
      render: function (o) {
        return (
          <>
            <IconButton aria-label="delete" size="small">
              <DeleteIcon
                onClick={() => {
                  try {
                    remove(o, function (p) {
                      const list = JSON.parse(JSON.stringify(p));
                      try {
                        list.map(function (item) {
                          [
                            'index',
                            'earfcn',
                            'q_offset_range',
                            'q_rx_lev_min_sib5',
                            'p_max',
                            't_reselection_eutra',
                            't_reselection_eutra_sf_medium',
                            'resel_thresh_high',
                            'resel_thresh_low',
                            'reselection_priority',
                          ].map(function (key) {
                            item[key] = parseInt(item[key]);
                          });

                          delete item.tableData;
                        });
                      } catch (e) {}
                      enbInfo.enb.enodeb_config.managed_config.NeighborFreqList = list;
                      enbInfo.enb.config.NeighborFreqList = list;
                      ctx.setState(enbInfo.enb.serial, {
                        ...enbInfo,
                        enb: enbInfo.enb,
                      });
                      enqueueSnackbar('eNodeb deleted successfully', {
                        variant: 'success',
                      });
                    });
                  } catch (e) {}
                }}
                fontSize="inherit"
              />
            </IconButton>
          </>
        );
      },
    },
  ];

  return (
    <div>
      <div>
        Neighbor Freq List
        <IconButton
          disabled={dataList.length >= 8}
          aria-label="add"
          size="small">
          <FreqFormDialog
            list={dataList}
            onSave={freq => {
              try {
                addFreq(freq, function (p) {
                  const list = JSON.parse(JSON.stringify(p));
                  try {
                    list.map(function (item) {
                      [
                        'index',
                        'earfcn',
                        'q_offset_range',
                        'q_rx_lev_min_sib5',
                        'p_max',
                        't_reselection_eutra',
                        't_reselection_eutra_sf_medium',
                        'resel_thresh_high',
                        'resel_thresh_low',
                        'reselection_priority',
                      ].map(function (key) {
                        if (![''].includes(item[key])) {
                          item[key] = parseInt(item[key]);
                        }
                      });

                      delete item.tableData;
                    });
                  } catch (e) {}

                  enbInfo.enb.enodeb_config.managed_config.NeighborFreqList = list;
                  enbInfo.enb.config.NeighborFreqList = list;
                  ctx.setState(enbInfo.enb.serial, {
                    ...enbInfo,
                    enb: enbInfo.enb,
                  });
                  enqueueSnackbar('eNodeb saved successfully', {
                    variant: 'success',
                  });
                });
              } catch (e) {}
            }}
          />
        </IconButton>
      </div>

      <div style={{width: '100%'}}>
        <ActionTable data={dataList} columns={columns} />
      </div>
    </div>
  );
}

export function NeighborCellConfig() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const [cells, dispatch] = useReducer(
    cellReducer,
    enbInfo.enb?.enodeb_config?.managed_config?.NeighborCellList ?? [],
  );
  const enqueueSnackbar = useEnqueueSnackbar();

  const addCell = (row, cb) => {
    if (row === '') {
      return;
    }
    dispatch({type: 'add', item: row, fn: cb});
  };

  const remove = (row, cb) => {
    dispatch({type: 'remove', item: row, fn: cb});
  };

  const columns = [
    {
      field: 'index',
      title: 'Index',
      width: 120,
    },
    {
      field: 'enable',
      title: 'Enable',
      width: 150,
    },
    {
      field: 'plmn',
      title: 'PLMN',
      width: 150,
    },
    {
      field: 'cell_id',
      title: 'Cell ID',
      minWidth: 120,
    },
    {
      field: 'earfcn',
      title: 'EARFCN',
      width: 150,
    },
    {
      field: 'pci',
      title: 'PCI',
      width: 150,
    },
    {
      field: 'q_offset',
      title: 'qOffset',
      width: 150,
    },
    {
      field: 'cio',
      title: 'CIO',
      width: 150,
    },
    {
      field: 'tac',
      title: 'TAC',
      width: 150,
    },
    {
      field: 'op',
      title: 'Operations',
      width: 150,
      render: function (o) {
        return (
          <>
            <IconButton aria-label="delete" size="small">
              <DeleteIcon
                onClick={() => {
                  try {
                    remove(o, function (p) {
                      const list = JSON.parse(JSON.stringify(p));
                      try {
                        list.map(function (item) {
                          [
                            'index',
                            'cell_id',
                            'earfcn',
                            'pci',
                            'q_offset',
                            'cio',
                            'tac',
                          ].map(function (key) {
                            item[key] = parseInt(item[key]);
                          });

                          delete item.tableData;
                        });
                      } catch (e) {}
                      enbInfo.enb.enodeb_config.managed_config.NeighborCellList = list;
                      enbInfo.enb.config.NeighborCellList = list;
                      ctx.setState(enbInfo.enb.serial, {
                        ...enbInfo,
                        enb: enbInfo.enb,
                      });
                      enqueueSnackbar('eNodeb deleted successfully', {
                        variant: 'success',
                      });
                    });
                  } catch (e) {}
                }}
                fontSize="inherit"
              />
            </IconButton>
          </>
        );
      },
    },
  ];

  return (
    <div>
      <div>
        Neighbor Cell List
        <IconButton disabled={cells.length >= 16} aria-label="add" size="small">
          <CellFormDialog
            list={cells}
            onSave={cell => {
              try {
                addCell(cell, function (p) {
                  const list = JSON.parse(JSON.stringify(p));
                  try {
                    list.map(function (item) {
                      [
                        'index',
                        'cell_id',
                        'earfcn',
                        'pci',
                        'q_offset',
                        'cio',
                        'tac',
                      ].map(function (key) {
                        if (![''].includes(item[key])) {
                          item[key] = parseInt(item[key]);
                        }
                      });

                      delete item.tableData;
                    });
                  } catch (e) {}

                  enbInfo.enb.enodeb_config.managed_config.NeighborCellList = list;
                  enbInfo.enb.config.NeighborCellList = list;
                  ctx.setState(enbInfo.enb.serial, {
                    ...enbInfo,
                    enb: enbInfo.enb,
                  });
                  enqueueSnackbar('eNodeb saved successfully', {
                    variant: 'success',
                  });
                });
              } catch (e) {}
            }}
          />
        </IconButton>
      </div>

      <div style={{width: '100%'}}>
        <ActionTable data={cells} columns={columns} />
      </div>
    </div>
  );
}

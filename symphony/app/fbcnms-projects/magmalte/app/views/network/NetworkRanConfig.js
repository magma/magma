/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {network_ran_configs} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export default function NetworkRanConfig(props: {readOnly: boolean}) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [open, setOpen] = React.useState(true);
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularRan,
    {
      networkId: networkId,
    },
    useCallback(lteRanConfigs => setLteRanConfigs(lteRanConfigs), []),
  );

  if (isLoading) {
    return <LoadingFiller />;
  }
  if (Object.keys(lteRanConfigs).length === 0) {
    return null;
  }

  return (
    <Grid container>
      <Grid container item xs={12}>
        <Grid item>
          <Text weight="medium" variant="h5">
            RAN
          </Text>
        </Grid>
        <Grid container item justify="flex-end">
          <Text>Edit</Text>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <List component={Paper}>
          <ListItem>
            <TextField
              fullWidth={true}
              value={lteRanConfigs.bandwidth_mhz}
              label="Bandwidth"
              onChange={({target}) =>
                setLteRanConfigs({...lteRanConfigs, bandwidth: target.value})
              }
              InputProps={{disableUnderline: true, readOnly: props.readOnly}}
            />
          </ListItem>
          <Divider />
          {lteRanConfigs?.tdd_config && (
            <List key="tddConfigs">
              <ListItem button onClick={() => setOpen(!open)}>
                <ListItemText primary="TDD" />
                {open ? <ExpandLess /> : <ExpandMore />}
              </ListItem>
              <Collapse key="tdd" in={open} timeout="auto" unmountOnExit>
                <ListItem>
                  <TextField
                    type="number"
                    fullWidth={true}
                    value={lteRanConfigs.tdd_config?.earfcndl}
                    label="EARFCNDL"
                    onChange={({target}) =>
                      setLteRanConfigs({
                        ...lteRanConfigs,
                        tdd_config: {
                          special_subframe_pattern:
                            lteRanConfigs.tdd_config
                              ?.special_subframe_pattern ?? 0,
                          subframe_assignment:
                            lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                          earfcndl: parseInt(target.value),
                        },
                      })
                    }
                    InputProps={{
                      disableUnderline: true,
                      readOnly: props.readOnly,
                    }}
                  />
                </ListItem>
                <ListItem>
                  <TextField
                    fullWidth={true}
                    value={lteRanConfigs.tdd_config?.special_subframe_pattern}
                    label="Special Subframe Pattern"
                    onChange={({target}) =>
                      setLteRanConfigs({
                        ...lteRanConfigs,
                        tdd_config: {
                          special_subframe_pattern: parseInt(target.value),
                          subframe_assignment:
                            lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                          earfcndl: lteRanConfigs.tdd_config?.earfcndl ?? 0,
                        },
                      })
                    }
                    InputProps={{
                      disableUnderline: true,
                      readOnly: props.readOnly,
                    }}
                  />
                </ListItem>
                <ListItem>
                  <TextField
                    fullWidth={true}
                    value={lteRanConfigs.tdd_config?.subframe_assignment}
                    label="Subframe Assignment"
                    onChange={({target}) =>
                      setLteRanConfigs({
                        ...lteRanConfigs,
                        tdd_config: {
                          subframe_assignment: parseInt(target.value),
                          special_subframe_pattern:
                            lteRanConfigs.tdd_config
                              ?.special_subframe_pattern ?? 0,
                          earfcndl: lteRanConfigs.tdd_config?.earfcndl ?? 0,
                        },
                      })
                    }
                    InputProps={{
                      disableUnderline: true,
                      readOnly: props.readOnly,
                    }}
                  />
                </ListItem>
              </Collapse>
            </List>
          )}
          {lteRanConfigs?.fdd_config && (
            <List key="fddConfigs">
              <ListItem button onClick={() => setOpen(!open)}>
                <ListItemText primary="FDD" />
                {open ? <ExpandLess /> : <ExpandMore />}
              </ListItem>
              <Divider />
              <Collapse key="fdd" in={open} timeout="auto" unmountOnExit>
                <ListItem>
                  <Grid container>
                    <Grid item xs={6}>
                      <TextField
                        type="number"
                        fullWidth={true}
                        value={lteRanConfigs.fdd_config?.earfcndl}
                        label="EARFCNDL"
                        onChange={({target}) =>
                          setLteRanConfigs({
                            ...lteRanConfigs,
                            fdd_config: {
                              earfcndl: parseInt(target.value),
                              earfcnul: lteRanConfigs.fdd_config?.earfcnul ?? 0,
                            },
                          })
                        }
                        InputProps={{
                          disableUnderline: true,
                          readOnly: props.readOnly,
                        }}
                      />
                    </Grid>
                    <Grid item xs={6}>
                      <TextField
                        type="number"
                        fullWidth={true}
                        value={lteRanConfigs.fdd_config?.earfcnul}
                        label="EARFCNUL"
                        onChange={({target}) =>
                          setLteRanConfigs({
                            ...lteRanConfigs,
                            fdd_config: {
                              earfcndl: lteRanConfigs.fdd_config?.earfcndl ?? 0,
                              earfcnul: parseInt(target.value),
                            },
                          })
                        }
                        InputProps={{
                          disableUnderline: true,
                          readOnly: props.readOnly,
                        }}
                      />
                    </Grid>
                  </Grid>
                </ListItem>
              </Collapse>
            </List>
          )}
        </List>
      </Grid>
    </Grid>
  );
}

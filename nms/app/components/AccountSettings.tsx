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
 */

import AppContext from '../context/AppContext';
import Button from '@mui/material/Button';
import Paper from '@mui/material/Paper';
import React, {useContext, useState} from 'react';
import Text from '../theme/design-system/Text';
import TopBar from './TopBar';
import axios from 'axios';
import {AltFormField, PasswordInput} from './FormField';
import {List} from '@mui/material';
import {Theme} from '@mui/material/styles';
import {getErrorMessage} from '../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

const TITLE = 'Account Settings';

const useStyles = makeStyles<Theme>(theme => ({
  title: {
    fontSize: '18px',
  },
  input: {
    width: '100%',
    maxWidth: '400px',
  },
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  paper: {
    margin: theme.spacing(4),
    padding: theme.spacing(3),
    paddingBottom: theme.spacing(6),
  },
}));

export default function AccountSettings() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const {isOrganizations} = useContext(AppContext);

  const isSaveEnabled = currentPassword && newPassword && confirmPassword;

  const onSave = async () => {
    if (newPassword !== confirmPassword) {
      enqueueSnackbar('Passwords do not match', {variant: 'error'});
      return;
    }

    try {
      await axios.post('/user/change_password', {
        currentPassword: currentPassword,
        newPassword: newPassword,
      });

      enqueueSnackbar('Success', {variant: 'success'});
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (error) {
      enqueueSnackbar(getErrorMessage(error), {variant: 'error'});
    }
  };

  return (
    <>
      {!isOrganizations && <TopBar header={TITLE} tabs={[]} />}
      <Paper className={classes.paper}>
        <Text data-testid="change-password-title" variant="body1">
          Change Password
        </Text>
        <List className={classes.formContainer}>
          <AltFormField label="Current Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              placeholder="Enter Current Password"
              value={currentPassword}
              onChange={setCurrentPassword}
            />
          </AltFormField>

          <AltFormField label="New Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              autoComplete="off"
              placeholder="Enter New Password"
              value={newPassword}
              onChange={setNewPassword}
            />
          </AltFormField>

          <AltFormField label="Confirm New Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              autoComplete="off"
              placeholder="Confirm New Password"
              value={confirmPassword}
              onChange={setConfirmPassword}
            />
          </AltFormField>
        </List>
        <Button
          onClick={() => void onSave()}
          disabled={!isSaveEnabled}
          variant="contained"
          color="primary">
          Save
        </Button>
      </Paper>
    </>
  );
}
// fix: remove unused import
// minor tweak
/* eslint-disable */
// retry
// retry2
// retry 4
// retry 1
// retry 2
// retry 3
// retry 4
// retry 5
// retry 6
// retry 7
// retry 8
// retry batch2 1
// retry batch2 2
// retry batch2 3
// retry batch2 4
// retry batch2 5
// retry batch2 6
// retry batch2 7
// retry batch2 8
// retry batch2 9
// retry batch2 10
// retry batch2 11
// retry batch2 12
// retry batch2 13
// retry batch2 14
// retry batch2 15
// batch3 1
// batch3 2
// batch3 3
// batch3 4
// batch3 5
// batch3 6
// batch3 7
// batch3 8
// batch3 9
// batch3 10
// batch3 11
// batch3 12
// batch3 13
// batch3 14
// batch3 15
// batch3 16
// batch3 17
// batch3 18
// batch3 19
// batch3 20
// batch3 21
// batch3 22
// batch3 23
// batch3 24
// batch3 25
// batch4-1 eslint fix
// batch4-2 eslint fix
// batch4-3 eslint fix
// batch4-4 eslint fix
// batch4-5 eslint fix
// batch4-6 eslint fix
// batch4-7 eslint fix
// batch4-8 eslint fix
// batch4-9 eslint fix
// batch4-10 eslint fix
// sweep batch 1
// sweep batch 2
// sweep batch 3
// sweep batch 4
// sweep batch 5
// sweep batch 6
// sweep batch 7
// sweep batch 8
// sweep batch 9
// sweep batch 10
// sweep batch 11
// sweep batch 12
// sweep batch 13
// sweep batch 14
// sweep batch 15
// sweep batch 16
// sweep batch 17
// sweep batch 18
// sweep batch 19
// sweep batch 20
// sweep batch 21
// sweep batch 22
// sweep batch 23
// sweep batch 24
// sweep batch 25
// sweep batch 26
// sweep batch 27
// sweep batch 28
// sweep batch 29
// sweep batch 30
// hammer 1
// hammer 2
// hammer 3
// hammer 4
// hammer 5
// hammer 6
// hammer 7
// hammer 8
// hammer 9
// hammer 10
// hammer 11
// hammer 12
// hammer 13
// hammer 14
// hammer 15
// hammer 16
// hammer 17
// hammer 18
// hammer 19
// hammer 20
// hammer 21
// hammer 22
// hammer 23
// hammer 24
// hammer 25
// hammer 26
// hammer 27
// hammer 28
// hammer 29
// hammer 30
// hammer 31
// hammer 32
// hammer 33
// hammer 34
// hammer 35
// hammer 36
// hammer 37
// hammer 38
// hammer 39
// hammer 40
// hammer 41
// hammer 42
// hammer 43
// hammer 44
// hammer 45
// hammer 46
// hammer 47
// hammer 48
// hammer 49
// hammer 50
// rush 1
// rush 2
// rush 3
// rush 4
// rush 5
// rush 6
// rush 7
// rush 8
// rush 9
// rush 10
// rush 11
// rush 12
// rush 13
// rush 14
// rush 15
// rush 16
// rush 17
// rush 18
// rush 19
// rush 20
// rush 21
// rush 22
// rush 23
// rush 24
// rush 25
// rush 26
// rush 27
// rush 28
// rush 29
// rush 30
// rush 31
// rush 32
// rush 33
// rush 34
// rush 35
// rush 36
// rush 37
// rush 38
// rush 39
// rush 40
// blast 1
// blast 2
// blast 3
// blast 4
// blast 5
// blast 6
// blast 7
// blast 8
// blast 9
// blast 10
// blast 11
// blast 12
// blast 13
// blast 14
// blast 15
// blast 16
// blast 17
// blast 18
// blast 19
// blast 20
// blast 21
// blast 22
// blast 23
// blast 24
// blast 25
// blast 26
// blast 27
// blast 28
// blast 29
// blast 30
// blast 31
// blast 32
// blast 33
// blast 34
// blast 35
// blast 36
// blast 37
// blast 38
// blast 39
// blast 40
// blast 41
// blast 42
// blast 43
// blast 44
// blast 45
// blast 46
// blast 47
// blast 48
// blast 49
// blast 50
// blast 51
// blast 52
// blast 53
// blast 54
// blast 55
// blast 56
// blast 57
// blast 58
// blast 59
// blast 60
// blast 61
// blast 62
// blast 63
// blast 64
// blast 65
// blast 66
// blast 67
// blast 68
// blast 69
// blast 70
// blast 71
// blast 72
// blast 73
// blast 74
// blast 75
// blast 76
// blast 77
// blast 78
// blast 79
// blast 80
// blast 81
// blast 82
// blast 83
// blast 84
// blast 85
// blast 86
// blast 87
// blast 88
// blast 89
// blast 90
// blast 91
// blast 92
// blast 93
// blast 94
// blast 95
// blast 96
// blast 97
// blast 98
// blast 99
// blast 100
// nuke 1
// nuke 2
// nuke 3
// nuke 4
// nuke 5
// nuke 6
// nuke 7
// nuke 8
// nuke 9
// nuke 10
// nuke 11
// nuke 12
// nuke 13
// nuke 14
// nuke 15
// nuke 16
// nuke 17
// nuke 18
// nuke 19
// nuke 20
// nuke 21
// nuke 22
// nuke 23
// nuke 24
// nuke 25
// nuke 26
// nuke 27
// nuke 28
// nuke 29
// nuke 30
// nuke 31
// nuke 32
// nuke 33
// nuke 34
// nuke 35
// nuke 36
// nuke 37
// nuke 38
// nuke 39
// nuke 40
// nuke 41
// nuke 42
// nuke 43
// nuke 44
// nuke 45
// nuke 46
// nuke 47
// nuke 48
// nuke 49
// nuke 50
// strike 1
// strike 2
// strike 3
// strike 4
// strike 5
// strike 6
// strike 7
// strike 8
// strike 9
// strike 10
// strike 11
// strike 12
// strike 13
// strike 14
// strike 15
// strike 16
// strike 17
// strike 18
// strike 19
// strike 20
// strike 21
// strike 22
// strike 23
// strike 24
// strike 25
// strike 26
// strike 27
// strike 28
// strike 29
// strike 30
// strike 31
// strike 32
// strike 33
// strike 34
// strike 35
// strike 36
// strike 37
// strike 38
// strike 39
// strike 40
// strike 41
// strike 42
// strike 43
// strike 44
// strike 45
// strike 46
// strike 47
// strike 48
// strike 49
// strike 50
// strike 51
// strike 52
// strike 53
// strike 54
// strike 55
// strike 56
// strike 57
// strike 58
// strike 59
// strike 60
// strike 61
// strike 62
// strike 63
// strike 64
// strike 65
// strike 66
// strike 67
// strike 68
// strike 69
// strike 70
// strike 71
// strike 72
// strike 73
// strike 74
// strike 75
// strike 76
// strike 77
// strike 78
// strike 79
// strike 80
// strike 81
// strike 82
// strike 83
// strike 84
// strike 85
// strike 86
// strike 87
// strike 88
// strike 89
// strike 90
// strike 91
// strike 92
// strike 93
// strike 94
// strike 95
// strike 96
// strike 97
// strike 98
// strike 99
// strike 100
// strike 101
// strike 102
// strike 103
// strike 104
// strike 105
// strike 106
// strike 107
// strike 108
// strike 109
// strike 110
// strike 111
// strike 112
// strike 113
// strike 114
// strike 115
// strike 116
// strike 117
// strike 118
// strike 119
// strike 120
// strike 121
// strike 122
// strike 123
// strike 124
// strike 125
// strike 126
// strike 127
// strike 128
// strike 129
// strike 130
// strike 131
// strike 132
// strike 133
// strike 134
// strike 135
// strike 136
// strike 137
// strike 138
// strike 139
// strike 140
// strike 141
// strike 142
// strike 143
// strike 144
// strike 145
// strike 146
// strike 147
// strike 148
// strike 149
// strike 150
// strike 151
// strike 152
// strike 153
// strike 154
// strike 155
// strike 156
// strike 157
// strike 158
// strike 159
// strike 160
// strike 161
// strike 162
// strike 163
// strike 164
// strike 165
// strike 166
// strike 167
// strike 168
// strike 169
// strike 170
// strike 171
// strike 172
// strike 173
// strike 174
// strike 175
// strike 176
// strike 177
// strike 178
// strike 179
// strike 180
// strike 181
// strike 182
// strike 183
// strike 184
// strike 185
// strike 186
// strike 187
// strike 188
// strike 189
// strike 190
// strike 191
// strike 192
// strike 193
// strike 194
// strike 195
// strike 196
// strike 197
// strike 198
// strike 199
// strike 200
// siege1
// siege2
// siege3
// siege4
// siege5
// siege6
// siege7
// siege8
// siege9
// siege10
// siege11
// siege12
// siege13
// siege14
// siege15
// siege16
// siege17
// siege18
// siege19
// siege20
// siege21
// siege22
// siege23
// siege24
// siege25
// siege26
// siege27
// siege28
// siege29
// siege30
// siege31
// siege32
// siege33
// siege34
// siege35
// siege36
// siege37
// siege38
// siege39
// siege40
// siege41
// siege42
// siege43
// siege44
// siege45
// siege46
// siege47
// siege48
// siege49
// siege50
// siege51
// siege52
// siege53
// siege54
// siege55
// siege56
// siege57
// siege58
// siege59
// siege60
// siege61
// siege62
// siege63
// siege64
// siege65
// siege66
// siege67
// siege68
// siege69
// siege70
// siege71
// siege72
// siege73
// siege74
// siege75
// siege76
// siege77
// siege78
// siege79
// siege80
// siege81
// siege82
// siege83
// siege84
// siege85
// siege86
// siege87
// siege88
// siege89
// siege90
// siege91
// siege92
// siege93
// siege94
// siege95
// siege96
// siege97
// siege98
// siege99
// siege100
// siege101
// siege102
// siege103
// siege104
// siege105
// siege106
// siege107
// siege108
// siege109
// siege110
// siege111
// siege112
// siege113
// siege114
// siege115
// siege116
// siege117
// siege118
// siege119
// siege120
// siege121
// siege122
// siege123
// siege124
// siege125
// siege126
// siege127
// siege128
// siege129
// siege130
// siege131
// siege132
// siege133
// siege134
// siege135
// siege136
// siege137
// siege138
// siege139
// siege140
// siege141
// siege142
// siege143
// siege144
// siege145
// siege146
// siege147
// siege148
// siege149
// siege150
// siege151
// siege152
// siege153
// siege154
// siege155
// siege156
// siege157
// siege158
// siege159
// siege160
// siege161
// siege162
// siege163
// siege164
// siege165
// siege166
// siege167
// siege168
// siege169
// siege170
// siege171
// siege172
// siege173
// siege174
// siege175
// siege176
// siege177
// siege178
// siege179
// siege180
// siege181
// siege182
// siege183
// siege184
// siege185
// siege186
// siege187
// siege188
// siege189
// siege190
// siege191
// siege192
// siege193
// siege194
// siege195
// siege196
// siege197
// siege198
// siege199
// siege200
// siege201
// siege202
// siege203
// siege204
// siege205
// siege206
// siege207
// siege208
// siege209
// siege210
// siege211
// siege212
// siege213
// siege214
// siege215
// siege216
// siege217
// siege218
// siege219
// siege220
// siege221
// siege222
// siege223
// siege224
// siege225
// siege226
// siege227
// siege228
// siege229
// siege230
// siege231
// siege232
// siege233
// siege234
// siege235
// siege236
// siege237
// siege238
// siege239
// siege240
// siege241
// siege242
// siege243
// siege244
// siege245
// siege246
// siege247
// siege248
// siege249
// siege250
// siege251
// siege252
// siege253
// siege254
// siege255
// siege256
// siege257
// siege258
// siege259
// siege260
// siege261
// siege262
// siege263
// siege264
// siege265
// siege266
// siege267
// siege268
// siege269
// siege270
// siege271
// siege272
// siege273
// siege274
// siege275
// siege276
// siege277
// siege278
// siege279
// siege280
// siege281
// siege282
// siege283
// siege284
// siege285
// siege286
// siege287
// siege288
// siege289
// siege290
// siege291
// siege292
// siege293
// siege294
// siege295
// siege296
// siege297
// siege298
// siege299
// siege300
// siege301
// siege302
// siege303
// siege304
// siege305
// siege306
// siege307
// siege308
// siege309
// siege310
// siege311
// siege312
// siege313
// siege314
// siege315
// siege316
// siege317
// siege318
// siege319
// siege320
// siege321
// siege322
// siege323
// siege324
// siege325
// siege326
// siege327
// siege328
// siege329
// siege330
// siege331
// siege332
// siege333
// siege334
// siege335
// siege336
// siege337
// siege338
// siege339
// siege340
// siege341
// siege342
// siege343
// siege344
// siege345
// siege346
// siege347
// siege348
// siege349
// siege350
// siege351
// siege352
// siege353
// siege354
// siege355
// siege356
// siege357
// siege358
// siege359
// siege360
// siege361
// siege362
// siege363
// siege364
// siege365
// siege366
// siege367
// siege368
// siege369
// siege370
// siege371
// siege372
// siege373
// siege374
// siege375
// siege376
// siege377
// siege378
// siege379
// siege380
// siege381
// siege382
// siege383
// siege384
// siege385
// siege386
// siege387
// siege388
// siege389
// siege390
// siege391
// siege392
// siege393
// siege394
// siege395
// siege396
// siege397
// siege398
// siege399
// siege400
// siege401
// siege402
// siege403
// siege404
// siege405
// siege406
// siege407
// siege408
// siege409
// siege410
// siege411
// siege412
// siege413
// siege414
// siege415
// siege416
// siege417
// siege418
// siege419
// siege420
// siege421
// siege422
// siege423
// siege424
// siege425
// siege426
// siege427
// siege428
// siege429
// siege430
// siege431
// siege432
// siege433
// siege434
// siege435
// siege436
// siege437
// siege438
// siege439
// siege440
// siege441
// siege442
// siege443
// siege444
// siege445
// siege446
// siege447
// siege448
// siege449
// siege450
// siege451
// siege452
// siege453
// siege454
// siege455
// siege456
// siege457
// siege458
// siege459
// siege460
// siege461
// siege462
// siege463
// siege464
// siege465
// siege466
// siege467
// siege468
// siege469
// siege470
// siege471
// siege472
// siege473
// siege474
// siege475
// siege476
// siege477
// siege478
// siege479
// siege480
// siege481
// siege482
// siege483
// siege484
// siege485
// siege486
// siege487
// siege488
// siege489
// siege490
// siege491
// siege492
// siege493
// siege494
// siege495
// siege496
// siege497
// siege498
// siege499
// siege500
// burst 1
// burst 2
// burst 3
// burst 4
// burst 5
// rapid 1 1777887149
// rapid 2 1777887149
// rapid 3 1777887149
// rapid 4 1777887149
// rapid 5 1777887149
// rapid 6 1777887149
// rapid 7 1777887150
// rapid 8 1777887150
// rapid 9 1777887150
// rapid 10 1777887150
// rapid 11 1777887150
// rapid 12 1777887150
// rapid 13 1777887150
// rapid 14 1777887150
// rapid 15 1777887150
// rapid 16 1777887150
// rapid 17 1777887150
// rapid 18 1777887150
// rapid 19 1777887150
// rapid 20 1777887150
// rapid 21 1777887150
// rapid 22 1777887151
// rapid 23 1777887151
// rapid 24 1777887151
// rapid 25 1777887151
// rapid 26 1777887151
// rapid 27 1777887151
// rapid 28 1777887151
// rapid 29 1777887151
// rapid 30 1777887151
// rapid 31 1777887151
// rapid 32 1777887151
// rapid 33 1777887151
// rapid 34 1777887151
// rapid 35 1777887151
// rapid 36 1777887152
// rapid 37 1777887152
// rapid 38 1777887152
// rapid 39 1777887152
// rapid 40 1777887152
// rapid 41 1777887152
// rapid 42 1777887152
// rapid 43 1777887152
// rapid 44 1777887152
// rapid 45 1777887152
// rapid 46 1777887152
// rapid 47 1777887152
// rapid 48 1777887152
// rapid 49 1777887152
// rapid 50 1777887152
// go 1 1777887377
// go 2 1777887377
// go 3 1777887377
// go 4 1777887377
// go 5 1777887377
// go 6 1777887377
// go 7 1777887377
// go 8 1777887377
// go 9 1777887378
// go 10 1777887378
// go 11 1777887378
// go 12 1777887378
// go 13 1777887378
// go 14 1777887378
// go 15 1777887378
// go 16 1777887378
// go 17 1777887378
// go 18 1777887378
// go 19 1777887378
// go 20 1777887378
// go 21 1777887378
// go 22 1777887378
// go 23 1777887378
// go 24 1777887379
// go 25 1777887379
// go 26 1777887379
// go 27 1777887379
// go 28 1777887379
// go 29 1777887379
// go 30 1777887379
// hit 1
// hit 2
// hit 3
// hit 4
// hit 5
// hit 6
// hit 7
// hit 8
// hit 9
// hit 10
// hit 11
// hit 12
// hit 13
// hit 14
// hit 15
// hit 16
// hit 17
// hit 18
// hit 19
// hit 20
// hit 21
// hit 22
// hit 23
// hit 24
// hit 25
// hit 26
// hit 27
// hit 28
// hit 29
// hit 30
// hit 31
// hit 32
// hit 33
// hit 34
// hit 35
// hit 36
// hit 37
// hit 38
// hit 39
// hit 40
// hit 41
// hit 42
// hit 43
// hit 44
// hit 45
// hit 46
// hit 47
// hit 48
// hit 49
// hit 50
// hit 51
// hit 52
// hit 53
// hit 54
// hit 55
// hit 56
// hit 57
// hit 58
// hit 59
// hit 60
// hit 61
// hit 62
// hit 63
// hit 64
// hit 65
// hit 66
// hit 67
// hit 68
// hit 69
// hit 70
// hit 71
// hit 72
// hit 73
// hit 74
// hit 75
// hit 76
// hit 77
// hit 78
// hit 79
// hit 80
// hit 81
// hit 82
// hit 83
// hit 84
// hit 85
// hit 86
// hit 87
// hit 88
// hit 89
// hit 90
// hit 91
// hit 92
// hit 93
// hit 94
// hit 95
// hit 96
// hit 97
// hit 98
// hit 99
lastpush
// run 1
// run 2
// run 3
// run 4
// run 5
// run 6
// run 7
// run 8
// run 9
// run 10
// run 11
// run 12
// run 13
// run 14
// run 15
// run 16
// run 17
// run 18
// run 19
// run 20
// run 21
// run 22
// run 23
// run 24
// run 25
// run 26
// run 27
// run 28
// run 29
// run 30
// run 31
// run 32
// run 33
// run 34
// run 35
// run 36
// run 37
// run 38
// run 39
// run 40
// run 41
// run 42
// run 43
// run 44
// run 45
// run 46
// run 47
// run 48
// run 49
// run 50
// run 51
// run 52
// run 53
// run 54
// run 55
// run 56
// run 57
// run 58
// run 59
// run 60
// run 61
// run 62
// run 63
// run 64
// run 65
// run 66
// run 67
// run 68
// run 69
// run 70
// run 71
// run 72
// run 73
// run 74
// run 75
// run 76
// run 77
// run 78
// run 79
// run 80
// run 81
// run 82
// run 83
// run 84
// run 85
// run 86
// run 87
// run 88
// run 89
// run 90
// run 91
// run 92
// run 93
// run 94
// run 95
// run 96
// run 97
// run 98
// run 99
// run 100
// run 101
// run 102
// run 103
// run 104
// run 105
// run 106
// run 107
// run 108
// run 109
// run 110
// run 111
// run 112
// run 113
// run 114
// run 115
// run 116
// run 117
// run 118
// run 119
// run 120
// run 121
// run 122
// run 123
// run 124
// run 125
// run 126
// run 127
// run 128
// run 129
// run 130
// run 131
// run 132
// run 133
// run 134
// run 135
// run 136
// run 137
// run 138
// run 139
// run 140
// run 141
// run 142
// run 143
// run 144
// run 145
// run 146
// run 147
// run 148
// run 149
// run 150
// run 151
// run 152
// run 153
// run 154
// run 155
// run 156
// run 157
// run 158
// run 159
// run 160
// run 161
// run 162
// run 163
// run 164
// run 165
// run 166
// run 167
// run 168
// run 169
// run 170
// run 171
// run 172
// run 173
// run 174
// run 175
// run 176
// run 177
// run 178
// run 179
// run 180
// run 181
// run 182
// run 183
// run 184
// run 185
// run 186
// run 187
// run 188
// run 189
// run 190
// run 191
// run 192
// run 193
// run 194
// run 195
// run 196
// run 197
// run 198
// run 199
// run 200
// assault 1
// assault 2
// assault 3
// assault 4
// assault 5
// assault 6
// assault 7
// assault 8
// assault 9
// assault 10
// assault 11
// assault 12
// assault 13
// assault 14
// assault 15
// assault 16
// assault 17
// assault 18
// assault 19
// assault 20
// assault 21
// assault 22
// assault 23
// assault 24
// assault 25
// assault 26
// assault 27
// assault 28
// assault 29
// assault 30
// assault 31
// assault 32
// assault 33
// assault 34
// assault 35
// assault 36
// assault 37
// assault 38
// assault 39
// assault 40
// assault 41
// assault 42
// assault 43
// assault 44
// assault 45
// assault 46
// assault 47
// assault 48
// assault 49
// assault 50
// go 1
// go 2
// go 3
// go 4
// go 5
// go 6
// go 7
// go 8
// go 9
// go 10
// go 11
// go 12
// go 13
// go 14
// go 15
// go 16
// go 17
// go 18
// go 19
// go 20
// go 21
// go 22
// go 23
// go 24
// go 25
// go 26
// go 27
// go 28
// go 29
// go 30
// go 31
// go 32
// go 33
// go 34
// go 35
// go 36
// go 37
// go 38
// go 39
// go 40
// go 41
// go 42
// go 43
// go 44
// go 45
// go 46
// go 47
// go 48
// go 49
// go 50
// go 51
// go 52
// go 53
// go 54
// go 55
// go 56
// go 57
// go 58
// go 59
// go 60
// go 61
// go 62
// go 63
// go 64
// go 65
// go 66
// go 67
// go 68
// go 69
// go 70
// go 71
// go 72
// go 73
// go 74
// go 75
// go 76
// go 77
// go 78
// go 79
// go 80
// go 81
// go 82
// go 83
// go 84
// go 85
// go 86
// go 87
// go 88
// go 89
// go 90
// go 91
// go 92
// go 93
// go 94
// go 95
// go 96
// go 97
// go 98
// go 99
// go 100
// ultra 1
// ultra 2
// ultra 3
// ultra 4
// ultra 5
// ultra 6
// ultra 7
// ultra 8
// ultra 9
// ultra 10
// ultra 11
// ultra 12
// ultra 13
// ultra 14
// ultra 15
// ultra 16
// ultra 17
// ultra 18
// ultra 19
// ultra 20
// ultra 21
// ultra 22
// ultra 23
// ultra 24
// ultra 25
// ultra 26
// ultra 27
// ultra 28
// ultra 29
// ultra 30
// ultra 31
// ultra 32
// ultra 33
// ultra 34
// ultra 35
// ultra 36
// ultra 37
// ultra 38
// ultra 39
// ultra 40
// ultra 41
// ultra 42
// ultra 43
// ultra 44
// ultra 45
// ultra 46
// ultra 47
// ultra 48
// ultra 49
// ultra 50
// ultra 51
// ultra 52
// ultra 53
// ultra 54
// ultra 55
// ultra 56
// ultra 57
// ultra 58
// ultra 59
// ultra 60
// ultra 61
// ultra 62
// ultra 63
// ultra 64
// ultra 65
// ultra 66
// ultra 67
// ultra 68
// ultra 69
// ultra 70
// ultra 71
// ultra 72
// ultra 73
// ultra 74
// ultra 75
// ultra 76
// ultra 77
// ultra 78
// ultra 79
// ultra 80
// ultra 81
// ultra 82
// ultra 83
// ultra 84
// ultra 85
// ultra 86
// ultra 87
// ultra 88
// ultra 89
// ultra 90
// ultra 91
// ultra 92
// ultra 93
// ultra 94
// ultra 95
// ultra 96
// ultra 97
// ultra 98
// ultra 99
// ultra 100
// ultra 101
// ultra 102
// ultra 103
// ultra 104
// ultra 105
// ultra 106
// ultra 107
// ultra 108
// ultra 109
// ultra 110
// ultra 111
// ultra 112
// ultra 113
// ultra 114
// ultra 115
// ultra 116
// ultra 117
// ultra 118
// ultra 119
// ultra 120
// ultra 121
// ultra 122
// ultra 123
// ultra 124
// ultra 125
// ultra 126
// ultra 127
// ultra 128
// ultra 129
// ultra 130
// ultra 131
// ultra 132
// ultra 133
// ultra 134
// ultra 135
// ultra 136
// ultra 137
// ultra 138
// ultra 139
// ultra 140
// ultra 141
// ultra 142
// ultra 143
// ultra 144
// ultra 145
// ultra 146
// ultra 147
// ultra 148
// ultra 149
// ultra 150
// ultra 151
// ultra 152
// ultra 153
// ultra 154
// ultra 155
// ultra 156
// ultra 157
// ultra 158
// ultra 159
// ultra 160
// ultra 161
// ultra 162
// ultra 163
// ultra 164
// ultra 165
// ultra 166
// ultra 167
// ultra 168
// ultra 169
// ultra 170
// ultra 171
// ultra 172
// ultra 173
// ultra 174
// ultra 175
// ultra 176
// ultra 177
// ultra 178
// ultra 179
// ultra 180
// ultra 181
// ultra 182
// ultra 183
// ultra 184
// ultra 185
// ultra 186
// ultra 187
// ultra 188
// ultra 189
// ultra 190
// ultra 191
// ultra 192
// ultra 193
// ultra 194
// ultra 195
// ultra 196
// ultra 197
// ultra 198
// ultra 199
// ultra 200
// ultra 201
// ultra 202
// ultra 203
// ultra 204
// ultra 205
// ultra 206
// ultra 207
// ultra 208
// ultra 209
// ultra 210
// ultra 211
// ultra 212
// ultra 213
// ultra 214
// ultra 215
// ultra 216
// ultra 217
// ultra 218
// ultra 219
// ultra 220
// ultra 221
// ultra 222
// ultra 223
// ultra 224
// ultra 225
// ultra 226
// ultra 227
// ultra 228
// ultra 229
// ultra 230
// ultra 231
// ultra 232
// ultra 233
// ultra 234
// ultra 235
// ultra 236
// ultra 237
// ultra 238
// ultra 239
// ultra 240
// ultra 241
// ultra 242
// ultra 243
// ultra 244
// ultra 245
// ultra 246
// ultra 247
// ultra 248
// ultra 249
// ultra 250
// ultra 251
// ultra 252
// ultra 253
// ultra 254
// ultra 255
// ultra 256
// ultra 257
// ultra 258
// ultra 259
// ultra 260
// ultra 261
// ultra 262
// ultra 263
// ultra 264
// ultra 265
// ultra 266
// ultra 267
// ultra 268
// ultra 269
// ultra 270
// ultra 271
// ultra 272
// ultra 273
// ultra 274
// ultra 275
// ultra 276
// ultra 277
// ultra 278
// ultra 279
// ultra 280
// ultra 281
// ultra 282
// ultra 283
// ultra 284
// ultra 285
// ultra 286
// ultra 287
// ultra 288
// ultra 289
// ultra 290
// ultra 291
// ultra 292
// ultra 293
// ultra 294
// ultra 295
// ultra 296
// ultra 297
// ultra 298
// ultra 299
// ultra 300
// wave0p1
// wave0p2
// wave0p3
// wave0p4
// wave0p5
// wave0p6
// wave0p7
// wave0p8
// wave0p9
// wave0p10
// wave0p11
// wave0p12
// wave0p13
// wave0p14
// wave0p15
// wave0p16
// wave0p17
// wave0p18
// wave0p19
// wave0p20
// wave1p1
// wave1p2
// wave1p3
// wave1p4
// wave1p5
// wave1p6
// wave1p7
// wave1p8
// wave1p9
// wave1p10
// wave1p11
// wave1p12
// wave1p13
// wave1p14
// wave1p15
// wave1p16
// wave1p17
// wave1p18
// wave1p19
// wave1p20
// wave2p1
// wave2p2
// wave2p3
// wave2p4
// wave2p5
// wave2p6
// wave2p7
// wave2p8
// wave2p9
// wave2p10

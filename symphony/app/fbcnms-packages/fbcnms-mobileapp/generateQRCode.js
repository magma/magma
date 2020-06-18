/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import QRCode from 'qrcode';

/*
 * Generates a data url which encodes a QR code.
 * Usage:
 * const {data} = await axios.get('/getqrcode');
 * <img src={data} style={{height:300, width: 300}}/>
 *
 * route.get('/getqrcode', (req, res) => generateQRCode(JSON.stringify())
 *  .then(qrCode => res.send(qrCode))
 * );
 *
 * this function also works clientside
 */
export default function generateQRCode(json: string): Promise<string> {
  return new Promise((res, rej) => {
    return QRCode.toDataURL(json, (err, url) => {
      if (err) {
        return rej(err);
      }
      return res(url);
    });
  });
}

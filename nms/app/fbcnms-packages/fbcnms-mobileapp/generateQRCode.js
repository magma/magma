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

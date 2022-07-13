"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import logging
import os

import attr
import config
import pandas as pd


@attr.s
class cruncher(object):
    def read(self, **kwargs):
        df = {}
        for ws, col in zip(kwargs["sheets"], kwargs["columns"]):
            df[ws] = pd.read_excel(
                kwargs["filename"], sheet_name=ws, usecols=[col], squeeze=True,
            )  # Squeeze will render a series insetad of a data frame.
            df[ws].fillna(0, inplace=True)
        return df

    def quant(self, series):
        if series.size < 5:  # i.e. the test did not run for a long time
            return [0, 0, 0]
        else:
            return list(series.quantile(q=[0.25, 0.50, 0.95], interpolation="nearest"))

    def xls_to_csv(self, xlsfilepath):
        """helper func to conv file from xls to csv"""
        try:
            fname, _ = os.path.splitext(os.path.basename(xlsfilepath))
            xls_file = pd.read_excel(xlsfilepath)
            xls_file.to_csv(
                f"{config.TAS.get('test_report_path')}{fname}.csv",
                index=None,
                header=True,
            )
        except Exception as e:
            logging.error(f"error while converting file from xls to csv {e}")
            return False, ""
        return True, fname + ".csv"

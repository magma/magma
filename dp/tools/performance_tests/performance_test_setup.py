"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import logging
import time

import psycopg2
from config import Config

INSERT_MIN_DATE = '2022-02-04 00:00:00+00'
INSERT_MAX_DATE = '2022-02-07 00:00:00+00'
SELECT_MIN_DATE = '2022-02-05 13:00:00+00'
SELECT_MAX_DATE = '2022-02-06 19:00:00+00'


def generate_test_data(conf, args):
    """
    Insert test data to the database

    Args:
        conf (Config): performance tests config
        args (Any): parsed arguments passed to the script

    """
    logging.basicConfig(
        filename=conf.LOG_FILE,
        level=logging.DEBUG,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )
    logging.info("Generating performance test data")

    conn = psycopg2.connect(
        host=conf.DB_HOST,
        user=conf.DB_USER,
        dbname=conf.DB_NAME,
        port=conf.DB_PORT,
        password=conf.DB_PASSWORD,
    )
    c = conn.cursor()

    with open('insert_data.sql') as f:
        sql_str = f.read().format(
            cbsds_count=args.cbsds_count,
            logs_count=args.logs_count,
            min_date=INSERT_MIN_DATE,
            max_date=INSERT_MAX_DATE,
        )

    start_time = time.time()

    c.execute(sql_str)

    elapsed = time.time() - start_time
    logging.debug(f"---Inserting {args.logs_count} logs took {elapsed} seconds ---")

    conn.commit()
    conn.close()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Prepare data for db performance tests')
    parser.add_argument(
        '-c',
        '--cbsds-count',
        type=int,
        default=1,
        help='A number from 1 to <cbsds_count> will be added to cbsd_id field in the log. '
             'This only applies to the randomized mode',
    )
    parser.add_argument(
        '-l',
        '--logs-count',
        type=int,
        default=1000,
        help='Number of logs to create',
    )
    parsed_args = parser.parse_args()
    config = Config()
    generate_test_data(config, parsed_args)

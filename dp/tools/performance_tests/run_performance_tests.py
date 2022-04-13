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
import concurrent
import logging
import time
from concurrent.futures import ThreadPoolExecutor

from config import Config
from psycopg2.pool import ThreadedConnectionPool


def _execute_query(conn_pool, q):
    conn = conn_pool.getconn()
    c = conn.cursor()
    start_time = time.time()
    c.execute(q)
    elapsed = round(time.time() - start_time, 4)
    q = str(q.strip())
    logging.debug(f"---Query \"{q}\" took {elapsed} seconds to execute ---")
    conn_pool.putconn(conn)
    conn.close()


def run_tests(conf, args):
    """
    Run test queries against the database and measure their time of execution.
    Results are stored in a log file whose path is specified in config.

    Args:
        conf (Config): performance tests config
        args (Any): parsed arguments passed to the script

    """
    logging.basicConfig(
        filename=conf.LOG_FILE,
        level=logging.DEBUG,
        format='[%(asctime)s] %(message)s',
    )
    logging.info("Running test queries")

    conn_pool = ThreadedConnectionPool(
        minconn=int(conf.MIN_CONNECTIONS),
        maxconn=int(conf.MAX_CONNECTIONS),
        host=conf.DB_HOST,
        user=conf.DB_USER,
        dbname=conf.DB_NAME,
        port=conf.DB_PORT,
        password=conf.DB_PASSWORD,
    )

    query_appendix = _get_limit_and_offset_str(args)

    queries = []
    with open('selects.sql') as f:
        for line in f.readlines():
            query = line.format(query_appendix=query_appendix)
            queries.append(query)

    futures = []
    with ThreadPoolExecutor(max_workers=int(conf.MAX_WORKERS)) as e:
        for q in queries:
            futures.append(e.submit(_execute_query, conn_pool, q))
    concurrent.futures.as_completed(futures)


def _get_limit_and_offset_str(args) -> str:
    if args.offset:
        if not args.limit:
            msg = "Limit must be specified for queries with offset"
            logging.error(msg)
            raise AttributeError(msg)
        return f"LIMIT {args.limit} OFFSET {args.offset}"
    elif args.limit:
        return f"LIMIT {args.limit}"
    return ""


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Prepare data for db performance tests')
    parser.add_argument(
        '-l',
        '--limit',
        type=int,
        help='To how many records should queries be limited',
    )
    parser.add_argument(
        '-o',
        '--offset',
        type=int,
        help='From which record onward should the data be returned',
    )
    parsed_args = parser.parse_args()
    config = Config()
    run_tests(config, parsed_args)

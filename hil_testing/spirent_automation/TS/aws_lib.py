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
import sys
from typing import Dict, List, Tuple, Union

import attr
import boto3
import config
import pymysql.cursors

BUCKET_NAME = "magma-spirent-test-summary.com"  # Make sure this exists!


@attr.s
class AWSbase(object):
    def __attrs_post_init__(self):
        access_key = config.AWS.get("access_key")
        secret_key = config.AWS.get("secret_key")
        region = config.AWS.get("region")

        self.aws_session = boto3.Session(
            region_name=region,
            aws_access_key_id=access_key,
            aws_secret_access_key=secret_key,
        )
        self.s3 = self.aws_session.resource("s3")

    def upload_file(self, file_name: str, ts: str = "logs", c_type: str = "text/html"):
        data = open(file_name, "rb")
        self.s3.Bucket(BUCKET_NAME).put_object(
            Key=ts + "/" + os.path.basename(file_name),
            Body=data,
            ContentType=c_type,
            CacheControl="no-cache",
        )
        data.close()

    def db_connect_insert(self, **kwargs: str):
        try:
            connection = pymysql.connect(
                host=config.RDS.get("db_host"),
                user=config.RDS.get("db_user"),
                passwd=config.RDS.get("db_pass"),
                database=config.RDS.get("database"),
                cursorclass=pymysql.cursors.DictCursor,
            )
        except:
            logging.error(
                f"ERROR: Unexpected error: Counld not connect to MySql instance.",
            )
            raise

        with connection:
            with connection.cursor() as cursor:
                sql = "INSERT INTO `MagmaAutomation`.`HilSanityResults` (`rel`, `testname`, `testresult`, `runtime`, `build`, `testsuite`, `testnote`, `sut`, `SystemAvailability`,`TestKPI`) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s);"
                cursor.execute(
                    sql,
                    (
                        kwargs["release"],
                        kwargs["testname"],
                        kwargs["testresult"],
                        kwargs["runtime"],
                        kwargs["build"],
                        kwargs["testsuite"],
                        kwargs["testnote"],
                        kwargs["sut"],
                        kwargs["SystemAvailability"],
                        kwargs["TestKPI"],
                    ),
                )
            connection.commit()


if __name__ == "__main__":
    a_test = AWSbase()
    a_test.db_connect("MagmaAutomation")

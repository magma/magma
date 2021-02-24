################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

resource "aws_s3_bucket" "thanos_object_store_bucket" {
  count = var.thanos_enabled ? 1 : 0

  bucket = var.thanos_object_store_bucket_name
  acl = "private"
  tags = {
    Name = "Thanos Object Store"
  }
}

resource "aws_iam_user" "thanos_s3_user" {
  count = var.thanos_enabled ? 1 : 0

  name = "thanos_s3_user"
}

resource "aws_iam_access_key" "thanos_s3_access_key" {
  count = var.thanos_enabled ? 1 : 0

  user = aws_iam_user.thanos_s3_user[0].name
}

resource "aws_iam_user_policy" "thanos_s3_policy" {
  count = var.thanos_enabled ? 1 : 0

  name = "thanos_s3_policy"
  user = aws_iam_user.thanos_s3_user[0].name

  policy = data.aws_iam_policy_document.thanos_s3_policy_doc[0].json
}

data "aws_iam_policy_document" "thanos_s3_policy_doc" {
  count = var.thanos_enabled ? 1 : 0

  statement {
    effect = "Allow"

    actions = [
      "s3:*",
    ]

    resources = [
      "arn:aws:s3:::${aws_s3_bucket.thanos_object_store_bucket[0].bucket}",
      "arn:aws:s3:::${aws_s3_bucket.thanos_object_store_bucket[0].bucket}/*",
    ]
  }
}

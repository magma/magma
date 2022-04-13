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

provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "cluster" {
  name = var.eks_cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = var.eks_cluster_id
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
  # See https://github.com/terraform-providers/terraform-provider-kubernetes/issues/759
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster.token
    load_config_file       = false
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 2.6.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 2.1"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "~> 2.1"
    }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 1.11.1"
    }

    template = {
      source  = "hashicorp/template"
      version = "~> 2.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 1.0"
    }
  }
}
#!/bin/bash
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

exit_on_error() {
    local message=$1
    if [[ $? != 0 ]]; then
        error_exit message
    fi
}

error_exit() {
    local message=$1
    echo "[FATAL] $message" 1>&2
    exit 1
}

AWS_METADATA_IP=169.254.169.254
REGION=`curl http://$AWS_METADATA_IP/latest/dynamic/instance-identity/document | jq -r '.region'`
AVAILABILITY_ZONE=`curl http://$AWS_METADATA_IP/latest/dynamic/instance-identity/document | jq -r '.availabilityZone'`

# Params:
# 1: Name of ebs_volume
# 2: Directory to be mounted
# 3: Device to be attached
attach_ebs() {
    EBS_NAME=$1
    MOUNT_DIR=$2
    DEVICE=$3
    MOUNT_DEVICE=$4

    INSTANCE_ID=$(curl -s http://$AWS_METADATA_IP/latest/meta-data/instance-id)
    echo INSTANCE_ID: "$INSTANCE_ID"

    IS_ALREADY_ATTACHED=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=tag:Name,Values=$EBS_NAME Name=availability-zone,Values=$AVAILABILITY_ZONE Name=attachment.instance-id,Values=$INSTANCE_ID --query 'Volumes[*].[VolumeId, State==`in-use`]' --output text | grep True | awk '{print $1}' | head -n 1)
    echo IS_ALREADY_ATTACHED:

    if [[ "$IS_ALREADY_ATTACHED" ]]; then
        echo Returning from attach_ebs
        return 0
    fi

    # getting available ebs volume-id
    EBS_VOLUME=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=tag:Name,Values=$EBS_NAME Name=availability-zone,Values=$AVAILABILITY_ZONE --query 'Volumes[*].[VolumeId, State==`available`]' --output text  | grep True | awk '{print $1}' | head -n 1)
    #check if there are available ebs volume

    if [[ -z "$EBS_VOLUME" ]]; then
        # See if the EBS volume is still attached to an instance
        ATTACHED_EBS_VOLUME=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=tag:Name,Values=$EBS_NAME Name=availability-zone,Values=$AVAILABILITY_ZONE --query 'Volumes[*].[VolumeId, State==`in-use`]' --output text | grep True | awk '{print $1}' | head -n 1)
        if [[ -n "$ATTACHED_EBS_VOLUME" ]]; then

            # detach volume if it is attached
            aws ec2 detach-volume --region $REGION --volume-id $ATTACHED_EBS_VOLUME

            sleep 10

            RETRY_LIMIT=5
            RETRIES=0
            EBS_VOLUME=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=tag:Name,Values=$EBS_NAME Name=availability-zone,Values=$AVAILABILITY_ZONE --query 'Volumes[*].[VolumeId, State==`available`]' --output text  | grep True | awk '{print $1}' | head -n 1)
            # Allow timed retries to find the now unattached volume
            while [[ -z "$EBS_VOLUME" ]]; do
                echo "retries = $RETRIES"
                RETRIES=$((RETRIES + 1))
                if [[ $RETRIES -ge $RETRY_LIMIT ]]; then
                    error_exit "Could not find available EBS within retry limit"
                fi

                echo "No available ebs volumes found, sleeping 10 seconds and retrying..."
                sleep 10
                EBS_VOLUME=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=tag:Name,Values=$EBS_NAME Name=availability-zone,Values=$AVAILABILITY_ZONE --query 'Volumes[*].[VolumeId, State==`available`]' --output text  | grep True | awk '{print $1}' | head -n 1)
            done
        else
            error_exit "could not find ebs volume"
        fi
    fi

    # attach ebs
    sudo aws ec2 attach-volume --region $REGION --volume-id $EBS_VOLUME --instance-id $INSTANCE_ID --device $DEVICE

    sleep 10
    IS_ATTACHED=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=attachment.instance-id,Values=$INSTANCE_ID Name=volume-id,Values=$EBS_VOLUME --query 'Volumes[*].[VolumeId]' --output text)

    # Allow timed retries for EBS volume to attach
    RETRIES=0
    while [[ -z "$IS_ATTACHED" ]]; do
        RETRIES=$((RETRIES + 1))
        if [[ $RETRIES -ge $RETRY_LIMIT ]]; then
            error_exit "Could not find available EBS within retry limit"
        fi

        echo "Volume not attached yet, sleeping 10 seconds and retrying..."
        sleep 10
        IS_ATTACHED=$(sudo aws ec2 describe-volumes --region $REGION --filters Name=attachment.instance-id,Values=$INSTANCE_ID Name=volume-id,Values=$EBS_VOLUME --query 'Volumes[*].[VolumeId]' --output text)
    done

    # Make the device a filesystem device if it isn't already
    DEVICE_TYPE=$(sudo file -s $MOUNT_DEVICE | awk '{print $2}')
    if [[ "$DEVICE_TYPE" == "data" ]]; then
        sudo mkfs -t ext4 $MOUNT_DEVICE
    fi

    # mount ebs to specified directory
    sudo mkdir -p $MOUNT_DIR
    sudo mount $MOUNT_DEVICE $MOUNT_DIR
    exit_on_error "Error mounting $DEVICE ($MOUNT_DEVICE) to $MOUNT_DIR"
}

PROMETHEUS_DATA_DIR="/prometheusData"
PROMETHEUS_CONFIG_DIR="/configs/prometheus"

# Attach prometheus config volume
attach_ebs "orc8r-prometheus-configs" $PROMETHEUS_CONFIG_DIR "/dev/xvdg" "/dev/nvme1n1"
sudo chmod -R 777 $PROMETHEUS_CONFIG_DIR

# Attach prometheus data volume
attach_ebs "orc8r-prometheus-data" $PROMETHEUS_DATA_DIR "/dev/xvdh" "/dev/nvme2n1"
sudo chmod -R 777 $PROMETHEUS_DATA_DIR
